FROM golang as builder
RUN go get github.com/afex/hystrix-go/hystrix
RUN go get github.com/gin-gonic/gin
RUN go get github.com/globalsign/mgo
RUN go get github.com/go-chi/chi
RUN go get github.com/gofrs/uuid
RUN go get github.com/gorilla/websocket
RUN go get github.com/streadway/amqp
RUN go get github.com/gofiber/fiber
RUN go get github.com/google/wire/cmd/wire
RUN go get github.com/TonPC64/gomon

RUN apt-get update
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y cmake git build-essential libtool autoconf automake pkg-config unzip libkrb5-dev

# install libzmq dependencies
RUN cd /tmp && git clone git://github.com/jedisct1/libsodium.git && cd libsodium && git checkout e2a30a && ./autogen.sh && ./configure && make check && make install && ldconfig
RUN cd /tmp && git clone --depth 1 git://github.com/zeromq/libzmq.git && cd libzmq && ./autogen.sh && ./configure &&  make install && ldconfig
RUN cd /tmp && git clone -b v4.3.0 https://github.com/zeromq/cppzmq.git && cd cppzmq && mkdir build && cd build && cmake .. &&  make install && ldconfig
RUN apt-get install libczmq-dev -y

RUN go get github.com/zeromq/goczmq
RUN go get github.com/joho/godotenv

# resume go-operations
COPY common /go/src/sample_rmq/common
COPY gateway /go/src/sample_rmq/gateway


WORKDIR /go/src/sample_rmq/gateway

CMD gomon
