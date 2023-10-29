FROM alpine:latest

RUN mkdir /app

COPY productsApp /app

CMD [ "/app/productsApp"]