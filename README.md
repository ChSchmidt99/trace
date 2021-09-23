## Path Tracer presented in bachelor's thesis "Progressive BVH Refinement in Interactive Ray Tracing" by Christian Schmidt

Prototype implementation of the presented approaches, might contain bugs and temporary code. Tests were removed to make project less cluttered. Place this project directory in your GOPATH, more information on https://golang.org/doc/gopath_code. 
### Overview
'./cmd' contains usage examples
'./pkg/pt' contains the path tracing engine
'./pkg/benchmark' contains benchmarking code 
'./pkg/demoscenes' contains a few demo scenes
'./pkg/interactive' contains a rudimentary interactive application implementation

### Benchmarking:
Benchmarking code is contained in './pkg/benchmark/benchmark_test.go'. To run benchmarks reliably, install 'bench' package and start the daemon using <code>sudo -b bench -daemon</code> (linux only). From the benchmark directory './pkg/benchmark/benchmark_test.go' run <code>bench</code> to execute all tests. For further information, refer to https://github.com/golang-design/bench. 

### Image rendering:
An example for image rendering is provided in './cmd/image'. From the directory, it can be executed using <code>go run .</code>. 
A heatmap of the constructed BVH can be created the same way using './cmd/heatmap'

### Window mode: 
An example for running interactive scenes is provided in './cmd/interactive'. From the directory, it can be executed using <code>go run .</code>. It is possible, that a local installation of SDL2 is necessary to create and write to windows.

Press w,a,s,d to move and 'q' to quit. Note that the interactive part is a very rough implementation and not at all representative of a proper implementation. The performance of interactive applications is also limited and does not utilize the full CPUs capacity.

'cmd/heatmap_interactive' contains an example on how to render an interactive BVH heatmap. Again, use <code>go run .</code> to execute. 