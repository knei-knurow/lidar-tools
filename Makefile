all: transmitter receiver

TRANSMITTER := ./cmd/transmitter
RECEIVER := ./cmd/receiver

transmitter: $(TRANSMITTER)/transmitter.go
	go build $(TRANSMITTER)/transmitter.go $(TRANSMITTER)/cloud.go

receiver: $(RECEIVER)/receiver.go
	go build $(RECEIVER)/receiver.go

install:
	cp ./lidar-tx /usr/local/bin
	cp ./lidar-rx /usr/local/bin

clean:
	rm lidar-tx
	rm lidar-rx
