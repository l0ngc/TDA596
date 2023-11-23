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
	mapTasks    []MapTask
	reduceTasks []ReduceTask
	nReduce     int
}

// Map Task
type MapTask struct {
	FileId   int
	FileName string
	Done     bool
}

type ReduceTask struct {
	ReduceID int
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
	if len(c.mapTasks) == 0 {
		// return errors.New("no more tasks available")
		return nil
	}
	for _, task := range c.mapTasks {
		if !task.Done {
			fmt.Println(task.FileId, task.FileName, task.Done)
			reply.FileId = task.FileId
			reply.FileName = task.FileName
			fmt.Println("Assigned task:", reply.FileName)
			return nil
		}
	}
	// return errors.New("no pending tasks available")
	return nil
}

func (c *Coordinator) GetReduceTask(args *ExampleArgs, reply *ExampleReply) error {
	// 确保还有任务可分配，然后返回reduceIdx
	if len(c.reduceTasks) == 0 {
		return errors.New("no more tasks available")
	}
	for _, task := range c.reduceTasks {
		if !task.Done {
			reply.Y = task.ReduceID
			return nil
		}
	}
	return errors.New("no pending tasks available")
}

func (c *Coordinator) DoneMapTask(args *ExampleArgs, reply *ExampleReply) error {
	doneFileId := args.X
	fmt.Println("Done Map task:", doneFileId)
	for i, task := range c.mapTasks {
		if task.FileId == doneFileId {
			c.mapTasks[i].Done = true
			reply.Y = 0
			return nil
		}
	}
	reply.Y = 1
	return nil
}

func (c *Coordinator) DoneReduceTask(args *ExampleArgs, reply *ExampleReply) error {
	doneReduceId := args.X
	fmt.Println("Done reduce task:", doneReduceId)
	for i, task := range c.reduceTasks {
		if task.ReduceID == doneReduceId {
			c.reduceTasks[i].Done = true
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
func (c *Coordinator) GetNMap(args *ExampleArgs, reply *ExampleReply) error {
	// return the num of undone tasks
	num := 0
	for _, task := range c.mapTasks {
		if !task.Done {
			num++
		}
	}
	reply.Y = num
	return nil
}

func (c *Coordinator) GetStatus(args *TaskAsk, reply *CoorStatusReply) error {
	// 确保还有任务可分配
	// NReduce           int
	// NFile             int
	// NUndoneMapTask    int
	// NUndoneReduceTask int
	reply.NReduce = c.nReduce
	reply.NFile = len(c.mapTasks)
	reply.NUndoneMapTask = 0
	reply.NUndoneReduceTask = 0
	for _, task := range c.mapTasks {
		if !task.Done {
			reply.NUndoneMapTask++
		}
	}
	for _, task := range c.reduceTasks {
		if !task.Done {
			reply.NUndoneReduceTask++
		}
	}
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
		nReduce:     nReduce,
		mapTasks:    make([]MapTask, len(files)),
		reduceTasks: make([]ReduceTask, nReduce),
	}
	// init mapTasks
	for i, file := range files {
		c.mapTasks[i] = MapTask{FileId: i, FileName: file, Done: false}
	}
	// init reduceTasks
	for i := 0; i < nReduce; i++ {
		c.reduceTasks[i] = ReduceTask{ReduceID: i, Done: false}
	}

	c.server() // Assuming this is a method that starts the server.
	return &c
}
