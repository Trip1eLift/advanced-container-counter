FROM golang:1.21-alpine3.19

RUN apk add curl
RUN apk add --update nodejs npm
RUN npm install -g nodemon

COPY ecs-container-metadata-mock-x.json /opt/ecs/metadata/cluster-x/ecs-container-metadata.json
COPY ecs-container-metadata-mock-y.json /opt/ecs/metadata/cluster-y/ecs-container-metadata.json

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

EXPOSE 8000

HEALTHCHECK CMD curl -sf http://localhost:8000/health || exit 1 --health-retries=3 --health-timeout=3s --health-interval=30s --health-start-period=1s