FROM chromedp/headless-shell

COPY xunya-legion .

EXPOSE 9099

ENTRYPOINT []

CMD ["./xunya-legion"]
