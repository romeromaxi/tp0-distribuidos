import logging
from common.courier import Courier
from common.utils import Bet
from common.utils import store_bets
from common.enums import MessageTypes
from common.constants import BET_FIELDS_NUMBER, MESSAGE_DELIMITER

class ClientHandler:
    def __init__(self, client_sock):
        self._id = None
        self._addr = client_sock.getpeername()
        self._courier = Courier(client_sock)
        
        self._is_running = True
        
    def run(self):
        while self._is_running:
            type_message = self._courier.recvTypeMessage()
        
            if type_message == MessageTypes.CONNECTION.value:
                self.__process_connection_message()
            elif type_message == MessageTypes.BET.value:
                self.__process_bet_message()
            elif type_message == MessageTypes.BETS_BATCH.value:
                self.__process_bets_batch_message()
            elif type_message == MessageTypes.END_BETS.value:
                self.__process_end_bets()
            else:
                logging.error("action: recv_message | error: Type Message not found")
        
        self.close()
        
    def close(self):
        self._courier.close()
        
    def __process_connection_message(self):
        payload_message = self._courier.recvPayloadMessage()
        
        if self._id is None:
            self._id = payload_message
        
    def __process_bet_message(self):
        payload_message = self._courier.recvPayloadMessage()
        
        fields_payload = payload_message.split(MESSAGE_DELIMITER)
        bet = Bet(self._id, fields_payload[0], fields_payload[1], fields_payload[2], fields_payload[3], fields_payload[4])
        store_bets([bet])
        
        logging.info(f"action: apuesta_almacenada | result: success | dni: {bet.document} | numero: {bet.number}")
        
        self._courier.sendResponseMessage(MessageTypes.OK_RESPONSE.value)
        
    def __process_bets_batch_message(self):
        payload_message = self._courier.recvPayloadMessage()
        
        fields_payload = payload_message.split(MESSAGE_DELIMITER)
        number_of_bets = int(fields_payload[0])
        
        try:
            bets = []
                    
            for i in range(number_of_bets):
                bet = Bet(
                    self._id,
                    fields_payload[1 + i * BET_FIELDS_NUMBER],
                    fields_payload[2 + i * BET_FIELDS_NUMBER],
                    fields_payload[3 + i * BET_FIELDS_NUMBER],
                    fields_payload[4 + i * BET_FIELDS_NUMBER],
                    fields_payload[5 + i * BET_FIELDS_NUMBER]
                )
                bets.append(bet)
            
            store_bets(bets)
            
            logging.info(f"action: apuesta_recibida | result: success | cantidad: {number_of_bets}")
            self._courier.sendResponseMessage(MessageTypes.OK_RESPONSE.value)
        
        except:
            logging.info(f"action: apuesta_recibida | result: fail | cantidad: {number_of_bets}")
            self._is_running = False
            self._courier.sendResponseMessage(MessageTypes.NO_OK_RESPONSE.value)
        
    def __process_end_bets(self):
        logging.info(f"action: fin_apuestas | result: success | client: {self._id}")
        self._is_running = False
    
        