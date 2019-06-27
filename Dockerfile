FROM chromedp/headless-shell

COPY $appname .

EXPOSE 9099
CMD ["/$appname"]
