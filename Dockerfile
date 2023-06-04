FROM golang:latest

RUN mkdir -p /app
COPY src/* /app/.
RUN cd /app ; go mod download
RUN go install github.com/codegangsta/gin@latest
