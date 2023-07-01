package utils

import (
	"net"
)

func ConnWrite(conn net.Conn, msg string) error {
	_, err := conn.Write([]byte(msg))
	return err
}
