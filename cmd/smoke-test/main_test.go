package main

import (
	"bufio"
	"net"
	"testing"
)

func TestEcho(t *testing.T) {
	client, server := net.Pipe()
	go echo(server)
	scanner := bufio.NewScanner(client)

	client.Write([]byte("echo\n"))
	scanner.Scan()
	assertResponseBody(t, scanner.Text(), "echo")

	client.Write([]byte("hello world\n"))
	scanner.Scan()
	assertResponseBody(t, scanner.Text(), "hello world")

	server.Close()
	client.Write([]byte("hello?\n"))
	scanner.Scan()
	got := len(scanner.Text())
	want := 0
	if got != want {
		t.Errorf("response length is wrong, got %d, want %d", got, want)
	}

	client.Close()
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}
