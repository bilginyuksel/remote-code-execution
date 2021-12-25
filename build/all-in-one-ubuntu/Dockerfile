FROM ubuntu:latest

RUN apt-get update -y && \
    apt-get upgrade -y && \
    apt-get install -y wget

# install golang
RUN wget https://dl.google.com/go/go1.17.5.linux-arm64.tar.gz && \
    tar -xvf go1.17.5.linux-arm64.tar.gz && \
    mv go /usr/local && \
    rm -rf go1.17.5.linux-arm64.tar.gz

# install python
RUN apt-get install -y python3

# install java
RUN apt-get install -y openjdk-8-jdk

# install nodejs
RUN apt-get install -y nodejs

# install c++
RUN apt-get install -y g++