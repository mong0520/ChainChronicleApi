FROM golang:latest

WORKDIR $GOPATH/src/mong0520/ChainChronicleApi
COPY . $GOPATH/src/mong0520/ChainChronicleApi
RUN GO111MODULE=on go build

EXPOSE 5000
ENTRYPOINT ["./ChainChronicleApi"]