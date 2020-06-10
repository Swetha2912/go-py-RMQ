from abc import ABC, abstractmethod


class IDataBroker(ABC):

    @abstractmethod
    def connect(self, context, data):
        pass

    @abstractmethod
    def register_endpoints(self, context):
        pass

    @abstractmethod
    def listen(self, context):
        pass
