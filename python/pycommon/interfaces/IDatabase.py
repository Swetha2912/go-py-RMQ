from abc import ABC, abstractmethod


class IDatabase(ABC):

    @abstractmethod
    def get_collection(self, name):
        pass

    def connect(self, context):
        pass
