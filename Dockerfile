FROM golang:1.21-alpine

WORKDIR /app

COPY . .

RUN go mod download
RUN go build .

FROM scratch 

COPY --from=0 /app/jellyexporter .

EXPOSE 8080

CMD ["/jellyexporter"]