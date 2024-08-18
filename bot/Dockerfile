FROM golang:1.21

WORKDIR /app

COPY . .

RUN go mod tidy

ENV PROD=1

CMD ["go", "run", "."]
