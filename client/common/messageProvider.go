package common

import (
	"fmt"
	"strconv"
	"strings"
)

type Message struct {
	Type    string
	Payload string
}

type MessageWinnerResponse struct {
	HasResult       bool
	NumberOfWinners int
	DNIWinners      []string
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

func GetBetsBatchMessage(betsInBatch [][]string) Message {
	number_of_bets := len(betsInBatch)
	strBuilder := strings.Builder{}
	number_of_bets_msg := fmt.Sprintf("%v", number_of_bets)
	strBuilder.WriteString(number_of_bets_msg)

	for _, row := range betsInBatch {
		betPayload := FormatBet(row[0], row[1], row[2], row[3], row[4])
		strBuilder.WriteString(DATA_DELIMITER)
		strBuilder.WriteString(betPayload)
	}

	return Message{
		Type:    MESSAGE_TYPE_BETS_BATCH,
		Payload: strBuilder.String(),
	}
}

func GetEndBetsMessage() Message {
	return Message{
		Type:    MESSAGE_TYPE_END_BETS,
		Payload: "",
	}
}

func GetObtainWinnersMessage() Message {
	return Message{
		Type:    MESSAGE_TYPE_GET_WINNERS,
		Payload: "",
	}
}

func GetWinnersResponseByPayload(payload string) (MessageWinnerResponse, error) {
	splitPayload := strings.Split(payload, DATA_DELIMITER)

	numberOfWinners, err := strconv.Atoi(splitPayload[0])
	if err != nil {
		return MessageWinnerResponse{}, err
	}

	winners := splitPayload[1:]

	return MessageWinnerResponse{
		HasResult:       true,
		NumberOfWinners: numberOfWinners,
		DNIWinners:      winners,
	}, nil
}
