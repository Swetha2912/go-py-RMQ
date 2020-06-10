from pycommon.interfaces.IRequest import IRequest
import pika
import json
import jsonpickle


class RMQRequest(IRequest):
    def __init__(self, rmq_connection, rmq_channel, matched_route, correlation_id, body, reply_to):
        self.rmq_connection = rmq_connection
        self.rmq_channel = rmq_channel
        self.payload_body = None
        self.matched_route = matched_route
        self.correlation_id = correlation_id
        self.body = body
        self.reply_to = reply_to

    def get_body(self):
        return self.body

    def reply_back(self, data):
        self.rmq_channel.basic_publish(exchange="", routing_key=str(self.reply_to),
                                       properties=pika.BasicProperties(correlation_id=self.correlation_id),
                                       body=data)
