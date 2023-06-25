# Benchmarks of Go serialization methods

[https://github.com/kubeservice-stack/serialization-benchmarks](https://github.com/kubeservice-stack/serialization-benchmarks)

## Tested serialization methods

- encoding/gob
- encoding/json
- encoding/xml
- [github.com/vmihailenco/msgpack/v5](https://github.com/vmihailenco/msgpack)
- [labix.org/v2/mgo/bson](https://github.com/valyala/fastjson)
- [github.com/valyala/fastjson](https://github.com/valyala/fastjson)
- [github.com/json-iterator/go](https://github.com/json-iterator/go)
- [github.com/kubeservice-stack/kspack-go](https://github.com/kubeservice-stack/kspack-go)

## Running the benchmarks

```bash
go get -u -t
go test -bench='.*' ./
```

## Data 

```
type A struct {
    Name     string
    BirthDay time.Time
    Phone    string
    Siblings int
    Spouse   bool
    Money    float64
}
```

## Results

2023-02-09 Results with Go 1.19.4, MacPro M2 OS 
```
goos: darwin
goarch: amd64
pkg: github.com/kubeservice-stack/serialization-benchmarks
cpu: VirtualApple @ 2.50GHz
Benchmark_Bson_Marshal-8         	 2272702	       499.2 ns/op	       110.0 B/serial	     376 B/op	      10 allocs/op
Benchmark_Bson_Unmarshal-8       	 1617357	       716.2 ns/op	       110.0 B/serial	     224 B/op	      19 allocs/op
Benchmark_FastJson_Marshal-8     	 4016576	       296.9 ns/op	       133.8 B/serial	     504 B/op	       6 allocs/op
Benchmark_FastJson_Unmarshal-8   	 1527366	       778.7 ns/op	       133.8 B/serial	    1704 B/op	       9 allocs/op
Benchmark_Gob_Marshal-8          	  467408	      2475 ns/op	       163.6 B/serial	    1616 B/op	      35 allocs/op
Benchmark_Gob_Unmarshal-8        	   95322	     12386 ns/op	       163.6 B/serial	    7688 B/op	     207 allocs/op
Benchmark_JsonIter_Marshal-8     	 1941692	       591.5 ns/op	       148.7 B/serial	     216 B/op	       3 allocs/op
Benchmark_JsonIter_Unmarshal-8   	 1307601	       871.8 ns/op	       148.7 B/serial	     359 B/op	      14 allocs/op
Benchmark_Json_Marshal-8         	 1406618	       850.8 ns/op	       148.6 B/serial	     208 B/op	       2 allocs/op
Benchmark_Json_Unmarshal-8       	  494430	      2153 ns/op	       148.6 B/serial	     399 B/op	       9 allocs/op
Benchmark_xxxPack_Marshal-8      	 3624224	       312.8 ns/op	       119.0 B/serial	     344 B/op	       6 allocs/op
Benchmark_xxxPack_Unmarshal-8    	 2070897	       574.7 ns/op	       119.0 B/serial	     112 B/op	       3 allocs/op
Benchmark_Msgpack_Marshal-8      	 2961548	       401.9 ns/op	        92.00 B/serial	     264 B/op	       4 allocs/op
Benchmark_Msgpack_Unmarshal-8    	 2209035	       536.9 ns/op	        92.00 B/serial	     160 B/op	       4 allocs/op
PASS
ok  	github.com/kubeservice-stack/serialization-benchmarks	23.928s

```

2023-02-12 Results with Go 1.19.4, MacPro M2 OS 
```
goos: darwin
goarch: amd64
pkg: github.com/kubeservice-stack/serialization-benchmarks
cpu: VirtualApple @ 2.50GHz
Benchmark_Bson_Marshal-8         	 2397604	       495.9 ns/op	       110.0 B/serial	     376 B/op	      10 allocs/op
Benchmark_Bson_Unmarshal-8       	 1670956	       713.6 ns/op	       110.0 B/serial	     224 B/op	      19 allocs/op
Benchmark_FastJson_Marshal-8     	 3992888	       305.9 ns/op	       133.7 B/serial	     504 B/op	       6 allocs/op
Benchmark_FastJson_Unmarshal-8   	 1548738	       775.3 ns/op	       133.8 B/serial	    1704 B/op	       9 allocs/op
Benchmark_Gob_Marshal-8          	  453764	      2532 ns/op	       163.6 B/serial	    1616 B/op	      35 allocs/op
Benchmark_Gob_Unmarshal-8        	   96286	     12224 ns/op	       163.6 B/serial	    7688 B/op	     207 allocs/op
Benchmark_JsonIter_Marshal-8     	 2015677	       590.0 ns/op	       148.7 B/serial	     216 B/op	       3 allocs/op
Benchmark_JsonIter_Unmarshal-8   	 1393594	       858.8 ns/op	       148.7 B/serial	     359 B/op	      14 allocs/op
Benchmark_Json_Marshal-8         	 1441102	       833.8 ns/op	       148.7 B/serial	     208 B/op	       2 allocs/op
Benchmark_Json_Unmarshal-8       	  553723	      2141 ns/op	       148.6 B/serial	     399 B/op	       9 allocs/op
Benchmark_xxxPack_Marshal-8      	 3709597	       317.6 ns/op	       119.0 B/serial	     344 B/op	       6 allocs/op
Benchmark_xxxPack_Unmarshal-8    	 2075898	       575.5 ns/op	       119.0 B/serial	     112 B/op	       3 allocs/op
Benchmark_Msgpack_Marshal-8      	 2868603	       414.2 ns/op	        92.00 B/serial	     264 B/op	       4 allocs/op
Benchmark_Msgpack_Unmarshal-8    	 2157040	       553.0 ns/op	        92.00 B/serial	     160 B/op	       4 allocs/op
PASS
ok  	github.com/kubeservice-stack/serialization-benchmarks	27.058s
```