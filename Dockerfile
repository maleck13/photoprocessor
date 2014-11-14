FROM google/golang

WORKDIR /gopath/src/app
ADD . /gopath/src/app/
RUN apt-get install  libexif-dev libexif12 -y
RUN mkdir -p /var/log/photoprocessor
RUN mkdir -p /etc/photoprocessor
RUN mkdir -p /opt/data/pictures
RUN mkdir -p /opt/data/completedPics
RUN mkdir -p /opt/data/thumbs
RUN go get app
RUN ls -al /gopath/bin
ENV PHOTO_PROC_CONF /etc/photoprocessor/conf.json
COPY ./conf.json /etc/photoprocessor/
CMD []
ENTRYPOINT ["/gopath/bin/app"]