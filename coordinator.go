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
	mapTasks    []MapTask
	reduceTasks []ReduceTask
	//the number of reduce tasks to be assigned
	nReduce int
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

// get map tasks from coodinator
func (c *Coordinator) GetMapTask(args *TaskAsk, reply *MapTaskReply) error {
	// make sure there are available tasks to be assigned
	if len(c.mapTasks) == 0 {
		// return errors.New("no more tasks available")
		return nil
	}
	// assign map tasks in order
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

// get reduce tasks from coordinator
func (c *Coordinator) GetReduceTask(args *ExampleArgs, reply *ExampleReply) error {
	// make sure there are available tasks to be assigned and return reduceIdx
	if len(c.reduceTasks) == 0 {
		return errors.New("no more tasks available")
	}

	// assign reduce tasks in order
	for _, task := range c.reduceTasks {
		if !task.Done {
			reply.Y = task.ReduceID
			return nil
		}
	}
	return errors.New("no pending tasks available")
}

// check if all map tasks are done
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

// check if all reduce tasks are done
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

// get the number of reduce tasks
func (c *Coordinator) GetNReduce(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = c.nReduce
	return nil
}

// get the number of map tasks
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
	// make sure there are available tasks to be assigned
	// NReduce           int
	// NFile             int
	// NUndoneMapTask    int
	// NUndoneReduceTask int
	reply.NReduce = c.nReduce
	reply.NFile = len(c.mapTasks)
	reply.NUndoneMapTask = 0
	reply.NUndoneReduceTask = 0
	// count the numbers of done map tasks
	for _, task := range c.mapTasks {
		if !task.Done {
			reply.NUndoneMapTask++
		}
	}

	// count the numbers of done reduce tasks
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
	if c.NUndoneReduceTask == c.nReduce {
		ret = true
	}
	return ret
}

// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to be assigned.
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
