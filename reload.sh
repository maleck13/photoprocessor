#!/bin/sh

docker stop photoprocessor;
docker rm photoprocessor
docker build -t photoprocessor .
docker run -it -v '/etc/photoprocessor:/etc/photoprocessor' -v '/opt/data/pictures:/opt/data/pictures' -v '/var/log/photoprocessor:/var/log/photoprocessor' -v '/opt/data/completedPics:/opt/data/completedPics' -v '/opt/data/thumbs:/opt/data/thumbs'  --name="photoprocessor" photoprocessor
