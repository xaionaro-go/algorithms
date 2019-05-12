Stupid solution. I think we can just mathematically prove that the solution exists only for small `n` and just print the
solution ;)
```
goos: linux
goarch: amd64
pkg: github.com/xaionaro-go/algorithms/simple_tasks/facebook/permutation
BenchmarkFindPermutations1-8   	10000000	       145 ns/op	     112 B/op	       5 allocs/op
BenchmarkFindPermutations2-8   	 5000000	       257 ns/op	     288 B/op	       7 allocs/op
BenchmarkFindPermutations3-8   	 2000000	      1069 ns/op	     592 B/op	      13 allocs/op
BenchmarkFindPermutations4-8   	  500000	      3690 ns/op	     784 B/op	      13 allocs/op
BenchmarkFindPermutations5-8   	   30000	     41590 ns/op	     849 B/op	      13 allocs/op
BenchmarkFindPermutations6-8   	    3000	    535642 ns/op	    1224 B/op	      15 allocs/op
BenchmarkFindPermutations7-8   	     300	   5942213 ns/op	    3377 B/op	      32 allocs/op
BenchmarkFindPermutations8-8   	      20	  75325611 ns/op	   18444 B/op	      99 allocs/op
BenchmarkFindPermutations9-8   	       1	1054364625 ns/op	   26688 B/op	      24 allocs/op
PASS
ok  	github.com/xaionaro-go/algorithms/simple_tasks/facebook/permutation	16.624s
```