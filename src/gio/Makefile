
all: gio gioc

gio: $(wildcard *.go) message.pb.go
	go install gio

gioc: $(wildcard gioc/*.go) message.pb.go
	go install gio/gioc

%.pb.go: %.proto
	protoc --go_out=. $<

test: gio
	gio
