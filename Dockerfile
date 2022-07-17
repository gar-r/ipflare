FROM golang:alpine AS builder
WORKDIR /ipflare
COPY . .
RUN go build

FROM alpine
COPY --from=builder /ipflare/ipflare /usr/bin/ipflare
ENTRYPOINT [ "/usr/bin/ipflare" ]
CMD [ "-h" ]