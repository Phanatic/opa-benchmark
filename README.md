## Benchmarking Open Policy Agent for use within Vitess. 

I made this repository to hold benchmark tests for using Open Policy Agent (OPA) for authorizing vitess queries.

### Latest results
``` bash
go test -bench=.   -benchtime=1s -benchmem
goos: linux
goarch: amd64
pkg: github.com/planetscale/opa
cpu: AMD Ryzen 7 2700 Eight-Core Processor
BenchmarkOPA_PrepareAndEvaluate-16                                      	     498	   2410280 ns/op	  483325 B/op	   11493 allocs/op
BenchmarkOPA_EvaluateAllowAll_UsingPreparedQuery-16                     	   28084	     45568 ns/op	   10274 B/op	     219 allocs/op
BenchmarkOPA_EvaluateUsingPreparedQuery-16                              	    7344	    167886 ns/op	   31255 B/op	     693 allocs/op
BenchmarkOPA_Evaluate_PreparedQuery_PreparedInput_NoColumns-16          	   27241	     43746 ns/op	   10506 B/op	     203 allocs/op
BenchmarkOPA_Evaluate_PreparedQuery_PreparedInput_SingleColumn-16       	   11400	    111760 ns/op	   22072 B/op	     458 allocs/op
BenchmarkOPA_Evaluate_PreparedQuery_PreparedInput_MultipleColumns-16    	    7416	    162198 ns/op	   30228 B/op	     660 allocs/op
BenchmarkOPA_Evaluate_PreparedQuery_EvaledInput_MultipleColumns-16      	    7501	    151957 ns/op	   30211 B/op	     659 allocs/op
Benchmark_Vitess_Classic_Authorized-16                                  	 8110430	       147.7 ns/op	      32 B/op	       1 allocs/op
Benchmark_Vitess_Classic_Prepared_Authorized-16                         	106023502	        13.88 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/planetscale/opa	16.961s
```

### Running this benchmark locally 
``` bash
    go test -bench=.   -benchtime=1s -benchmem
```

