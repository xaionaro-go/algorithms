This is quite precise approximate method.

```
goos: linux
goarch: amd64
pkg: github.com/xaionaro-go/algorithms/tsp/solver/approximate
BenchmarkSolver_cities4-8     	   10000	    198104 ns/op	   50597 B/op	    1001 allocs/op
BenchmarkSolver_cities8-8     	    1000	   1982926 ns/op	  255250 B/op	    5649 allocs/op
BenchmarkSolver_cities10-8    	     500	   3726642 ns/op	  448612 B/op	    9682 allocs/op
BenchmarkSolver_cities12-8    	     300	   4400313 ns/op	  462419 B/op	    9976 allocs/op
BenchmarkSolver_cities15-8    	     300	   5955272 ns/op	  597451 B/op	   12317 allocs/op
BenchmarkSolver_cities20-8    	     100	  11328068 ns/op	  999726 B/op	   19970 allocs/op
BenchmarkSolver_cities25-8    	     100	  17323736 ns/op	 1183328 B/op	   22322 allocs/op
BenchmarkSolver_cities30-8    	      50	  30308630 ns/op	 1625266 B/op	   29500 allocs/op
BenchmarkSolver_cities40-8    	      30	  52896895 ns/op	 2362477 B/op	   40323 allocs/op
BenchmarkSolver_cities50-8    	      20	  95037529 ns/op	 3005639 B/op	   48016 allocs/op
BenchmarkSolver_cities60-8    	      10	 167432783 ns/op	 4470933 B/op	   66257 allocs/op
BenchmarkSolver_cities70-8    	       3	 420310924 ns/op	 5229098 B/op	   75730 allocs/op
BenchmarkSolver_cities80-8    	       3	 398814056 ns/op	 6774104 B/op	   96187 allocs/op
BenchmarkSolver_cities90-8    	       5	 318029991 ns/op	 9084654 B/op	  125932 allocs/op
BenchmarkSolver_cities95-8    	       3	 482644849 ns/op	 9069986 B/op	  123611 allocs/op
BenchmarkSolver_cities100-8   	       1	1486157700 ns/op	 9674392 B/op	  127757 allocs/op
BenchmarkSolver_cities120-8   	       2	 652659302 ns/op	13919412 B/op	  175752 allocs/op
BenchmarkSolver_cities150-8   	       1	1131565454 ns/op	20058240 B/op	  235491 allocs/op
BenchmarkSolver_cities200-8   	       1	3397649655 ns/op	32800144 B/op	  352725 allocs/op
BenchmarkSolver_cities250-8   	       1	4721367822 ns/op	49688552 B/op	  504508 allocs/op
BenchmarkSolver_cities350-8   	       1	7055837754 ns/op	92295560 B/op	  854285 allocs/op
BenchmarkSolver_cities500-8   	       1	21497862494 ns/op	177679624 B/op	 1555707 allocs/op
PASS
ok  	github.com/xaionaro-go/algorithms/tsp/solver/approximate	76.967s
```

Keep in mind:
* For every city it solves 8 different tasks (so it shows performance in 8 times less than real one)
* The most of the performance is utilized by normalization process (I was solving more common problem than classic TSP: there're just random one-directional routes from random cities to random cities; so for example costs of travels from specific cities to specific cities are unknown).
* There's no simplifications in the problem (like clusters of cities). Just a cloud of cities in non-metric space.
* Theoretically it should calculate quite precise.

Performance optimizations:
* Remove useless routes and duplicate paths.
* Pre-calculation of costs from every city to every city (this part is the most CPU utilizing one)
* Cache results while/for the pre-calculation (to do not do full search on every calculation).
* Boundaries (for example an estimation of the maximal cost by a stupid solver) to reduce useless branching (of the bruteforcing process)
* Sort routes by some stupid heuristics to reduce useless branching (of the bruteforcing process)
* Ant-like markers (and re-sort routes) to reduce useless branching (of the bruteforcing process)
* An approach for global/full path is done by heuristics (which chooses the cheapest route to a city the most far from the finish-city)
* Bruteforce only on not more than 6 sequential cities (to improve local solutions only).

Precision optimizations:
* Sort routes by some stupid heuristics to pick the most useful routes.
* Ant-like markers (and re-sort routes) to pick the most useful routes.
* Pick optimal solution for every sequential 6 cities (by bruteforce).
* Optimization of local solutions using the same heuristics method which was used for the global solution (for example it helps with problems like "drilling route").

Obvious things to improve:
* Parallelize (it uses only one logical CPU core now).
* More boundaries
* Clean up, optimize simple things (like do no reallocate memory when it's not required).
* Restructurize the code, add comments
* Visualize the process