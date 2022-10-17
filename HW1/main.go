package main

import (
	"encoding/gob"
	"fmt"
	"math"
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

func rpcServer(mySlice IntSlice, number int) {
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
	fmt.Printf("Sending %d numbers took %d milliseconds for the rpcServer\n", number, endTime-startTime)
}

func rpcClient(number int) {
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
	fmt.Printf("Sending %d numbers took %d milliseconds for the rpcClient\n", number, endTime-startTime)
}

func gobServer(mySlice IntSlice, number int) {
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
	fmt.Printf("Sending %d numbers took %d milliseconds with printing in gobClient\n", number, endTime-startTime)
}

func gobClient(number int) {
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
	fmt.Printf("Sending %d numbers took %d milliseconds use gobClient\n", number, endTime-startTime)
}

func localfileServer(mySlice IntSlice, number int) {
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
	fmt.Printf("Sending %d numbers took %d milliseconds via localfileServer", number, endTime-startTime)
}

func localfileClient(number int) {
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
	fmt.Printf("Sending %d numbers took %d milliseconds via local fileClient\n", number, endTime-startTime)
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
	value := os.Args[2]
	max := math.MaxInt
	currentNumbers := 1000 //This is the minimum number we want to start testing
	for j := 1; j < 10; j++ {
		currentNumbers = int(math.Min(float64(currentNumbers*2), float64(max)))
		var mySlice IntSlice = make([]int, currentNumbers, currentNumbers)
		for i := 0; i < len(mySlice); i++ {
			mySlice[i] = i + i%10 //just to avoid any potential trivial encodings
		}
		if value == "listen" {
			switch num {
			case 1:
				rpcServer(mySlice, currentNumbers)
			case 2:
				gobServer(mySlice, currentNumbers)
			case 3:
				localfileServer(mySlice, currentNumbers)
			}
		} else if value == "send" {
			//startTime := time.Now().UnixMilli()
			switch num {
			case 1:
				rpcClient(currentNumbers)
			case 2:
				gobClient(currentNumbers)
			case 3:
				localfileClient(currentNumbers)
			}
		}
	}
}
