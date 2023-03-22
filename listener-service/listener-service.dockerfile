# build tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY listenerService /app

# TODO: remove wildcard for cors before end of development
ENTRYPOINT ["app/listenerService"]