FROM golang:1.19.4-alpine3.16 as build-env

RUN mkdir /app
WORKDIR /app

# Install ca-certificates and set timezone
RUN apk update && apk add --no-cache ca-certificates tzdata

# Set the timezone to Asia/Jakarta
RUN ln -sf /usr/share/zoneinfo/Asia/Jakarta /etc/localtime && echo "Asia/Jakarta" > /etc/timezone

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/app

FROM scratch

COPY --from=build-env /go/bin/app /go/bin/app
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-env /usr/share/zoneinfo/Asia/Jakarta /usr/share/zoneinfo/Asia/Jakarta
COPY --from=build-env /etc/localtime /etc/localtime
COPY --from=build-env /etc/timezone /etc/timezone

ENTRYPOINT ["/go/bin/app"]
