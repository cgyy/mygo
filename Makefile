
all: pwgen gomemcached

pwgen:
	go install $@

gomemcached:
	go install $@
