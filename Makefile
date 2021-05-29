all: receiver servoctl sync transmitter

RECEIVER := ./cmd/receiver
SERVOCTL:= ./cmd/servoctl
SYNC := ./cmd/sync
TRANSMITTER := ./cmd/transmitter

receiver: $(RECEIVER)/receiver.go
	go build $(RECEIVER)/receiver.go

servoctl: $(SERVOCTL)/servoctl.go
	go build $(SERVOCTL)/servoctl.go

sync: $(SYNC)/sync.go
	go build $(SYNC)/sync.go $(SYNC)/servo.go $(SYNC)/accelerometer.go

transmitter: $(TRANSMITTER)/transmitter.go
	go build $(TRANSMITTER)/transmitter.go $(TRANSMITTER)/cloud.go

install:
	cp ./receiver /usr/local/bin
	cp ./servoctl /usr/local/bin
	cp ./sync /usr/local/bin
	cp ./transmitter /usr/local/bin

clean:
	rm -f receiver servoctl sync transmitter
