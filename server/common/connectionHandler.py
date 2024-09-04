import socket
import logging

from common.exceptions import ClientConnectionClosedException

class ConnectionHandler:
    def __init__(self, client_sock: socket):
        self._socket = client_sock

    def send(self, data):
        size_to_send = len(data)
        total_bytes_sent = 0
        
        while total_bytes_sent < size_to_send:
            bytes_sent = self._socket.send(data[total_bytes_sent:])

            if bytes_sent == 0:
                raise ClientConnectionClosedException()

            total_bytes_sent += bytes_sent

    def recv(self, size_to_recv):
        total_bytes_recv = 0
        data = b''
        
        while total_bytes_recv < size_to_recv:
            data_chunk = self._socket.recv(size_to_recv - total_bytes_recv)
            
            if data_chunk == b'':
                raise ClientConnectionClosedException()
            
            total_bytes_recv += len(data_chunk)
            data += data_chunk
            
        return data
    
    def close(self):
        self._socket.close()
        