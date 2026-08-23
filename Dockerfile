FROM golang:1.23 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /out/pg-lattice-proxy ./cmd/pg-lattice-proxy
FROM gcr.io/distroless/static-debian12
COPY --from=build /out/pg-lattice-proxy /pg-lattice-proxy
COPY configs /configs
EXPOSE 6432 8080
ENTRYPOINT ["/pg-lattice-proxy", "-config", "/configs/config.yaml"]
