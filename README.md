# Overview
Go katas, package experiment, unit testing and memory management analysis. 

# Package Tests in the katas module
## basics
- shapes_test.go > TestShapeStruct()
``` BASH
go test -v -run TestShapeStruct ./basics
```

## memory
- maps_leak_test.go > TestMapsLeak()
``` BASH
go test -v -run TestMapsLeak ./memory
```
- slice_leak_test.go > TestSliceLeak()
``` BASH
go test -v -run TestSliceLeak ./memory
```

## leetcode
- bfs_test.go > TestMaxLevelSum() 
``` BASH 
go test -v -run TestMaxLevelSum ./leetcode
```
- dfs_test.go > TestMaxDepth()
``` BASH 
go test -v -run TestMaxDepth ./leetcode
```

