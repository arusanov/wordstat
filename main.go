package main

import (
    "flag"
    "sync"
    "sync/atomic"
    "net/http"
    "log"
    "strconv"
    "encoding/json"
)

var countWords = make(map[string]uint64)
var countLetters = make(map[string]uint64)
var countTotal uint64 = 0
var countRWLock sync.RWMutex

func main() {
    tcpPortPtr := flag.Int("tcp", 5555, "server tcp port")
    httpPortPtr := flag.Int("http", 8080, "server http port")
    flag.Parse()
    wordChan := make(chan string, 32)
    go statWords(wordChan)
    go tcpListen(*tcpPortPtr, wordChan)

    http.HandleFunc("/stat", handleHttpStat)
    log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*httpPortPtr), nil))
}

func makeResponse(topCount int) map[string]interface{} {
    res:=map[string]interface{} {
        "count":countTotal,
    }
    countRWLock.RLock()
    defer countRWLock.RUnlock()

    res["top_"+strconv.Itoa(topCount)+"_words"] = sortedKeys(countWords)[:min(topCount,len(countWords))]
    res["top_"+strconv.Itoa(topCount)+"_letters"] = sortedKeys(countLetters)[:min(topCount,len(countLetters))]
    return res
}

func handleHttpStat(w http.ResponseWriter, r *http.Request) {
    topNCount,_:= strconv.Atoi(r.URL.Query().Get("N"))
    if topNCount==0 {
        topNCount = 5
    }
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    json.NewEncoder(w).Encode(makeResponse(topNCount))
}

func incCount(word string) {
    countTotal = atomic.AddUint64(&countTotal, 1)

    countRWLock.Lock()
    defer countRWLock.Unlock()

    countWords[word]++
    for _, r := range word {
        countLetters[string(r)]++
    }
}

func statWords(wordsChan chan string) {
    for {
        incCount(<-wordsChan)
    }
}

func min(v1 int, v2 int) int {
    if v1<=v2 {
        return v1
    }
    return v2
}