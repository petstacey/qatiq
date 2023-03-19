# build tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY authenticationService /app

ENTRYPOINT ["app/authenticationService", "-trusted-origins='*'"]