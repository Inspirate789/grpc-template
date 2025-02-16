Load testing:
```
rm ./data/data.db
go run ./cmd/app
k6 run -u 100 -d 1m ./test/load/write_events.js
k6 run -u 100 -d 1m ./test/load/read_events.js
```