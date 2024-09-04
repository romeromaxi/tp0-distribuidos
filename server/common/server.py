import socket
import logging
import signal
import threading

from common.utils import load_bets, store_bets, has_won, Bet

from common.clientHandler import ClientHandler
from common.exceptions import ClientConnectionClosedException

class Server:
    def __init__(self, port, listen_backlog, number_of_agencies):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        
        self._sigterm_received = False
        signal.signal(signal.SIGTERM, self.__graceful_shutdown)
        
        self._number_of_agencies = number_of_agencies
        self._agencies_completed = set()
        
        self._agencies_threads = []
        self._bets_files_lock = threading.Lock()
        self._agencies_completed_lock = threading.Lock()

    def run(self):
        while not self._sigterm_received:
            try:
                client_sock = self.__accept_new_connection()
                client_thread = threading.Thread(target=self.__handle_client_connection, args=(client_sock,))
                client_thread.start()
                self._agencies_threads.append((client_thread, client_sock))
            except OSError as e:
                if (not self._sigterm_received):
                    raise

    def __handle_client_connection(self, client_sock):
        try:
            client = ClientHandler(
                client_sock, self.__handle_store_bets, self.__handle_end_agency, self.__handle_get_winners
            )
            client.run()
        except ClientConnectionClosedException as e:
            pass
        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
        finally:
            client.close()

    def __accept_new_connection(self):
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c

    def __handle_end_agency(self, agency_id: int):
        with self._agencies_completed_lock:
            if agency_id not in self._agencies_completed:
                self._agencies_completed.add(agency_id)
                
                if len(self._agencies_completed) == self._number_of_agencies:
                    logging.info("action: sorteo | result: success")
                    
                return True
            return False
    
    def __handle_store_bets(self, bets: list[Bet]):
        with self._bets_files_lock:
            store_bets(bets)
    
    def __handle_get_winners(self, agency_id: int):
        with self._agencies_completed_lock:
            if len(self._agencies_completed) != self._number_of_agencies:
                return (False, [])
        
        agency_winners = []
        with self._bets_files_lock:
            for bet in load_bets():
                if bet.agency == agency_id and has_won(bet):
                    agency_winners.append(bet)
                
        return (True, agency_winners)
    
    def __graceful_shutdown(self, signum, frame): 
        for agency_thread,client_handler in self._agencies_threads:
            client_handler.close() 
            agency_thread.join() 
        
        self._sigterm_received = True
        self._server_socket.close()
        logging.info('action: graceful_shutdown | result: success')
        