FROM alpine
MAINTAINER ash@the-rebellion.net

RUN apk --update add bash

RUN mkdir -p /app/bin
COPY app/ /app/

WORKDIR /app

CMD [ "/app/bin/forego", "start" ]
