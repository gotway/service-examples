FROM golang:1.18.3-alpine3.16 AS builder

ARG SERVICE
RUN test -n "$SERVICE" || (echo "Build argument \"SERVICE\" not set" && false)
ENV WORKDIR /go/src/services
RUN mkdir -p ${WORKDIR}
WORKDIR ${WORKDIR}

RUN apk update && \
  apk add --no-cache --update make bash git ca-certificates && \
  update-ca-certificates

COPY . .

RUN make build-$SERVICE

FROM alpine:3.16.0

COPY --from=builder /go/src/services/bin/app /app

CMD [ "/app" ]