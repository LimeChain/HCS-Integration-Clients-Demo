package functions

type APIResponse struct {
	Status bool   `json:"status"`
	Error  string `json:"error,omitempty"`
}

type incommingUpdates struct {
	OrderItems int `json:"newOrderItems"`
	Rfp        int `json:"newRFPs"`
	Msa        int `json:"newMSAs"`
	POs        int `json:"newPOs"`
}

type updateResponse struct {
	APIResponse
	Updates incommingUpdates `json:"updates,omitempty"`
}

type sentProposalsResponse struct {
	APIResponse
	ProposalsSent int `json:"proposalsSent"`
}

type IntegrationNodeAPIResponse struct {
	Status bool   `json:"status"`
	Error  string `json:"error,omitempty"`
}

// RFP

type RFP struct {
	RFPId      string `json:"rfpId" bson:"rfpId"`
	SupplierId string `json:"supplierId" bson:"supplierId"`
	BuyerId    string `json:"buyerId" bson:"buyerId"`
	Items      []Item `json:"items" bson:"items"`
}

type Item struct {
	OrderItemId int     `json:"orderItemId" bson:"orderItemId"`
	SKUBuyer    string  `json:"skuBuyer" bson:"skuBuyer"`
	SKUSupplier string  `json:"skuSupplier" bson:"skuSupplier"`
	Quantity    int     `json:"quantity" bson:"quantity"`
	Unit        string  `json:"unit" bson:"unit"`
	SinglePrice float32 `json:"singlePrice" bson:"singlePrice"`
	TotalValue  float32 `json:"totalValue" bson:"totalValue"`
	Currency    string  `json:"currency" bson:"currency"`
}

type allRFPsResponse struct {
	IntegrationNodeAPIResponse
	RFPs []*RFP `json:"rfps"`
}

// Proposal

type Proposal struct {
	ProposalId      string       `json:"proposalId" bson:"proposalId"`
	SupplierId      string       `json:"supplierId" bson:"supplierId"`
	BuyerId         string       `json:"buyerId" bson:"buyerId"`
	ReferencedRfpId string       `json:"referencedRfpId" bson:"referencedRfpId"`
	PriceScales     []PriceScale `json:"priceScales" bson:"priceScales"`
}

type PriceScale struct {
	Sku          ProposalSku `json:"sku" bson:"sku"`
	QuantityFrom int         `json:"quantityFrom" bson:"quantityFrom"`
	QuantityTo   int         `json:"quantityTo" bson:"quantityTo"`
	SinglePrice  float32     `json:"singlePrice" bson:"singlePrice"`
	Unit         string      `json:"unit" bson:"unit"`
	Currency     string      `json:"currency" bson:"currency"`
}

type ProposalSku struct {
	ProductName       string `json:"productName" bson:"productName"`
	BuyerProductId    string `json:"buyerProductId" bson:"buyerProductId"`
	SupplierProductId string `json:"supplierProductId" bson:"supplierProductId"`
}

type storedProposalsResponse struct {
	IntegrationNodeAPIResponse
	Proposals []*Proposal `json:"proposals"`
}

type CreateProposal struct {
	ProposalId      string       `json:"proposalId" bson:"proposalId"`
	SupplierId      string       `json:"supplierId" bson:"supplierId"`
	BuyerId         string       `json:"buyerId" bson:"buyerId"`
	ReferencedRfpId string       `json:"referencedRfpId" bson:"referencedRfpId"`
	PriceScales     []priceScale `json:"priceScales" bson:"priceScales"`
}

type priceScale struct {
	Sku          proposalSku `json:"sku" bson:"sku"`
	QuantityFrom int         `json:"quantityFrom" bson:"quantityFrom"`
	QuantityTo   int         `json:"quantityTo" bson:"quantityTo"`
	SinglePrice  float32     `json:"singlePrice" bson:"singlePrice"`
	Unit         string      `json:"unit" bson:"unit"`
	Currency     string      `json:"currency" bson:"currency"`
}

type proposalSku struct {
	ProductName       string `json:"productName" bson:"productName"`
	BuyerProductId    string `json:"buyerProductId" bson:"buyerProductId"`
	SupplierProductId string `json:"supplierProductId" bson:"supplierProductId"`
}

type createProposalsResponse struct {
	IntegrationNodeAPIResponse
	ProposalId string `json:"proposalId,omitempty"`
}

// Contracts

type UnsignedContract struct {
	ContractId           string `json:"contractId" bson:"contractId"`
	SupplierId           string `json:"supplierId" bson:"supplierId"`
	BuyerId              string `json:"buyerId" bson:"buyerId"`
	ReferencedProposalId string `json:"referencedProposalId" bson:"referencedProposalId"`
}

type Contract struct {
	UnsignedContract  `json:"unsignedContract" bson:"unsignedContract"`
	BuyerSignature    string `json:"buyerSignature" bson:"buyerSignature"`
	SupplierSignature string `json:"supplierSignature" bson:"supplierSignature"`
	DLTAnchored       bool   `json:"DLTAnchored" bson:"DLTAnchored"`
	DLTProof          string `json:"DLTProof" bson:"DLTProof"`
}

type storedContractsResponse struct {
	IntegrationNodeAPIResponse
	Contracts []*Contract `json:"contracts"`
}

// POs

type UnsignedPurchaseOrder struct {
	PurchaseOrderId      string      `json:"purchaseOrderId" bson:"purchaseOrderId"`
	SupplierId           string      `json:"supplierId" bson:"supplierId"`
	BuyerId              string      `json:"buyerId" bson:"buyerId"`
	ReferencedProposalId string      `json:"referencedProposalId" bson:"referencedProposalId"`
	ReferencedContractId string      `json:"referencedContractId" bson:"referencedContractId"`
	OrderItems           []OrderItem `json:"orderItems" bson:"orderItems"`
}

type OrderItem struct {
	OrderItemId int     `json:"orderItemId" bson:"orderItemId"`
	SKUBuyer    string  `json:"skuBuyer" bson:"skuBuyer"`
	SKUSupplier string  `json:"skuSupplier" bson:"skuSupplier"`
	Quantity    int     `json:"quantity" bson:"quantity"`
	Unit        string  `json:"unit" bson:"unit"`
	SinglePrice float32 `json:"singlePrice" bson:"singlePrice"`
	TotalValue  float32 `json:"totalValue" bson:"totalValue"`
	Currency    string  `json:"currency" bson:"currency"`
}

type PurchaseOrder struct {
	UnsignedPurchaseOrder `json:"unsignedPurchaseOrder" bson:"unsignedPurchaseOrder"`
	BuyerSignature        string `json:"buyerSignature" bson:"buyerSignature"`
	SupplierSignature     string `json:"supplierSignature" bson:"supplierSignature"`
	DLTAnchored           bool   `json:"DLTAnchored" bson:"DLTAnchored"`
	DLTProof              string `json:"DLTProof" bson:"DLTProof"`
}

type storedPOsResponse struct {
	IntegrationNodeAPIResponse
	PurchaseOrders []*PurchaseOrder `json:"contracts"`
}
