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
	fileTasks := CallTask()
	for _, fileName := range fileTasks {
		fmt.Println(fileName)
	}
	intermediate := []KeyValue{}
	for _, filename := range fileTasks {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("cannot open %v", filename)
		}
		content, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatalf("cannot read %v", filename)
		}
		file.Close()
		// fmt.Println(string(content))
		kva := mapf(filename, string(content))
		intermediate = append(intermediate, kva...)
	}

	sort.Sort(ByKey(intermediate))
	outputFile, err := os.Create("doodle") // Change "doodle" to your desired output file name
	if err != nil {
		log.Fatalf("cannot create output file")
	}
	defer outputFile.Close()

	enc := json.NewEncoder(outputFile)
	for _, kv := range intermediate {
		err := enc.Encode(&kv)
		if err != nil {
			log.Fatalf("error encoding JSON: %v", err)
		}
	}
	// oname := "mr-out-1"
	// ofile, _ := os.Create(oname)
	// i := 0
	// for i < len(intermediate) {
	// 	j := i + 1
	// 	for j < len(intermediate) && intermediate[j].Key == intermediate[i].Key {
	// 		j++
	// 	}
	// 	values := []string{}
	// 	for k := i; k < j; k++ {
	// 		values = append(values, intermediate[k].Value)
	// 	}
	// 	output := reducef(intermediate[i].Key, values)

	// 	// this is the correct format for each line of Reduce output.
	// 	fmt.Fprintf(ofile, "%v %v\n", intermediate[i].Key, output)

	// 	i = j
	// }

	// ofile.Close()
}

// example function to show how to make an RPC call to the coordinator.
//
// the RPC argument and reply types are defined in rpc.go.
func CallTask() []string {
	// declare an argument structure.
	args := TaskAsk{}
	args.X = 1
	// declare a reply structure.
	reply := TaskReply{}

	ok := call("Coordinator.GetTask", &args, &reply)
	if ok {
		// reply.Y should be 100.
		fmt.Println(reply.FileNames)
	} else {
		fmt.Printf("call failed!\n")
	}
	return reply.FileNames
}

func CallExample() {

	// declare an argument structure.
	args := ExampleArgs{}

	// // fill in the argument(s).
	args.X = 99

	// // declare a reply structure.
	reply := ExampleReply{}

	// send the RPC request, wait for the reply.
	// the "Coordinator.Example" tells the
	// receiving server that we'd like to call
	// the Example() method of struct Coordinator.
	ok := call("Coordinator.Example", &args, &reply)
	if ok {
		// reply.Y should be 100.
		fmt.Printf("reply.Y %v\n", reply.Y)
	} else {
		fmt.Printf("call failed!\n")
	}
}

// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
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
