from abc import ABC, abstractmethod


class IRequest(ABC):

    @abstractmethod
    def get_body(self):
        pass

    @abstractmethod
    def reply_back(self, context):
        pass
