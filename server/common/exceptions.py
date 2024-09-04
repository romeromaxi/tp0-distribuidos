class ClientConnectionClosedException(Exception):
    def __init__(self, message="The client closed the connection"):
        self.message = message
        super().__init__(self.message)