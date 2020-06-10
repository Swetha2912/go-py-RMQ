FROM python

RUN python3 -m pip install --upgrade pip

#update and install tooling
RUN apt-get update\
    && apt-get install -y git \
    && apt-get install -y wget \
    && apt-get install -y unzip \
    && apt-get install -y curl \
    && apt-get install -y libglib2.0-0 \
    && apt-get install -y libsm6 libxext6 libxrender-dev \
    && apt-get install -y nano

RUN pip install pika \
    && pip install jsonpickle

COPY . /python/sample

WORKDIR /python/sample
CMD python3 index.py



