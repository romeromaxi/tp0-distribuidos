package common

import "fmt"

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

func FormatBet(name string, surname string, dni string, birth string, number string) string {
	return fmt.Sprintf("%s%s%s%s%s%s%s%s%s", name, DATA_DELIMITER, surname, DATA_DELIMITER, dni, DATA_DELIMITER, birth, DATA_DELIMITER, number)
}

func GetBetMessage(name string, surname string, dni string, birth string, number string) Message {
	betPayload := FormatBet(name, surname, dni, birth, number)

	return Message{
		Type:    MESSAGE_TYPE_BET,
		Payload: betPayload,
	}
}
