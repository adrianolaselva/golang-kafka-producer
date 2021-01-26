

### Build de producer
```sh
go build -a -ldflags="-s -w" -o kafka-producer
```

### Example producer
```sh
./kafka-producer produce -f config.yaml
```

### Example consumer with kafka CLI
```sh
kafka-console-consumer \
    --group example-group \
    --bootstrap-server kafka-1:29092 \
    --from-beginning \
    --topic events-streaming
```
