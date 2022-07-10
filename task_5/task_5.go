package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"io/ioutil"
	"os"
	"strconv"
	proto2 "vezdecodeProject_task_5/proto"
)

func remove(slice [][]byte, s int) [][]byte {
	return append(slice[:s], slice[s+1:]...)
}

func comparePb(files [][]byte, m proto.Message) {
	for i, file := range files {
		err := proto.Unmarshal(file, m)
		if err == nil {
			files = remove(files, i)
			fmt.Println(protojson.Format(m.(protoreflect.ProtoMessage)) + "\n")
			break
		}
	}
}

func main() {

	var files [][]byte
	for i := 1; i <= 4; i++ {
		bz, err := ioutil.ReadFile("./pb/example" + strconv.Itoa(i) + ".pb")
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		files = append(files, bz)
	}

	var arr []proto.Message
	arr = append(arr, &proto2.Teams{})
	arr = append(arr, &proto2.Person{})
	arr = append(arr, &proto2.Points{})
	arr = append(arr, &proto2.Cities{})
	arr = append(arr, &proto2.Names{})

	for _, msg := range arr {
		comparePb(files, msg)
	}

}
