# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.0.3] - 2021-08-12
### Added
- Axis-aligned bounding box
- Bounding volume hierarchy traversal
- Morton code computation 
- Basic LBVH construction 
- .obj file parser
- light emitting materials 

### Changed
- partially fixed scene node transformation
- split up primitive and tracable
- transforming normal using transposed inverse of transformation matrix
- fixed ApproxZero bug

## [0.0.2] - 2021-07-15
### Changed
- Using Ray call by value instead of pointer for performance increase
- Reusing Rays when scattering and casting new rays from camera for performance increase
- Return scatter result as struct for performance increase
- Hit records are now a c style out variable to save allocs
- render worker channel is now buffered 
- changed RandomInUnitSphere implementation to more efficient approach

## [0.0.1] - 2021-07-15
### Added 
- Very simple scene graph
- Minimal renderer
- Sphere and Triangle primitives
- Diffuse, reflective and refractive material
- Render multiple samples per pixel
- Basic math operations including vectors, matricies and quanternions

### Removed
- All math operations that directly add to struct, as benchmarking showed that the performance is worse than creating new values
- ApproxEqual in favor of ApproxZero, which is more efficient
