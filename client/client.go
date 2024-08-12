package main

import (
        "archive/zip"
        "flag"
        "fmt"
        "io"
        "net"
        "os"
        "path/filepath"
        "strconv"
        "strings"
)

const (
        MinChunkSize  = 1024 * 32
        MaxChunkSize  = 1024 * 1024 * 16
        ProgressWidth = 50
)

func main() {
        zipPath := flag.String("path", "", "Directory path to compress into a ZIP file")
        output := flag.String("output", "", "Output ZIP file name (optional)")
        filePath := flag.String("file", "", "File path to send (use after compression)")
        serverIP := flag.String("ip", "localhost", "Server IP address")
        serverPort := flag.String("port", "8080", "Server port")
        flag.Parse()

        if *zipPath != "" {
                zipFileName, err := compressDirectory(*zipPath, *output)
                if err != nil {
                        fmt.Println("Failed to compress directory:", err)
                        return
                }
                fmt.Println("Directory compressed to:", zipFileName)
                *filePath = zipFileName
        }

        if *filePath != "" {
                err := transferFile(*serverIP, *filePath, *serverPort)
                if err != nil {
                        fmt.Println("Failed to transfer file:", err)
                }
        } else if *zipPath == "" {
                fmt.Println("No file specified for transfer.")
        }
}

func compressDirectory(dirPath, outputFileName string) (string, error) {
        if outputFileName == "" {
                outputFileName = filepath.Base(dirPath) + ".zip"
        }
        zipFile, err := os.Create(outputFileName)
        if err != nil {
                return "", err
        }
        defer zipFile.Close()

        zipWriter := zip.NewWriter(zipFile)
        defer zipWriter.Close()

        err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
                if err != nil {
                        return err
                }
                relPath := strings.TrimPrefix(path, filepath.Dir(dirPath)+string(os.PathSeparator))
                if info.IsDir() {
                        return nil
                }
                file, err := os.Open(path)
                if err != nil {
                        return err
                }
                defer file.Close()

                writer, err := zipWriter.Create(relPath)
                if err != nil {
                        return err
                }
                _, err = io.Copy(writer, file)
                return err
        })

        if err != nil {
                return "", err
        }
        return outputFileName, nil
}

func transferFile(serverIP, filePath, serverPort string) error {
        conn, err := net.Dial("tcp", serverIP+":"+serverPort)
        if err != nil {
                return fmt.Errorf("error connecting to server: %w", err)
        }
        defer conn.Close()

        fileName := filepath.Base(filePath)
        fileSize, err := getFileSize(filePath)
        if err != nil {
                return fmt.Errorf("failed to get file size: %w", err)
        }

        resume := false
        if _, err := os.Stat(fileName); err == nil {
                resume = true
        }

        info := fmt.Sprintf("%s|%d|%t", fileName, fileSize, resume)
        _, err = conn.Write([]byte(info))
        if err != nil {
                return fmt.Errorf("failed to send file info: %w", err)
        }

        offset := int64(0)
        if resume {
                offsetBuf := make([]byte, 256)
                n, err := conn.Read(offsetBuf)
                if err != nil {
                        return fmt.Errorf("failed to read resume offset: %w", err)
                }
                offset, _ = strconv.ParseInt(string(offsetBuf[:n]), 10, 64)
        }

        err = sendFile(conn, filePath, fileSize, offset)
        if err != nil {
                return fmt.Errorf("failed to send file: %w", err)
        }

        fmt.Println("\nFile transfer completed successfully.")
        return nil
}

func sendFile(conn net.Conn, filePath string, fileSize, offset int64) error {
        file, err := os.Open(filePath)
        if err != nil {
                return err
        }
        defer file.Close()

        buf := make([]byte, getOptimalChunkSize(fileSize))

        file.Seek(offset, 0)

        var sentBytes = offset

        for {
                n, err := file.Read(buf)
                if err == io.EOF {
                        break
                }
                if err != nil {
                        return err
                }

                _, err = conn.Write(buf[:n])
                if err != nil {
                        return err
                }

                sentBytes += int64(n)
                printProgress(sentBytes, fileSize)
        }

        return nil
}

func getFileSize(filePath string) (int64, error) {
        fileInfo, err := os.Stat(filePath)
        if err != nil {
                return 0, err
        }
        return fileInfo.Size(), nil
}

func getOptimalChunkSize(fileSize int64) int {
        switch {
        case fileSize < 1024*1024*100:
                return MinChunkSize
        case fileSize < 1024*1024*500:
                return 1024 * 128
        case fileSize < 1024*1024*1024:
                return 1024 * 512
        default:
                return MaxChunkSize
        }
}

func printProgress(current, total int64) {
        progress := float64(current) / float64(total)
        bars := int(progress * ProgressWidth)
        fmt.Printf("\r[%-*s] %.2f%%", ProgressWidth, strings.Repeat("=", bars)+">", progress*100)
}
