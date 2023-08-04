FROM golang:alpine AS build 
WORKDIR /go/src/consumption-ms

COPY . .
RUN go build -o /go/bin/consumption-ms cmd/api/main.go

EXPOSE 8080
FROM scratch
COPY --from=build /go/bin/consumption-ms /go/bin/consumption-ms
ENTRYPOINT [ "/go/bin/consumption-ms" ]

