clear
go build 
chmod a+x ./sniffer
./sniffer -snifferTTY="/dev/ttyUSB0" -senderTTY="/dev/ttyUSB1" 