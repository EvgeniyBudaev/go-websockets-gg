package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"net/http"
	"time"
)

// Server - веб сервер
type Server struct {
	// conns - карта соединений
	conns map[*websocket.Conn]bool
}

// NewServer - конструктор сервера
func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

// handleWSOrderBook - слушатель событий заказов книги
func (s *Server) handleWSOrderBook(ws *websocket.Conn) {
	fmt.Println("new incoming connection from client to order book feed:", ws.RemoteAddr())
	for {
		payload := fmt.Sprintf("orderbook data -> %d\n", time.Now().UnixNano())
		ws.Write([]byte(payload))
		time.Sleep(time.Second * 2)
	}
}

// handleWS - обработчик соединений
func (s *Server) handleWS(ws *websocket.Conn) {
	// соединение поступило от клиента
	fmt.Println("new incoming connection from client:", ws.RemoteAddr())
	// TODO: добавить mutex
	// поддержка соединения
	s.conns[ws] = true
	// считываем данные с каждого цикла
	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	// buf - создание буфера
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			// если конец файла, то прерываем соединение
			if err == io.EOF {
				break
			}
			fmt.Println("read error: ", err)
			continue
		}
		// msg - сообщение из буфера и байтов
		msg := buf[:n]
		// сообщаем всем о появлении нового сообщения
		s.broadcast(msg)
	}
}

func (s *Server) broadcast(b []byte) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("write error: ", err)
			}
		}(ws)
	}
}

func main() {
	// создаем сервер
	server := NewServer()
	// обработчик маршрутов
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.Handle("/orderbookfeed", websocket.Handler(server.handleWSOrderBook))
	// слушаем сервер
	http.ListenAndServe(":3000", nil)
}
