# Raspberry Pi Build Light Indicator

## Cross Compiling on Windows

### GNU Toolchain

http://gnutoolchains.com/raspberry/tutorial/

### Go Environment Variables

GOOS=linux
GOARM=7
GOARCH=arm
CC=arm-linux-gnueabihf-gcc
CGO_ENABLED=1

### HID ( USB ) read/write access for non root users ( in my case for user pi on an raspberry pi 2 running nodered )

Taken from the following link (thanks @somebuddy87):
https://flows.nodered.org/node/node-red-contrib-usbhid

The Pd-extended [hid] object allows you to access Human Interface Devices such as mice, keyboards, and joysticks. However, in most Linux distributions, these devices are setup to where they cannot be read/written directly by Pd unless you run it as root.

Running a non-system process as root is considered a security risk, so an alternative is to change the permissions of the input devices so that pd can read/write them.

```
sudo mkdir -p /etc/udev/rules.d
sudo nano /etc/udev/rules.d/85-pure-data.rules
```
Now add the following rules to /etc/udev/rules.d/85-pure-data.rules:

```
SUBSYSTEM=="usb", GROUP="input", MODE="777"
```

Then create an "input" group and add yourself to it:

```
sudo groupadd -f input
sudo gpasswd -a YOURUSERNAME input
```
Reboot your machine for the rules to take effect.

Your nodejs / nodered has now FULL ACCESS !! to you usb devides. Feel free to adjust the permissions to fit your needs.
