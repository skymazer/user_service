# GRPC-demo
# build command : docker build . -t grpc-demo/server
# run command : docker run -it grpc-demo/server
FROM golang:latest

# Install grpc
RUN go get -u google.golang.org/grpc && \
    go get -u github.com/golang/protobuf/protoc-gen-go && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Install protoc and zip system library
RUN apt-get update && apt-get install -y zip && \
    mkdir /opt/protoc && cd /opt/protoc && wget https://github.com/protocolbuffers/protobuf/releases/download/v3.7.0/protoc-3.7.0-linux-x86_64.zip && \
    unzip protoc-3.7.0-linux-x86_64.zip

ENV PATH=$PATH:$GOPATH/bin:/opt/protoc/bin

RUN mkdir  /app
COPY app /app
COPY ./users.proto /app/proto/
RUN cd /app/proto && \
    protoc --go-grpc_out=. --go_out=. users.proto

RUN cd /app && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates bash
WORKDIR /root/
COPY --from=0 /app/app ./
COPY ./ops/wait-for-it.sh ./wait-for-it.sh
RUN chmod +x ./wait-for-it.sh
ENTRYPOINT ./wait-for-it.sh -t 30 db:5432 -- ./app