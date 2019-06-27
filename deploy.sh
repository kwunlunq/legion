#!/bin/bash

branch=$1
server=reg.paradise-soft.com.tw:5000
appname=xunya-legion

docker pull $server/$appname:$branch
docker run --rm -p 9099:9099 --name $appname $server/$appname:$branch
