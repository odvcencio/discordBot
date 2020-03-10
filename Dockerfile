###############
# Build Stage #
###############
FROM golang:1.14-alpine as builder

# Installing dependencies
RUN apk add --no-cache \
      git \
      ca-certificates \
      curl \
      tzdata

# Setting work directory
WORKDIR ${GOPATH}/src/discordBot

# Copying rest of the code and building
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./discordBot

#################
# Release Stage #
#################
FROM scratch

# Copying files and folders from builder stage
COPY --from=builder /go/src/discordBot/discordBot /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo/

# Setting docker command
CMD ["./discordBot"]
