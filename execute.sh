clear
go build -o sniffer
chmod a+x ./sniffer
./sniffer -snifferTTY="/dev/ttyUSB0" -senderTTY="/dev/ttyUSB1" #-DemoMode=true