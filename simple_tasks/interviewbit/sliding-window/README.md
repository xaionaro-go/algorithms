```
goos: linux
goarch: amd64
pkg: github.com/xaionaro-go/algorithms/simple_tasks/interviewbit/sliding-window
BenchmarkMy100000w20000-8         	     100	  13522274 ns/op	 4977715 B/op	  100004 allocs/op
BenchmarkProposed100000w20000-8   	       1	1329149128 ns/op	 3687680 B/op	      29 allocs/op
BenchmarkMy10000w2000-8           	    2000	   1068395 ns/op	  500256 B/op	   10004 allocs/op
BenchmarkProposed10000w2000-8     	     100	  13236690 ns/op	  287992 B/op	      19 allocs/op
BenchmarkMy1000w200-8             	   20000	     86226 ns/op	   50208 B/op	    1004 allocs/op
BenchmarkProposed1000w200-8       	   10000	    140053 ns/op	   16376 B/op	      11 allocs/op
BenchmarkMy100w20-8               	  200000	      6191 ns/op	    5184 B/op	     104 allocs/op
BenchmarkProposed100w20-8         	  500000	      2660 ns/op	    2040 B/op	       8 allocs/op
BenchmarkMy10w2-8                 	 3000000	       580 ns/op	     560 B/op	      14 allocs/op
BenchmarkProposed10w2-8           	 3000000	       601 ns/op	     560 B/op	      14 allocs/op
PASS
ok  	github.com/xaionaro-go/algorithms/simple_tasks/interviewbit/sliding-window	17.691s
```
