package main

import (
	"bufio"
	"log"
	"net"

	"github.com/usrme/protohackers"
)

func main() {
	log.Fatal(protohackers.ListenAndAccept(5000, echo))
}

func echo(c net.Conn) error {
	defer c.Close()

	lines := bufio.NewReader(c)

	for {
		line, err := lines.ReadString('\n')
		if err != nil {
			return err
		}

		c.Write([]byte(line))
	}
}
