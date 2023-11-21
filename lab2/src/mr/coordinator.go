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
	// take out the N, N = fileNums tasks form c.taskFiles
	reply.FileNames = c.taskFiles[0]
	fmt.Println(reply.FileNames)
	// c.taskFiles = c.taskFiles[1:]
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

	return ret
}

// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{}

	c.taskFiles = files

	c.server()
	return &c
}

// I want to first print out the files
// here it could split files to workers, then ask the worker to read the content and do the work
// so as to reduce the time of transfering files from coordinator to workers
// intermediate := []string{}
// for _, filename := range files {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		log.Fatalf("cannot open %v", filename)
// 	}
// 	content, err := ioutil.ReadAll(file)
// 	if err != nil {
// 		log.Fatalf("cannot read %v", filename)
// 	}
// 	file.Close()
// 	intermediate = append(intermediate, string(content))
// }
// fmt.Println(intermediate)

// func MakeCoordinator(files []string, nReduce int) *Coordinator {
//     c := Coordinator{}

//     // 初始化协调器的字段
//     c.mapTasks = make([]Task, len(files))
//     c.reduceTasks = make([]Task, nReduce)
//     c.mapDone = make([]bool, len(files))
//     c.reduceDone = make([]bool, nReduce)
//     c.nReduce = nReduce

//     // 设置 Map 任务的状态为未完成
//     for i := range c.mapDone {
//         c.mapDone[i] = false
//     }

//     // 设置 Reduce 任务的状态为未完成
//     for i := range c.reduceDone {
//         c.reduceDone[i] = false
//     }

//     // 设置 Map 任务和 Reduce 任务的数量
//     c.totalMapTasks = len(files)
//     c.totalReduceTasks = nReduce

//     // 启动 RPC 服务器
//     c.server()

//     return &c
// }
