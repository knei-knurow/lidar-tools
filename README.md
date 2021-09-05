# lidar-tools

This repository is a part of a bigger *lidar* project and contains some tools which allow to control our lidar setup.

## Programs

### sync

The main program which synchronizes lidar and AVR controlling accelerometer with gyroscope and servo. Computed point clouds are sent via stdout.

**Output:**

```
...
X Y Z
X Y Z
...
```

where *X, Y, Z* are floating point numbers representing single cartesian points of scanned point cloud. 

### servoctl

Controls servo rotating the axis on which lidar is mounted.

`$ ./servoctl --port /dev/tty.usbserial-14220`

### receiver

Enables transmitting data from [lidar-scan](https://github.com/knei-knurow/lidar-scan) over the network using UDP.

`$ ./receiver --port /dev/ttyUSB0 | lidar-tx --address 192.168.1.1 --port 8080`

### transmitter

Enables receiving data from [lidar-scan](https://github.com/knei-knurow/lidar-scan) over the network using UDP.

`$ ./transmitter --port 8080 | lidar-vis -s`

### scan-dummy

- Genereate dummy data to imitate the original lidar-scan output.

  `$ ./scan-dummy`

  Time differences are represended as milliseconds (like the original one). `stout` is used for data output. More detaled info in [lidar-scan repository](https://github.com/knei-knurow/lidar-scan#point-cloud-output). 

  Example output structure:

  ```
  # a comment
  ! 0 0
  120  100
  240  100
  360  100
  ! 1 500
  120  200
  240  200
  360  200
  ! 2 500
  120  300
  240  300
  360  300
  ! 3 500
  ```
