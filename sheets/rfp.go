package functions

import (
	"encoding/json"
	"fmt"
	"google.golang.org/api/sheets/v4"
	"log"
	"net/http"
	"strconv"
)

func getOrderItemsCount(srv *sheets.Service, spreadsheetID string) (count int, err error) {
	readRange := "Order_Items!A2:C"
	respRead, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		return 0, err
	}

	return len(respRead.Values), nil
}

func appendOrderItem(srv *sheets.Service, spreadsheetID string, rfpId string, o Item) (err error) {
	var appendValues sheets.ValueRange

	appendValues.Values = append(appendValues.Values, []interface{}{rfpId, strconv.Itoa(o.OrderItemId), o.SKUBuyer, o.SKUSupplier, strconv.Itoa(o.Quantity), o.Unit, o.Currency})

	_, err = srv.Spreadsheets.Values.Append(spreadsheetID, "Order_Items!A2:A", &appendValues).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		log.Fatalf("Unable to write orderItems data to sheet: %v", err)
		return err
	}
	return nil
}

func appendRfp(srv *sheets.Service, spreadsheetID string, rfp RFP) (err error) {
	var appendValues sheets.ValueRange

	fmt.Println(rfp.RFPId)

	appendValues.Values = append(appendValues.Values, []interface{}{rfp.RFPId, "USMF - Contoso Entertainment System USA", strconv.Itoa(len(rfp.Items))})

	_, err = srv.Spreadsheets.Values.Append(spreadsheetID, "RFPS!A2", &appendValues).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		log.Fatalf("Unable to write RFP data to sheet: %v", err)
		return err
	}
	return nil
}

func insertRFP(srv *sheets.Service, spreadsheetID string, rfp RFP) (orderItemsIncrease int, err error) {

	for _, o := range rfp.Items {
		err = appendOrderItem(srv, spreadsheetID, rfp.RFPId, o)
		if err != nil {
			return 0, err
		}
	}

	err = appendRfp(srv, spreadsheetID, rfp)
	if err != nil {
		return 0, err
	}

	return len(rfp.Items), nil
}

func updateIncommingRFPs(srv *sheets.Service, spreadsheetID string) (orderItemsIncrease, rfpsIncrease int, err error) {
	readRange := "RFPS!A2:A"
	respRead, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		return 0, 0, err
	}

	resp, err := http.Get(proxySprintf("/rfp"))
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)

	var response allRFPsResponse
	err = decoder.Decode(&response)
	if err != nil {
		return 0, 0, err
	}

	savedRfpsCount := len(respRead.Values)
	totalRfpsCount := len(response.RFPs)

	fmt.Println(savedRfpsCount, totalRfpsCount)

	if totalRfpsCount > savedRfpsCount {
		orderItemsIncrease = 0
		newRfps := response.RFPs[savedRfpsCount:]
		for _, rfp := range newRfps {
			newOrderCount, err := insertRFP(srv, spreadsheetID, *rfp)
			if err != nil {
				return 0, 0, err
			}
			orderItemsIncrease += newOrderCount
		}
		return orderItemsIncrease, totalRfpsCount - savedRfpsCount, nil
	}

	return 0, 0, nil
}
