## Embedded databases comparison

This project was developed to do a performance comparison of various embedded databases used for caching application. 
The selection of the databases was based on the status (active/deactive), application and performance.

This repository contains embedded store abstraction and implementations for Go.


#### Databases

- [badgerdb](https://github.com/dgraph-io/badger)
- [bigcache](https://github.com/allegro/bigcache)
- [moss](https://github.com/couchbase/moss)
- [nutsdb](https://github.com/xujiajun/nutsdb)
- [pudge](https://github.com/recoilme/pudge)
- [ristretto](https://github.com/dgraph-io/ristretto)
- In-memory (with Go map)

#### Interface
All the databases implement this interface.
```go
type Store interface {
    Set(k string, v interface{}) error
    Get(k string, v interface{}) (found bool, err error)
    Delete(k string) error
    Close() error
}
```

#### Benchmarks

| Database | counter | ns/op | B/op |allocs/op | counter | ns/op | B/op |allocs/op | counter | ns/op | B/op |allocs/op |
| --- | --- |--- | --- |--- | --- |--- |--- | --- |--- | --- |--- |--- | --- |--- | --- |
| ristretto	| 1366906	|767	  |545	 |5	   |143102	|7808	     |1004	  |20	 |1645538 |692	       |143	    |4
| bigcache	| 377283	|2670	  |748	 |5	   |200302	|6082	     |602	  |19	 |5419498 |246	       |16	    |2
| In-Memory	| 511095	|2677	  |538	 |5	   |157450	|6418	     |371	  |17	 |3667006 |332	       |80	    |4
| moss	    | 79590	    |14346	  |4994	 |54   |7750376	|166	     |1	      |0	 |9174933 |150	       |90	    |0
| pudge	    | 27457	    |48396	  |850	 |22   |117904	|8647	     |645	  |21	 |3468951 |301	       |40	    |4
| badgerdb	| 7477	    |168751	  |15363 |79   |124051	|10153	     |1791	  |28	 |7509	  |185560      | 14536	|77
| nutsdb	| 108	    |11160554 |47596 |722  |1	    |10198412296 |4357920 |67041 |1	      |10100557116 |8148224	|125802