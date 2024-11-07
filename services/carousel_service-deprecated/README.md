# TODO
Implement HB handler with further chnage of carousel status to offlin if HB is missed
Implement Message convertation between broker and carousel parts

# Environment

Runc carusel simulater. It replies with Ack event
```sh
  cd clousel\carousel_simulator
  go run .\cmd\carousel_simulator\carousel_simulator.go
```

Run hartbeat script preventing 'offline' status
```sh
  cd clousel/services/carousel_service/scripts/mqtt
  bash hb.sh
```
