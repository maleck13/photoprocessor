FROM google/golang

WORKDIR /gopath/src/app
RUN apt-get install  libexif-dev libexif12 -y
RUN mkdir -p /var/log/photoprocessor
ADD . /gopath/src/app/
RUN go get app
RUN ls -al /gopath/bin
ENV PHOTO_PROC_CONF /etc/photoprocessor/conf.json
CMD []
ENTRYPOINT ["/gopath/bin/app"]