# lidar-tools

This project contains currently 3 programs:

- **lidar-tx**

  Enables transmitting data from [lidar-scan](https://github.com/knei-knurow/lidar-scan)
  over the network using UDP.

  `$ lidar-scan --port /dev/ttyUSB0 | lidar-tx --address 192.168.1.1 --port 8080`

- **lidar-rx**

  Enables receiving data from [lidar-scan](https://github.com/knei-knurow/lidar-scan)
  over the network using UDP.

  `$ lidar-rx --port 8080 | lidar-vis -s`

- **lidar-servo**

  `$ lidar-servo --port /dev/tty.usbserial-14220`
