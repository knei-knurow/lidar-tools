all: transmitter receiver

transmitter: cmd/transmitter/transmitter.go
	go build -o lidar-tx cmd/transmitter/transmitter.go

receiver: cmd/receiver/receiver.go
	go build -o lidar-rx cmd/receiver/receiver.go

install:
	cp ./lidar-tx /usr/local/bin
	cp ./lidar-rx /usr/local/bin

clean:
	rm lidar-tx
	rm lidar-rx
