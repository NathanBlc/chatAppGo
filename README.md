# chatAppGo

A simple go application to exchange messages between clients. All clients
connect to a server, then the user can send message through a TCP socket,
the server will broadcast the message to all other clients.

## How to use it

Launch the server using:

```bash
go run server.go <port>
```

Launch the client using:

```bash
go run client.go <address:port>
```
