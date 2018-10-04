FROM golang:1.11-stretch

ENV HOME /home/jenkins

# User configuration for proper interopability with the CI container's environment
RUN groupadd -g 10000 jenkins
RUN useradd -c "Jenkins user" -d $HOME -u 10000 -g 10000 -m jenkins

# apt dependencies
RUN apt-get update
RUN apt-get install -y unzip

# Protobuf compiler
RUN wget https://github.com/protocolbuffers/protobuf/releases/download/v3.6.1/protoc-3.6.1-linux-x86_64.zip
RUN unzip protoc-3.6.1-linux-x86_64.zip -d /opt/protoc
RUN rm protoc-3.6.1-linux-x86_64.zip
RUN ln -s /opt/protoc/bin/protoc /usr/bin/protoc
RUN ln -s /opt/protoc/include/google /usr/include/google

# Build dependencies
RUN go get -u -v github.com/golang/protobuf/protoc-gen-go
RUN go get -u -v golang.org/x/lint/golint
RUN go get -u -v google.golang.org/grpc

# Prepare environment for non-root user
RUN chown -R jenkins:jenkins $HOME $GOPATH
