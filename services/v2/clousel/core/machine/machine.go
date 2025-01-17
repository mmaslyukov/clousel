package machine

import (
	"clousel/lib/fault"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Machine struct {
	// cfg    IBusinessConfigAdapter
	repoGame  IMachineRepoGameAdapter
	repoGen   IMachineRepoGeneralAdapter
	client    IMachineClientAdapter
	log       *zerolog.Logger
	barouness IMachineBarounessAdapter
}

func MachineCreate(
	repoGame IMachineRepoGameAdapter,
	repoGen IMachineRepoGeneralAdapter,
	client IMachineClientAdapter,
	cfg IMachineConfigAdapter,
	ipc IMachineIpcAdapter,
	log *zerolog.Logger,
) *Machine {
	b := BarounessCreatePartial(log, cfg, ipc)
	m := &Machine{
		repoGame:  repoGame,
		repoGen:   repoGen,
		client:    client,
		log:       log,
		barouness: b,
	}

	b.InjectAndFinish(m).Init()
	go b.Run()
	return m
}

func (m *Machine) SaveNewMachineEntry(machId uuid.UUID, companyId uuid.UUID, gameCost int) fault.IError {
	return m.repoGen.SaveNewMachineEntry(machId, companyId, gameCost)
}

func (m *Machine) ReadMachineEntriesBySelector(selector IMachineSelector) (entries []*MachineEntry, err fault.IError) {
	if selector.MachId() != nil {
		if entry, e := m.repoGen.ReadMachineById(*selector.MachId()); e == nil {
			entries = append(entries, entry)
		} else {
			err = fault.New(EMachineWrongSelector).Msgf("%s", e.Error())
		}
	} else if selector.CompanyId() != nil && selector.Status() != nil {
		if pentries, e := m.repoGen.ReadMachinesByStatus(*selector.CompanyId(), *selector.Status()); e == nil {
			entries = pentries
		} else {
			err = fault.New(EMachineWrongSelector).Msgf("%s", e.Error())
		}
	} else if selector.CompanyId() != nil {
		if pentries, e := m.repoGen.ReadMachinesByCompanyId(*selector.CompanyId()); e == nil {
			entries = pentries
		} else {
			err = fault.New(EMachineWrongSelector).Msgf("%s", e.Error())
		}
	}
	return entries, err
}

func (m *Machine) ChangeGameCost(machId uuid.UUID, gameCost int) fault.IError {
	return m.repoGen.UpdateMachineCost(machId, gameCost)
}

func (m *Machine) ChangeSubscriptionFee(machId uuid.UUID, fee int) fault.IError {
	return m.repoGen.UpdateFee(machId, fee)
}

func (m *Machine) PlayRequest(machId uuid.UUID, userId uuid.UUID) (eventId *uuid.UUID, err fault.IError) {
	const fn = "Core.Machine.PlayRequest"
	var entryMachine *MachineEntry
	gameId := uuid.New()

	for ok := true; ok; ok = false {
		if entryMachine, err = m.repoGen.ReadMachineById(machId); err != nil {
			m.log.Error().Msgf("%s: Fail to read machine with id: '%s'", fn, machId)
			break
		}

		if entryMachine.Status != MachineStatusOnline {
			err = fault.New(EMachineStatus).Msgf("Machine '%s' is not '%s', but it is '%s'", machId, MachineStatusOnline, entryMachine.Status)
			m.log.Error().Msgf("%s: Fail to read machine with id: '%s'", fn, machId)
			break
		}

		if games, e := m.repoGame.ReadGamesByMachineIdWithStatus(machId, GameStatusPending); e != nil || len(games) != 0 {
			if e != nil {
				m.log.Error().Msgf("%s: Fail to read pending games for machine: '%s'", fn, machId)
			} else if len(games) > 0 {
				m.log.Warn().Msgf("%s: Already has a pending games for machine: '%s'", fn, machId)
				eventId = &games[0].EventId
			}
			break
		}
		var can bool
		if can, err = m.client.IsCanPay(userId, entryMachine.GameCost); err != nil {
			m.log.Error().Msgf("%s: Fail to check balance for user '%s'", fn, userId)
			err = fault.New(EMachineReadUser).Err(err)
			break
		}
		if !can {
			err = fault.New(EMachineNoTickets).Msgf("Not enough tickets")
			m.log.Error().Msgf("%s: User '%s' %s", fn, userId, err.Error())
			break
		}

		if err = m.repoGame.SaveNewGameEntry(gameId, machId, userId, entryMachine.GameCost); err != nil {
			m.log.Error().Msgf("%s: Fail to save new game entry for machine: '%s'", fn, machId)
			break
		}

		eventId = &gameId
		m.barouness.SendGameRequestEvent(GameEvent{
			EventId: gameId,
			MachId:  machId,
		})
	}

	return eventId, err
}

func (m *Machine) PollRequestStatus(eventId uuid.UUID) (GameStatus, fault.IError) {
	if entry, err := m.repoGame.ReadGamesByGameId(eventId); err == nil {
		return entry.Status, err
	} else {
		return GameStatusIvalid, err
	}
}

func (m *Machine) GameStartConfirm(gameId uuid.UUID) (err fault.IError) {
	for ok := true; ok; ok = false {
		var game *GameEntry
		if game, err = m.repoGame.ReadGamesByGameId(gameId); err != nil {
			break
		}
		if err = m.client.ApplyGameCost(game.EventId, game.UserId, game.GameCost); err != nil {
			break
		}
		if err = m.repoGame.UpdateGameStatus(gameId, GameStatusCompleted); err != nil {
			break
		}
	}
	return err
}

func (m *Machine) GameStartFailed(gameId uuid.UUID) fault.IError {
	return m.repoGame.UpdateGameStatus(gameId, GameStatusFailed)
}

func (m *Machine) MachineUpdateStatus(machId uuid.UUID, status MachineStatus) fault.IError {
	return m.repoGen.UpdateMachineStatus(machId, status)
}

func (m *Machine) ReadPendingGames() (gameEvents []*GameEvent, err fault.IError) {
	var gamePending []*GameEntry
	if gamePending, err = m.repoGame.ReadGamesByStatus(GameStatusPending); err == nil {
		for _, p := range gamePending {
			gameEvents = append(gameEvents, p.toGameEvent())
		}
	}
	return gameEvents, err
}
func (m *Machine) ReadExpiredGames(ts time.Time) (gameEvents []*GameEvent, err fault.IError) {
	var gamePending []*GameEntry
	if gamePending, err = m.repoGame.ReadGamesByStatusAndTime(GameStatusPending, ts); err == nil {
		for _, p := range gamePending {
			gameEvents = append(gameEvents, p.toGameEvent())
		}
	}
	return gameEvents, err
}
