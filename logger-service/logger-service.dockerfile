# build tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY loggerService /app

ENTRYPOINT ["app/loggerService", "-trusted-origins='*'"]