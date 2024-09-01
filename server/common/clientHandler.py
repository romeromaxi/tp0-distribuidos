import logging
from common.courier import Courier

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
            logging.info(f'Type Message: {type_message}')
        
            if type_message == "CONN":
                self.__process_connection_message()
            elif type_message == "BET":
                self.__process_bet_message()
            else:
                logging.debug("Type Message not found")
        
        self._courier.close()
        
    def __process_connection_message(self):
        payload_message = self._courier.recvPayloadMessage()
        logging.info(f'Payload: {payload_message}')
        
        if self._id is None:
            self._id = payload_message
        
    def __process_bet_message(self):
        payload_message = self._courier.recvPayloadMessage()
        logging.info(f'Payload: {payload_message} | clienteId {self._id}')
        
        self._courier.sendResponseMessage("OK")
        
        self._is_running = False
    
    
        