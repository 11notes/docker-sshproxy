# :: Builder
        FROM golang:1.14-alpine as builder

        RUN apk update \
                && apk upgrade \
                && apk add --no-cache \
                ca-certificates \
                && update-ca-certificates 2>/dev/null
        RUN apk add google-authenticator git gcc libc-dev linux-pam-dev

        ADD . /go/src/github.com/tg123/sshpiper/
        WORKDIR /go/src/github.com/tg123/sshpiper/sshpiperd
        RUN go build -ldflags "$(/go/src/github.com/tg123/sshpiper/sshpiperd/ldflags.sh)" -tags pam -o /go/bin/sshpiperd

# :: Header
        FROM alpine:3.12

# :: Run
        RUN apk update \
                && apk upgrade \
                && apk add --no-cache \
                ca-certificates \
                && update-ca-certificates 2>/dev/null
                
        RUN apk add google-authenticator

        RUN mkdir /etc/ssh/

        # :: docker copy source
                COPY src /

        # :: docker start and health script
                RUN chmod +x \
                        /usr/local/bin/entrypoint.sh \
                        /usr/local/bin/healthcheck.sh

        COPY --from=builder /go/bin/sshpiperd /
        EXPOSE 2222

# :: Monitor
        HEALTHCHECK CMD /usr/local/bin/healthcheck.sh || exit 1

# :: Start
        ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
        CMD ["/sshpiperd", "daemon"]