package main

import (
    "strings"
    "html/template"
    "fmt"
    "net/http"
    "os"
    "encoding/json"
)

type LogFile struct {
    Name string `json:"name"`
    Path string `json:"path"`
    Lines []string
}

type LogList struct {
    Name string `json:"name"`
    Logs []LogFile `json:"logs"`
}

type Data struct {
    List []LogList
    Log LogFile
}

var logList []LogList
var logMap = make(map[string]LogFile)

func check(e error) {
    if e!= nil {
        panic(e)
    }
}

func reverse(numbers []string) {
    for i, j := 0, len(numbers)-1; i < j; i, j = i+1, j-1 {
        numbers[i], numbers[j] = numbers[j], numbers[i]
    }
}

func getLog(log LogFile) []string {
    dat, err := os.ReadFile(log.Path)
    check(err)

    sdat := string(dat)
    sdat = strings.TrimSuffix(sdat,"\n")
    sdat = strings.ReplaceAll(sdat,"\000","")

    lines := strings.Split(sdat,"\n")
    reverse(lines)

    return lines
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    log := r.URL.Query().Get("log")

    var lines []string
    var l LogFile

    if val, ok := logMap[log]; ok {
        lines = getLog(val)
        l = LogFile {
            Name: val.Name,
            Path: val.Path,
            Lines: lines,
        }
    }

    data := Data {
        List: logList,
        Log: l,
    }

    fmt.Printf("%s\n", logMap[log])
    t, _ := template.ParseFiles("templates/index.html")
    t.Execute(w, data)
}

func main() {
    dat, err := os.ReadFile("config.json")
    check(err)

    json.Unmarshal(dat, &logList) 

    for _, logEntry := range logList {
        for _, logFile := range logEntry.Logs {
            logMap[logFile.Name] = logFile
        }
    }

    http.HandleFunc("/", indexHandler)
    fmt.Println("View logs on port 8090.")
    http.ListenAndServe(":8090", nil)
}
