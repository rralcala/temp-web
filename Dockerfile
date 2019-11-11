FROM golang:1.13-alpine

COPY temp-web /

CMD ["/temp-web"]