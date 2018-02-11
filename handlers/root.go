package handlers

import (
	"net/http"
	"google.golang.org/appengine"
	gae_log "google.golang.org/appengine/log"
	"html/template"
	"calendar-synch/logic"
)

func Root(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	gae_log.Debugf(ctx, "Logging test 2")

	visits, err := logic.QueryVisits(ctx, 100)
	if err != nil {
		gae_log.Debugf(ctx, "Querying vists failed")
	}

	if err := guestbookTemplate.Execute(w, visits); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var guestbookTemplate = template.Must(template.New("book").Parse(`
<html>
  <head>
    <title>Visits</title>
  </head>
  <body>
    {{range .}}
      <p>From {{.UserIP}}</p>
      {{with .Body}}
        <p><b>{{.}}</b> :</p>
      {{else}}
        <p>Empty request or error on parsing</p>
      {{end}}
      <pre>{{.Timestamp}}</pre>
    {{end}}
  </body>
</html>
`))
