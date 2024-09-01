import logging
from common.connectionHandler import ConnectionHandler

BYTES_SIZE_MESSAGE_TYPE = 4
BYTES_SIZE_PAYLAOD_SIZE = 4

class Courier:
    def __init__(self, client_sock):
        self._addr = client_sock.getpeername()
        self._conn = ConnectionHandler(client_sock)
        
    def sendResponseMessage(self, message_type: str):
        type_to_send = bytearray(BYTES_SIZE_MESSAGE_TYPE)
        type_to_send[:len(message_type)] = message_type.encode('utf-8')
        self._conn.send(type_to_send)
        
    def recvTypeMessage(self) -> str:
        msg_type_bytes = self._conn.recv(BYTES_SIZE_MESSAGE_TYPE)        
        return msg_type_bytes.rstrip(b'\x00').decode('utf-8')
        
    def recvPayloadMessage(self) -> str:
        size_payload_bytes = self._conn.recv(BYTES_SIZE_PAYLAOD_SIZE)
        size_payload = int.from_bytes(size_payload_bytes, byteorder='big', signed=False)
        
        payload_bytes = self._conn.recv(size_payload)
        return payload_bytes.decode('utf-8')
    
    def close(self):
        self._conn.close() 