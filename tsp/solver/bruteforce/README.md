goos: linux
goarch: amd64
pkg: github.com/xaionaro-go/algorithms/tsp/solver/bruteforce
BenchmarkSolver_cities4-8    	   20000	     76713 ns/op	   25232 B/op	     563 allocs/op
BenchmarkSolver_cities6-8    	    2000	    735415 ns/op	   34520 B/op	     649 allocs/op
BenchmarkSolver_cities8-8    	     100	  23317775 ns/op	   49239 B/op	     815 allocs/op
BenchmarkSolver_cities10-8   	       3	 351991471 ns/op	   65128 B/op	     926 allocs/op
BenchmarkSolver_cities11-8   	       1	3399775263 ns/op	   82432 B/op	    1042 allocs/op
BenchmarkSolver_cities12-8   	       1	12928480922 ns/op	   92592 B/op	    1104 allocs/op
BenchmarkSolver_cities13-8   	       1	18555609246 ns/op	   86632 B/op	    1062 allocs/op
PASS
ok  	github.com/xaionaro-go/algorithms/tsp/solver/bruteforce	43.231s
