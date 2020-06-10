import json


class Response:
    def __init__(self, http_code, data, msg=None, error=None, errors=None):
        self.http_code = http_code
        self.msg = msg
        self.data = data
        self.error = error
        self.errors = errors

    def toJson(self):
        return json.dumps(self, default=lambda o: o.__dict__)
