FROM golang:1.18

WORKDIR kalpsdk
COPY . .
RUN go install -v golang.org/x/tools/cmd/godoc@latest
RUN go mod vendor

EXPOSE 8080
CMD /go/bin/godoc -http=:8080
