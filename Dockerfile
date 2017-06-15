FROM ubuntu:16.10

RUN apt-get update
RUN apt-get install -y ca-certificates

COPY user-service /user-service

ENTRYPOINT ["./user-service"]
