# benchmark result

```bash
goos: linux
goarch: amd64
pkg: experiment/benchmark/assign_interface
BenchmarkValueMethod_AssignValue-8       	52575579	        22.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkValueMethod_AssignPointer-8     	51119427	        22.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkPointerMethod_AssignPointer-8   	599274757	         1.95 ns/op	       0 B/op	       0 allocs/op
PASS
```