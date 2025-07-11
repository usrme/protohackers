package protohackers

import (
	"fmt"
	"net"
)

type ConnHandler func(net.Conn) error

func ListenAndAccept(port int, handler ConnHandler) error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("can't listen on %d/tcp: %s", port, err)
	}

	fmt.Println("listening on port", port)

	for {
		conn, err := l.Accept()
		if err != nil {
			return fmt.Errorf("accept: %w", err)
		}

		fmt.Println("connection from", conn.RemoteAddr())

		go handler(conn)
	}
}
