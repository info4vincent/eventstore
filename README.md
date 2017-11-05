# Installation
sudo apt-get install libzmq-dev
go get github.com/info4vincent/eventstore

# Required installation of 3th party tools:
go get github.com/gin-gonic/gin
go get github.com/boltdb/bolt/...
go get github.com/evnix/boltdbweb
cd $GOPATH/src/github.com/evnix/boltdbweb
go build boltdbweb.go

## libsodium
wget https://github.com/jedisct1/libsodium/releases/download/1.0.3/libsodium-1.0.3.tar.gz
tar -zxvf libsodium-1.0.3.tar.gz
cd libsodium-1.0.3/
./configure
make
sudo make install 
 
## ZeroMQ 4.1.3
wget http://download.zeromq.org/zeromq-4.1.3.tar.gz
tar -zxvf zeromq-4.1.3.tar.gz
cd zeromq-4.1.3/
./configure
make
sudo make install
sudo ldconfig

# To browse what is stored in the db..
boltdbweb --db-name=myevents.db --static-path=/home/$USER/go/bin/
