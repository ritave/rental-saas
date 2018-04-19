FROM golang:1.10.1

WORKDIR /go/src/rental-saas
COPY . .

RUN go get -d -v ./...; exit 0
RUN go install -v ./...; exit 0
RUN go build -o standalone/main.exe standalone/main.go

WORKDIR standalone

EXPOSE 8080
CMD ["./main.exe"]