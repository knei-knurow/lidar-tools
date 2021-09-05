.PHONY: all receiver servoctl sync transmitter scandummy clean

all: receiver servoctl sync transmitter scandummy

RECEIVER := ./receiver
SERVOCTL:= ./servoctl
SYNC := ./sync
TRANSMITTER := ./transmitter
SCAN_DUMMY := ./scan-dummy

receiver: $(RECEIVER)/receiver.go
	go build $(RECEIVER)/receiver.go

servoctl: $(SERVOCTL)/servoctl.go
	go build $(SERVOCTL)/servoctl.go

sync: $(SYNC)/sync.go
	go build $(SYNC)/sync.go $(SYNC)/servo.go $(SYNC)/accelerometer.go $(SYNC)/lidar.go $(SYNC)/data-buffer.go $(SYNC)/fusion.go $(SYNC)/process.go 

transmitter: $(TRANSMITTER)/transmitter.go
	go build $(TRANSMITTER)/transmitter.go $(TRANSMITTER)/cloud.go

scandummy: $(SCAN_DUMMY)/scan-dummy.go
	go build  $(SCAN_DUMMY)/scan-dummy.go

install:
	cp ./receiver /usr/local/bin
	cp ./servoctl /usr/local/bin
	cp ./sync /usr/local/bin
	cp ./transmitter /usr/local/bin

clean:
	rm -f receiver servoctl sync transmitter scan-dummy
