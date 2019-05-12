Solution for: https://www.careercup.com/question?id=5723406763294720

* `T: O(num_digits)`
* `M: O(1)`

```
goos: linux
goarch: amd64
pkg: github.com/xaionaro-go/algorithms/simple_tasks/facebook/sum_linked_lists
BenchmarkNode_Add_10e5-8    	50000000	        39.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkNode_Add_10e10-8   	20000000	        68.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkNode_Add_10e15-8   	20000000	        97.5 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/xaionaro-go/algorithms/simple_tasks/facebook/sum_linked_lists	5.615s
```