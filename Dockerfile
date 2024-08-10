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
RUN mkdir /app
COPY --from=builder /app/build/meilisitemap /app
RUN chmod +x /app/meilisitemap
CMD ["./app/meilisitemap",  "-config", "/etc/meilisitemap/config.yml"]