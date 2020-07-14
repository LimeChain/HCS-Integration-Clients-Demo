package functions

import (
	"fmt"
	"github.com/go-chi/render"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
	"io/ioutil"
	"log"
	"net/http"
)

var proxyURL = "http://35.202.31.77:8182"

func proxySprintf(pattern string, a ...interface{}) string {
	return fmt.Sprintf(proxyURL+pattern, a...)
}

func getSheetsService() (*sheets.Service, error) {
	data, err := ioutil.ReadFile("./credentials.json")
	if err != nil {
		return nil, err
	}
	conf, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, err
	}

	client := conf.Client(context.Background())

	return sheets.New(client)
}

func UpdateIncoming(w http.ResponseWriter, r *http.Request) {

	srv, err := getSheetsService()

	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	spreadsheetID := "1Z_DonR4P5T5xyjKODgDyyF3BgQG7eObodl8d1IWxS6s"

	orderItemsIncrease, rfpsIncrease, err := updateIncommingRFPs(srv, spreadsheetID)
	if err != nil {
		log.Fatalf("response %v\n", err)
		render.JSON(w, r, updateResponse{APIResponse{false, err.Error()}, incommingUpdates{0, 0, 0, 0}})
		return
	}

	contractIncrease, err := updateIncommingContracts(srv, spreadsheetID)
	if err != nil {
		log.Fatalf("response %v\n", err)
		render.JSON(w, r, updateResponse{APIResponse{false, err.Error()}, incommingUpdates{orderItemsIncrease, rfpsIncrease, 0, 0}})
		return
	}

	poIncrease, err := updateIncommingPOs(srv, spreadsheetID)
	if err != nil {
		log.Fatalf("response %v\n", err)
		render.JSON(w, r, updateResponse{APIResponse{false, err.Error()}, incommingUpdates{orderItemsIncrease, rfpsIncrease, contractIncrease, 0}})
		return
	}

	render.JSON(w, r, updateResponse{APIResponse{true, ""}, incommingUpdates{orderItemsIncrease, rfpsIncrease, contractIncrease, poIncrease}})

}

func SendProposals(w http.ResponseWriter, r *http.Request) {

	srv, err := getSheetsService()

	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	spreadsheetID := "1Z_DonR4P5T5xyjKODgDyyF3BgQG7eObodl8d1IWxS6s"

	sentProposals, err := sendOutgoingProposals(srv, spreadsheetID)
	if err != nil {
		log.Fatalf("response %v\n", err)
		render.JSON(w, r, sentProposalsResponse{APIResponse{false, err.Error()}, 0})
		return
	}

	render.JSON(w, r, sentProposalsResponse{APIResponse{true, ""}, sentProposals})
}
