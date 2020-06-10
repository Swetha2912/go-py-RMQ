from abc import ABC, abstractmethod

class IBroker(ABC):

    @abstractmethod
    def connect(self):
        pass

    def register_endpoints(self, context):
        pass

    def listen(self, patterns):
        pass
