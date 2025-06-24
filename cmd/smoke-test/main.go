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

func echo(c net.Conn) error {
	defer c.Close()

	_, err := io.Copy(c, c)
	return err
}
