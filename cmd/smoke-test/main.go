package main

import (
	"io"
	"log"
	"net"

	"github.com/usrme/protohackers"
)

func main() {
	log.Fatal(protohackers.ListenAndAccept(5000, echo))
}

func echo(conn net.Conn) error {
	defer conn.Close()

	_, err := io.Copy(conn, conn)
	return err
}
