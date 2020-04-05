FROM golang:alpine3.11 AS builder
WORKDIR /satelit-scraper/
COPY . .
RUN docker/scripts/build-project.sh

FROM alpine:latest
WORKDIR /satelit-scraper/
EXPOSE 9050
COPY --from=builder /satelit-scraper/satelit-scraper.tar.gz /satelit-scraper/
RUN ["tar", "-xvf", "satelit-scraper.tar.gz"]
ENTRYPOINT ["/satelit-scraper/entry.sh"]
HEALTHCHECK --interval=1m --start-period=10s CMD ps cax | grep satelit-scraper >/dev/null || exit 1
