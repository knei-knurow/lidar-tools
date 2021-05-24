all: receiver servo transmitter

RECEIVER := ./cmd/receiver
SERVO := ./cmd/servo
TRANSMITTER := ./cmd/transmitter

receiver: $(RECEIVER)/receiver.go
	go build $(RECEIVER)/receiver.go

servo: $(SERVO)/servo.go
	go build $(SERVO)/servo.go

transmitter: $(TRANSMITTER)/transmitter.go
	go build  $(TRANSMITTER)/transmitter.go $(TRANSMITTER)/cloud.go

install:
	cp ./lidar-rx /usr/local/bin
	cp ./lidar-servo /usr/local/bin
	cp ./lidar-tx /usr/local/bin

clean:
	rm -f lidar-rx lidar-servo lidar-tx
