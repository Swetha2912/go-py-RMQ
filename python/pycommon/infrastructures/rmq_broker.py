from pycommon.interfaces.IDataBroker import IDataBroker
import pika
import json
import sys
import traceback

from pycommon.infrastructures.rmq_request import RMQRequest


class RMQBroker(IDataBroker):
    def __init__(self, rmq_host, rmq_user, rmq_pass, rmq_port, default_exchange):
        self.rmq_host = rmq_host
        self.rmq_user = rmq_user
        self.rmq_pass = rmq_pass
        self.rmq_port = rmq_port
        self.default_exchange = default_exchange
        self._credentials = pika.PlainCredentials(self.rmq_user, self.rmq_pass)
        self._parameters = pika.ConnectionParameters(self.rmq_host, self.rmq_port, '/', self._credentials)
        self.rmq_conn = pika.BlockingConnection(self._parameters)
        self.rmq_channel = self.rmq_conn.channel()
        self.rmq_conn.close()
        self.endpoint_handlers = None

    def connect(self):
        self.rmq_conn = pika.BlockingConnection(self._parameters)
        self.rmq_channel = self.rmq_conn.channel()
        self.rmq_channel.exchange_declare(exchange=self.default_exchange, exchange_type="topic")

    def register_endpoints(self, endpoint_handlers):
        self.endpoint_handlers = endpoint_handlers
        pass

    def listen(self, patterns):
        result = self.rmq_channel.queue_declare(queue='', exclusive=True)
        for pattern in patterns:
            self.rmq_channel.queue_bind(exchange=self.default_exchange, queue=result.method.queue, routing_key=pattern)
        self.rmq_channel.basic_consume(queue=result.method.queue, on_message_callback=self.listener, auto_ack=True)

        try:
            print("listening ..")
            self.rmq_channel.start_consuming()
        except Exception as e:
            traceback.print_exc(file=sys.stdout)

    def listener(self, ch, method, props, req_body):
        for pattern in self.endpoint_handlers:
            if method.routing_key == pattern:
                handler = self.endpoint_handlers[pattern]
                req = RMQRequest(rmq_channel=self.rmq_channel, rmq_connection=self.rmq_conn,
                                 correlation_id=props.correlation_id, matched_route=pattern, body=req_body,
                                 reply_to=props.reply_to)
                handler(req)
                return

        print("invalid routing_key for the request : ", method.routing_key)
