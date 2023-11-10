package main

import (
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

const css = "style.css"

func checkSuffix(str string) bool {
    return strings.HasSuffix(str, ".jpg")   ||
    strings.HasSuffix(str, ".jpeg")         ||
    strings.HasSuffix(str, ".png")          ||
    strings.HasSuffix(str, ".webp")         ||
    strings.HasSuffix(str, ".gif")
}

func getMode() int {
    query_string := os.Getenv("QUERY_STRING")
    if query_string == "" {
        return -1
    }
    query := strings.Split(query_string, "&")
    query = strings.Split(query[0], "=")
    tmp, err := strconv.Atoi(query[1])
    if query[0] == "mode" && err == nil {
        return tmp
    }
    return -1
}

func openDir(input string) []fs.DirEntry {
    dir, err := os.Open(input)
    files, err := dir.ReadDir(-1)
    if err != nil {
        panic(err)
    }
    dir.Close()
    return files
}

func main() {
    mode := getMode()

    var files []fs.DirEntry
    var url string

    switch mode {
    case 1:
        fmt.Printf("Content-Type: text/html; charset=utf-8\r\n\r\n")
        fmt.Printf("<!DOCTYPE html><html><head><title>Wallpapers</title>")
        style, err := os.ReadFile(css)
        if err == nil {
            fmt.Printf("<style type=\"text/css\">%s</style>", string(style))
        }

        fmt.Printf("<body><div class=\"frame\">")

        files = openDir("./pics/")
        for _, v := range files {
            name := v.Name()
            if checkSuffix(name) {
                fmt.Printf("<p><img src=\"https://dl.mahi.ro/Wallpapers/%s\" width=900 class=\"image\" \\>", name)
            }
        }
        fmt.Printf("</div>")
        fmt.Printf("</body></html>")
        return
    case 2:
        files = openDir("./bgs/")
        url = "mahi.ro/images/bgs"
    default:
        files = openDir("./pics/")
        url = "dl.mahi.ro/Wallpapers"
    }
    for i:=0; i < 20; i++ { // 20 retries should be enough
        randomfile := files[rand.Intn(len(files))].Name()
        if checkSuffix(randomfile) {
            fmt.Printf("Status:302 \nLocation: https://%s/%s\r\n\r\n", url, randomfile)
            return
        }
    }
    fmt.Printf("Status: 302 \nLocation: https://mahi.ro/50x.html")
}
