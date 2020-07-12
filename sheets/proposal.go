package functions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"google.golang.org/api/sheets/v4"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func getSKU(skuID, counterpartyId string, skus [][]interface{}) (sku proposalSku, err error) {
	for _, row := range skus {
		if row[0] != skuID {
			continue
		}

		return proposalSku{fmt.Sprintf("%s", row[2]), counterpartyId, fmt.Sprintf("%s", row[0])}, nil
	}

	return sku, fmt.Errorf("Could not find sku with id %v", skuID)
}

func getPriceScale(priceScaleID string, priceScales [][]interface{}, skus [][]interface{}) (p priceScale, err error) {
	for _, row := range priceScales {
		if row[0] != priceScaleID {
			continue
		}

		sku, err := getSKU(fmt.Sprintf("%s", row[1]), fmt.Sprintf("%s", row[7]), skus)
		if err != nil {
			return p, err
		}

		quantityFrom, err := strconv.Atoi(fmt.Sprintf("%s", row[2]))
		if err != nil {
			return p, err
		}

		quantityTo, err := strconv.Atoi(fmt.Sprintf("%s", row[3]))
		if err != nil {
			return p, err
		}

		price, err := strconv.ParseFloat(fmt.Sprintf("%s", row[5]), 32)
		if err != nil {
			return p, err
		}

		p = priceScale{sku, quantityFrom, quantityTo, float32(price), fmt.Sprintf("%s", row[4]), fmt.Sprintf("%s", row[6])}

		return p, nil
	}
	return p, fmt.Errorf("Could not find price scale with id %v", priceScaleID)
}

func createProposal(proposal []interface{}, priceScales [][]interface{}, skus [][]interface{}) (b *CreateProposal, err error) {
	scaleIds := strings.Split(fmt.Sprintf("%s", proposal[3]), ",")
	proposalPriceScales := make([]priceScale, len(scaleIds))
	for i, k := range scaleIds {
		proposalPriceScales[i], err = getPriceScale(k, priceScales, skus)
		if err != nil {
			return b, err
		}
	}
	b = &CreateProposal{fmt.Sprintf("%s", proposal[0]), "ACME GSheet", fmt.Sprintf("%s", proposal[1]), fmt.Sprintf("%s", proposal[2]), proposalPriceScales}
	return b, nil
}

func sendProposal(proposal *CreateProposal) error {
	jsonValue, _ := json.Marshal(proposal)

	fmt.Println(string(jsonValue))

	resp, err := http.Post(proxySprintf("/proposal"), "application/json", bytes.NewBuffer(jsonValue))
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Posting proposal returned status code %v", resp.Status)
	}
	return nil
}

func markProposalSent(proposalID string, srv *sheets.Service, spreadsheetID string) error {
	proposalIDInt, err := strconv.Atoi(proposalID)
	if err != nil {
		return err
	}

	proposalRow := proposalIDInt + 1
	updateRange := fmt.Sprintf("Proposals!E%v:E%v", proposalRow, proposalRow)

	var updateValues sheets.ValueRange

	updateValues.Values = append(updateValues.Values, []interface{}{"Yes"})

	_, err = srv.Spreadsheets.Values.Update(spreadsheetID, updateRange, &updateValues).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		log.Fatalf("Unable to write update proposal data to sheet: %v", err)
		return err
	}
	return nil

}

func sendOutgoingProposals(srv *sheets.Service, spreadsheetID string) (sentProposals int, err error) {
	skusReadRange := "SKU!A2:Z"
	skusRead, err := srv.Spreadsheets.Values.Get(spreadsheetID, skusReadRange).Do()
	if err != nil {
		return sentProposals, err
	}

	priceScalesRange := "Proposal_Tiers!A2:Z"
	priceScalesRead, err := srv.Spreadsheets.Values.Get(spreadsheetID, priceScalesRange).Do()
	if err != nil {
		return sentProposals, err
	}

	proposalsReadRange := "Proposals!A2:Z"
	proposalsRead, err := srv.Spreadsheets.Values.Get(spreadsheetID, proposalsReadRange).Do()
	if err != nil {
		return sentProposals, err
	}

	resp, err := http.Get(proxySprintf("/proposal"))
	if err != nil {
		return sentProposals, err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var response storedProposalsResponse
	err = decoder.Decode(&response)
	if err != nil {
		return sentProposals, err
	}

	savedProposalsCount := len(proposalsRead.Values)
	sentProposalsCount := len(response.Proposals)

	if savedProposalsCount > sentProposalsCount {
		log.Printf("Ok, lets start sending")
		newProposals := proposalsRead.Values[sentProposalsCount:]
		fmt.Println(newProposals)
		for _, row := range newProposals {
			if row[4] != "No" {
				continue
			}

			proposal, err := createProposal(row, priceScalesRead.Values, skusRead.Values)
			if err != nil {
				return sentProposals, err
			}

			sendProposal(proposal)
			if err != nil {
				return sentProposals, err
			}
			err = markProposalSent(proposal.ProposalId, srv, spreadsheetID)
			if err != nil {
				return sentProposals, err
			}

			sentProposals++
		}
	}

	return sentProposals, nil
}
