FROM devopsworks/golang-upx:1.16 as builder

WORKDIR /app
COPY go.mod /app
COPY go.sum /app

CMD ["/bin/bash"]
ADD . /app/
RUN --mount=type=cache,target=/go/pkg/mod go mod download

RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build  CGO_ENABLED=1  go install main.go
RUN strip --strip-unneeded /go/bin/main
RUN /usr/local/bin/upx -9 /go/bin/main


# this is where the application actually runs
FROM ubuntu:20.04

RUN apt-get update
RUN apt-get install -y --no-install-recommends --fix-missing ca-certificates tzdata
RUN apt-get autoremove
RUN apt-get clean && rm -rf /tmp/* /var/tmp/* /var/lib/apt/archive/* /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /go/bin/main worker
RUN mkdir etc
RUN mkdir templates
RUN mkdir static
COPY etc/conf.yaml /app/etc/
COPY templates/* /app/templates/
COPY static/* /app/static/
COPY chat.db /app/

CMD ["/app/worker"]
