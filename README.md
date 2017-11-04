# Installation
sudo apt-get install libzmq-dev
go get github.com/info4vincent/eventstore

# Required installation of 3th party tools:
go get github.com/gin-gonic/gin
go get github.com/boltdb/bolt/...
go get github.com/evnix/boltdbweb
cd $GOPATH/src/github.com/evnix/boltdbweb
go build boltdbweb.go
