Stupid solution. I think we can just mathematically prove that the solution exists only for small `n` and just print the
solution ;)
```
goos: linux
goarch: amd64
pkg: github.com/xaionaro-go/algorithms/simple_tasks/facebook/permutation
BenchmarkFindPermutations1-8   	10000000	       147 ns/op	     112 B/op	       5 allocs/op
BenchmarkFindPermutations2-8   	 5000000	       295 ns/op	     288 B/op	       7 allocs/op
BenchmarkFindPermutations3-8   	 1000000	      1098 ns/op	     592 B/op	      13 allocs/op
BenchmarkFindPermutations4-8   	  300000	      4375 ns/op	     896 B/op	      15 allocs/op
BenchmarkFindPermutations5-8   	   30000	     49574 ns/op	     848 B/op	      13 allocs/op
BenchmarkFindPermutations6-8   	    3000	    585202 ns/op	    1241 B/op	      15 allocs/op
BenchmarkFindPermutations7-8   	     200	   6763583 ns/op	   10385 B/op	      76 allocs/op
BenchmarkFindPermutations8-8   	      20	  85533026 ns/op	   66060 B/op	     329 allocs/op
BenchmarkFindPermutations9-8   	       1	1220408866 ns/op	   26688 B/op	      24 allocs/op
PASS
ok  	github.com/xaionaro-go/algorithms/simple_tasks/facebook/permutation	14.922s
```