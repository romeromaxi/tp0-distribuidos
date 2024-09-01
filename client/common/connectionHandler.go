package common

import (
	"net"
)

type ConnectionHandler struct {
	conn net.Conn
}

func NewConnectionHandler(conn net.Conn) *ConnectionHandler {
	ConnectionHandler := &ConnectionHandler{
		conn: conn,
	}
	return ConnectionHandler
}

func (mc *ConnectionHandler) send(data []byte) error {
	sizeToSend := len(data)

	totalBytesSent := 0
	for totalBytesSent < sizeToSend {
		bytesSent, err := mc.conn.Write(data[totalBytesSent:])

		if err != nil {
			return err
		}

		totalBytesSent += bytesSent
	}

	return nil
}

func (mc *ConnectionHandler) recv(sizeToRecv int) ([]byte, error) {
	log.Debugf("Se desea recibir %v", sizeToRecv)
	data := make([]byte, sizeToRecv)

	totalBytesRecv := 0
	for totalBytesRecv < sizeToRecv {
		bytesRecv, err := mc.conn.Read(data[totalBytesRecv:])

		log.Debugf("Se recibe %v | %v", bytesRecv, err)

		if err != nil {
			return nil, err
		}

		totalBytesRecv += bytesRecv
	}

	return data, nil
}

func (c *ConnectionHandler) close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
