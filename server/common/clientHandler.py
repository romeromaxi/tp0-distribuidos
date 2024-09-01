import logging
from common.courier import Courier
from common.utils import Bet
from common.utils import store_bets
from common.enums import MessageTypes

BYTES_SIZE_MESSAGE_TYPE = 4
BYTES_SIZE_PAYLAOD_SIZE = 4

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
            else:
                logging.error("action: recv_message | error: Type Message not found")
        
        self._courier.close()
        
    def __process_connection_message(self):
        payload_message = self._courier.recvPayloadMessage()
        
        if self._id is None:
            self._id = payload_message
        
    def __process_bet_message(self):
        payload_message = self._courier.recvPayloadMessage()
        
        fields_payload = payload_message.split("|")
        bet = Bet(self._id, fields_payload[0], fields_payload[1], fields_payload[2], fields_payload[3], fields_payload[4])
        store_bets([bet])
        
        logging.info(f"action: apuesta_almacenada | result: success | dni: {bet.document} | numero: {bet.number}")
        
        self._courier.sendResponseMessage(MessageTypes.OK_RESPONSE.value)
        
        self._is_running = False
    
    
        