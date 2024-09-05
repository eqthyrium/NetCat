# NET-CAT

## About

This code allows you to create a group chat using the nc (netcat) utility.

## Usage

### Start:
```sh
./TCPChat 
```

If you do not specify a port, port 8989 will be used by default. If you want to use a different port, specify it as follows:

```sh
go run . $port
```

### Connection 

Run in a different terminal to connect:

```sh
nc localhost 8989
```
or a specified port:
```sh
nc localhost $port
```

Enter the correct nickname and you will be in chat.

Note: The maximum name and message length is 1024 characters and only ASCII characters 32 through 126 are accepted.

## Installation

1. Clone the repository:
```bash
git clone git@git.01.alem.school:eqthyrium/net-cat.git
```

