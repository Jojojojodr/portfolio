#!/bin/sh
set -e

# static/img is bind-mounted from the host. When the host folder is empty (e.g. a fresh
# checkout or a folder Docker just created), fill it with the images shipped in the image.
IMG_DIR=/app/static/img
DEFAULT_IMG_DIR=/app/defaults/img

mkdir -p "$IMG_DIR"
if [ -z "$(ls -A "$IMG_DIR")" ]; then
    echo "$IMG_DIR is empty, copying default images"
    cp -a "$DEFAULT_IMG_DIR/." "$IMG_DIR/"
fi

exec "$@"
