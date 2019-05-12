Solution for: [https://www.careercup.com/question?id=5638939143045120](https://www.careercup.com/question?id=5638939143045120)
```
goos: linux
goarch: amd64
pkg: github.com/xaionaro-go/algorithms/simple_tasks/facebook/3_coins
BenchmarkFindAllCoinSums_dynamic_100000-8       	     200	   9926046 ns/op	 1010730 B/op	     580 allocs/op
BenchmarkFindAllCoinSums_increment_100000-8     	    2000	    753929 ns/op	    7052 B/op	      31 allocs/op
BenchmarkFindAllCoinSums_dynamic_1000000-8      	      10	 121426441 ns/op	 8926533 B/op	    8070 allocs/op
BenchmarkFindAllCoinSums_increment_1000000-8    	     200	   7343438 ns/op	    7034 B/op	      31 allocs/op
BenchmarkFindAllCoinSums_dynamic_10000000-8     	       1	2402371661 ns/op	126067672 B/op	   76875 allocs/op
BenchmarkFindAllCoinSums_increment_10000000-8   	      20	  71781408 ns/op	    7108 B/op	      31 allocs/op
PASS
ok  	github.com/xaionaro-go/algorithms/simple_tasks/facebook/3_coins	12.234s
```
