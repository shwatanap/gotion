REPOSITORY_NAME=gotion-server
LOCATION_NAME=asia-northeast1
GCP_PROJECT_NAME=gotion-395708

DOCKER_IMAGE_NAME=${LOCATION_NAME}-docker.pkg.dev/${GCP_PROJECT_NAME}/${REPOSITORY_NAME}/images

all: build push

build:
	docker build --tag ${DOCKER_IMAGE_NAME} -f ./Dockerfile.prod .

push:
	docker push ${DOCKER_IMAGE_NAME}
