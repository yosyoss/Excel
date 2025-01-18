package model

type (
	Parameter struct {
		Start string
	}

	TaxData struct {
		SlipType    string  `json:"Slip Type"`
		StmtNo      string  `json:"Stmt.No."`
		FinalTaxAmt float64 `json:"Final Tax Amt"`
	}

	Response struct {
		SlipType    string `json:"Slip Type"`
		StmtNo      string `json:"Stmt.No."`
		FinalTaxAmt string `json:"Final Tax Amt"`
	}
)
