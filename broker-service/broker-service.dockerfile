# build tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY brokerService /app

ENTRYPOINT ["app/brokerService", "-trusted-origins='*'"]