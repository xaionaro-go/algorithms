Stupid solution. I think we can just mathematically prove that the solution exists only for `n==3` and just print the
solution ;)
```
goos: linux
goarch: amd64
pkg: github.com/xaionaro-go/algorithms/simple_tasks/facebook/permutation
BenchmarkFindPermutations1-8    	10000000	       140 ns/op	     112 B/op	       5 allocs/op
BenchmarkFindPermutations2-8    	 5000000	       260 ns/op	     288 B/op	       7 allocs/op
BenchmarkFindPermutations3-8    	 2000000	       722 ns/op	     496 B/op	      11 allocs/op
BenchmarkFindPermutations5-8    	  200000	      7998 ns/op	     848 B/op	      13 allocs/op
BenchmarkFindPermutations7-8    	   10000	    256649 ns/op	    1424 B/op	      17 allocs/op
BenchmarkFindPermutations10-8   	      50	  38400056 ns/op	    2656 B/op	      23 allocs/op
PASS
ok  	github.com/xaionaro-go/algorithms/simple_tasks/facebook/permutation	11.505s
```
