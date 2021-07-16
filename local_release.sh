#! /bin/bash

export VERSION=v5.3.0-release                                            
export BUILD_IMAGE_BASE_NAME=registry.gitlab.com/gridworkz/kato
./release.sh all push
