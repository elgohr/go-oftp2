package main

import (
	"bufio"
	"github.com/elgohr/go-oftp2/oftp2"
	"github.com/stretchr/testify/require"
	"net"
	"os"
	"testing"
)

func TestListener(t *testing.T) {
	c := make(chan os.Signal, 1)
	p, err := NewListener(c)
	require.NoError(t, err)
	go p.Listen()

	addr, err := net.ResolveTCPAddr("tcp", "localhost:3305")
	require.NoError(t, err)
	conn, err := net.DialTCP("tcp", nil, addr)
	require.NoError(t, err)

	reader := bufio.NewReader(conn)

	t.Run("<-SSRM-", func(t *testing.T) {
		con, err := reader.ReadBytes('\r')
		sth := "\x10\x00\x00\x17" // https://tools.ietf.org/html/rfc5024#section-8.1
		require.NoError(t, err)
		require.Equal(t, sth+"IODETTE FTP READY "+oftp2.CarriageReturn, string(con))
	})

	t.Run("-SSID->", func(t *testing.T) {
		_, err = conn.Write([]byte("test\n"))
		require.NoError(t, err)
	})

	t.Run("<-SSID-", func(t *testing.T) {
		con, err := reader.ReadBytes('\r')
		sth := "\x10\x00\x00\x17" // https://tools.ietf.org/html/rfc5024#section-8.1
		require.NoError(t, err)
		require.Equal(t, sth+"IODETTE FTP READY "+oftp2.CarriageReturn, string(con))
	})

}
