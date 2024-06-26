FROM golang:1.14.0-alpine3.11 AS builder

# Add timezone information.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create appuser.
ENV USER=appuser
ENV UID=10001

RUN adduser \
	--disabled-password \
	--gecos "" \
	--home "/nonexistent" \
	--shell "/sbin/nologin" \
	--no-create-home \
	--uid "${UID}" \
	"${USER}"

WORKDIR $GOPATH/src/alextanhongpin/app/

COPY go.mod go.mod 

RUN go mod download
RUN go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/app

FROM scratch

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /go/bin/app /go/bin/app

USER appuser:appuser

CMD ["/go/bin/app"]
