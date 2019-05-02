```
goos: linux
goarch: amd64
pkg: github.com/xaionaro-go/algorithms/tsp/solver/bruteforce
BenchmarkSolver_cities4-8          20000             81795 ns/op           25234 B/op        563 allocs/op
BenchmarkSolver_cities6-8           2000            661759 ns/op           34649 B/op        652 allocs/op
BenchmarkSolver_cities8-8            100          23593019 ns/op           49415 B/op        816 allocs/op
BenchmarkSolver_cities10-8             5         336393035 ns/op           64155 B/op        918 allocs/op
BenchmarkSolver_cities11-8             1        2996688701 ns/op           81600 B/op       1032 allocs/op
BenchmarkSolver_cities12-8             1        12468774186 ns/op          92168 B/op       1098 allocs/op
BenchmarkSolver_cities13-8             1        19351610248 ns/op          94912 B/op       1106 allocs/op
PASS
ok      github.com/xaionaro-go/algorithms/tsp/solver/bruteforce 44.760s
```
