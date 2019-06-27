FROM golang:1.12.1 AS builder

ARG branch
RUN echo $branch

ENV appname=xunya-legion
ENV GO111MODULE=on
RUN mkdir /$appname
WORKDIR /$appname
RUN ls -l

COPY . .

RUN ls -l

RUN go mod download &&\
    echo "go mod Completed"

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $appname &&\
    ls -l


FROM chromedp/headless-shell

COPY --from=builder /$appname/$appname /$appname

CMD ["/$appname/$appname"]
