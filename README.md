Description
-----------
This repository contains code that allows geophone clients to connect to a server
on a specified port. 

Installation
------------
If you have golang installed, all that is needed is to just run 
```
go get github.com/geophones/communications-server
cd $GOPATH/src/github.com/geophones/communications-server
go build
```

Running
-------
You can run the server by calling the executable with the port you want to
run the server on as an argument, e.g.
```
./communications-server 8888
```

You can connect to the server with basically anything that speaks TCP, e.g.
netcat: 
```
nc localhost 8888 
```
