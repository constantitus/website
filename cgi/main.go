package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const file = "anime.norg"
const css = "style.css"

func parseList() (
    titles []string,
    ratings []string,
    reviews []string,
    linecount int,
    ) {
    bytes, err := os.ReadFile(file)
    if err != nil {
        panic(err)
    }
    tmp_data := string(bytes)
    // we are only using \n, so we get rid of \r
    tmp_data = strings.ReplaceAll(tmp_data, "\r", "")
    tmp_data = strings.TrimSuffix(tmp_data, "\n")
    tmp_data = strings.TrimPrefix(tmp_data, "* ")
    tmp_data += "\n\r   " // avoid out of range errors

    data := make([]string, 1)
    data_i := 0
    for i:=0; i < len(tmp_data); {
        if tmp_data[i] == '\n' {
            skip_space:
            if tmp_data[i+1] == '-' {
                i++
            } else if tmp_data[i+1] == ' ' {
                i++
                goto skip_space;
            } else if tmp_data[i+1] == '*' {
                if tmp_data[i+2] == '*' {
                    i++
                }
                i+=3
            } else {
                goto exit_if
            }
            data = append(data, "")
            data_i++
            continue
        }
        exit_if:
        data[data_i] += string(tmp_data[i])
        i++
    }

    state := 0
    for _, v := range data {
        switch state {
        case 0:
            titles = append(titles, v)
            state++
        case 1:
            ratings = append(ratings, v)
            state++
        case 2:
            reviews = append(reviews, v)
            state = 0
        }
    }

    linecount = len(titles)
    return titles, ratings, reviews, linecount
}

func getTitleIndex() int {
    query_string := os.Getenv("QUERY_STRING")
    if query_string == "" {
        return -1
    }
    query := strings.Split(query_string, "&")
    query = strings.Split(query[0], "=")
    tmp, err := strconv.Atoi(query[1])
    if query[0] == "title_index" && err == nil {
        return tmp
    }
    return -1
}

func main() {
    // print header
    fmt.Printf("Content-Type: text/html; charset=utf-8\r\n\r\n")
    title, rating, review, linecount := parseList()

    fmt.Printf("<!DOCTYPE html><html><head><title>Anime list</title>")
    style, err := os.ReadFile(css)
    if err == nil {
        fmt.Printf("<style type=\"text/css\">%s</style>", string(style))
    }
    fmt.Printf("<body><div class=\"frame\">")

    title_index := getTitleIndex()

    for i := 0; i < linecount; i++ {
        fmt.Printf("<p><a name=\"number%d\" href=\"https://mahi.ro/animelist/anim?title_index=%d#number%d\" target=\"iframe_a\">%d. <b>%s</b></a> %s\n",
            i, i, i, i+1, title[i], rating[i])
        if title_index == i {
            fmt.Printf("<p>%s\n", review[i])
        }

    }
    fmt.Printf("</body></html>")
}
