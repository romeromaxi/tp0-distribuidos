package common

import (
	"encoding/binary"
	"net"
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

func (c *Courier) SendMessage(message_type string, message string) error {
	typeToSend := make([]byte, BYTES_SIZE_MESSAGE_TYPE)
	copy(typeToSend, []byte(message_type))

	log.Debugf("action: send_message_type | result: in_process | msg: %v",
		typeToSend,
	)

	err := c.conn.send(typeToSend)

	log.Debugf("action: send_message_type | result: success | msg: %v",
		typeToSend,
	)

	if err != nil {
		log.Debugf("action: send_message_type | result: failure | msg: %v",
			typeToSend,
		)
		return err
	}

	payload := []byte(message)
	sizePayload := len(payload)
	sizePayloadBytes := make([]byte, BYTES_SIZE_PAYLAOD_SIZE)
	binary.BigEndian.PutUint32(sizePayloadBytes, uint32(sizePayload))

	finalPayload := append(sizePayloadBytes, payload...)

	log.Debugf("action: send_message | result: in_process | msg: %v",
		finalPayload,
	)

	err = c.conn.send(finalPayload)

	if err != nil {
		log.Debugf("action: send_message | result: failure | msg: %v",
			finalPayload,
		)
		return err
	}

	log.Debugf("action: send_message | result: success | msg: %v",
		finalPayload,
	)

	return nil
}

func (c *Courier) RecvTypeMessage() (string, error) {
	msgTypeBytes, err := c.conn.recv(BYTES_SIZE_MESSAGE_TYPE)
	if err != nil {
		return "", err
	}

	return string(msgTypeBytes), nil
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
