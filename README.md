# flow
A data channel creation and chaining CLI tool for directing data flow between processes, files, sockets and serial ports.

## Quick Start

Install the tool via command line:

> `go install github.com/Kilemonn/flow@latest`

## Usage

Below is general usage outlines for the main use cases for using this tool.

### Flow

To initiate and apply a configuration from a `yaml` file to direct data between different components on a machine.

To run:
> flow -f ./connection.yaml config-apply

This scenario requires a yaml file to be defined that holds the `node`, their `connections` and `settings`.
A simple sample is below, which will copy all data from the "input.txt" file into the "output.txt" file then from that file to `stdout`.

```yaml
connections:
  - readerid: "InputFile"
    writerid: "OutputFile"
  - readerid: "OutputFile"
    writerid: "stdout"
nodes:
    files:
      - id: "InputFile"
        path: "input.txt"
      - id: "OutputFile"
        path: "output.txt"
settings:
    timeout: 5
```

#### Connections

The `connections` defines a list of pairs of `readerid` and `writerid` pairs that references the defined `nodes` by `id` and outlines that data will be read from the `readerid` `node` and written to the `writerid` `node`.
**NOTE: `stdin` and `stdout` are always accessible as a `readerid` or `writerid` without having to define any `nodes`.**

##### Multiple Writers with the same Reader

When the same `readerid` is defined multiple times, its data will be written to **each** configured `writerid` that it is paired with. Data is essentially duplicated and written to each defined `writer`.

#### Nodes

The `nodes` contains several categories: `files`, `sockets`, `ports`, and `ipcs` that defines the underlying object and its `id`.

#### Files

A File is used to define a file in the system that you wish to read from or write to.

The `files` structure requires two properties:
- `id` used to identify the `node` itself
- `path` the path to the file
- `trunc` (optional) determines whether the file should be truncated **before writing** by default this is false and file changes are appended

```yaml
...
nodes:
  files:
    - id: "InputFile"
      path: "input.txt"
      trunc: false # optional
...
```

#### Ports

A Port is used to define a connected **Serial Port** that you wish to read from or write to.

`TODO: Confirm Ports configuration and example config`

The `ports` structure requires three properties:
- `id` used to identify the `node` itself
- `channel` "COM4" or /dev/tty1
- `mode` which contains the following properties:
    - `baudrate` The serial port bitrate
    - `databits` 
    - `parity` ?
    - `stopbits` ?
    - `initialstatusbits` ?

```yaml
...
nodes:
  ports:
    - id: "Port1"
      channel: "COM4"
      mode:
        baudrate: 9600
        databits: ""
        parity: ""
        stopbits: ""
        initialstatusbits: ""
...
```

#### Sockets

A Socket is used to define a TCP or UDP **Socket**, its address and port that it wants to send to or listen and read from.

The `sockets` structure requires four properties:
- `id` used to identify the `node` itself
- `protocol`, either `TCP` or `UDP`
- `address` this is either the address to **listen** on (if this is being used as a `reader`) or to **send** to (if this is being used as a `writer`)
- `port`, this is either the port to **listen** on (if this is being used as a `reader`) or to **send** to (if this is being used as a `writer`)

When used as a `reader` the socket will accept any incoming connection and immediately read it and forward data to the configured `writers` defined as a `connection`. All data will be read from a socket before attempting to read the next, however the order that data is read and from which socket cannot be guaranteed.

```yaml
...
nodes:
  sockets:
    - id: "TCP-Socket"
      protocol: "TCP"
      address: "127.0.0.1"
      port: 57132
...
```

#### IPC

A IPC (Inter-Process Communication Socket) is used to define an **IPC Socket** channel that you wish to to send to or listen and read from.

The `ipcs` struct requires four properties:
- `id` used to identify the `node` itself
- `channel` the socket channel to **reader** from (if this is being used as a `reader`) or to **send** to (if this is being used as a `writer`). This uses underlying **unix** sockets to communicate between processes on the device

Similarly to the `socket` configuration, when multiple incoming connections are configured, the data order cannot be guaranteed, but reading from all connections until they reach EOF is guaranteed.

```yaml
...
nodes:
  ipcs:
    - id: "IPC-1"
      channel: "channel1"
...
```

#### Settings

The Settings contains general configuration settings, if omitted the flow configuration itself will run indefinitely (Ctrl + C is your friend here).

The properties available in **Settings** are:
- `timeout` - this is a timeout in **seconds** indicating how long the flow configuration should wait before ending. The entire time must elapse **without** any new data being available in **any** reader. In other words, once all readers have no more new data to read from for **timeout** amount of seconds then the application will close all readers and writers and exit.

### Interactive Serial

In progress...
