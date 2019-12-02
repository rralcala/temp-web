FROM golang:1.13-alpine

COPY temp-web /
COPY home.html /

CMD ["/temp-web"]