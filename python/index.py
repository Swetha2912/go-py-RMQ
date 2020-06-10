import json
import os
from pycommon.infrastructures.rmq_broker import RMQBroker
from pycommon.interfaces.IRequest import IRequest
from pycommon.interfaces.IResponse import Response

print("i am in main")

def sample(data):
    body = data.get_body()
    response_bytes = bytes("hello world", 'utf-8')
    data.reply_back(response_bytes)

if __name__ == '__main__':
    rmqBroker = RMQBroker(rmq_host=os.getenv("RMQ_HOST"), rmq_user=os.getenv("RMQ_USER"), rmq_pass=os.getenv("RMQ_PASS"), rmq_port=os.getenv("PORT"),
                          default_exchange='sample_rmq')
    rmqBroker.connect()
    rmqBroker.register_endpoints({'sample.data': sample})
    rmqBroker.listen(['sample.*'])
