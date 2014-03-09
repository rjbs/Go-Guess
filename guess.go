package main

import (
  "fmt"
  "math/rand"
  "net"
  "strconv"
  "strings"
  "time"
)

type Game struct {
  conn net.Conn
  ident string
}

func (game *Game) kill() {
  err := game.conn.Close()
  if err != nil {
    game.log("error closing game connection: ", err.Error())
  }
}

func (game *Game) send(s string) (err error) {
  _, err = game.conn.Write([]byte(s))
  return err
}

func (game *Game) recv() (s string, err error) {
  var b [512]byte;

  game.conn.SetReadDeadline( time.Now().Add( time.Duration(10) * time.Second ) )
  c, err := game.conn.Read(b[:])
  if err != nil {
    return "", err
  }

  return strings.TrimSpace(string(b[:c])), nil
}

func (game *Game) log(s string, a ...interface{}) {
  fmt.Printf("[%s] ", game.ident)
  fmt.Printf(s, a...)
  fmt.Println("")
}

func (game *Game) play() {
  n := rand.Intn(100) + 1

  game.log("beginning game; answer = %d", n)

  err := game.send("enter your username\n")

  username, err := game.recv()

  if err != nil {
    game.log("error reading username: %s", err.Error())
    game.kill()
    return
  }

  game.log("got username: %s", username)

  game.mainLoop(n)

  game.kill()
}

func (game *Game) mainLoop(n int) {
  err := game.send("target number is 1 <= n <= 100; enter a guess\n")

  if err != nil {
    game.log("error sending: %s", err.Error())
    return
  }

  done := false;

  for ! done {
    guess_str, err := game.recv()

    if err != nil {
      game.log("error reading guess: %s", err.Error())
      return
    }

    guess, err := strconv.Atoi( guess_str )

    if err != nil {
      game.log("error with atoi: %s", err.Error())
      continue
    }

    game.log("guess: %d", guess)

    if guess > n {
      err = game.send("Too high.\r\n")
    } else if guess < n {
      err = game.send("Too low.\r\n")
    } else {
      err = game.send("GOT IT!\r\n")
      done = true
    }

    if err != nil {
      game.log("error sending result: %s", err.Error())
    }
  }
}

func playGame(game Game) { game.play() }

func main() {
  rand.Seed( time.Now().UTC().UnixNano())

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
    go playGame( Game{ conn: conn, ident: conn.RemoteAddr().String() } )
  }
}

