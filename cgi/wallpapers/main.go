package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

const css = "style.css"

func CheckSuffix(str string) bool {
    return strings.HasSuffix(str, ".jpg")   ||
    strings.HasSuffix(str, ".jpeg")         ||
    strings.HasSuffix(str, ".png")          ||
    strings.HasSuffix(str, ".webp")         ||
    strings.HasSuffix(str, ".gif")
}

func get_mode() int {
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

func main() {
    dir, err := os.Open("./pics/")
    files, err := dir.ReadDir(-1)
    if err != nil {
        panic(err)
    }
    dir.Close()

    mode := get_mode()

    if mode == 1 {
        fmt.Printf("Content-Type: text/html; charset=utf-8\r\n\r\n")
        fmt.Printf("<!DOCTYPE html><html><head><title>Wallpapers</title>")
        style, err := os.ReadFile(css)
        if err == nil {
            fmt.Printf("<style type=\"text/css\">%s</style>", string(style))
        }

        fmt.Printf("<body><div class=\"frame\">")


        for _, v := range files {
            name := v.Name()
            if CheckSuffix(name) {
                fmt.Printf("<p><img src=\"https://dl.mahi.ro/Wallpapers/%s\" width=900 class=\"image\" \\>", name)
            }
        }
        fmt.Printf("</div>")
        fmt.Printf("</body></html>")
    } else {
        // fmt.Printf("https://dl.mahi.ro/Wallpapers/%s", files[rand.Intn(len(files))].Name())
        fmt.Printf("Status:302 \nLocation: https://dl.mahi.ro/Wallpapers/%s\r\n\r\n", files[rand.Intn(len(files))].Name())
    }
}
