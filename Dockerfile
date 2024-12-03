FROM golang:1.23.3

WORKDIR /app

ADD ./app /app

RUN go get
#RUN go build -o bin .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server

EXPOSE 9000

CMD [ "/server" ]