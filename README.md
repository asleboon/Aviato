![UiS](http://www.ux.uis.no/~telea/uis-logo-en.png)

# Lab 7: Processing Channel Zaps

| Lab 7:	 | Processing Channel Zaps	 |
| -------------------- | ------------------------------------- |
| Subject: | DAT320 Operating Systems |
| Deadline: | Nov 19 2015 23:00 |
| Expected effort: | 15-20 hours |
| Grading: | Graded |
| Submission: | Group |


### Table of Contents

1. [Introduction](https://github.com/uis-dat320-fall17/glabs/blob/master/lab7/README.md#introduction)
2. [Collecting Channel Zaps](https://github.com/uis-dat320-fall17/glabs/blob/master/lab7/README.md#collecting-channel-zaps)
3. [Traffic Generator](https://github.com/uis-dat320-fall17/glabs/blob/master/lab7/README.md#traffic-generator)
4. [Building a Zap Event Processing Server](https://github.com/uis-dat320-fall17/glabs/blob/master/lab7/README.md#building-a-zap-event-processing-server)
5. [Publish/Subscribe RPC Client and Server](https://github.com/uis-dat320-fall17/glabs/blob/master/lab7/README.md#publishsubscribe-grpc-client-and-server)
6. [Lab Approval](https://github.com/uis-dat320-fall17/glabs/blob/master/lab7/README.md#lab-approval)


## Introduction

This is the main project in this course. It will take you through some interesting challenges,
and hopefully you will be able to use much of the stuff that you have learnt in the previous lab
exercises. **The project must be written in Go.**

**Heads up!** Before you begin, you should read through the whole document. This will hopefully 
help you plan your design, so that you can separate code into separate files and make
packages and separate structs and so forth that will help you design a good piece of software.

**Another Heads up!** You are expected to use the provided template code. Remove `TODO`
comments when you have implemented your solution. When submitting, make sure that your
program compiles and runs. If have trouble with some of tasks below: (1) Try to get help
during the lab exercise hours; (2) If you are unable to get help in time for the deadline, simply
comment the non-compiling code, and provide a commit comment describing the problem. This
comment should be separate from the first line of the final commit message, which should say:

`username lab7 submission`


## Collecting Channel Zaps

Imagine that you are working for ZapBox, an Internet and Cable service provider. ZapBox has
deployed a huge number of set-top boxes at customer homes that allows them to watch TV
over a fiber optic cable. The TV signal is distributed to customers based on a multicast stream
for each available TV channel in ZapBox’s channel portfolio. Recently, ZapBox commissioned
a software update on their deployed set-top boxes. After this software update, the set-top box
will send a UDP message to a server every time a user changes the channel on their set-top
box. In addition to channel changes, a few other items of interest may also be sent. Thus, a
message sent by a set-top box may contain information about either channel changes, volume,
mute status, or HDMI status. The content depends on the actions of the different TV viewers.
Below is shown a few samples of the message format:
```
2013/07/20, 21:56:13, 252.126.91.56, HDMI_Status: 0
2013/07/20, 21:56:55, 111.229.208.129, MAX, Viasat 4
2013/07/20, 21:57:48, 98.202.244.97, FEM, TVNORGE
2013/07/20, 21:57:44, 12.23.36.158, Canal 9, MAX
2013/07/20, 21:57:46, 81.187.186.219, TV2 Bliss, TV2 Zebra
2013/07/20, 21:57:42, 61.77.4.101, TV2 Film, TV2 Bliss
2013/07/20, 21:57:42, 203.124.29.72, Volume: 50
2013/07/20, 21:57:42, 203.124.29.72, Mute_Status: 0
```
Each line above represents an event, triggered by a single TV viewer’s action, either to
change the channel on their set-top box, or adjust the volume and so forth. These set-top
box events are sent in text format shown above. The fields are separated by comma and have
the meaning shown in the table below. Note that the message format with 5 fields represents
channel change events, while a message with only 4 fields contains a status change in the 4th
field, and no 5th field.


| Field No. | Field Name | Description |
| --------- | ---------- | ----------- |
| 1 | Date | The date that the event was sent. |
| 2 | Time | The time that the event was sent. |
| 3 | IP | The IPv4 address of the sending set-top box unit. |
| 4 | ToChan | The new channel of the set-top box. |
| 5 | FromChan | The previous channel of the set-top box. |
| 4 | StatusChange | A change in status on the set-top box. |

A StatusChange may contain one of the following entries:


| StatusChange | Value range | Description |
| ------------ | ----------- | ----------- |
| Volume: | 0-100 | The volume setting on the set-top box. |
| Mute Status: | 0/1 | The mute setting on the set-top box. |
| HDMI Status: | 0/1 | The HDMI status of the set-top box indicates whether or not a TV is connected to the set-top box and powered on. |


## Traffic Generator

For the purposes of this lab project, we have built a traffic generator to simulate the set-top
box events generated by ZapBox’s customer set-top boxes. The traffic generator resends set-top
box events loaded from a large dataset obtained from real traffic. The IP addresses have been
scrambled and do not represent a real set-top box. The traffic generator works by synchronizing
the timestamp obtained from the dataset with the local clock on the simulator machine. The
date is not synchronized.

In a real deployment, the traffic would typically be sent from set-top boxes using UDP
and received at a single UDP server, where the data can be processed. However, to make the
simulator scale to multiple receiver groups (you the students), we have instead set up the traffic
generator on a single machine multicasting each set-top box event to a single multicast address.


## Building a Zap Event Processing Server

The objective of this part is to develop a UDP multicast server that will process the events that
are sent by the set-top box clients (in our case the traffic generator). The server can run on one
of the machines in the Linux lab. Your server should be able to receive UDP packets from the
traffic generator using multicast address and port:

`224.0.1.130:10000`

Note that since the traffic generator is continuously sending out a stream of zap events, it
may be difficult to work with this part of the lab on your own machine. The multicast stream
is only available on the subnet of the Linux lab. It is therefore recommended that you work on
the lab machines, either physically or remotely using ssh.

**Tasks:**

`a.` (5 points) Build a UDP zapserver that listens to the IP multicast address and port number 
specified above. Your server *must not* echo anything back (respond) to the traffic
generator. Your server should only receive zap events in a loop. In this task, the server
only needs to print to the console whatever it receives from the traffic generator. Hint:
`net.ListenMulticastUDP`.

`b.` (10 points) Develop a data structure for storing individual zap events. The struct must
contain all the necessary fields to store channel changes (ignore storing status changes for
now). The test cases in `chzap test.go` should pass. The main task here is to implement the
constructor `NewSTBEvent()` which can be used by your server when it receives a zap event.
In addition the struct should have the following methods. See the template in `chzap.go`.

| Method | Description |
| ------ | ----------- |
| NewSTBEvent() | Returns one of three items depending on the input string, either a channel zap event, a status change event, or an error. |
| String() string | Return a string representation of your struct. |
| Duration(provided ChZap) time.Duration | Return the duration between two zap events: the receiving zap event and the provided event. |

Hints: `time.Time package`, Methods: `time.Parse()`, `strings.Split()`, `strings.TrimSpace()`, 
Layout: `const timeLayout = "2006/01/02, 15:04:05"`

`c.` (10 points) The next task is to use `zlog/simplelogger.go` (available on github) to store the
channel changes received on your zapserver.

`	1.` Use the API of the simple logger to compute the number of viewers on `NRK1` periodically,
once every second. Print the output to the console.

`	2.` Implement the same for `TV2 Norge`. They should both be printed to the console on
a separate line. Measure the time it takes to compute the Viewers() function using
`TimeElapsed()`.

**Optional task:** In the tasks above, you may observe different outputs at different times
of the day, reflecting the actual number of viewers that were actively changing channels at
your current time of day (on a previous date). Study the output at different times of the
day, perhaps coinciding with well-known TV programs on the two channels in question.
Document and explain what you observe. Is there any correlation between the data for
NRK1 and TV2?

`d.` (5 points) Take note of the measurements obtained for the `Viewers()` function over time.
What does these results show? What could be the cause of the observed problem?

`e.` (10 points) Implement a function that can compute a list of the top-10 channels. Call this
function periodically, once every second. Hint: Results returned from the `ChannelViewers()`
method defined in the `ZapLogger` interface can be sorted.

*Note that the underlying data structure used so far precludes an efficient implementation.*

`f.` (20 points) Implement a new data structure that avoids the problems that you should have
identified with the simple slice-based storage solution. Implement the data structure so that
it can support you with keeping track of the top-10 list of channels. Your implementation
must adhere to the `ZapLogger` interface. Hint: You do not need to store all the zap events
to compute the number of viewers for each channel.

##Publish/Subscribe gRPC Client and Server

`a.` (20 points) Associated with your UDP zapserver, implement a gRPC-based server (called
a publisher) that takes `Subscribe()` requests from external clients wishing to subscribe to
a stream of viewership statistics (the top-10 list). A subscriber client must include the
refresh rate in the subscribe request.

* The `RefreshRate` specific how often a subscriber wishes to be notified.
* The gRPC server is serving statistics based on zap events from the zap storage, while
continuously updating the server’s storage (the state of the server).
* The gRPC server and the zapserver part receiving zap events should be implemented
as separate goroutines, preferably in separate files.
* Assume that the refresh rate is one second or more.

Implement the publisher gRPC server and the corresponding subscriber gRPC client. The
client should display the viewership updates as they are received from the gRPC server,
leaving the refresh rate handling to the server.

This section requires three things: protobuf, go support files for protobuf, and grpc. 
The lab computers already have protobuf installed. 
If you wish to install it on your own computer, you will need 
[version 3.0.0 or higher](https://github.com/google/protobuf/releases). *Note that some Linux
distributions may only have an old version of protobuf and the protoc compiler.*
To check which protoc version is installed on your system, you can do `protoc --version`.
It should be 3.0 or higher.

The support files and grpc should already be installed after lab 3, but in case you need them.
```
go get -u github.com/golang/protobuf/proto
go get -u github.com/golang/protobuf/protoc-gen-go

go get -u google.golang.org/grpc
```

There is a sample subscribe.proto file provided for you. You will need to add variables
to the messages. Note the keyword `stream` in the file. Streaming is needed for the
pub/sub scenario.

Once you have added the variables, you must generate a go file by using
the following command:
```
protoc --go_out=plugins=grpc:. subscribe.proto
```

`	1.` How would you characterize the access pattern to the server’s state? That is, what is the
relationship between reads and writes to the server's state (storage datastructure). Suggest to
make a drawing of the system architecture, illustrating the relationship between the gRPC clients
and the server (viewed as a combination of the gRPC server and zapserver) and the STBs and the server.

`	2.` With this access pattern (workload) in mind. How would you protect the server’s state
to avoid returning a statistics computation that is incorrect or otherwise malformed?

`b.` (20 points) Now we want to analyze the duration between channel change clicks. To do
that, we need to store the previous zap event for each IP, so that you can use the `Duration()`
method that you developed earlier. You will need to extend your new data structure or
add another data structure for storing these durations. Also, extend the Subscription
struct with an additional field to select the type of statistics the subscription refers to,
either viewership or duration statistics. Whatever statistics is chosen, your publisher sends
publications to subscribers at the specified refresh rate.

`c.` (20 points) **Extra:** Profile the data structure implemented in Part 2 (a). Implement a
data structure that better supports the workload experienced by the zapserver. Profile the
new data structure and compare the results to the one implemented in (a).


## Lab Approval

To have your lab assignment approved, you must come to the lab during lab hours
and present your solution. This lets you present the thought process behind
your solution, and gives us more information for grading purposes. When you are
ready to show your solution, reach out to a member of the teaching staff. It
is expected that you can explain your code and show how it works. You may show
your solution on a lab workstation or your own computer. The results from
Autograder will also be taken into consideration when approving a lab. At least
60% of the Autograder tests should pass for the lab to be approved. A lab needs
to be approved before Autograder will provide feedback on the next lab
assignment.

Also see the [Grading and Collaboration
Policy](https://github.com/uis-dat320-fall17/course-info/blob/master/policy.md)
document for additional information.
