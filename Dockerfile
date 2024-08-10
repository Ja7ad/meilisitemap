FROM golang:1.22-alpine AS builder
LABEL authors="Javad Rajabzadeh"

RUN apk add make

RUN mkdir /app
WORKDIR /app
COPY . /app
RUN make build

FROM alpine
RUN apk --no-cache add ca-certificates tzdata
RUN mkdir /etc/meilisitemap/
COPY --from=builder /app/build/meilisitemap /usr/local/bin
RUN chmod +x /usr/local/bin
CMD ["meilisitemap",  "-config", "/etc/meilisitemap/config.yml"]