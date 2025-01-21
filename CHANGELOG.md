<a name="unreleased"></a>
## [2.6.0](https://github.com/noandrea/geo2tz/compare/v2.5.0...v2.6.0) (2025-01-21)


### Features

* update tzdata to 2025a ([a92f1b2](https://github.com/noandrea/geo2tz/commit/a92f1b2acc8815d178ca098ec4bf0e902ae6d702))

## [2.5.0](https://github.com/noandrea/geo2tz/compare/v2.4.0...v2.5.0) (2024-09-16)


### Features

* update tzdata to 2024b ([86f0b97](https://github.com/noandrea/geo2tz/commit/86f0b97d607152974ac6fe2685837b9df6a14d77))

## [2.4.0](https://github.com/noandrea/geo2tz/compare/v2.3.0...v2.4.0) (2024-06-29)


### Features

* improve geojson parser ([370804c](https://github.com/noandrea/geo2tz/commit/370804cb5143a262557cd08495f84724518dea19))
* remove unnecessary Size function ([e3acbd5](https://github.com/noandrea/geo2tz/commit/e3acbd57d638f84e2177fd64bac0a4eedd47198d))

## [2.3.0](https://github.com/noandrea/geo2tz/compare/v2.2.0...v2.3.0) (2024-06-25)


### Features

* fail fast if a lat/lng is not found ([37cff76](https://github.com/noandrea/geo2tz/commit/37cff762145f621110e0cb39399583148b2d53cb))
* improve precision ([a194645](https://github.com/noandrea/geo2tz/commit/a1946451417d87feed0fa4fcf3f353c8f2cdfa09))

## [2.2.0](https://github.com/noandrea/geo2tz/compare/v2.1.5...v2.2.0) (2024-06-24)


### Features

* improve timezone lookup and add tests ([5c00920](https://github.com/noandrea/geo2tz/commit/5c00920c91c19729e9282bbf90b95e7661045e49))
* use an in memory rtree to find the timezone ([561e500](https://github.com/noandrea/geo2tz/commit/561e500fe7cc88c6f9f33c659620efef434bbead))

## [2.1.5](https://github.com/noandrea/geo2tz/compare/v2.1.4...v2.1.5) (2024-06-16)


### Bug Fixes

* github release jobs ([e4b720a](https://github.com/noandrea/geo2tz/commit/e4b720a6b3d2d16904665f9d6c92130320958a12))

## [2.1.4](https://github.com/noandrea/geo2tz/compare/v2.1.3...v2.1.4) (2023-06-23)


### Bug Fixes

* invalid docker build ([12d5d0d](https://github.com/noandrea/geo2tz/commit/12d5d0d2790ee8dc64b0c4b9901ca9b04b25795b))

## [2.1.3](https://github.com/noandrea/geo2tz/compare/v2.1.2...v2.1.3) (2023-06-23)


### Bug Fixes

* docker image won't build ([f13e809](https://github.com/noandrea/geo2tz/commit/f13e80991a92ea5e1208ae2466c39272446f348b))

## [2.1.2](https://github.com/noandrea/geo2tz/compare/v2.1.1...v2.1.2) (2023-06-23)


### Bug Fixes

* **build:** the docker image is broken ([4cf592b](https://github.com/noandrea/geo2tz/commit/4cf592bd09b85f69248b2e5b1f04af1ee53cca83))

## [2.1.1](https://github.com/noandrea/geo2tz/compare/v2.1.0...v2.1.1) (2023-06-23)


### Bug Fixes

* release workflow ([3658395](https://github.com/noandrea/geo2tz/commit/3658395c33531db79404eb87ee59515433b22458))

## [2.1.0](https://github.com/noandrea/geo2tz/compare/v2.0.0...v2.1.0) (2023-06-23)


### Features

* add script to update tzdata ([67d34db](https://github.com/noandrea/geo2tz/commit/67d34dbc7f78910530ae3a9c354b8f527ee259da))
* update tzdata to 2023b ([43f9260](https://github.com/noandrea/geo2tz/commit/43f9260affdd4a9c1907f68592b4138d3e665fc3))


### Bug Fixes

* error detecting non-existing files ([992a2ea](https://github.com/noandrea/geo2tz/commit/992a2ea6966b0a93761a48c3b99c3388909e01a4))
* **security:** update dependencies ([33fed50](https://github.com/noandrea/geo2tz/commit/33fed50b189893d7d06bbba5cd442f58afa93813))

## [Unreleased]


<a name="v2.0.0"></a>
## [v2.0.0] - 2022-07-31
### Chore
- improve readme, remove unused deps, improve logging
- update dependencies and go to v1.18

### Feat
- remove support for boltdb
- **data:** update tz source to 2021c version

### Test
- improve test coverage

### BREAKING CHANGE

the boltdb options are not available anymore, only the
in-memory db is supported


<a name="v1.0.0"></a>
## [v1.0.0] - 2021-11-28
### Chore
- upgrade go version and dependencies

### Docs
- update docs

### Feat
- use in memory db instead of bolt
- add server request rate limit


<a name="v0.4.2"></a>
## [v0.4.2] - 2021-09-13
### Build
- add multiple automated build options for docker


<a name="0.4.1"></a>
## [0.4.1] - 2021-09-13
### Build
- publish image on release

### Docs
- update docker image url

### Fix
- go mod tidy


<a name="0.4.0"></a>
## [0.4.0] - 2021-09-13
### Feat
- add support for memory shapefile
- add support for memory shapefile

### Fix
- **ci:** install linter in CI


<a name="0.3.1"></a>
## [0.3.1] - 2020-11-16
### Fix
- docker-image and release process fix


<a name="0.3.0"></a>
## [0.3.0] - 2020-11-16
### Chore
- move to github packages (from docker hub)
- set go version to 1.15

### Feat
- update tzdata to v2020d

### Fix
- linter warnings

### Misc
- hide echo startup banner


<a name="0.2.2"></a>
## [0.2.2] - 2020-06-29
### Build
- add deepsource, move to noandrea org on docker hub

### Fix
- **deepsource:** add error handling
- **deepsource:** simplify comparison


<a name="0.2.1"></a>
## [0.2.1] - 2020-05-17
### Doc
- fix README


<a name="0.2.0"></a>
## [0.2.0] - 2020-05-17
### Feat
- add authorization support


<a name="0.1.0"></a>
## [0.1.0] - 2020-05-06

<a name="0.0.0"></a>
## 0.0.0 - 2020-05-06
### Build
- add changelog support
- better integration with docker-hub
- fix issue with linting
- add travis config
- add dockeringore and gitignore

### Doc
- minor fixes in documentation
- add details in readme
- fix badges URLs
- add readme

### Feat
- first commit

### Test
- add tests for coordinate parsing


[Unreleased]: https://github.com/noandrea/geo2tz/compare/v2.0.0...HEAD
[v2.0.0]: https://github.com/noandrea/geo2tz/compare/v1.0.0...v2.0.0
[v1.0.0]: https://github.com/noandrea/geo2tz/compare/v0.4.2...v1.0.0
[v0.4.2]: https://github.com/noandrea/geo2tz/compare/0.4.1...v0.4.2
[0.4.1]: https://github.com/noandrea/geo2tz/compare/0.4.0...0.4.1
[0.4.0]: https://github.com/noandrea/geo2tz/compare/0.3.1...0.4.0
[0.3.1]: https://github.com/noandrea/geo2tz/compare/0.3.0...0.3.1
[0.3.0]: https://github.com/noandrea/geo2tz/compare/0.2.2...0.3.0
[0.2.2]: https://github.com/noandrea/geo2tz/compare/0.2.1...0.2.2
[0.2.1]: https://github.com/noandrea/geo2tz/compare/0.2.0...0.2.1
[0.2.0]: https://github.com/noandrea/geo2tz/compare/0.1.0...0.2.0
[0.1.0]: https://github.com/noandrea/geo2tz/compare/0.0.0...0.1.0
