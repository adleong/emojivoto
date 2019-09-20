#!/bin/bash

make build
tag="git-$(git rev-parse --short HEAD)"
docker tag buoyantio/emojivoto-web:v8 buoyantio/emojivoto-web:$tag
docker tag buoyantio/emojivoto-emoji-svc:v8 buoyantio/emojivoto-emoji-svc:$tag
docker tag buoyantio/emojivoto-voting-svc:v8 buoyantio/emojivoto-voting-svc:$tag
docker push buoyantio/emojivoto-voting-svc:$tag
docker push buoyantio/emojivoto-emoji-svc:$tag
docker push buoyantio/emojivoto-web:$tag
