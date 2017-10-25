FROM golang:1.9.1

EXPOSE 4273

WORKDIR /go/src/github.com/8-team/bacotto
ADD . .

RUN go-wrapper download
RUN go-wrapper install

RUN go get -u github.com/pilu/fresh

CMD bacotto
