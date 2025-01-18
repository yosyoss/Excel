package model

type (
	Parameter struct {
		Start string
	}

	TaxData struct {
		SlipType      string  `json:"Slip Type"`
		VendorInvNo   string  `json:"Vendor Inv No."`
		CustomerName  string  `json:"Customer Name"`
		StmtNo        string  `json:"Stmt.No."`
		LocCur        string  `json:"Loc Cur"`
		Proxy         string  `json:"Proxy Y/N"`
		FinalAmt      float64 `json:"Final Amt"`
		FinalTaxAmt   float64 `json:"Final Tax Amt"`
		FinalTotalAmt float64 `json:"Final Total Amt"`
		AccRefNo      string  `json:"Acc Ref No."`
	}

	Response struct {
		No             string `json:"No"`
		Office         string `json:"Office"`
		SlipType       string `json:"Account Code"`
		BusinessType   string `json:"Business Type"`
		Payment        string `json:"Payment"`
		PaymentDate    string `json:"Payment date"`
		OperatingMonth string `json:"Operating Month"`
		VendorInvNo    string `json:"Vendor Inv No."`
		CustomerName   string `json:"Vendor name (User)"`
		VendorInvDate  string `json:"Vendor Inv Date"`
		StmtNo         string `json:"STP"`
		LocCur         string `json:"Cur"`
		FinalAmt       string `json:"Gross"`
		FinalTaxAmt    string `json:"Tax"`
		FinalTotalAmt  string `json:"Amoun"`
		AccRefNo       string `json:"Acc Ref No"`
		Remark         string `json:"Remark"`
		DatePayment    string `json:"Date Payment"`
	}
)
