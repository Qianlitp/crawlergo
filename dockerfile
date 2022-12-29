## Build
FROM golang:1.16-buster AS build

WORKDIR /app

RUN apt-get update && apt-get install unzip && rm -rf /var/lib/apt/lists/*
COPY ./ ./
RUN make build
RUN chmod +x ./get_chrome.sh && ./get_chrome.sh


## Deploy
FROM ubuntu:18.04

WORKDIR /

COPY --from=build /app/bin/crawlergo /crawlergo
COPY --from=build /app/latest/ /chrome/
RUN apt-get update && apt-get install -yq --no-install-recommends \
     libasound2 libatk1.0-0 libc6 libcairo2 libcups2 libdbus-1-3 \
     libexpat1 libfontconfig1 libgcc1 libgconf-2-4 libgdk-pixbuf2.0-0 libglib2.0-0 libgtk-3-0 libnspr4 \
     libpango-1.0-0 libpangocairo-1.0-0 libstdc++6 libx11-6 libx11-xcb1 libxcb1 libgbm1 \
     libxcursor1 libxdamage1 libxext6 libxfixes3 libxi6 libxrandr2 libxrender1 libxss1 libxtst6 libnss3 \
     && rm -rf /var/lib/apt/lists/*

ENTRYPOINT ["/crawlergo", "-c", "/chrome/chrome"]