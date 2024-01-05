Инициализация зависимостей
```
go mod init github.com/EvgeniyBudaev/go-websockets-gg
```

Сборка
```
go build -v ./cmd/
```

Удаление неиспользуемых зависимостей
```
go mod tidy -v
```

WebSocket
```
go get -u golang.org/x/net/websocket
```

Тестирование в консоли в браузере
```
let socket = new WebSocket("ws://localhost:3000/ws")
socket.onmessage = (event) => { console.log("received from the server: ", event.data) }
socket.send("hello from client")
```