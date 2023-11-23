package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import (
	"os"
	"strconv"
)

//
// example to show how to declare the arguments
// and reply for an RPC.
//

type ExampleArgs struct {
	X int
}

type ExampleReply struct {
	Y int
}

// Add your RPC definitions here.
type TaskAsk struct {
	X int // task nums
}

type TaskReply struct {
	FileName string
}

type MapTaskReply struct {
	FileId   int
	FileName string
}

type ReduceTaskReply struct {
	FileId int
}

type CoorStatusReply struct {
	NReduce           int
	NFile             int
	NUndoneMapTask    int
	NUndoneReduceTask int
}

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	s := "/var/tmp/5840-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
