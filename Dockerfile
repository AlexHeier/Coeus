FROM ubuntu

RUN apt-get update && apt-get upgrade -y

RUN mkdir /server
RUN mkdir /server/dashboard

WORKDIR /server

COPY ./dashboard/index.html /server/dashboard/index.html
COPY ./server /server/server
COPY ./.env /server/.env

EXPOSE 8081

RUN useradd app
USER app