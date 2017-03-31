## CLI

```
go build
./follow_along -follower lukewendling -followee potus44

```

## Redis jobs integration

```
cd svc
go run main.go
# hmset 1 id <screename> state new
# lpush genie:followfinder 1
# ... wait ...
# hgetall 1
# check for error, state field
```