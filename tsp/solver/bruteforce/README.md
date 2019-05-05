```
goos: linux
goarch: amd64
pkg: github.com/xaionaro-go/algorithms/tsp/solver/bruteforce
BenchmarkSolver_cities4-8    	   20000	     87274 ns/op	   34342 B/op	     962 allocs/op
BenchmarkSolver_cities6-8    	    5000	    265602 ns/op	   63926 B/op	    1681 allocs/op
BenchmarkSolver_cities8-8    	     200	   6481720 ns/op	  119171 B/op	    2830 allocs/op
BenchmarkSolver_cities10-8   	      20	  88749914 ns/op	  194497 B/op	    4191 allocs/op
BenchmarkSolver_cities11-8   	       2	 844651113 ns/op	  256940 B/op	    5025 allocs/op
BenchmarkSolver_cities12-8   	       1	2448357414 ns/op	  308104 B/op	    5986 allocs/op
BenchmarkSolver_cities13-8   	       1	5137915905 ns/op	  396208 B/op	    7228 allocs/op
PASS
ok  	github.com/xaionaro-go/algorithms/tsp/solver/bruteforce	17.965s
```
