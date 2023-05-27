FROM golang:latest

RUN mkdir -p /app
COPY src/* /app/.
RUN cd /app ; go mod download
