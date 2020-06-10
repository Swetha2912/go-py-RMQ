from multiprocessing import Queue
from interfaces.IDataBroker import IDataBroker
from multiprocessing import Queue


class QueueBroker(IDataBroker):
    def __init__(self):
        self.queue = Queue()

    def get_data(self, context):
        return self.queue.get(block=True)

    def get_frame(self, context):
        return self.queue.get(block=True)

    def re_connect(self, context, data):
        pass

    def put_data(self, context, data):
        self.queue.put(data)

    def put_frame(self, context, data):
        self.queue.put(data)
