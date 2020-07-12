package functions

import (
	"encoding/json"
	"google.golang.org/api/sheets/v4"
	"log"
	"net/http"
	"strconv"
)

func appendPOOrderItem(srv *sheets.Service, spreadsheetID string, poId string, o OrderItem) (err error) {
	var appendValues sheets.ValueRange

	appendValues.Values = append(appendValues.Values, []interface{}{poId, strconv.Itoa(o.OrderItemId), o.SKUBuyer, o.SKUSupplier, strconv.Itoa(o.Quantity), o.Unit, o.Currency})

	_, err = srv.Spreadsheets.Values.Append(spreadsheetID, "PO_Items!A2:A", &appendValues).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		log.Fatalf("Unable to write orderItems data to sheet: %v", err)
		return err
	}
	return nil
}

func insertPO(srv *sheets.Service, spreadsheetID string, po PurchaseOrder) (err error) {

	var appendValues sheets.ValueRange

	appendValues.Values = append(appendValues.Values, []interface{}{po.PurchaseOrderId, po.ReferencedContractId, po.BuyerId, strconv.Itoa(len(po.OrderItems)), po.DLTAnchored, po.DLTProof})

	_, err = srv.Spreadsheets.Values.Append(spreadsheetID, "POs!A2:A", &appendValues).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		log.Fatalf("Unable to write PO data to sheet: %v", err)
		return err
	}

	for _, o := range po.OrderItems {
		err = appendPOOrderItem(srv, spreadsheetID, po.PurchaseOrderId, o)
		if err != nil {
			return err
		}
	}

	return nil
}

func updateIncommingPOs(srv *sheets.Service, spreadsheetID string) (POIncrease int, err error) {
	readRange := "POs!A2:A"
	respRead, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		return 0, err
	}

	resp, err := http.Get(proxySprintf("/po"))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var response storedPOsResponse
	err = decoder.Decode(&response)
	if err != nil {
		return 0, err
	}

	savedPOCount := len(respRead.Values)
	totalPOCount := len(response.PurchaseOrders)

	if totalPOCount > savedPOCount {
		log.Printf("Ok, lets start inserting POs")
		newPOs := response.PurchaseOrders[savedPOCount:]
		for _, po := range newPOs {
			err := insertPO(srv, spreadsheetID, *po)
			if err != nil {
				return 0, err
			}
		}
		return totalPOCount - savedPOCount, nil
	}

	return 0, nil
}
