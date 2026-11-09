# ---------- Stage 1: compile the Go server ----------
FROM golang:1.27-alpine AS build

# Same label as the final image, so leftovers from this stage can be cleaned up by label.
LABEL project="ascii-art-web-dockerize"

WORKDIR /src

# Only the Go files are needed to compile, so editing HTML or CSS
# doesn't invalidate this layer.
COPY go.mod ./
COPY *.go ./

# CGO_ENABLED=0 gives a static binary that runs on plain Alpine.
# -s -w strips debug info to make the binary smaller.
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o server .

# ---------- Stage 2: the small image that actually runs ----------
FROM alpine:3.24

LABEL org.opencontainers.image.title="ascii-art-web" \
      org.opencontainers.image.description="Web server that turns text into ASCII art banners" \
      org.opencontainers.image.authors="Soul and Mo" \
      org.opencontainers.image.version="1.0.0" \
      project="ascii-art-web-dockerize"

# bash lets you run `docker exec -it <container> /bin/bash`.
# The server runs as its own user instead of root.
RUN apk add --no-cache bash \
    && adduser -D -H app

WORKDIR /app

# The binary from stage 1, plus the files the server reads at runtime.
COPY --from=build /src/server ./server
COPY templates ./templates
COPY static ./static
COPY standard.txt shadow.txt thinkertoy.txt ./

USER app

EXPOSE 8080

CMD ["./server"]
