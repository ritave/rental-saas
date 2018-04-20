FROM golang:1.10.1

WORKDIR /go/src/rental-saas/service
COPY . .

RUN go get -d -v ./...; exit 0
RUN go install -v ./...; exit 0
RUN go build -o main.exe main.go

#EXPOSE 8080
CMD ["./main.exe"]