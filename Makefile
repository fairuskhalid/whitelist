.PHONY: all server plugin clean

default: all

all: server plugin

server: 
	CGO_ENABLED=0 go build  -o wlserver -a -installsuffix cgo ./server/server.go

plugin:
	CGO_ENABLED=0 go build -o wlplugin -a -installsuffix cgo ./broker/wl-broker.go

clean:
	rm wlserver
	rm wlplugin
