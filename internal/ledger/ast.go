package ledger

import (
	"time"

	"github.com/shopspring/decimal"
)

const FormatVersion = "qstock_0.1.0"

// Position references a source line in the ledger file.
// Every parsed entity keeps its line number so validation
// errors can point at the offending line.
type Position struct {
	File string
	Line int
}

type TxType string

// TxType enumerates supported transaction kinds for v0.1.
const (
	TxBuy      TxType = "buy"
	TxSell     TxType = "sell"
	TxDividend TxType = "dividend"
	TxSplit    TxType = "split"
)

// DividendMode describes whether a dividend was paid in
// cash or reinvested.
type DividendMode string

const (
	DividendCash     DividendMode = "cash"
	DividendReinvest DividendMode = "reinvest"
)

type Portfolio struct {
	Name        string
	DisplayName string
	Pos         Position
}

type Account struct {
	Name        string
	DisplayName string
	Portfolio   string
	Currency    string
	Pos         Position
}

// Transaction is the union of all v0.1 transaction kinds.
// Fields not relevant to a given kind are left zero.
type Transaction struct {
	Date     time.Time
	Type     TxType
	Symbol   string
	Quantity decimal.Decimal
	Price    decimal.Decimal
	Amount   decimal.Decimal
	Currency string
	Account  string
	Fee      decimal.Decimal

	Mode           DividendMode
	ReinvestShares decimal.Decimal
	ReinvestPrice  decimal.Decimal

	SplitNum decimal.Decimal
	SplitDen decimal.Decimal

	Note string
	Pos  Position
}

// Ledger is the parsed in-memory representation of a
// .qstock file.
type Ledger struct {
	Format       string
	Currency     string
	Portfolios   []Portfolio
	Accounts     []Account
	Transactions []Transaction
}
