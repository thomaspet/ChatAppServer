FROM golang

WORKDIR /app

ADD ./ /app/

RUN go get -u github.com/cosmtrek/air

CMD [ "air" ]