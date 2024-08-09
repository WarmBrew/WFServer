# WFServer - A Simple and Efficient File Transfer Server and Client

**WFServer** is a lightweight file transfer toolset consisting of a server and client that allows for easy and reliable file transfers between systems. The tool supports features such as directory compression, resuming interrupted transfers, and is cross-platform, making it suitable for various operating systems.

## Features

### WFServer

- **File Reception**: Listens on a specified port and receives files sent from the WFclient.
- **Resume Transfers**: Supports resuming file transfers from where they left off in case of interruptions.
- **File Size and Name Display**: Displays the name and size of the received file after the transfer is complete.
- **Cross-Platform**: Can run on multiple operating systems, including Linux and Windows.

### WFClient

- **File Transfer**: Sends files from the local machine to the WFServer.
- **Directory Compression**: Compresses a specified directory into a ZIP file before transferring it.
- **Resume Transfers**: Supports resuming file transfers from the point of interruption, ensuring data integrity and efficiency.
- **Progress Display**: Shows real-time progress of file transfers, including the percentage completed.

## How to Compile

### Prerequisites

- Go 1.16 or higher installed on your system.

### Compilation Steps

#### WFServer

1. **Clone the Repository**:

   ```
   git clone https://github.com/WarmBrew/wfserver.git
   cd wfserver
   ```

2. **Compile the Server**:

   ```
   go build -o wfserver WFserver.go
   ```

#### WFClient

1. **Clone the Repository** (if not already done):

   ```
   git clone https://github.com/yourusername/wfserver.git
   cd wfserver
   ```

2. **Compile the Client**:

   ```
   go build -o wfclient WFclient.go
   ```

3. **Static Linking (Optional)**: If you need a statically linked binary (useful for portability across different Linux distributions), use the following command:

   ```
   CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o wfclient WFclient.go
   ```

4. **Cross-Compilation (Optional)**: To compile the binary for a different platform:

   ```
   GOOS=windows GOARCH=amd64 go build -o wfclient.exe WFclient.go  # For Windows
   GOOS=linux GOARCH=amd64 go build -o wfclient WFclient.go        # For Linux
   ```

## How to Use

### Running WFServer

1. Start the WFServer on the desired port:

   ```
   ./wfserver -port 59999
   ```

2. The server will now listen on the specified port, ready to receive files.

### Using WFClient

1. **Transfer a Single File**:

   ```
   ./wfclient -ip <Server_IP_Address> -port 59999 -file /path/to/file.txt
   ```

2. **Compress and Transfer a Directory**:

   ```
   ./wfclient -ip <Server_IP_Address> -port 59999 -path /path/to/directory -output compressed.zip
   ```

3. **Command-Line Arguments**:

   - `-ip`: The IP address of the remote server.
   - `-port`: The port number on the remote server.
   - `-file`: The path to the file you want to transfer.
   - `-path`: The directory path you want to compress and transfer (optional).
   - `-output`: The name of the output ZIP file after compression (optional). If not provided, the directory name will be used.

### Example Workflow

1. Start the server on your remote machine:

   ```
   ./wfserver -port 59999
   ```

2. On your local machine, transfer a file or directory:

   ```
   ./wfclient -ip 192.168.1.10 -port 59999 -file /path/to/file.txt
   or
   ./wfclient -ip 192.168.1.10 -port 59999 -path /root/xxxx/ -output xxx.zip 
   ```

## Compatibility

- **Operating Systems**: Linux, Windows
- **Dependencies**: Go 1.16 or higher

## Contributing

Contributions are welcome! Please fork this repository, make your changes, and submit a pull request. If you find any issues, feel free to open an issue on GitHub.

## License

This project is licensed under the MIT License - see the LICENSE file for details.