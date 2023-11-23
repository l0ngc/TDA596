package mr

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"sort"
)

// Map functions return a slice of KeyValue.
type KeyValue struct {
	Key   string
	Value string
}
type ByKey []KeyValue

func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

// main/mrworker.go calls this function.
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {
	// Your worker implementation here.
	// uncomment to send the Example RPC to the coordinator.
	// CallExample()
	corrStatus := CallCoorStatus()
	mapTask := CallMapTask()
	currentFileName := mapTask.FileName
	// if map is not finished
	fmt.Println(corrStatus.NUndoneMapTask)

	if corrStatus.NUndoneMapTask > 0 {
		intermediate := []KeyValue{}
		// read file
		file, err := os.Open(currentFileName)
		if err != nil {
			log.Fatalf("cannot open %v", currentFileName)
		}
		defer file.Close()

		content, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatalf("cannot read %v", currentFileName)
		}
		// mapf
		kva := mapf(currentFileName, string(content))
		intermediate = append(intermediate, kva...)
		// set up intermediate files
		reduceOutputFiles := make(map[int]*os.File)
		encoders := make(map[int]*json.Encoder)

		for i := 0; i < corrStatus.NReduce; i++ {
			filename := fmt.Sprintf("intermediate_%d_%d.json", mapTask.FileId, i)
			outputFile, err := os.Create(filename)
			if err != nil {
				log.Fatalf("cannot create output file %v", filename)
			}
			defer outputFile.Close()
			reduceOutputFiles[i] = outputFile

			encoders[i] = json.NewEncoder(outputFile)
		}
		// write to intermediate files
		for _, kv := range intermediate {
			reduceTaskNumber := ihash(kv.Key) % corrStatus.NReduce
			enc := encoders[reduceTaskNumber]

			if err := enc.Encode(&kv); err != nil {
				log.Fatalf("error encoding JSON: %v", err)
			}
		}
		// done map task
		isDone := CallMapDone(mapTask.FileId)
		if isDone == 0 {
			fmt.Println("Map Task Done")
		}
	} else {
		// start reduce task
		reduceTaskId := CallReduceTask()
		fmt.Println("Start Reduce Task", reduceTaskId)
		intermediate := loadIntermediateData(reduceTaskId, corrStatus.NFile)
		sort.Sort(ByKey(intermediate))
		oname := fmt.Sprintf("mr-out-%d", reduceTaskId)
		ofile, _ := os.Create(oname)
		i := 0
		for i < len(intermediate) {
			j := i + 1
			for j < len(intermediate) && intermediate[j].Key == intermediate[i].Key {
				j++
			}
			values := []string{}
			for k := i; k < j; k++ {
				values = append(values, intermediate[k].Value)
			}
			output := reducef(intermediate[i].Key, values)

			// this is the correct format for each line of Reduce output.
			fmt.Fprintf(ofile, "%v %v\n", intermediate[i].Key, output)

			i = j
		}
		ofile.Close()
		CallReduceDone(reduceTaskId)
	}
}

func loadIntermediateData(reduceTaskNumber int, nFiles int) []KeyValue {
	intermediate := []KeyValue{}

	for i := 0; i < nFiles; i++ {
		filename := fmt.Sprintf("intermediate_%d_%d.json", i, reduceTaskNumber)
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("cannot open intermediate file %v", filename)
		}

		dec := json.NewDecoder(file)
		for {
			var kv KeyValue
			if err := dec.Decode(&kv); err != nil {
				break
			}
			intermediate = append(intermediate, kv)
		}

		file.Close()
	}

	return intermediate
}

// example function to show how to make an RPC call to the coordinator.
//
// the RPC argument and reply types are defined in rpc.go.
func CallReduceTask() int {
	args := ExampleArgs{}
	args.X = 99
	reply := ExampleReply{}
	ok := call("Coordinator.GetReduceTask", &args, &reply)
	if ok {
		// reply.Y should be 100.
		fmt.Printf("N Reduce is %v\n", reply.Y)
	} else {
		fmt.Printf("call failed!\n")
	}
	return reply.Y
}

func CallCoorStatus() CoorStatusReply {
	// declare an argument structure.
	args := TaskAsk{}
	args.X = 1
	// declare a reply structure.
	reply := CoorStatusReply{}

	ok := call("Coordinator.GetStatus", &args, &reply)
	if ok {
		// reply.Y should be 100.
		fmt.Println(reply.NReduce, reply.NFile, reply.NUndoneMapTask, reply.NUndoneReduceTask)
		return reply
	} else {
		fmt.Printf("call failed!\n")
		return CoorStatusReply{}
	}
}

func CallMapTask() MapTaskReply {
	// declare an argument structure.
	args := TaskAsk{}
	args.X = 1
	// declare a reply structure.
	reply := MapTaskReply{}

	ok := call("Coordinator.GetMapTask", &args, &reply)
	if ok {
		// reply.Y should be 100.
		fmt.Println(reply.FileName)
		return reply
	} else {
		fmt.Printf("call failed!\n")
		return MapTaskReply{}
	}
}
func CallMapDone(taskId int) int {
	// send text id means it is done
	args := ExampleArgs{}
	args.X = taskId
	reply := ExampleReply{}
	ok := call("Coordinator.DoneMapTask", &args, &reply)
	if ok {
		// reply.Y should be 100.
		fmt.Printf("reply.Y %v\n", reply.Y)
	} else {
		fmt.Printf("call failed!\n")
	}
	return reply.Y
}

func CallReduceDone(taskId int) int {
	// send text id means it is done
	args := ExampleArgs{}
	args.X = taskId
	reply := ExampleReply{}
	ok := call("Coordinator.DoneReduceTask", &args, &reply)
	if ok {
		// reply.Y should be 100.
		fmt.Printf("reply.Y %v\n", reply.Y)
	} else {
		fmt.Printf("call failed!\n")
	}
	return reply.Y
}

func CallNReduce() int {
	// declare an argument structure.
	args := ExampleArgs{}

	// // fill in the argument(s).
	args.X = 99

	// // declare a reply structure.
	reply := ExampleReply{}

	ok := call("Coordinator.GetNReduce", &args, &reply)
	if ok {
		// reply.Y should be 100.
		fmt.Printf("N Reduce is %v\n", reply.Y)
	} else {
		fmt.Printf("call failed!\n")
	}
	return reply.Y
}

func CallNMap() int {
	// declare an argument structure.
	args := ExampleArgs{}

	// // fill in the argument(s).
	args.X = 99

	// // declare a reply structure.
	reply := ExampleReply{}

	ok := call("Coordinator.GetNMap", &args, &reply)
	if ok {
		// reply.Y should be 100.
		fmt.Printf("N Reduce is %v\n", reply.Y)
	} else {
		fmt.Printf("call failed!\n")
	}
	return reply.Y
}

func CallExample() {

	args := ExampleArgs{}
	args.X = 99
	reply := ExampleReply{}
	ok := call("Coordinator.Example", &args, &reply)
	if ok {
		// reply.Y should be 100.
		fmt.Printf("reply.Y %v\n", reply.Y)
	} else {
		fmt.Printf("call failed!\n")
	}
}
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := coordinatorSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}

// // test json read
// inputFile, err := os.Open("doodle") // Change "doodle" to your output file name
// if err != nil {
// 	log.Fatalf("cannot open output file for reading")
// }
// defer inputFile.Close()

// kva := []KeyValue{}
// dec := json.NewDecoder(inputFile)
// for {
// 	var kv KeyValue
// 	if err := dec.Decode(&kv); err != nil {
// 		break
// 	}
// 	kva = append(kva, kv)
// }

// // read the intermediate files and reduce them
// for reduceTaskNumber := 0; reduceTaskNumber < NReduce; reduceTaskNumber++ {
// 	intermediate := loadIntermediateData(reduceTaskNumber)
// 	sort.Sort(ByKey(intermediate))
// 	oname := fmt.Sprintf("mr-out-%d", reduceTaskNumber)
// 	ofile, _ := os.Create(oname)
// 	i := 0
// 	for i < len(intermediate) {
// 		j := i + 1
// 		for j < len(intermediate) && intermediate[j].Key == intermediate[i].Key {
// 			j++
// 		}
// 		values := []string{}
// 		for k := i; k < j; k++ {
// 			values = append(values, intermediate[k].Value)
// 		}
// 		output := reducef(intermediate[i].Key, values)

// 		// this is the correct format for each line of Reduce output.
// 		fmt.Fprintf(ofile, "%v %v\n", intermediate[i].Key, output)

// 		i = j
// 	}
// 	ofile.Close()
// }
