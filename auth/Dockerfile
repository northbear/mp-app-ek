FROM docker.io/golang:1.18 AS buildenv

COPY . src/app
WORKDIR src/app

RUN make

FROM docker.io/alpine:3.16
RUN apk add --no-cache libc6-compat

WORKDIR /root
COPY --from=buildenv /go/src/app/build/auth /usr/local/bin/

EXPOSE 10080/tcp

CMD '/usr/local/bin/auth'
