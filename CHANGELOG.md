<a name="unreleased"></a>
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
