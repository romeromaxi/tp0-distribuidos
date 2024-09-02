package common

import (
	"fmt"
	"net"
)

type MessageHandler struct {
	clientId string
	courier  *Courier
}

func NewMessageHandler(id string, conn net.Conn) *MessageHandler {
	MessageHandler := &MessageHandler{
		clientId: id,
		courier:  NewCourier(conn),
	}
	return MessageHandler
}

func (mh *MessageHandler) SendConnectionMessage() error {
	connection_message := GetConnectionMessage(mh.clientId)
	err := mh.courier.SendMessage(connection_message)

	if err != nil {
		log.Errorf("action: send_message | result: fail | client_id: %v | CONN error: %v",
			mh.clientId,
			err,
		)
		return err
	}

	return nil
}

func (mh *MessageHandler) SendBetsBatchMessageAndRecv(betsInBatch [][]string) error {
	number_of_bets := len(betsInBatch)
	nbets_message := GetBetsBatchMessage(betsInBatch)

	err := mh.courier.SendMessage(nbets_message)

	if err != nil {
		return err
	}

	recv_msg_type, err := mh.courier.RecvTypeMessage()

	if err != nil {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
			mh.clientId,
			err,
		)
		return err
	}

	log.Infof("action: receive_message | result: success | client_id: %v | msg: %v",
		mh.clientId,
		recv_msg_type,
	)

	if recv_msg_type != MESSAGE_TYPE_OK_RESPONSE {
		nok_reponse := fmt.Errorf("action: apuesta_enviada | result: fail | client_id: %v | cantidad: %v",
			mh.clientId,
			number_of_bets,
		)
		return nok_reponse
	}

	return nil
}

func (mh *MessageHandler) SendEndBetsMessage() error {
	end_message := GetEndBetsMessage()
	err := mh.courier.SendMessage(end_message)

	if err != nil {
		log.Errorf("action: send_message | result: fail | client_id: %v | END error: %v",
			mh.clientId,
			err,
		)
		return err
	}

	return nil
}

func (mh *MessageHandler) Close() {
	if mh.courier != nil {
		mh.courier.Close()
	}
}
