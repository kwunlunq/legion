#!/bin/bash

docker build . -f Dockerfile.base -t reg.paradise-soft.com.tw:5000/base-legion:latest

docker push reg.paradise-soft.com.tw:5000/base-legion:latest

docker rmi reg.paradise-soft.com.tw:5000/base-legion:latest

docker build . -f Dockerfile.build.base -t reg.paradise-soft.com.tw:5000/base-build-legion:latest

docker push reg.paradise-soft.com.tw:5000/base-build-legion:latest

docker rmi reg.paradise-soft.com.tw:5000/base-build-legion:latest
