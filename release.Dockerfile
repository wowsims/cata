# syntax=docker/dockerfile:1


FROM debian:12-slim

ARG TARGETARCH

WORKDIR /cata

COPY wowsimcata-${TARGETARCH}-linux ./wowsimcata

EXPOSE 3333

CMD ["./wowsimcata", "--host", "0.0.0.0:3333", "--launch=false"]
