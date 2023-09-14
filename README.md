# Networked Sort in Go

A multi-node sorting program with sockets and Go’s net package. The main objective of this project was to get familiarized with basic socket level programming, handling servers and clients, and some popular concurrency control methods that golang provides natively.

## Description

The program will read files consisting of zero or more records. A record is a 100 byte binary key-value pair, consisting of a 10-byte key and a 90-byte value. The sort is ascending, meaning that the output should have the record with the smallest key first, then the second-smallest, etc.

We will have multiple servers running. Each server has its own single input file. At the end of the netsort process, each server will have a portion of the sorted output dataset. The distributed sorting program will concurrently run on all servers. In particular, this netsort program has the ability to sort a data set that is larger than any one server can hold. The program does the following:

* read the input file,
* appropriately partition the records,
* send relevant records to peer servers,
* receive records that belong to you from peers,
* sort the merged data (a fraction of records from its own input list augmented with a fraction of records from each of the other servers),
* write the sorted data to a single per-server output file

In the end, each server will produce a sorted output file with only records belonging to (based on the data partition algorithm below) that server.  If one were to concatenate the output files from server 0 then server 1 etc then there would be a single big sorted file.

## Assumptions

#### Number of servers: 

The number of servers running the distributed sorting program is a power of 2 (so 2, 4, 8 servers, etc). The set of servers is defined in a configuration file (see one of the testcases) specified on the command line.  There are at most 16 servers.

#### Data partition algorithm: 
Given a record with a 10 byte key and assuming 2^N servers, we would use the most significant N bits to map this record to the appropriate server. For example, in a system with 4 servers, if we encounter a record with a key starting with 1101… , it would belong to server 3.

## Network protocol

Each netsort instance will send (and receive) 100-byte records over the network socket. We send a one-byte boolean value called stream_complete (0 = false, 1 = true) indicating whether all of the data is sent, followed by the 10-byte key, then the 90-byte value.  If stream_complete is false, then the contents of the 100 bytes after it represent a valid record. If stream_complete if true, then the contents of the 100 bytes after it are not defined and no more records will be sent over the TCP socket.

![protocol_img](https://i.imgur.com/JgLjt9F.png)

## Input and output

Assume that the number of servers is N=4 (it could be any power of 2 up to 16).  Each server starts with its own input file, so server 0 has input-0.dat, server 1 has input-1.dat, etc up to input-3.dat.  After the sort is finished, then server 0 should have output-0.dat, server 1 should have output-1.dat, etc up to output-3.dat.

output-0.dat should contain a sorted set of records that all begin with bits 00. output-1.dat should contain a sorted set of records that all begin with bits 01, output-2 should have bits 10, and output-3 should have bits 11.

Thus if you concatenated output-0.dat, output-1.dat, output-2.dat, then output-3.dat together into a single big file output.dat, then output.dat would be a sorted set of key-value pairs and should contain the data that was in all the input files.

## Example

![example_img](https://i.imgur.com/lxrP7PF.png)

## Utility scripts

Read about the scripts [here](https://github.com/frank-chris/sorting-in-go#utility-scripts)

## Building

```
go build -o netsort netsort.go
```

## Running

```
netsort <serverId> <inputFilePath> <outputFilePath> <configFilePath>
```

* serverId: integer-valued id starting from 0 that specifies which server YOU are. For example, if there are 4 servers, valid values of this field would be {0, 1, 2, 3}.
* inputFilePath: input file
* outputFilePath: output file
* configFilePath: path to the config file

Edit and use `src/run-demo.sh` to run the program with the test cases by starting several servers at the same time (the script provided is for 4 servers).

Make it executable
```
chmod 777 src/run-demo.sh
```

Run it
```
./src/run-demo.sh
```

## Verifying

Concatenate all the input files into a file called INPUT

```
$ cp input-0.dat INPUT
$ cat input-1.dat >> INPUT
$ cat input-2.dat >> INPUT
$ cat input-3.dat >> INPUT
```

Concatenate all the output files into a file called OUTPUT

```
$ cp output-0.dat OUTPUT
$ cat output-1.dat >> OUTPUT
$ cat output-2.dat >> OUTPUT
$ cat output-3.dat >> OUTPUT
```

Sort INPUT and compare it with OUTPUT

```
$ utils/{architecture}/bin/showsort INPUT | sort > REF_OUTPUT
$ diff REF_OUTPUT OUTPUT
```

The diff command should not produce any output.


