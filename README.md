This is example of trace propagation over 2 services

To check the system out, launch jaeger

```bash
podman run --rm --name jaeger   -e COLLECTOR_ZIPKIN_HOST_PORT=:9411   -p 5775:5775/udp   -p 6831:6831/udp   -p 6832:6832/udp   -p 5778:5778   -p 16686:16686   -p 14268:14268   -p 14250:14250   -p 9411:9411   jaegertracing/all-in-one:1.22
```

Then run both services

```bash
go run consumer/main.go
```

```bash
go run producer/main.go
```

And fire up a GET request to consumer service

```bash
http localhost:1323
```

Go to [jaerer ui](http://localhost:16686) and search for traces
