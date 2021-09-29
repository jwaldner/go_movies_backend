# syntax=docker/dockerfile:1
FROM golang:1.16.8-alpine3.14
RUN addgroup app && adduser -S -G app app
WORKDIR /app/home
#RUN npm install
COPY . .
EXPOSE 4000
CMD go build -o movies cmd/api/*.go && ./movies
 





