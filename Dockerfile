FROM golang:1.20 as bob


WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/app


FROM gcr.io/distroless/static-debian11 

COPY --from=bob  /go/bin/app /
CMD ["/app"]



