# Sockunx

Golang unix-socket wrapper

## Server

### Running server

```go
server, e := sockunx.NewServer("/path/to/your/socks.sock", 512)
if e != nil {
    log.Fatal(e)
}
defer func() {
    log.Println("Shutting down...")
    server.Stop()
}()

// server.RegisterHandler(handler.Index) // just in case you need to implement handler
e = server.Run()
if e != nil {
    log.Fatal("Error while running server", e.Error())
}
```

### Handler

```go
var Index sockunx.Handler = func(request string) (response interface{}, e error) {
		return strings.ToUpper(request), nil
}
```

### Client

```go
client, e := sockunx.NewClient(*socketPath)
if e != nil {
    log.Fatal(e)
}

for i := 0; i < 10; i++ {
    response, e := client.Send(`{"id":"one","from":0,"to":15,"fizz":"zzif","buzz":"zzub"}\n`)
    if e != nil {
        log.Println("error : ", e.Error())
        continue
    }
    fmt.Println(">>>", response)
}
```
