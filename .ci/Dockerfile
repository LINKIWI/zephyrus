FROM docker.internal.kevinlin.info/infra/ci-base:0.2.5

# apt dependencies
RUN sudo apt-get update
RUN sudo apt-get install -y unzip

# Protobuf compiler
RUN wget https://github.com/protocolbuffers/protobuf/releases/download/v3.6.1/protoc-3.6.1-linux-x86_64.zip
RUN sudo unzip protoc-3.6.1-linux-x86_64.zip -d /opt/protoc
RUN rm protoc-3.6.1-linux-x86_64.zip
RUN sudo ln -s /opt/protoc/bin/protoc /usr/bin/protoc
RUN sudo ln -s /opt/protoc/include/google /usr/include/google

# Build dependencies
RUN go get -u -v github.com/golang/protobuf/protoc-gen-go
RUN go get -u -v golang.org/x/lint/golint
RUN go get -u -v google.golang.org/grpc
