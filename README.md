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

#### Benchmark results
<table class="tg">
<thead>
  <tr>
    <th class="tg-bobw" align=Center>Database</th>
    <th class="tg-amwm" align=Center>Counter</th>
    <th class="tg-amwm" align=Center>ns/op</th>
    <th class="tg-amwm" align=Center>B/op </th>
    <th class="tg-amwm" align=Center>allocs/op</th>
    <th class="tg-amwm" align=Center>Counter</th>
    <th class="tg-amwm" align=Center>ns/op</th>
    <th class="tg-amwm" align=Center>B/op </th>
    <th class="tg-amwm" align=Center>allocs/op</th>
    <th class="tg-amwm" align=Center>Counter</th>
    <th class="tg-amwm" align=Center>ns/op</th>
    <th class="tg-amwm" align=Center>B/op </th>
    <th class="tg-amwm" align=Center>allocs/op</th>
  </tr>
</thead>
<tbody>
  <tr>
    <td class="tg-5frq">operation</td>
    <td class="tg-5frq" colspan="4" align=Center>Write</td>
    <td class="tg-5frq" colspan="4" align=Center>Read</td>
    <td class="tg-5frq" colspan="4" align=Center>Delete</td>
  </tr>
  <tr>
    <td class="tg-8d8j" align=Center>ristretto</td>
    <td class="tg-8d8j" align=Center>1366906</td>
    <td class="tg-8d8j" align=Center>767</td>
    <td class="tg-8d8j" align=Center>545</td>
    <td class="tg-8d8j" align=Center>5</td>
    <td class="tg-8d8j" align=Center>143102</td>
    <td class="tg-8d8j" align=Center>7808</td>
    <td class="tg-8d8j" align=Center>1004</td>
    <td class="tg-8d8j" align=Center>20</td>
    <td class="tg-8d8j" align=Center>1645538</td>
    <td class="tg-8d8j" align=Center>692</td>
    <td class="tg-8d8j" align=Center>143</td>
    <td class="tg-8d8j" align=Center>4</td>
  </tr>
  <tr>
    <td class="tg-8d8j" align=Center>bigcache</td>
    <td class="tg-8d8j" align=Center>377283</td>
    <td class="tg-8d8j" align=Center>2670</td>
    <td class="tg-8d8j" align=Center>748</td>
    <td class="tg-8d8j" align=Center>5</td>
    <td class="tg-8d8j" align=Center>200302</td>
    <td class="tg-8d8j" align=Center>6082</td>
    <td class="tg-8d8j" align=Center>602</td>
    <td class="tg-8d8j" align=Center>19</td>
    <td class="tg-8d8j" align=Center>5419498</td>
    <td class="tg-8d8j" align=Center>246</td>
    <td class="tg-8d8j" align=Center>16</td>
    <td class="tg-8d8j" align=Center>2</td>
  </tr>
  <tr>
    <td class="tg-8d8j" align=Center>In-Memory</td>
    <td class="tg-8d8j" align=Center>511095</td>
    <td class="tg-8d8j" align=Center>2677</td>
    <td class="tg-8d8j" align=Center>538</td>
    <td class="tg-8d8j" align=Center>5</td>
    <td class="tg-8d8j" align=Center>157450</td>
    <td class="tg-8d8j" align=Center>6418</td>
    <td class="tg-8d8j" align=Center>371</td>
    <td class="tg-8d8j" align=Center>17</td>
    <td class="tg-8d8j" align=Center>3667006</td>
    <td class="tg-8d8j" align=Center>332</td>
    <td class="tg-8d8j" align=Center>80</td>
    <td class="tg-8d8j" align=Center>4</td>
  </tr>
  <tr>
    <td class="tg-8d8j" align=Center>moss</td>
    <td class="tg-8d8j" align=Center>79590</td>
    <td class="tg-8d8j" align=Center>14346</td>
    <td class="tg-8d8j" align=Center>4994</td>
    <td class="tg-8d8j" align=Center>54</td>
    <td class="tg-8d8j" align=Center>7750376</td>
    <td class="tg-8d8j" align=Center>166</td>
    <td class="tg-8d8j" align=Center>1</td>
    <td class="tg-8d8j" align=Center>0</td>
    <td class="tg-8d8j" align=Center>9174933</td>
    <td class="tg-8d8j" align=Center>150</td>
    <td class="tg-8d8j" align=Center>90</td>
    <td class="tg-8d8j" align=Center>0</td>
  </tr>
  <tr>
    <td class="tg-8d8j" align=Center>pudge</td>
    <td class="tg-8d8j" align=Center>27457</td>
    <td class="tg-8d8j" align=Center>48396</td>
    <td class="tg-8d8j" align=Center>850</td>
    <td class="tg-8d8j" align=Center>22</td>
    <td class="tg-8d8j" align=Center>117904</td>
    <td class="tg-8d8j" align=Center>8647</td>
    <td class="tg-8d8j" align=Center>645</td>
    <td class="tg-8d8j" align=Center>21</td>
    <td class="tg-8d8j" align=Center>3468951</td>
    <td class="tg-8d8j" align=Center>301</td>
    <td class="tg-8d8j" align=Center>40</td>
    <td class="tg-8d8j" align=Center>4</td>
  </tr>
  <tr>
    <td class="tg-8d8j" align=Center>badgerdb</td>
    <td class="tg-8d8j" align=Center>7477</td>
    <td class="tg-8d8j" align=Center>168751</td>
    <td class="tg-8d8j" align=Center>15363</td>
    <td class="tg-8d8j" align=Center>79</td>
    <td class="tg-8d8j" align=Center>124051</td>
    <td class="tg-8d8j" align=Center>10153</td>
    <td class="tg-8d8j" align=Center>1791</td>
    <td class="tg-8d8j" align=Center>28</td>
    <td class="tg-8d8j" align=Center>7509</td>
    <td class="tg-8d8j" align=Center>185560</td>
    <td class="tg-8d8j" align=Center>14536</td>
    <td class="tg-8d8j" align=Center>77</td>
  </tr>
  <tr>
    <td class="tg-8d8j" align=Center>nutsdb</td>
    <td class="tg-8d8j" align=Center>108</td>
    <td class="tg-8d8j" align=Center>11160554</td>
    <td class="tg-8d8j" align=Center>47596</td>
    <td class="tg-8d8j" align=Center>722</td>
    <td class="tg-8d8j" align=Center>1</td>
    <td class="tg-8d8j" align=Center>10198412296</td>
    <td class="tg-8d8j" align=Center>4357920</td>
    <td class="tg-8d8j" align=Center>67041</td>
    <td class="tg-8d8j" align=Center>1</td>
    <td class="tg-8d8j" align=Center>10100557116</td>
    <td class="tg-8d8j" align=Center>8148224</td>
    <td class="tg-8d8j" align=Center>125802</td>
  </tr>
</tbody>
</table>



