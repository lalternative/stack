# core image — distroless Go binary
FROM golang:1.26-bookworm AS build
WORKDIR /src
COPY go.work go.work.sum ./
COPY apps/core ./apps/core
RUN cd apps/core && CGO_ENABLED=0 go build -o /out/core .

FROM gcr.io/distroless/base-debian12
COPY --from=build /out/core /core
COPY apps/core/migrations /migrations
# Run as distroless' non-root user (UID 65532); :4100 is >1024 so no
# CAP_NET_BIND_SERVICE needed, and /migrations is read-only.
USER nonroot
EXPOSE 4100
ENTRYPOINT ["/core"]
