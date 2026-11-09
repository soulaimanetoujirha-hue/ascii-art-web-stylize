# ascii-art-web-dockerize

A Go web server that turns text into ASCII art, packaged in a Docker container.
Only standard Go packages are used.

## Run with Docker

```bash
chmod +x docker.sh   # once
./docker.sh          # build the image and start the container
```

Then open http://localhost:8080.

| Command | What it does |
|---|---|
| `./docker.sh` | Build the image, replace the old container, clean up old images |
| `./docker.sh stop` | Stop and remove the container |
| `./docker.sh clean` | Remove the container, the image and leftovers |

### By hand (what the script does)

```bash
docker image build -f Dockerfile -t ascii-art-web-docker .
docker images
docker container run -p 8080:8080 --detach --name dockerize ascii-art-web-docker
docker ps -a
docker exec -it dockerize /bin/bash   # then: ls -l, exit
```

## Run without Docker

```bash
go run .
```

## How the Docker setup works

- **Multi-stage build:** Go compiles the server in `golang:1.27-alpine`. Only the
  binary and the files it reads (`templates/`, `static/`, the three banners) are
  copied into a small `alpine:3.24` image.
- **Metadata:** the image has labels for title, description, authors, version
  and project. The container gets the project label too.
- **Garbage collection:** the build stage is thrown away. The script removes the
  old container and this project's dangling images (found by their label).
- **Safety:** the server runs as a normal user, not root. `bash` is installed so
  you can look inside with `docker exec`.

## Authors

Soul and Mo
