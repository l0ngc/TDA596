package mr

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

type Coordinator struct {
	// Your definitions here.
	taskFiles []MapTask
	nReduce   int
}

// Map Task
type MapTask struct {
	FileId   int
	FileName string
	Done     bool
}

// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 100
	return nil
}

func (c *Coordinator) GetMapTask(args *TaskAsk, reply *MapTaskReply) error {
	// 确保还有任务可分配
	if len(c.taskFiles) == 0 {
		return errors.New("no more tasks available")
	}
	// for _, task := range c.taskFiles {
	// 	fmt.Println(task.FileId, task.FileName, task.Done)
	// }
	// 分配第一个未完成的任务
	for _, task := range c.taskFiles {
		if !task.Done {
			fmt.Println(task.FileId, task.FileName, task.Done)
			reply.FileId = task.FileId
			reply.FileName = task.FileName
			fmt.Println("Assigned task:", reply.FileName)
			return nil
		}
	}

	return errors.New("no pending tasks available")
}

func (c *Coordinator) DoneMapTask(args *ExampleArgs, reply *ExampleReply) error {
	doneFileId := args.X
	fmt.Println("Done task:", doneFileId)
	for i, task := range c.taskFiles {
		if task.FileId == doneFileId {
			c.taskFiles[i].Done = true
			reply.Y = 0
			return nil
		}
	}
	reply.Y = 1
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

	return ret
}

// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{
		nReduce:   nReduce,
		taskFiles: make([]MapTask, len(files)),
	}
	// init taskFiles
	for i, file := range files {
		c.taskFiles[i] = MapTask{FileId: i, FileName: file, Done: false}
	}

	c.server() // Assuming this is a method that starts the server.
	return &c
}

// func MakeCoordinator(files []string, nReduce int) *Coordinator {
// 	c := Coordinator{}
// 	c.nReduce = nReduce
// 	c.taskFiles = files

// 	c.server()
// 	return &c
// }
