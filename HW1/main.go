package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"time"
)

type IntSlice []int

var rpcPort string = "8081"

func (mySlice IntSlice) Values(ignore int, sameSlice *IntSlice) error {
	*sameSlice = mySlice
	return nil
}

func rpcServer(mySlice IntSlice) {
	rpc.Register(mySlice)
	l, err := net.Listen("tcp", ":"+rpcPort)
	if err != nil {
		fmt.Println("Unable to connect to listener")
		return
	}
	startTime := time.Now().UnixMilli()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Unable to accept connection")
			return
		}
		rpc.ServeConn(conn) //note that this uses gob behind the scenes
	}
	endTime := time.Now().UnixMilli()
	fmt.Printf("This took %d milliseconds for the rpcServer\n", endTime-startTime)
}

func rpcClient() {
	var mySlice IntSlice
	client, err := rpc.Dial("tcp", "127.0.0.1:"+rpcPort)
	if err != nil {
		fmt.Println("Unable to connect to server")
		return
	}
	startTime := time.Now().UnixMilli()
	args := 0 //new(struct{})
	err = client.Call("IntSlice.Values", args, &mySlice)
	if err != nil {
		fmt.Println("Unable to send data")
		fmt.Println(err)
		return
	}
	endTime := time.Now().UnixMilli()
	fmt.Printf("This took %d milliseconds for the rpcClient\n", endTime-startTime)
}

func gobServer(mySlice IntSlice) {
	rpc.Register(mySlice)
	l, err := net.Listen("tcp", ":"+rpcPort)
	if err != nil {
		fmt.Println("Unable to connect to listener")
		return
	}
	startTime := time.Now().UnixMilli()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Unable to accept connection")
			return
		}
		encoder := gob.NewEncoder(conn)
		encoder.Encode(mySlice)
	}
	endTime := time.Now().UnixMilli()
	fmt.Printf("This took %d milliseconds with printing in gobClient\n", endTime-startTime)
}

func gobClient() {
	var mySlice IntSlice
	client, err := net.Dial("tcp", "127.0.0.1:"+rpcPort)
	if err != nil {
		fmt.Println("Unable to connect to server")
		return
	}
	decoder := gob.NewDecoder(client)
	startTime := time.Now().UnixMilli()
	err = decoder.Decode(&mySlice)
	if err != nil {
		fmt.Println("Unable to receieve data")
		fmt.Println(err)
		return
	}
	endTime := time.Now().UnixMilli()
	fmt.Printf("This took %d milliseconds use gobClient\n", endTime-startTime)
}

func localfileServer(mySlice IntSlice) {
	startTime := time.Now().UnixMilli()
	file, err := os.Create(os.TempDir() + "/hw1_file")
	if err != nil {
		fmt.Println("Failed to open file")
		return
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	encoder.Encode(mySlice)
	endTime := time.Now().UnixMilli()
	fmt.Printf("This took %d milliseconds via localfileServer", endTime-startTime)
}

func localfileClient() {
	startTime := time.Now().UnixMilli()
	var mySlice IntSlice
	file, err := os.Open(os.TempDir() + "/hw1_file")
	if err != nil {
		fmt.Println("Failed to open file")
		return
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&mySlice)
	if err != nil {
		fmt.Println("Unable to read data")
		fmt.Println(err)
		return
	}
	endTime := time.Now().UnixMilli()
	fmt.Printf("This took %d milliseconds via local fileClient\n", endTime-startTime)
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Incorrect format: ./main [1-3] [send/listen]")
		return
	}
	num, err := strconv.Atoi(os.Args[1])
	if err != nil || num > 3 || num < 1 {
		fmt.Println("Incorrect format: pass a number between 1 and 3")
		return
	} else if os.Args[2] != "send" && os.Args[2] != "listen" {
		fmt.Println("Incorrect format: 2nd arg should be either send or listen")
		return
	}
	var mySlice IntSlice = make([]int, 1000000, 1000000)
	for i := 0; i < len(mySlice); i++ {
		mySlice[i] = i + i%10 //just to avoid any potential trivial encodings
	}
	if os.Args[2] == "listen" {
		switch num {
		case 1:
			rpcServer(mySlice)
		case 2:
			gobServer(mySlice)
		case 3:
			localfileServer(mySlice)
		}
	} else if os.Args[2] == "send" {
		//startTime := time.Now().UnixMilli()
		switch num {
		case 1:
			rpcClient()
		case 2:
			gobClient()
		case 3:
			localfileClient()
		}
		/*endTime := time.Now().UnixMilli()
		fmt.Printf("This took %d milliseconds fucking main\n", endTime-startTime)*/
	}
}
