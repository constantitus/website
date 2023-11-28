package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const file = "anime.norg"
const css = "style.css"

type fields_t struct {
    titles []string
    ratings []string
    reviews []string
    linecount int
}

func main() {
    PrintHtml(ParseList())
}

func ParseList() (fields fields_t) {
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
                i+=2
            } else if tmp_data[i+1] == ' ' {
                i++
                goto skip_space;
            } else if tmp_data[i+1] == '*' {
                if tmp_data[i+2] == '*' {
                    i++
                }
                i+=3
            } else {
                data[data_i] += "<p>"
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
            fields.titles = append(fields.titles, v)
            state++
        case 1:
            fields.ratings = append(fields.ratings, v)
            state++
        case 2:
            fields.reviews = append(fields.reviews, v)
            state = 0
        }
    }

    fields.linecount = len(fields.titles)
    return
}

func GetTitleByIndex() int {
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

func PrintHtml(fields fields_t) {
    fmt.Print("Content-Type: text/html; charset=utf-8\r\n\r\n")

    var page strings.Builder

    page.WriteString(`<!DOCTYPE html><html>
<head>
    <title>Anime list</title>`)
    style, err := os.ReadFile(css)
    if err == nil {
        page.WriteString(`
        <style type="text/css">`)
        page.WriteString(string(style))
        page.WriteString(`</style>`)
    }
    fmt.Print(`
</head>
<body>
    <div class="frame">`)

    title_index := GetTitleByIndex()

    for i := 0; i < fields.linecount; i++ {
        page.WriteString(`<p class="title"><a name="number`)
        page.WriteString(fmt.Sprint(i))
        page.WriteString(`" href="https://mahi.ro/animelist/anim?title_index=`)
        page.WriteString(fmt.Sprint(i))
        page.WriteString(`#number`)
        page.WriteString(fmt.Sprint(i))
        page.WriteString(`" target="iframe_a">`)
        page.WriteString(fmt.Sprint(i+1))
        page.WriteString(`. <b>`)
        page.WriteString(fields.titles[i])
        page.WriteString(`</b></a> `)
        page.WriteString(fields.ratings[i])
        if fields.reviews[i] != " " { page.WriteString(" *") }
        if title_index == i {
            page.WriteString(`<p>`)
            page.WriteString(fields.reviews[i])
        }

    }
    page.WriteString(`
</body>
</html>`)
    fmt.Print(page.String())
}
