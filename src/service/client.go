package service

import (
	"google.golang.org/api/calendar/v3"
	"io/ioutil"
	"golang.org/x/oauth2/google"
	"golang.org/x/net/context"
	"log"
	"golang.org/x/oauth2"
	"net/http"
)

const serviceClientJsonLocation = "local/service_client.json"
const readWriteCalendars = "https://www.googleapis.com/auth/calendar"

func New() *calendar.Service {

	ctx := context.Background()

	b, err := ioutil.ReadFile(serviceClientJsonLocation)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.JWTConfigFromJSON(b, readWriteCalendars)
	if err != nil {
		log.Fatalf("Unable to parse service client secret file to config: %v", err)
	}
	// TODO tokens?
	client := config.Client(ctx)

	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve calendar Client %v", err)
	}
	return srv
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}
