WMBus-Sniffer-MUC
=================

This repository contains the source code of the demonstration tools used in the Black Hat '13 presentation "Energy fraud and orchestrated blackouts: Issues with wireless metering protocols (WM-BUS)" by Cyrill Brunschwiler

Please note that this code is mainly meant as proof of concept and therefore is far away from being perfect 'nor is the little included MBus library complete.

# Usage
## Requirements

- [AMB8465-M](http://amber-wireless.de/406-1-AMB8465-M.html) (Commander)
- [AMB8465-AT](http://amber-wireless.de/415-1-AMB8465-AT.html) (Sniffer)
- >= [Go 1.0](http://golang.org/)  
- MySQL database
- Chromium / Chrome
- Linux

## Setup
### Get third-party packages

Before compiling you need to ``go get`` the used packages:

- ``go get github.com/tarm/goserial``
- ``go get github.com/ziutek/serial``
- ``go get code.google.com/p/go.net/websocket``
- ``go get github.com/go-sql-driver/mysql``

### Setup MySQL database

After getting all required third-party packages you have to setup the MySQL database, this can be done using the shell:

- ``cd WMBus-Sniffer-MUC/``
- ``cat database.sql | mysql -u YourUsername -p``

### Compile the application

- ``go build -o sniffer``
- ``chmod a+x ./sniffer``

(or just run execute.sh)

### Setup the Sniffer and Commander

You have to enable the "CMD Output" (UART Settings) and set the Baud Rate to 9600 via the [Amber Wireless ACC](http://amber-wireless.de/files/acc.zip) software.

## Execution

The application supports multiple parameters:

| Parameter     | Default        | Description                                								|
| ------------- |----------------|--------------------------------------------------------------------------|
| snifferTTY    | /dev/ttyUSB0   | Mountpoint of sniffing device (AMB8465-AT) 								|
| senderTTY     | /dev/ttyUSB1   | Mountpoint of sending device (AMB-8465-M) 								|
| DBUser        | root           | Username of the DB user                   								|
| DBPass        | root           | Username of the DB user                    								|
| DBName        | capturedFrames | Name of the database                       								|
| DemoMode		| false 		 | Insert sended frames directly into the DB (in case your sender is defect)|

e.g. ``./sniffer -snifferTTY="/dev/ttyUSB0" -senderTTY="/dev/ttyUSB1"``

The sniffer is then listening on "http://localhost:80" and the MUC on "http://localhost:8080/webui" - please be advised that this has to be executed as root as the application is using a privileged port.