// Code generated by github.com/swaggest/json-cli v1.8.3, DO NOT EDIT.

// Package transaction contains JSON mapping structures.
package transaction

import (
	"github.com/google/uuid"
)

// Transaction structure is generated from "openapi.yaml#/components/schemas/Transaction".
type Transaction struct {
	// Format: uuid.
	// Required.
	ID uuid.UUID `json:"id" csv:"id"`
	// Format: uuid.
	// Required.
	UserID               uuid.UUID `json:"userId" csv:"userId"`
	Type                 string    `json:"type" csv:"type"`                 // Required.
	Amount               float64   `json:"amount" csv:"amount"`             // Required.
	CurrencyCode         string    `json:"currencyCode" csv:"currencyCode"` // Required.
	OriginalAmount       float64   `json:"originalAmount,omitempty" csv:"originalAmount"`
	OriginalCurrency     string    `json:"originalCurrency,omitempty" csv:"originalCurrency"`
	ExchangeRate         float64   `json:"exchangeRate,omitempty" csv:"exchangeRate"`
	MerchantCity         string    `json:"merchantCity,omitempty" csv:"merchantCity"`
	VisibleTS            int64     `json:"visibleTS" csv:"visibleTS"` // Required.
	Mcc                  int64     `json:"mcc,omitempty" csv:"mcc"`
	MccGroup             int64     `json:"mccGroup,omitempty" csv:"mccGroup"`
	MerchantName         string    `json:"merchantName,omitempty" csv:"merchantName"`
	Recurring            bool      `json:"recurring,omitempty" csv:"recurring"`
	PartnerBankName      string    `json:"partnerBankName,omitempty" csv:"partnerBankName"`
	PartnerBic           string    `json:"partnerBic,omitempty" csv:"partnerBic"`
	PartnerBcn           string    `json:"partnerBcn,omitempty" csv:"partnerBcn"`
	PartnerAccountIsSepa bool      `json:"partnerAccountIsSepa,omitempty" csv:"partnerAccountIsSepa"`
	PartnerName          string    `json:"partnerName,omitempty" csv:"partnerName"`
	// Format: uuid.
	// Required.
	AccountID           uuid.UUID `json:"accountId" csv:"accountId"`
	PartnerIban         string    `json:"partnerIban,omitempty" csv:"partnerIban"`
	PartnerAccountBan   string    `json:"partnerAccountBan,omitempty" csv:"partnerAccountBan"`
	Category            string    `json:"category" csv:"category"`       // Required.
	CardID              uuid.UUID `json:"cardId,omitempty" csv:"cardId"` // Format: uuid.
	ReferenceText       string    `json:"referenceText,omitempty" csv:"referenceText"`
	UserAccepted        int64     `json:"userAccepted,omitempty" csv:"userAccepted"`
	UserCertified       int64     `json:"userCertified" csv:"userCertified"`         // Required.
	Pending             bool      `json:"pending" csv:"pending"`                     // Required.
	Nature              string    `json:"transactionNature" csv:"transactionNature"` // Required.
	CreatedTS           int64     `json:"createdTS" csv:"createdTS"`                 // Required.
	MerchantCountry     int64     `json:"merchantCountry,omitempty" csv:"merchantCountry"`
	MerchantCountryCode int64     `json:"merchantCountryCode,omitempty" csv:"merchantCountryCode"`
	// Format: uuid.
	// Required.
	SmartLinkID    uuid.UUID `json:"smartLinkId" csv:"smartLinkId"`
	SmartContactID uuid.UUID `json:"smartContactId,omitempty" csv:"smartContactId"` // Format: uuid.
	// Format: uuid.
	// Required.
	LinkID       uuid.UUID `json:"linkId" csv:"linkId"`
	TxnCondition string    `json:"txnCondition,omitempty" csv:"txnCondition"`
	Confirmed    int64     `json:"confirmed" csv:"confirmed"` // Required.
}
