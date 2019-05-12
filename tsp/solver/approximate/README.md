This is quite precise approximate method.

```
goos: linux
goarch: amd64
pkg: github.com/xaionaro-go/algorithms/tsp/solver/approximate
BenchmarkSolver_cities4-8     	   10000	    171970 ns/op	   41153 B/op	     992 allocs/op
BenchmarkSolver_cities8-8     	    1000	   1877874 ns/op	  270129 B/op	    4916 allocs/op
BenchmarkSolver_cities10-8    	     500	   3536111 ns/op	  443241 B/op	    7712 allocs/op
BenchmarkSolver_cities12-8    	     300	   4168368 ns/op	  487095 B/op	    8461 allocs/op
BenchmarkSolver_cities15-8    	     300	   6145205 ns/op	  691901 B/op	   11313 allocs/op
BenchmarkSolver_cities20-8    	     100	  13309341 ns/op	 1433825 B/op	   20734 allocs/op
BenchmarkSolver_cities25-8    	     100	  18028180 ns/op	 1517034 B/op	   21541 allocs/op
BenchmarkSolver_cities30-8    	      50	  34088634 ns/op	 2687221 B/op	   34260 allocs/op
BenchmarkSolver_cities40-8    	      20	  68578759 ns/op	 5160452 B/op	   57729 allocs/op
BenchmarkSolver_cities50-8    	      20	 102253111 ns/op	 4785423 B/op	   54583 allocs/op
BenchmarkSolver_cities60-8    	       5	 206331794 ns/op	11582904 B/op	  104890 allocs/op
BenchmarkSolver_cities70-8    	       3	 453716559 ns/op	10964074 B/op	   96990 allocs/op
BenchmarkSolver_cities80-8    	       3	 477341039 ns/op	15288808 B/op	  129248 allocs/op
BenchmarkSolver_cities90-8    	       3	 521593311 ns/op	27950917 B/op	  201086 allocs/op
BenchmarkSolver_cities95-8    	       2	 606791975 ns/op	21001500 B/op	  162934 allocs/op
BenchmarkSolver_cities100-8   	       1	1588006830 ns/op	20248456 B/op	  158458 allocs/op
BenchmarkSolver_cities120-8   	       2	 843257428 ns/op	36470744 B/op	  234523 allocs/op
BenchmarkSolver_cities150-8   	       1	1458863518 ns/op	57951056 B/op	  326371 allocs/op
BenchmarkSolver_cities200-8   	       1	4978097888 ns/op	102456800 B/op	  495117 allocs/op
BenchmarkSolver_cities250-8   	       1	6951132766 ns/op	148696824 B/op	  670470 allocs/op
BenchmarkSolver_cities350-8   	       1	10941885671 ns/op	257584600 B/op	 1053924 allocs/op
BenchmarkSolver_cities500-8   	       1	44557145314 ns/op	746297392 B/op	 2165350 allocs/op
PASS
ok  	github.com/xaionaro-go/algorithms/tsp/solver/approximate	107.732s
```

Keep in mind:
* For every city it solves 8 different tasks (so it show performance in 8 times less than real)
* The most of the performance is utilized by normalization process (I was solving more common problem than classic TSP: there're just random one-directional routes from random cities to random cities; so for example costs of travels from specific cities to specific cities are unknown).
* There's no simplifications in the problem (like clusters of cities). Just a cloud of cities in non-metric space.
* Theoretically it should calculate quite precise.

Performance optimizations:
* Remove useless routes and duplicate paths.
* Pre-calculation of costs from every city to every city (this part is the most CPU utilizing one)
* Boundaries (for example an estimation of the maximal cost by a stupid solver) to reduce useless branching (of the bruteforcing process)
* Estimate maximal cost by a stupid solver 
* Sort routes by some stupid heuristics to reduce useless branching (of the bruteforcing process)
* Ant-like markers (and re-sort routes) to reduce useless branching (of the bruteforcing process)
* Global path approach is done by heuristics (choose the cheapest route to a city the most far from the finish-city)
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