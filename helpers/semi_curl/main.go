package main

import (
	"rental-saas/src/api_integration"
	"net/http"
	"fmt"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"strings"
	"encoding/json"
)

var _ = fmt.Printf
var _ = http.Get
var _ = rand.Int
var atoi = func(in string) (int) {out, err := strconv.Atoi(in); if err != nil {panic(err)}; return out}
var itoa = func(in int) (string) {return strconv.Itoa(in)}
var JsonToA = func(in interface{}) (string) {btz, _ := json.MarshalIndent(in, "", "  "); return string(btz)}

func main() {
	logrus.SetLevel(logrus.InfoLevel)
	p := api_integration.NewProvider()

	create := api_integration.Create2ActionRequest{
		ClientID: 11597,
		Address: api_integration.Address{
			Street: "Testowa",
			Zip:    "02-103",
			City:   "Testowo",
		},
		Frequency: 0,
		Start:     "2018-05-14 12:00:00",
		Length:    3.5,
		Zip:       "02-103",
		Chemicals: 1,
		Pets:      0,
		Eng:       1,
		Services:  []int{},
		Osource:   "A",
		Info:      "extra",
		CouponID:  123,
	}

	fmt.Printf("Trying to create: \n%s\n\n", JsonToA(create))

	createSuc, err := p.Create2(create)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	if len(createSuc.Cleaners) == 0 {
		fmt.Printf("No cleaners for address: %v\n", create)
		return
	}

	cleanersNeatly := make([]string, len(createSuc.Cleaners))
	for i, c := range createSuc.Cleaners {
		cleanersNeatly[i] = JsonToA(c)
	}
	fmt.Printf("Potential cleaners for order %s:\n%s\n\n", createSuc.OrderID, strings.Join(cleanersNeatly, "\n"))

	cleaner := createSuc.Cleaners[0]
	if len(createSuc.Cleaners) > 1 {
		fmt.Printf("Picking one at random.\n")
		randInt := rand.Intn(len(createSuc.Cleaners))
		cleaner = createSuc.Cleaners[randInt]
	}

	// saving order with many cleaners has no effect and only the first is chosen :|
	contSave := api_integration.Save2Request{
		OrderID: atoi(createSuc.OrderID),
		UserID: create.ClientID,
		Info: create.Info,
		Cleaners: []api_integration.Save2RequestCleaner{
			{ID: itoa(cleaner.ID), Stage: itoa(cleaner.Stage)},
		},
		Address: create.Address,
	}

	saveSuc, err := p.Save2(contSave)
	if err != nil {
		fmt.Printf("Saving failed: %s\n", err.Error())
		btz1, _ := json.Marshal(create)
		btz2, _ := json.Marshal(contSave)
		fmt.Printf("Create payload parsed: %s\n", string(btz1))
		fmt.Printf("Save payload parsed: %s\n", string(btz2))
		return
	}

	fmt.Printf("Server responded with: %s\n\n", saveSuc.Status)

	fmt.Printf("Activating order %d\n\n", contSave.OrderID)

	activateResp, err := p.Activate(contSave.OrderID)
	if err != nil || activateResp.Message != "" {
		fmt.Printf("Activating failed: %s\n", err.Error())
		return
	}

	fmt.Printf("Activating succeeded\n\n")

	fmt.Printf("Editing order %d (but we're really not changing anything)\n\n", contSave.OrderID)

	/*
	passing empty length breaks everything so IT HAS TO BE EXACTLY THE SAME
	creation date and time can be left blank
	address and address_id can be empty too
	*/

	edit := api_integration.EditRequest{
		OrderID: contSave.OrderID,
		UserID: contSave.UserID,
		Length: create.Length,
	}

	noopEdit, err := p.Edit(edit)
	if err != nil {
		fmt.Printf("Editing failed: %s", err.Error())
		return
	}

	btz, _ := json.MarshalIndent(noopEdit, "", "  ")
	fmt.Printf("Result: %s\n\n", string(btz))

	fmt.Printf("Deleting the event for future re-use of the date\n")
	err = p.Cancel(api_integration.CancelRequest{
		OrderID: contSave.OrderID,
		Canceler: "client",
		Cycle: 1,
		Why: "because",
	})
	if err != nil {
		fmt.Printf("Server presumably rejected our request\n%s\n", err.Error())
		return
	}

	fmt.Println("Ta-dah!")
}

