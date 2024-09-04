import logging
from common.courier import Courier
from common.utils import Bet
from common.constants import MessageTypes, BET_FIELDS_NUMBER, MESSAGE_DELIMITER

class ClientHandler:
    def __init__(self, client_sock, store_bets_callback):
        self._id = None
        self._addr = client_sock.getpeername()
        self._courier = Courier(client_sock)
        
        self._is_running = True
        self._store_bets_callback = store_bets_callback
        
    def run(self):
        while self._is_running:
            type_message = self._courier.recvTypeMessage()
        
            if type_message == MessageTypes.CONNECTION.value:
                self.__process_connection_message()
            elif type_message == MessageTypes.BET.value:
                self.__process_bet_message()
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
        
        fields_payload = payload_message.split("|")
        bet = Bet(self._id, fields_payload[0], fields_payload[1], fields_payload[2], fields_payload[3], fields_payload[4])
        self._store_bets_callback([bet])
        
        logging.info(f"action: apuesta_almacenada | result: success | dni: {bet.document} | numero: {bet.number}")
        
        self._courier.sendResponseMessage(MessageTypes.OK_RESPONSE.value)
        
        self._is_running = False
    
    
        