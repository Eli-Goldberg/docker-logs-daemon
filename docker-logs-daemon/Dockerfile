FROM golang:1.13

RUN mkdir /app 
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

ADD . /app/  

RUN go build -o dist/main . 

EXPOSE 8080
CMD ["/app/dist/main"]
