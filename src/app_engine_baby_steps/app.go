package app_engine_baby_steps

import (
	"html/template"
	"net/http"
	"log"
	"fmt"
	"google.golang.org/appengine"
)

func init() {
	fmt.Println("Where are my fucking logs")
	log.Println("I seriously need them now")
	http.HandleFunc("/", root)
}

// [START func_root]
func root(w http.ResponseWriter, r *http.Request) {

	message := fmt.Sprintf("Is dev server? %b", appengine.IsDevAppServer())
	log.Println(message)

	if err := guestbookTemplate.Execute(w, message); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
// [END func_root]

var guestbookTemplate = template.Must(template.New("book").Parse(`
<html>
  <head>
    <title>Baby steps...</title>
  </head>
  <body>
    <pre>{{.}}</pre>
  </body>
</html>
`))
