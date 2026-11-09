#!/usr/bin/env bash
# Builds the image and (re)starts the container in one command.
#
#   ./docker.sh          build the image and start the container
#   ./docker.sh stop     stop and remove the container
#   ./docker.sh clean    remove the container, the image and leftovers
set -e

IMAGE="ascii-art-web-docker"
CONTAINER="dockerize"
PORT="8080"
LABEL="project=ascii-art-web-dockerize"

remove_container() {
    docker container rm --force "$CONTAINER" > /dev/null 2>&1 || true
}

# Garbage collection: when a new image takes the name, the old one becomes
# "dangling" (<none>). The label filter makes sure only this project's
# leftovers are removed, never your other images.
remove_leftovers() {
    docker image prune --force --filter "label=$LABEL" > /dev/null
}

case "${1:-start}" in
    start)
        docker image build -f Dockerfile -t "$IMAGE" .
        remove_container
        docker container run --detach --publish "$PORT:8080" --name "$CONTAINER" --label "$LABEL" "$IMAGE"
        remove_leftovers
        docker container ls --filter "name=$CONTAINER"
        echo "Running at http://localhost:$PORT"
        ;;
    stop)
        remove_container
        echo "Container $CONTAINER removed"
        ;;
    clean)
        remove_container
        docker image rm "$IMAGE" > /dev/null 2>&1 || true
        remove_leftovers
        echo "Container, image and leftovers removed"
        ;;
    *)
        echo "Usage: ./docker.sh [start|stop|clean]"
        exit 1
        ;;
esac
