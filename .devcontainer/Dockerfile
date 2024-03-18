FROM golang:1.22

RUN apt -y update

RUN apt -y install \
    systemd \
    git \
    tree \
    curl \
    nano \
    python3-pip \
    unzip \
    tar \
    gzip

RUN apt -y clean all

ENV PS1="[\u@\h \W]# "

COPY post_create /usr/local/bin/post_create

RUN post_create
