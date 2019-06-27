FROM chromedp/headless-shell

COPY $appname .

EXPOSE 9099

ENTRYPOINT []

CMD ["/$appname"]
