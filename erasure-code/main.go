package main

import (
	"fmt"
	"github.com/klauspost/reedsolomon"
	"io/ioutil"
	"os"
	"path"
)

var (
	encoder, _ = reedsolomon.New(5, 3)
	fileName   = "./test.txt"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Encode() {
	testData, err := ioutil.ReadFile(fileName)
	checkErr(err)
	encodeFile, err := encoder.Split(testData)
	checkErr(err)
	err = encoder.Encode(encodeFile)
	checkErr(err)
	err = os.MkdirAll("./", 0777)
	checkErr(err)

	for index, data := range encodeFile {
		err = ioutil.WriteFile(path.Join("./", fmt.Sprintf("encode_%d", index)), data, 0644)
		checkErr(err)
	}
	fmt.Println("encode successful")

}

func Decode() {
	shards := make([][]byte, 8)
	var missingShards []int
	for i := 0; i < 8; i++ {
		encodeFile := path.Join("./", fmt.Sprintf("encode_%d", i))
		data, err := ioutil.ReadFile(encodeFile)
		if err == nil {
			shards[i] = data
		} else if os.IsNotExist(err) {
			missingShards = append(missingShards, i)
			continue
		} else {
			panic(err)
		}

	}
	err := encoder.Reconstruct(shards)
	checkErr(err)

	f, err := os.OpenFile("./test.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	checkErr(err)
	dataSize := 0
	for i := 0; i < 5; i++ {
		dataSize += len(shards[i])
	}
	err = encoder.Join(f, shards, dataSize)
	checkErr(err)
	fmt.Println("decode successful")
}

func main() {
	//Encode()
	Decode()
}
