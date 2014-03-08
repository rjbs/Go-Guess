package main

import (
    "fmt"
    "math/rand"
    "net"
    "strconv"
    "strings"
    "time"
)

func playGame(conn net.Conn) {
    var b [512]byte;
    n := rand.Intn(100) + 1

    _, err := conn.Write([]byte("enter your username\n"))

    conn.SetReadDeadline( time.Now().Add( time.Duration(10) * time.Second ) )
    c, err := conn.Read(b[:])

    if err != nil {
        fmt.Println("error reading: ", err.Error())
        err = conn.Close()
        if err != nil {
          fmt.Println("Fatal error closing: ", err.Error())
        }
        return
    }

    username := string(b[:c])
    username = strings.TrimSpace(username)

    fmt.Println("bytes read: ", c)
    fmt.Println("username: ", username)

    if err != nil {
        fmt.Println("Fatal error writing: ", err.Error())
    }

    _, err = conn.Write([]byte("target number is 1 <= n <= 100; enter a guess\n"))

    if err != nil {
        fmt.Println("Fatal error writing: ", err.Error())
    }

    done := false;

    for ! done {
      conn.SetReadDeadline( time.Now().Add( time.Duration(10) * time.Second ) )
      c, err = conn.Read(b[:])
      if err != nil {
          fmt.Println("Fatal error reading guess: ", err.Error())
          done = true
          continue
      }

      guess, err := strconv.Atoi( strings.TrimSpace(string(b[:c])) )

      if err != nil {
          fmt.Println("Error with atoi: ", err.Error())
          continue
      }

      fmt.Println("guess:", guess)

      if guess > n {
        _, err = conn.Write([]byte("Too high.\r\n"))
      } else if guess < n {
        _, err = conn.Write([]byte("Too low.\r\n"))
      } else {
        _, err = conn.Write([]byte("GOT IT!\r\n"))
        done = true
      }
    }

    err = conn.Close()
    if err != nil {
        fmt.Println("Fatal error closing: ", err.Error())
    }
}

func main() {
    fmt.Println("start");
    ln, err := net.Listen("tcp", ":8080")

    if err != nil {
        fmt.Println("Couldn't listen: ", err.Error())
    }

    for {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println("Couldn't listen: ", err.Error())
            continue
        }
        go playGame(conn)
    }
}

