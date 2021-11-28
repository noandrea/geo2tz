############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
ARG DOCKER_TAG=0.0.0
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git curl unzip
# download tz data
WORKDIR /tz
# download the location file 
RUN curl -LO https://github.com/evansiroky/timezone-boundary-builder/releases/download/2020d/timezones-with-oceans.geojson.zip
RUN unzip timezones-with-oceans.geojson.zip 
# checkout the project 
WORKDIR /builder
COPY . .
# Fetch dependencies.
# Using go get.
RUN go get -d -v
# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /geo2tz -ldflags="-s -w -extldflags \"-static\" -X main.Version=$DOCKER_TAG"
# build the database
# RUN /geo2tz build --json "/tz/combined-with-oceans.json" --db=/timezone --type=boltdb
# build the memory database timezone.snap.json
RUN /geo2tz build --json "/tz/combined-with-oceans.json" --db=/timezone --type=memory
############################
# STEP 2 build a small image
############################
FROM scratch
# Copy our static executable.
# COPY --from=builder /timezone.snap.db /
COPY --from=builder /timezone.snap.json /
COPY --from=builder /geo2tz /
# Copy the temlates folder
# COPY templates /templates
# Run the hello binary.
ENTRYPOINT [ "/geo2tz" ]
CMD [ "start" ]
