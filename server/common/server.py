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
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """
        
        while not self._sigterm_received:
            try:
                client_sock = self.__accept_new_connection()
                client_thread = threading.Thread(target=self.__handle_client_connection, args=(client_sock,))
                client_thread.start()
                self._agencies_threads.append(client_thread)
            except OSError as e:
                if (not self._sigterm_received):
                    raise

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            client = ClientHandler(
                client_sock, self.__handle_store_bets, self.__handle_end_agency, self.__handle_get_winners
            )
            client.run()
        except ClientConnectionClosedException as e:
            pass
        except OSError as e:
            logging.critical(f"action: receive_message | result: fail | error: {e}")
        finally:
            client.close()

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c

    def __handle_end_agency(self, agency_id: int):
        self._agencies_completed_lock.acquire()
        
        if agency_id not in self._agencies_completed:
            self._agencies_completed.add(agency_id)
            
            if len(self._agencies_completed) == self._number_of_agencies:
                logging.info("action: sorteo | result: success")
            
            self._agencies_completed_lock.release()
        
            return True
        
        self._agencies_completed_lock.release()
        
        return False
    
    def __handle_store_bets(self, bets: list[Bet]):
        self._bets_files_lock.acquire()
        store_bets(bets)
        self._bets_files_lock.release()
        
    
    def __handle_get_winners(self, agency_id: int):
        self._bets_files_lock.acquire()
        if len(self._agencies_completed) != self._number_of_agencies:
            self._bets_files_lock.release()
            return (False, [])
        
        agency_winners = []
        for bet in load_bets():
            if bet.agency == agency_id and has_won(bet):
                agency_winners.append(bet)
                
        self._bets_files_lock.release()
            
        return (True, agency_winners)
    
    def __graceful_shutdown(self, signum, frame):
        self._sigterm_received = True
        self._server_socket.close()
        logging.info('action: graceful_shutdown | result: success')
        