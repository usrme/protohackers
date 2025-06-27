package main

import (
	"bufio"
	"encoding/json"
	"log"
	"math/big"
	"net"

	"github.com/usrme/protohackers"
)

func main() {
	log.Print(protohackers.ListenAndAccept(5000, handler))
}

type request struct {
	Method string   `json:"method"`
	Number *float64 `json:"number"`
}

type response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func isPrime(f float64) bool {
	if float64(int(f)) == f {
		return big.NewInt(int64(f)).ProbablyPrime(0)
	}
	return false
}

func isMalformed(req request) bool {
	if req.Number == nil {
		log.Println("no number, malformed request")
		return true
	}

	if req.Method != "isPrime" {
		return true
	}

	return false
}

func handler(conn net.Conn) error {
	defer conn.Close()

	log.Println("connection from:", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		var req request
		if err := json.Unmarshal([]byte(scanner.Text()), &req); err != nil {
			log.Println("malformed request", err)
			conn.Write([]byte("malformed request"))
			conn.Close()
			continue
		}

		if isMalformed(req) {
			log.Println("malformed request")
			conn.Write([]byte("malformed request"))
			conn.Close()
			continue
		}

		res := response{
			Method: "isPrime",
			Prime:  isPrime(*req.Number),
		}
		enc := json.NewEncoder(conn)
		if err := enc.Encode(res); err != nil {
			log.Printf("connection read: %v\n", err)
			continue
		}
	}
	return nil
}
