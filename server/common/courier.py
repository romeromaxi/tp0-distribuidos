import logging
from common.connectionHandler import ConnectionHandler
from common.constants import BYTES_SIZE_MESSAGE_TYPE, BYTES_SIZE_PAYLAOD_SIZE

class Courier:
    def __init__(self, client_sock):
        self._addr = client_sock.getpeername()
        self._conn = ConnectionHandler(client_sock)
        
    def sendResponseMessage(self, message_type: str):
        type_to_send = bytearray(BYTES_SIZE_MESSAGE_TYPE)
        type_to_send[:len(message_type)] = message_type.encode('utf-8')
        self._conn.send(type_to_send)
        
    def sendMessageWithPayload(self, message_type: str, message_payload: str):
        self.sendResponseMessage(message_type)
        payload = message_payload.encode('utf-8')
        size_payload = len(payload)
        size_payload_bytes = size_payload.to_bytes(BYTES_SIZE_PAYLAOD_SIZE, byteorder='big', signed=False)
        final_payload = size_payload_bytes + payload
        
        self._conn.send(final_payload)
        
        
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