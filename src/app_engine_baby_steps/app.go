package app_engine_baby_steps

import (
	"html/template"
	"net/http"
	"log"
	"fmt"
	"google.golang.org/appengine"
	gae_log "google.golang.org/appengine/log"
)

func init() {
	log.Println("Technology sucks")
	http.HandleFunc("/", root)
}

// [START func_root]
func root(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	message := fmt.Sprintf("Is dev server? %t", appengine.IsDevAppServer())
	log.Println(message)

	gae_log.Infof(c, message)

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
