# Geo2Tz

[![Build Status](https://travis-ci.com/noandrea/geo2tz.svg?branch=master)](https://travis-ci.com/noandrea/geo2tz) [![codecov](https://codecov.io/gh/noandrea/geo2tz/branch/master/graph/badge.svg)](https://codecov.io/gh/noandrea/geo2tz) [![GoDoc](https://godoc.org/github.com/noandrea/geo2tz?status.svg)](https://godoc.org/github.com/noandrea/distill) [![Go Report Card](https://goreportcard.com/badge/github.com/noandrea/geo2tz)](https://goreportcard.com/report/github.com/noandrea/geo2tz)


A self-host-able service to get the timezone given geo coordinates (lat/long)

It does it by exposing the library from https://github.com/evanoberholster/timezoneLookup 

Tz data comes from https://github.com/evansiroky/timezone-boundary-builder


## Motivations

Geo coordinates might be a sensible information to share in many context,
and I needed a self-hosted solution to ensure that coordinates where not leaked.


## API

the services exposes only one api:

```
GET /tz/${LAT}/${LON}
```

and returns a json reply (`http/200`):

```
{
    "tz": "${TIMEZONE}",
    "coords": {
        "lat": ${LAT},
        "lon": ${LON}
    }
}
```

or in case of errors (`http/4**`):

```
{
    "message": "${DESCRIPTION}
}
```

## Docker

Docker image is available at https://hub.docker.com/orgs/apeunit/repositories

```
docker run -p 2004:2004 apeunit/geo2tz
```

The image is built on scratch, the image size is ~76mb:

- ~11mb the application
- ~62mb the tz data 

## K8s

Kubernetes configuration example:

```
---
# Deployment
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  labels:
    app: geo2tz
  name: geo2tz
spec:
  replicas: 1
  revisionHistoryLimit: 3
  selector:
    matchLabels:
      app: geo2tz
  template:
    metadata:
      labels:
        app: geo2tz
    spec:
      containers:
      - env:
        # this is an example, a better alternative would be to use 
        # a k8s secret or config map
        - name: GEO2TZ_API_KEY
          value: asecretmaybebetter
        image: apeunit/geo2tz
        imagePullPolicy: Always
        name: geo2tz
        ports:
        - name: http
          containerPort: 2004
---
# Service 
# the service for the above deployment
apiVersion: v1
kind: Service
metadata:
  name: geo2tz-service
spec:
  type: ClusterIP
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: http
  selector:
    app: geo2tz

```


