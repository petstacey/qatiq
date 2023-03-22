# build tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY loggerService /app

# TODO: remove wildcard for cors before end of development
ENTRYPOINT ["app/loggerService", "-trusted-origins='*'"]