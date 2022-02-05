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

# Copy the grpc proto file and generate the go module
RUN mkdir  /app
COPY app /app
RUN cd /app/proto && \
    protoc --go-grpc_out=. --go_out=. users.proto

RUN cd /app && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /app/app ./
RUN chmod +x ./app
ENTRYPOINT ./app