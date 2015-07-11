FROM google/golang

WORKDIR /gopath/src/github.com/maleck13/photoProcessor
RUN apt-get install  libexif-dev libexif12 -y
RUN mkdir -p /var/log/photoprocessor
ADD . /gopath/src/github.com/maleck13/photoProcessor
RUN go get github.com/maleck13/photoProcessor
RUN ls -al /gopath/bin
ENV PHOTO_PROC_CONF /etc/photoprocessor/conf.json
EXPOSE 9002
CMD []
ENTRYPOINT ["/gopath/bin/photoProcessor"]