package mr

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

type Coordinator struct {
	// Your definitions here.
	taskFiles []string
	nReduce   int
}

// Your code here -- RPC handlers for the worker to call.

// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 100
	return nil
}
func (c *Coordinator) GetTask(args *TaskAsk, reply *TaskReply) error {
	// Your code here
	// start from Map, when all of the Map is finished, then turned into Reduce
	// now just put ev
	reply.FileNames = c.taskFiles
	fmt.Println(reply.FileNames)
	// c.taskFiles = c.taskFiles[1:]
	return nil
}

func (c *Coordinator) GetNReduce(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = c.nReduce
	return nil
}

// start a thread that listens for RPCs from worker.go
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
func (c *Coordinator) Done() bool {
	ret := false

	// Your code here.
	// 这里用来标记一下目前搞定的任务，我觉得应该是coordinator分布任务，然后worker完成任务，然后worker告诉coordinator完成了任务
	// 然后coordinator就把这个任务标记为完成了，然后等到所有任务都完成了，就返回true
	// 所以应该有一个标记，来标记要被map的任务，同时标记要被reduce的任务
	return ret
}

// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{}
	c.nReduce = nReduce
	c.taskFiles = files

	c.server()
	return &c
}
