import socket
import logging
import signal

from common.utils import store_bets, Bet

from common.clientHandler import ClientHandler
from common.exceptions import ClientConnectionClosedException

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        
        self._sigterm_received = False
        signal.signal(signal.SIGTERM, self.__graceful_shutdown)

    def run(self):
        while not self._sigterm_received:
            try:
                client_sock = self.__accept_new_connection()
                self.__handle_client_connection(client_sock)
            except OSError as e:
                if (not self._sigterm_received):
                    raise

    def __handle_client_connection(self, client_sock):
        try:
            client = ClientHandler(client_sock, self.__handle_store_bets)
            client.run()
        except ClientConnectionClosedException as e:
            pass
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            client.close()

    def __accept_new_connection(self):
        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c
    
    def __handle_store_bets(self, bets: list[Bet]):
        store_bets(bets)

    def __graceful_shutdown(self, signum, frame):
        self._sigterm_received = True
        self._server_socket.close()
        logging.info('action: graceful_shutdown | result: success')
        