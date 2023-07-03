FROM golang:1.20.4-buster AS builder

WORKDIR /app

COPY . ./
RUN go mod download

RUN CGO_ENABLED=0 go build -o /bin/app

FROM gcr.io/distroless/static-debian11

COPY --from=builder /bin/app /bin
COPY app.env .

EXPOSE 9536

ENTRYPOINT [ "/bin/app" ]