FROM ubuntu

RUN mkdir /server
COPY ./server /server

EXPOSE 32420

RUN useradd app
USER app

CMD ["/server"]