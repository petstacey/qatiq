# build tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY authenticationService /app

# TODO: remove wildcard for cors before end of development
ENTRYPOINT ["app/authenticationService", "-trusted-origins='*'"]