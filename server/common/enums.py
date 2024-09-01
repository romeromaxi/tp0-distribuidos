from enum import Enum

class MessageTypes(Enum):
    BET = "BET"
    BETS_BATCH = "NBET"
    CONNECTION = "CONN"
    END_BETS = "END"
    
    OK_RESPONSE = "OK"
    NO_OK_RESPONSE = "NOK"