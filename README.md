## Path Tracer presented in bachelor's thesis "Progressive BVH Refinement in Interactive Ray Tracing"
### by Christian Schmidt

Prototype implementation of the presented approaches, might contain bugs and temporary code. Tests were removed to make project less cluttered. 

### Benchmarking:
Benchmarking code is contained in './pkg/benchmark/benchmark_test.go'. To run benchmarks reliably, install 'bench' package and start the daemon using <code>sudo -b bench -daemon</code> (linux only). From the benchmark directory './pkg/benchmark/benchmark_test.go' run <code>bench</code> to execute all tests. For further information, refer to https://github.com/golang-design/bench. 

### Image rendering:
An example for image rendering is provided in './cmd/image'. From the directory, it can be executed using <code>go run .</code>. 
A heatmap of the constructed BVH can be created the same way using './cmd/heatmap'

### Window mode: 
An example for running interactive scenes is provided in './cmd/window'. From the directory, it can be executed using <code>go run .</code>. Press 'q' to quit the window mode. 