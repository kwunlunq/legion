FROM chromedp/headless-shell

ENV appname=xunya-legion

COPY $appname .

EXPOSE 9099

ENTRYPOINT []

CMD ["/$appname"]
