FROM python:3.9-rc-alpine3.10
RUN apk add --no-cache git
RUN git clone https://github.com/skyway/skyway-js-sdk.git /skyway-js-sdk

WORKDIR /skyway-js-sdk/examples
ARG APIKEY
RUN sed -i -e "s/<YOUR_KEY_HERE>/$APIKEY/g" /skyway-js-sdk/examples/_shared/key.js

CMD python -m http.server 8080
