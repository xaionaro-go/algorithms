```
goos: linux
goarch: amd64
pkg: github.com/xaionaro-go/algorithms/simple_tasks/interviewbit/sliding-window
BenchmarkMy100000w20000-8         	     100	  13691744 ns/op	 4977698 B/op	  100004 allocs/op
BenchmarkProposed100000w20000-8   	       1	1322544760 ns/op	 3687680 B/op	      29 allocs/op
BenchmarkMy10000w2000-8           	    2000	   1093007 ns/op	  500256 B/op	   10004 allocs/op
BenchmarkProposed10000w2000-8     	     100	  12762169 ns/op	  287994 B/op	      19 allocs/op
BenchmarkMy1000w200-8             	   20000	     85645 ns/op	   50208 B/op	    1004 allocs/op
BenchmarkProposed1000w200-8       	   10000	    138409 ns/op	   16376 B/op	      11 allocs/op
BenchmarkMy100w20-8               	  200000	      6334 ns/op	    5184 B/op	     104 allocs/op
BenchmarkProposed100w20-8         	  500000	      2626 ns/op	    2040 B/op	       8 allocs/op
BenchmarkMy10w2-8                 	 3000000	       598 ns/op	     560 B/op	      14 allocs/op
BenchmarkProposed10w2-8           	10000000	       179 ns/op	     248 B/op	       5 allocs/op
PASS
ok  	github.com/xaionaro-go/algorithms/simple_tasks/interviewbit/sliding-window	18.337s
```
