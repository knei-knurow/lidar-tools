# lidar-tools

Programs in this repository:

- **receiver**

  Enables transmitting data from [lidar-scan](https://github.com/knei-knurow/lidar-scan)
  over the network using UDP.

  `$ ./receiver --port /dev/ttyUSB0 | lidar-tx --address 192.168.1.1 --port 8080`

- **transmitter**

  Enables receiving data from [lidar-scan](https://github.com/knei-knurow/lidar-scan)
  over the network using UDP.

  `$ ./transmitter --port 8080 | lidar-vis -s`

- **servoctl**
  Controls servo rotating the axis on which lidar is mounted.

  `$ ./servoctl --port /dev/tty.usbserial-14220`

- **sync**
  Synchronizes scan data from lidar with acceleration and gyroscope data from
  MPU6050 accelerometer.

  `$ ./sync --port /dev/tty.usbserial-14220`

  TODO: add more description to `sync`
