FROM golang:1.10 as base
WORKDIR /go/src/chatbot
COPY . .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o main

FROM scratch
COPY --from=base /go/src/chatbot/main /chatbot
COPY --from=base /go/src/chatbot/index.html /index.html
COPY --from=base /go/src/chatbot/static /static/

EXPOSE 8080

CMD ["/chatbot"]
