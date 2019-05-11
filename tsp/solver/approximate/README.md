This is quite precise approximate method.

```
goos: linux
goarch: amd64
pkg: github.com/xaionaro-go/algorithms/tsp/solver/approximate
BenchmarkSolver_cities4-8    	   10000	    118574 ns/op	   39794 B/op	     894 allocs/op
BenchmarkSolver_cities8-8    	    2000	   1271148 ns/op	  198811 B/op	    3782 allocs/op
BenchmarkSolver_cities10-8   	    1000	   2463331 ns/op	  327866 B/op	    5980 allocs/op
BenchmarkSolver_cities12-8   	     500	   3474088 ns/op	  435976 B/op	    7113 allocs/op
BenchmarkSolver_cities15-8   	     300	   4642179 ns/op	  599064 B/op	    8368 allocs/op
BenchmarkSolver_cities20-8   	     100	  10288453 ns/op	 1268070 B/op	   13834 allocs/op
BenchmarkSolver_cities25-8   	     100	  15348473 ns/op	 1370101 B/op	   13537 allocs/op
BenchmarkSolver_cities30-8   	      50	  23542900 ns/op	 1911450 B/op	   16515 allocs/op
BenchmarkSolver_cities40-8   	      30	  47275640 ns/op	 4236715 B/op	   28784 allocs/op
BenchmarkSolver_cities50-8   	      10	 110631608 ns/op	 7182500 B/op	   39742 allocs/op
BenchmarkSolver_cities60-8   	      10	 191942666 ns/op	12534673 B/op	   57845 allocs/op
BenchmarkSolver_cities70-8   	       2	 550705384 ns/op	12568492 B/op	   56761 allocs/op
BenchmarkSolver_cities80-8   	       3	 479881046 ns/op	17616229 B/op	   70559 allocs/op
BenchmarkSolver_cities90-8   	       2	 514134093 ns/op	32496200 B/op	  107364 allocs/op
BenchmarkSolver_cities95-8   	       2	 909439964 ns/op	32929056 B/op	  106136 allocs/op
PASS
ok  	github.com/xaionaro-go/algorithms/tsp/solver/approximate	29.028s
```

Keep in mind:
* For every city it solves 8 different tasks (so it show performance in 8 times less than real)
* The most of the performance is utilized by normalization process (I was solving more common problem than classic TSP: there're just random one-directional routes from random cities to random cities; so for example costs of travels from specific cities to specific cities are unknown).
* There's no simplifications in the problem (like clusters of cities). Just a cloud of cities in non-metric space.
* Theoretically it should calculate quite precise.

Performance optimizations:
* Pre-calculation of costs from every city to every city (this part is the most CPU utilizing one)
* Boundaries (for example an estimation of the maximal cost by a stupid solver) to reduce useless branching (of the bruteforcing process)
* Estimate maximal cost by a stupid solver 
* Sort routes by some stupid heuristics to reduce useless branching (of the bruteforcing process)
* Ant-like markers (and re-sort routes) to reduce useless branching (of the bruteforcing process)
* Global path approach is done by heuristics (choose the cheapest route to a city the most far from the finish-city)
* Bruteforce only on length not more than 8 cities.

Precision optimizations:
* Sort routes by some stupid heuristics to pick the most useful routes.
* Ant-like markers (and re-sort routes) to pick the most useful routes.
* Pick optimal solution for every sequential 8 cities (by bruteforce).
* Optimize local solutions using the same heuristics method which was use for the global solution (for example it help with problems like "drilling route").
