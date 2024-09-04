from common.utils import Bet
from common.constants import MESSAGE_DELIMITER

def get_winners_payload_message(bets: list[Bet]) -> str:
    number_of_winners = len(bets)
    
    dni_list = [bet.document for bet in bets]
    dni_with_delimiter = MESSAGE_DELIMITER.join(dni_list)
    
    return f"{number_of_winners}{MESSAGE_DELIMITER}{dni_with_delimiter}"