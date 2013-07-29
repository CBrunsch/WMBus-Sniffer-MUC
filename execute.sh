clear
go build 
chmod a+x ./Downloads
./Downloads -snifferTTY="/dev/ttyUSB0" -senderTTY="/dev/ttyUSB1" #-DemoMode=true