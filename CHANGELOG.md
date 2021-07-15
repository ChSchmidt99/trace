# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- Scene graphs
- Minimal renderer
- Sphere and Triangle primitive
- Diffuse, reflective and refractive material

## [0.0.1] - 2021-07-13
### Added 
- This CHANGELOG file
- Copied math.go from old project and split up into vector.go, matrix.go and quanternion.go

### Removed
- All math operations that directly add to struct, as benchmarking showed that the performance is worse than creating new values
- ApproxEqual in favor of ApproxZero, which is more efficient
