FROM python:3.12.4

WORKDIR /code

COPY ./requirements.txt /code/requirements.txt

RUN pip install --no-cache-dir --upgrade -r /code/requirements.txt

COPY ./chat.py /code
COPY ./train.py /code
COPY ./intents-default.json /data/intents.json

RUN python train.py

EXPOSE 6969

CMD ["uvicorn", "chat:app", "--host", "0.0.0.0", "--port", "6969"]
