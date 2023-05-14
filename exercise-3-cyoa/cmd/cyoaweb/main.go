package main

import (
	"cyoa"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	filename := flag.String("file", "gopher.json", "the JSON file with CYOA story")
	port := flag.Int("port", 3000, "the port to start the CYOA web app on")
	flag.Parse()
	fmt.Printf("Using the story in %s\n", *filename)

	f, err := os.Open(*filename)
	if err != nil {
		// Makes more sense for printing a logical error than doing this
		// Usually a very bad idea in nested code
		panic(err)
	}

	story, err := cyoa.JsonStory(f)
	if err != nil {
		panic(err)
	}

	// %v prints verbose output (all the struct)
	// %+v prints the above, but also with field names
	// fmt.Printf("%+v", story)
	tpl := template.Must(template.New("").Parse(storyTmpl))
	h := cyoa.NewHandler(story,
		cyoa.WithTemplate(tpl),
		cyoa.WithPathFunc(pathFn),
	)
	mux := http.NewServeMux()
	mux.Handle("/story/", h)
	fmt.Printf("Starting the server on port %d\n", *port)

	// ListenAndServe keeps running until there is an error
	// in which log.Fatal should handle the loggin
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}

func pathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/story" || path == "/story/" {
		path = "/story/intro"
	}
	return path[len("/story/"):]
}

var storyTmpl = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Choose your own adventure</title>
</head>
<body>
    <h1>{{.Title}}</h1>
    {{range .Paragraphs}}
        <p>{{.}}</p>
    {{end}}
    <ul>
        {{range .Options}}
            <li><a href="/story/{{.Chapter}}">{{.Text}}</a></li>
        {{end}}
    </ul>
</body>
</html>
`
