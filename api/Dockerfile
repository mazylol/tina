FROM nvidia/cuda:12.5.0-runtime-ubuntu22.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y \
        git \
        python3-pip \
        python3-dev \
        python3-opencv \
        libglib2.0-0

WORKDIR /code

COPY ./requirements.txt /code/requirements.txt

RUN pip install --no-cache-dir --upgrade -r /code/requirements.txt

COPY ./chat.py /code
COPY ./train.py /code
COPY ./intents-default.json /data/intents.json

RUN python3 train.py

EXPOSE 6969

CMD ["uvicorn", "chat:app", "--host", "0.0.0.0", "--port", "6969"]
