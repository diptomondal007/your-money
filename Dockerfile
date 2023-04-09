ARG GO_VERSION=1.18

FROM golang:${GO_VERSION}-alpine as builder

WORKDIR /src
COPY ./ ./

RUN CGO_ENABLED=0 go build -mod=vendor -o /your-money .

FROM scratch AS final

COPY --from=builder your-money your-money

EXPOSE 8080
ENTRYPOINT ["/your-money"]