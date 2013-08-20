# qr_server2

A QR Server with groupcache.

## Test

Start two instances:

```
go build && ./qr_server2
go build && ./qr_server2 -port 8002 -addr :1719
```

Run tests:

```
./t.sh
```
