package functions

import (
	"encoding/json"
	"google.golang.org/api/sheets/v4"
	"log"
	"net/http"
)

func insertContract(srv *sheets.Service, spreadsheetID string, contract Contract) (err error) {

	var appendValues sheets.ValueRange

	appendValues.Values = append(appendValues.Values, []interface{}{contract.ContractId, contract.BuyerId, contract.SupplierId, contract.ReferencedProposalId, contract.DLTAnchored, contract.DLTProof})

	_, err = srv.Spreadsheets.Values.Append(spreadsheetID, "Contracts!A2:A", &appendValues).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		log.Fatalf("Unable to write MSA data to sheet: %v", err)
		return err
	}
	return nil
}

func updateIncommingContracts(srv *sheets.Service, spreadsheetID string) (contractsIncrease int, err error) {
	readRange := "Contracts!A2:A"
	respRead, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		return 0, err
	}

	resp, err := http.Get(proxySprintf("/contract"))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var response storedContractsResponse
	err = decoder.Decode(&response)
	if err != nil {
		return 0, err
	}

	savedMSACount := len(respRead.Values)
	totalMSACount := len(response.Contracts)

	log.Println(respRead.Values, savedMSACount, totalMSACount)

	if totalMSACount > savedMSACount {
		log.Printf("Ok, lets start inserting MSA")
		newContracts := response.Contracts[savedMSACount:]
		for _, contract := range newContracts {
			err := insertContract(srv, spreadsheetID, *contract)
			if err != nil {
				return 0, err
			}
		}
		return totalMSACount - savedMSACount, nil
	}

	return 0, nil
}
