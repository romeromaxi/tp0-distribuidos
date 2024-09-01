package common

import (
	"encoding/binary"
	"net"
	"strings"
)

const BYTES_SIZE_MESSAGE_TYPE = 4
const BYTES_SIZE_PAYLAOD_SIZE = 4

type Courier struct {
	conn *ConnectionHandler
}

func NewCourier(id string, conn net.Conn) *Courier {
	Courier := &Courier{
		conn: NewConnectionHandler(conn),
	}
	return Courier
}

func (c *Courier) SendTypeMessage(message_type string) error {
	typeToSend := make([]byte, BYTES_SIZE_MESSAGE_TYPE)
	copy(typeToSend, []byte(message_type))

	err := c.conn.send(typeToSend)

	if err != nil {
		log.Debugf("action: send_message_type | result: failure | msg: %v",
			message_type,
		)
		return err
	}

	return nil
}

func (c *Courier) SendMessage(message_type string, message string) error {
	err := c.SendTypeMessage(message_type)

	if err != nil {
		return err
	}

	payload := []byte(message)
	sizePayload := len(payload)
	sizePayloadBytes := make([]byte, BYTES_SIZE_PAYLAOD_SIZE)
	binary.BigEndian.PutUint32(sizePayloadBytes, uint32(sizePayload))

	finalPayload := append(sizePayloadBytes, payload...)

	err = c.conn.send(finalPayload)

	if err != nil {
		log.Debugf("action: send_message | result: failure | msg: %v",
			payload,
		)
		return err
	}

	return nil
}

func (c *Courier) RecvTypeMessage() (string, error) {
	msgTypeBytes, err := c.conn.recv(BYTES_SIZE_MESSAGE_TYPE)

	if err != nil {
		return "", err
	}

	msgType := string(msgTypeBytes)

	return strings.TrimRight(msgType, "\x00"), nil
}

func (c *Courier) RecvPayloadMessage() (string, error) {
	sizePayloadBytes, err := c.conn.recv(BYTES_SIZE_PAYLAOD_SIZE)
	if err != nil {
		return "", err
	}

	sizePayload := int(binary.BigEndian.Uint32(sizePayloadBytes))

	payloadBytes, err := c.conn.recv(sizePayload)

	if err != nil {
		return "", err
	}

	return string(payloadBytes), nil
}

func (c *Courier) Close() {
	if c.conn != nil {
		c.conn.close()
	}
}
