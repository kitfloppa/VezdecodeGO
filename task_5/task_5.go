package main

import (
	"github.com/golang/protobuf/proto"
	"os"
)

func main() {
	cit := &Cities{}

	f, _ := os.Open("./pb/example1.pb")

	bz := make([]byte, 1024)
	_, _ = f.Read(bz)

	proto.Unmarshal(bz, cit)

	i := 0
	_ = i
}
