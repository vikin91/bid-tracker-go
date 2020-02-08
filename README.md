# Bid-Tracker

## Assumptions

1. No support for DELETE on Bids (also on all other models)
2. Focus on so called _assignment functions_, i.e., functions implementing the tasks from the assignment definition:
  - Place bid
  - Get winning bid for item
  - Get all bids for item
  - Get all items on which user has bid
3. Where optimizing for time and space is in conflict, optimization for time takes priority

## Notes to reviewers

1. `/pkg/storage/storage.go` is not necessary now, but maybe useful for future performance comparison of various data structures
2. This project was started from my private API-template - this explains the folders structure. For this project, much simpler structure could be used, but I choose to stay with the template.

## Building, Running, Testing

### Quick start

```
git clone https://github.com/vikin91/bid-tracker-go
cd bid-tracker-go
make vendor
make build
make run-demo
```

Navigate to one of the following URLs:
- http://localhost:9000/api/v1/users
- http://localhost:9000/api/v1/items
- http://localhost:9000/api/v1/item/{itemID}/winner
- http://localhost:9000/api/v1/item/{itemID}/bids
- (POST) http://localhost:9000/api/v1/item/{itemID}/bids (to add bid)
- http://localhost:9000/api/v1/users/{userID}/items

Or see [swagger definition of API](swagger/api.yml)

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
