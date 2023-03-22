# build tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY brokerService /app

# TODO: remove wildcard for cors before end of development
ENTRYPOINT ["app/brokerService", "-trusted-origins='*'"]