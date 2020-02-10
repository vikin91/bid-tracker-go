# Bid-Tracker

Bid-tracker is a simple API program to track users' bids on items put on auctions.

This program has been implemented based on **assignment requirements** (a.k.a _Remote Task SDE version october 2019_) document that is treated as confidential.

## Assumptions

1. No support for DELETE on Bids (also on all other models)
2. Focus on so called _assignment functions_, i.e., functions implementing the tasks from the assignment definition:
  - Place bid
  - Get winning bid for item
  - Get all bids for item
  - Get all items on which user has bid
3. Where optimizing for time and space is in conflict, optimization for time takes priority.
4. Apply the simplest solutions that work - provided the limited time.

## Notes to Reviewers

1. `/pkg/storage/storage.go` is not necessary now, but maybe useful for future performance comparison of various data structures
2. This project was started from my private API-template - this explains the folders structure. For this project, much simpler structure could be used, but I choose to stay with the template.

### Choice of Data Structures

Provided the requirements (e.g., no delete operations), I stared with a simple data structure containing three maps:

```go
type MapBiddingSystem struct {
    Items map[uuid.UUID]*models.Item
    Users map[uuid.UUID]*models.User
    Bids  map[uuid.UUID]*models.Bid
}
```

Next, I implemented most of the tests and benchmarks - especially for the four required functions as specified in the assignment description.
As a next step, I optimized the code, to obtain `O(1)` in time and space on all four required functions.

This allowed me to simplify the structures to the following:

```go
type MapBiddingSystem struct {
    Items map[uuid.UUID]*models.Item
    Users map[uuid.UUID]*models.User
}

type User struct {
    Name       string
    bids       map[uuid.UUID]*Bid
    ItemsBid   []*Item
}

type Item struct {
    Name         string
    bids         []*Bid
    WinningBid   *Bid
    MaxBidAmount float64
}
```

### Concurrency Approach

Next, I wanted to ensure, that no race conditions exist when reading/writing users and items.
Thanks to the `map`, I need to care about locks only at the level of a single user or single item<sup>[1](#foot1)</sup>.
Next, I protect the following variables with `RWMutex`es:
1. `User.ItemsBid`
2. `User.bids` (despite not being in focus of the requirements)
3. `Item.WinningBid` and `Item.MaxBidAmount`
4. `Item.bids`

**Note**:
`User.bids` could be defined as:
`bids []*Bid` instead of `bids map[uuid.UUID]*Bid`,
because each bid is unique - but there is no requirement to optimize it further.

<a name="foot1">[1]</a>: Despite possible, I exclude here the possibility of  race condition between creating a user and using it - reason: not in the scope of the four functions required in the assignment.

## Building, Running, Testing

### Quick start

```
git clone https://github.com/vikin91/bid-tracker-go
cd bid-tracker-go

make help

make vendor
make build
make run-demo
```

You may change port with environment variable `BID_PORT`. Its default value is `9000`.

Navigate to one of the following URLs:
- http://localhost:9000/api/v1/user
- http://localhost:9000/api/v1/item
- http://localhost:9000/api/v1/item/{itemID}/winner
- http://localhost:9000/api/v1/item/{itemID}/bids
- (POST) http://localhost:9000/api/v1/item/{itemID}/bids (to add bid)
- http://localhost:9000/api/v1/user/{userID}/items

For other options see `make help`.

### Testing

```
make test      # run unit tests and code coverage
make test-race # run race detector
make bench     # run benchmarks for selected functions
```

### Benchmark results

Reference measurements - creation of objects
```
goos: darwin
goarch: amd64
pkg: github.com/vikin91/bid-tracker-go/pkg/models

Benchmark_Reference_NewBid-8    	 5112609	       237 ns/op	      16 B/op	       1 allocs/op
Benchmark_Reference_NewItem-8   	 5093360	       233 ns/op	      16 B/op	       1 allocs/op
Benchmark_Reference_NewUser-8   	 4921887	       247 ns/op	      16 B/op	       1 allocs/op
```

Performance measurements - functions required in the assignment definition

(Naming: `<benchmark_name>/<problem_size>-<GOPROC>`)
```
goos: darwin
goarch: amd64
pkg: github.com/vikin91/bid-tracker-go/pkg/storage

Benchmark_PlaceBid_OneUser_OneItem-8       	             2708265	       434 ns/op	     224 B/op	       4 allocs/op

Benchmark_PlaceBid_ManyUsers_ManyItems/1-8 	             1231968	       939 ns/op	     242 B/op	       2 allocs/op
Benchmark_PlaceBid_ManyUsers_ManyItems/2-8 	             1234220	       982 ns/op	     242 B/op	       2 allocs/op
Benchmark_PlaceBid_ManyUsers_ManyItems/4-8 	             1263750	       909 ns/op	     238 B/op	       2 allocs/op
Benchmark_PlaceBid_ManyUsers_ManyItems/8-8 	             1288119	       901 ns/op	     236 B/op	       2 allocs/op
Benchmark_PlaceBid_ManyUsers_ManyItems/16-8         	 1251896	       923 ns/op	     240 B/op	       2 allocs/op
Benchmark_PlaceBid_ManyUsers_ManyItems/32-8         	 1258936	       939 ns/op	     239 B/op	       2 allocs/op
Benchmark_PlaceBid_ManyUsers_ManyItems/64-8         	 1246604	       996 ns/op	     240 B/op	       2 allocs/op
Benchmark_PlaceBid_ManyUsers_ManyItems/128-8        	 1268413	       934 ns/op	     238 B/op	       2 allocs/op
Benchmark_PlaceBid_ManyUsers_ManyItems/256-8        	 1227532	       905 ns/op	     242 B/op	       2 allocs/op

Benchmark_GetWinningBid/1-8                         	12266350	        85.3 ns/op	       0 B/op	       0 allocs/op
Benchmark_GetWinningBid/2-8                         	12791896	        97.1 ns/op	       0 B/op	       0 allocs/op
Benchmark_GetWinningBid/4-8                         	12238524	       109 ns/op	       0 B/op	       0 allocs/op
Benchmark_GetWinningBid/8-8                         	11549455	       103 ns/op	       0 B/op	       0 allocs/op
Benchmark_GetWinningBid/16-8                        	11337864	        98.3 ns/op	       0 B/op	       0 allocs/op
Benchmark_GetWinningBid/32-8                        	11126266	       106 ns/op	       0 B/op	       0 allocs/op
Benchmark_GetWinningBid/64-8                        	11474744	       112 ns/op	       0 B/op	       0 allocs/op
Benchmark_GetWinningBid/128-8                       	11461794	       121 ns/op	       0 B/op	       0 allocs/op
Benchmark_GetWinningBid/256-8                       	10705972	       112 ns/op	       0 B/op	       0 allocs/op

Benchmark_GetBidsOnItem/1-8                         	23238231	        51.1 ns/op	       0 B/op	       0 allocs/op
Benchmark_GetBidsOnItem/2-8                         	23718804	        50.8 ns/op	       0 B/op	       0 allocs/op
Benchmark_GetBidsOnItem/4-8                         	22323140	        51.4 ns/op	       0 B/op	       0 allocs/op
Benchmark_GetBidsOnItem/8-8                         	22066209	        56.2 ns/op	       0 B/op	       0 allocs/op
Benchmark_GetBidsOnItem/16-8                        	20470885	        58.1 ns/op	       0 B/op	       0 allocs/op
Benchmark_GetBidsOnItem/32-8                        	19242194	        59.8 ns/op	       0 B/op	       0 allocs/op
Benchmark_GetBidsOnItem/64-8                        	20083909	        60.6 ns/op	       0 B/op	       0 allocs/op
Benchmark_GetBidsOnItem/128-8                       	18716521	        64.7 ns/op	       0 B/op	       0 allocs/op
Benchmark_GetBidsOnItem/256-8                       	15642604	        72.4 ns/op	       0 B/op	       0 allocs/op

Benchmark_GetItemsUserHasBid/1-8                    	27043590	        40.0 ns/op	       0 B/op	       0 allocs/op
Benchmark_GetItemsUserHasBid/2-8                    	24158804	        49.5 ns/op	       0 B/op	       0 allocs/op
Benchmark_GetItemsUserHasBid/4-8                    	21926752	        55.2 ns/op	       0 B/op	       0 allocs/op
Benchmark_GetItemsUserHasBid/8-8                    	23139181	        51.3 ns/op	       0 B/op	       0 allocs/op
Benchmark_GetItemsUserHasBid/16-8                   	21908314	        54.4 ns/op	       0 B/op	       0 allocs/op
Benchmark_GetItemsUserHasBid/32-8                   	21045462	        58.2 ns/op	       0 B/op	       0 allocs/op
Benchmark_GetItemsUserHasBid/64-8                   	20114014	        57.0 ns/op	       0 B/op	       0 allocs/op
Benchmark_GetItemsUserHasBid/128-8                  	19409780	        60.7 ns/op	       0 B/op	       0 allocs/op
Benchmark_GetItemsUserHasBid/256-8                  	17945068	        64.6 ns/op	       0 B/op	       0 allocs/op
```

## Swagger API definition

The API specification can be found in `/swagger/api.yml`. To preview the specification,
1. Add a [swagger viewer](https://marketplace.visualstudio.com/items?itemName=Arjun.swagger-viewer) to _VSCode_
1. Open the `yml` file
1. Open the preview with `SHIFT + OPTION + P`

**Important!**

Swagger API definition is **focused** on the API calls **required in the assignment definition**. Other API calls may not be specified.
