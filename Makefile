all: receiver servo transmitter

RECEIVER := ./cmd/receiver
SERVO := ./cmd/servo
TRANSMITTER := ./cmd/transmitter

receiver: $(RECEIVER)/receiver.go
	go build -o lidar-rx $(RECEIVER)/receiver.go

servo: $(SERVO)/servo.go
	go build -o lidar-servo $(SERVO)/servo.go

transmitter: $(TRANSMITTER)/transmitter.go
	go build -o lidar-tx $(TRANSMITTER)/transmitter.go $(TRANSMITTER)/cloud.go

install:
	cp ./lidar-rx /usr/local/bin
	cp ./lidar-servo /usr/local/bin
	cp ./lidar-tx /usr/local/bin

clean:
	rm -f lidar-rx lidar-servo lidar-tx
