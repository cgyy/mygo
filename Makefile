
all: mypw gomemcached golint

mypw:
	go install $@

gomemcached:
	go install $@

golint:
	go install github.com/golang/lint/golint
