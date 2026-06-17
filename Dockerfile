# core image — distroless Go binary
FROM golang:1.25-bookworm AS build
WORKDIR /src
COPY go.work go.work.sum ./
COPY apps/core ./apps/core
COPY apps/cli ./apps/cli
RUN cd apps/core && CGO_ENABLED=1 go build -o /out/core .

FROM gcr.io/distroless/base-debian12
COPY --from=build /out/core /core
COPY apps/core/migrations /migrations
ENV DUCKDB_PATH=/data/app.duckdb
EXPOSE 4100
ENTRYPOINT ["/core"]
