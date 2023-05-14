package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

// Put the HTML template in the source code because it's simpler (not to worry about packaging)
// and the template is relatively simple
// Backticks enable multiline strings (raw string)
var defaultHandlerTmpl = `
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
            <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
        {{end}}
    </ul>
</body>
</html>
`

// Using global var but it is not exported so does not matter
var tpl *template.Template

// This runs when the package is initialized, and run only once
// https://tutorialedge.net/golang/the-go-init-function/
func init() {
	// template.Must is used, because if the template fails to compile, then
	// it is not ready to ship. Template needs to compile to be useful, and there is no recovery otherwise
	tpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

// Method which should be used a lot with JSON stories, so makes sense to put in here
// instead of forcing the developers to implement logic
func JsonStory(r io.Reader) (Story, error) {
	// Marshal and unmarshal accepts input as byte slices
	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		// Makes more sense for printing a logical error than doing this
		// Usually a very bad idea in nested code
		// panic(err)
		return nil, err
	}
	// Note that the returned map does not have a guaranteed order, and
	// order changes everytime we run the code again
	return story, nil
}

// Use JSON-to-go here to generate
// https://mholt.github.io/json-to-go/
// The struct returned is broken into smaller struct, instead of inline struct in Chapter
type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

// Define Story so that developers don't have to think about Story as a map
// making life easier (!?)
type Story map[string]Chapter

// Functional options design pattern
// https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
type HandlerOption func(h *handler)

type handler struct {
	s      Story
	t      *template.Template
	pathFn func(r *http.Request) string
}

// Functional options of the same type share the same prefix, in this case
// "With", for possible clearer code
func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		// Same package so can access private fields
		h.t = t
	}
}

func WithPathFunc(fn func(r *http.Request) string) HandlerOption {
	return func(h *handler) {
		h.pathFn = fn
	}
}

// Go proverb "Accept interfaces, return structs."
// Be forgiving when you take in inputs from the user, let them do whatever
// but explicitly return type so that they can do whatever they want
func NewHandler(s Story, opts ...HandlerOption) http.Handler {
	// In this case we return the Handler interface, because
	// that is eventually what we want to use this function for
	h := handler{s, tpl, defaultPathFn}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

func defaultPathFn(r *http.Request) string {
	// TrimSpace just in case
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	// Get rid of the prefix, in particular the "/"
	return path[1:]
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFn(r)
	if chapter, ok := h.s[path]; ok {
		err := h.t.Execute(w, chapter)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found", http.StatusNotFound)
}
