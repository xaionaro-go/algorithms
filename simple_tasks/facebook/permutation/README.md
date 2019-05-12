Stupid solution. Also it's required to think if we can just mathematically prove that the solution exists only for
small `n` and just print the solution ;)

```
goos: linux
goarch: amd64
pkg: github.com/xaionaro-go/algorithms/simple_tasks/facebook/permutation
BenchmarkFindPermutations1-8    	10000000	       157 ns/op	     112 B/op	       5 allocs/op
BenchmarkFindPermutations2-8    	 5000000	       303 ns/op	     288 B/op	       7 allocs/op
BenchmarkFindPermutations3-8    	 2000000	       927 ns/op	     592 B/op	      13 allocs/op
BenchmarkFindPermutations4-8    	  500000	      2403 ns/op	     896 B/op	      15 allocs/op
BenchmarkFindPermutations5-8    	  200000	      6577 ns/op	     848 B/op	      13 allocs/op
BenchmarkFindPermutations6-8    	   30000	     50158 ns/op	    1216 B/op	      15 allocs/op
BenchmarkFindPermutations7-8    	    5000	    336486 ns/op	   10310 B/op	      76 allocs/op
BenchmarkFindPermutations8-8    	    1000	   2068261 ns/op	   64856 B/op	     329 allocs/op
BenchmarkFindPermutations9-8    	     100	  12793559 ns/op	    2287 B/op	      21 allocs/op
BenchmarkFindPermutations10-8   	      20	  91383175 ns/op	    4028 B/op	      24 allocs/op
BenchmarkFindPermutations11-8   	       2	 711335603 ns/op	10898304 B/op	   35636 allocs/op
BenchmarkFindPermutations12-8   	       1	5786154805 ns/op	70391264 B/op	  216358 allocs/op
PASS
ok  	github.com/xaionaro-go/algorithms/simple_tasks/facebook/permutation	26.241s
```