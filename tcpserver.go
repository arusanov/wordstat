package main
import (
    "strconv"
    "net"
    "log"
    "bufio"
    "strings"
)

func tcpListen(port int, wordsChan chan string) (error) {
    addr := ":"+strconv.Itoa(port)
    l, err := net.Listen("tcp", addr)
    if err != nil {
        log.Fatal(err)
    }
    defer l.Close()
    for {
        conn, err := l.Accept()
        if err != nil {
            log.Fatal(err)
        }
        go handleRequest(conn, wordsChan)
    }
}

func handleRequest(conn net.Conn, wordsChan chan string) {
    defer conn.Close()
    scanner := bufio.NewScanner(bufio.NewReader(conn))
    scanner.Split(bufio.ScanWords)
    for scanner.Scan() {
        wordsChan <- strings.ToLower(strings.Trim(scanner.Text(), ",;.:!?"))
    }
    if err := scanner.Err(); err != nil {
        log.Printf("Error reading input:", err.Error())
    }
}
