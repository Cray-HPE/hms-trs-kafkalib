# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2022-01-05

### Changed

- Now uses Confluent kafka Go package, which changes the interface slightly, requiring a major rev bump.

## [1.5.3] - 2021-08-10

### Changed

- Updated jenkinsfile and snyk. Added .github

## [1.5.2] - 2021-07-22

### Changed

- Updated github references

## [1.5.1] - 2021-07-22

### Changed

Add support for building within the CSM Jenkins.

## [1.5.0] - 2021-06-28

### Security

- CASMHMS-4898 - Updated base container images for security updates.

## [1.4.2] - 2021-05-05

### Changed

- Updated docker-compose files to pull images from Artifactory instead of DTR.

## [1.4.1] - 2021-04-20

### Changed

- Updated Dockerfiles to pull base images from Artifactory instead of DTR.

## [1.4.0] - 2021-01-26

### Changed

- Updated license in all source code.
- Corrected string creation error in library.go

## [1.3.0] - 2021-01-14

### Changed

- Updated license file.

## [1.2.2] - 2020-10-21

### Security

- CASMHMS-4105 - Updated base Golang Alpine image to resolve libcrypto vulnerability.

## [1.2.1] - 2020-08-12

### Changed

- CASMHMS-2979 - Updated hms-trs-kafkalib to use the latest trusted baseOS images.

## [1.2.0] - 2020-07-01

### Fixed

- Added product field to jenkins file to enable building in the pipeline

## [1.1.1] - 2020-02-20

- Enhanced logging, made builds work

## [1.1.0] - 2020-02-13

- Changed contract of the API, converted the logger to use logrus.

## [1.0.3] - 2020-02-06

- Changed contract Kafka Init.

## [1.0.2] - 2020-02-06

- Changed contract of library so topics can come and go as needed.

## [1.0.1] - 2020-01-13

- Implemented library functionality.  

## [1.0.0] - 2019-12-11

### Added

### Changed

### Deprecated

### Removed

### Fixed

### Security
