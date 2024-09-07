package core

type ICommandExecuter interface {
	Execute(ICommand)
}