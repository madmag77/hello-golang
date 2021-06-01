FROM golang:1.16 as build
COPY . /kvs
WORKDIR /kvs/src
RUN CGO_ENABLED=0 GOOS=linux go build -a -o kvs

FROM scratch
COPY --from=build /kvs/src/kvs .
EXPOSE 8080
CMD ["/kvs"]