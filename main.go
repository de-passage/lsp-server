package main

import (
  "fmt"
  "net"
)

func main() {
  startTcpServer()
}

func startTcpServer() {
  fmt.Println("Starting server...")
  listener, err := net.Listen("tcp", "localhost:8080")
  if err != nil {
    fmt.Println("Error:", err.Error())
    return
  }

  for {
    conn, err := listener.Accept()
    if err != nil {
      fmt.Println("Error on connection:", err.Error())
      return
    }

    go handleConnection(conn)
  }
}

func handleConnection(conn net.Conn) {
  defer conn.Close()

  reqBuffer := make([]byte, 4096)
  _, err := conn.Read(reqBuffer)

  if err != nil {
    fmt.Println("Error reading:", err.Error())
  }

  fmt.Println("Received request:", string(reqBuffer))

  _, err = conn.Write([]byte("Hello world"))
  if err != nil {
    fmt.Println("Error writing:", err.Error())
  }
}

