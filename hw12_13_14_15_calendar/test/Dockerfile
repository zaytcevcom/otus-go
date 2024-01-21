FROM golang:1.21

RUN mkdir -p /opt/integration_tests
WORKDIR /opt/integration_tests

COPY /test/go.mod .
COPY /test/go.sum .
RUN go mod download

COPY /test/. .
CMD ["go", "test"]
