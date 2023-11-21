*help me generate one readme*
I am a readme generator
## Table of Contents
* [Installation](#installation)
* [Usage](#usage)

I am here to think through the logic of current project

I need to write one Map Reduce for this project

Sequential Approach
```
$ cd ~/6.5840
$ cd src/main
$ go build -buildmode=plugin ../mrapps/wc.go
$ rm mr-out*
$ go run mrsequential.go wc.so pg*.txt
$ more mr-out-0
A 509
ABOUT 2
ACT 8
```
Basic idea here is, you write one plugin, it contains the right function, then use the main function to call it.

Distributed Approach
```
$ go build -buildmode=plugin ../mrapps/wc.go
$ rm mr-out* 
$ go run -race mrcoordinator.go pg-*.txt
$ go run -race mrworker.go wc.so
$ cat mr-out-* | sort | more
A 509
ABOUT 2
ACT 8
...
```

Target files
- mr/coordinator.go
- mr/worker.go
- mr/rpc.go.

We have given you a little code to start you off. 

The "main" routines for the coordinator and worker are in main/mrcoordinator.go and main/mrworker.go; don't change these files. You should put your implementation in mr/coordinator.go, mr/worker.go, and mr/rpc.go.