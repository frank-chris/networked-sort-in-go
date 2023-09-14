package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net"
	"os"
	"sort"
	"strconv"
	"time"

	"gopkg.in/yaml.v2"
)

type ServerConfigs struct {
	Servers []struct {
		ServerId int    `yaml:"serverId"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
	} `yaml:"servers"`
}

func readServerConfigs(configPath string) ServerConfigs {
	f, err := ioutil.ReadFile(configPath)

	if err != nil {
		log.Fatalf("could not read config file %s : %v", configPath, err)
	}

	scs := ServerConfigs{}
	err = yaml.Unmarshal(f, &scs)
	errorCheck("Error while trying to read config file", err, true)

	return scs
}

func errorCheck(msg string, err error, fatal bool) {
	if err != nil {
		if fatal {
			log.Fatal(msg, ":", err)
		} else {
			log.Println(msg, ":", err)
		}
	}
}

func sendRecord(host string, port string, record []byte, streamComplete bool) {
	var connection net.Conn
	var err error
	for i := 0; i < 10; i++ {
		address := host + ":" + port
		connection, err = net.Dial("tcp", address)
		if err == nil {
			break
		} else {
			log.Println("Error while dialing. Retrying. ", err)
			time.Sleep(time.Second / 20)
		}
	}
	defer connection.Close()

	if streamComplete {
		record = append([]byte{1}, record...)
	} else {
		record = append([]byte{0}, record...)
	}
	_, err = connection.Write(record)
	errorCheck("Error while sending record", err, false)
}

func readFromConnection(connection net.Conn, records chan<- []byte) {
	record := make([]byte, 101)
	n, err := connection.Read(record)
	errorCheck("Error while reading from connection", err, false)
	record = record[:n]
	records <- record
}

func listenForConnection(host string, port string, records chan<- []byte) {
	address := host + ":" + port
	listener, err := net.Listen("tcp", address)
	errorCheck("Error while trying to listen on "+address, err, false)
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		errorCheck("Error while trying to accept connection", err, false)
		go readFromConnection(connection, records)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) != 5 {
		log.Fatal("Usage : ./netsort {serverId} {inputFilePath} {outputFilePath} {configFilePath}")
	}

	// What is my serverId
	serverId, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Invalid serverId, must be an int %v", err)
	}
	fmt.Println("My server Id:", serverId)

	// Read server configs from file
	scs := readServerConfigs(os.Args[4])
	fmt.Println("Got the following server configs:", scs)

	records := make(chan []byte)
	for _, server := range scs.Servers {
		if serverId == server.ServerId {
			go listenForConnection(server.Host, server.Port, records)
			break
		}
	}
	time.Sleep(time.Second / 4)

	inputFilePath := os.Args[2]
	outputFilePath := os.Args[3]
	inputFile, err := os.Open(inputFilePath)
	errorCheck("Error while opening file", err, false)
	defer inputFile.Close()

	significantBits := int(math.Log2(float64(len(scs.Servers))))
	for {
		record := make([]byte, 100)
		n, err := inputFile.Read(record)
		errorCheck("Error while reading from input file", err, false)
		if n == 0 {
			break
		}

		destinationId := int(record[0] >> (8 - significantBits))

		for _, server := range scs.Servers {
			if destinationId == server.ServerId {
				sendRecord(server.Host, server.Port, record, false)
				break
			}
		}
	}

	record := make([]byte, 100)
	for _, server := range scs.Servers {
		sendRecord(server.Host, server.Port, record, true)
	}

	numServersCompleted := 0
	myRecords := [][]byte{}
	for numServersCompleted < len(scs.Servers) {
		record := <-records
		if record[0] == byte(1) {
			numServersCompleted += 1
		} else {
			myRecords = append(myRecords, record[1:])
		}
	}

	sort.Slice(myRecords, func(i, j int) bool {
		return bytes.Compare(myRecords[i][:10], myRecords[j][:10]) == -1
	})

	outputFile, err := os.Create(outputFilePath)
	errorCheck("Error while creating output file", err, false)
	defer outputFile.Close()

	for i := 0; i < len(myRecords); i++ {
		outputFile.Write(myRecords[i])
	}
	fmt.Println(len(myRecords))
}
