package common

import (
	"fmt"
	"strings"
)

const DATA_DELIMITER = "|"

const (
	MESSAGE_TYPE_BET        string = "BET"
	MESSAGE_TYPE_BETS_BATCH string = "NBET"
	MESSAGE_TYPE_CONNECTION string = "CONN"
	MESSAGE_TYPE_END_BETS   string = "END"

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

func FormatBet(name string, surname string, dni string, birth string, number string) string {
	return fmt.Sprintf("%s|%s|%s|%s|%s", name, surname, dni, birth, number)
}

func GetBetMessage(name string, surname string, dni string, birth string, number string) Message {
	betPayload := FormatBet(name, surname, dni, birth, number)

	return Message{
		Type:    MESSAGE_TYPE_BET,
		Payload: betPayload,
	}
}

func GetBetsBatchMessage(betsInBatch [][]string) Message {
	number_of_bets := len(betsInBatch)
	strBuilder := strings.Builder{}
	number_of_bets_msg := fmt.Sprintf("%v", number_of_bets)
	strBuilder.WriteString(number_of_bets_msg)

	for _, row := range betsInBatch {
		betPayload := FormatBet(row[0], row[1], row[2], row[3], row[4])
		strBuilder.WriteString("|")
		strBuilder.WriteString(betPayload)
	}

	return Message{
		Type:    MESSAGE_TYPE_BETS_BATCH,
		Payload: strBuilder.String(),
	}
}

func GetEndBetsMessage() Message {
	return Message{
		Type:    MESSAGE_TYPE_CONNECTION,
		Payload: "",
	}
}
