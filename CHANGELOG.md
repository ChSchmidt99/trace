# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- Implemented PHR BVH builder
- Benchmark package including all benchmarks conducted
- Added optional unlit renderer
- Implemented parallel bucket sort
- Optimizer that uses Grid Search and Bayesian Optimization to optimize PHR parameters 
- BVH cost
### Changed 
- Fixed bug in enclosing functions
- Parallelized PHR find initial Cut
- Caching subtree size
- Atomic counter for updateAABB
- Incrementally compute SAH
- Morton code generation now without table

## [0.0.3] - 2021-08-12
### Added
- Axis-aligned bounding box
- Bounding volume hierarchy traversal
- Morton code computation 
- Basic LBVH construction 
- .obj file parser
- Light emitting materials 

### Changed
- Partially fixed scene node transformation
- Split up primitive and tracable
- Transforming normal using transposed inverse of transformation matrix
- Fixed ApproxZero bug

## [0.0.2] - 2021-07-15
### Changed
- Using Ray call by value instead of pointer for performance increase
- Reusing Rays when scattering and casting new rays from camera for performance increase
- Return scatter result as struct for performance increase
- Hit records are now a c style out variable to save allocs
- Render worker channel is now buffered 
- Changed RandomInUnitSphere implementation to more efficient approach

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
