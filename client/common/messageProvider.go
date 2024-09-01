package common

import "fmt"

const DATA_DELIMITER = "|"

const (
	MESSAGE_TYPE_BET        string = "BET"
	MESSAGE_TYPE_CONNECTION string = "CONN"

	MESSAGE_TYPE_OK_RESPONSE string = "OK"
)

type Message struct {
	Type    string
	Payload string
}

func GetConnectionMessage(clientId string) Message {
	return Message{
		Type:    MESSAGE_TYPE_CONNECTION,
		Payload: clientId,
	}
}

func GetBetMessage(name string, surname string, dni string, birth string, number string) Message {
	betPayload := fmt.Sprintf("%s|%s|%s|%s|%s", name, surname, dni, birth, number)

	return Message{
		Type:    MESSAGE_TYPE_BET,
		Payload: betPayload,
	}
}
