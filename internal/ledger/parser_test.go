package ledger

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestParseSimpleLedger(t *testing.T) {
	testFile := fixturePath(t, "simple_001.qstock")

	l, err := ParseFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	if l.Format != FormatVersion {
		t.Fatalf("format = %q", l.Format)
	}

	if len(l.Accounts) != 1 {
		t.Fatalf("# accounts = %d", len(l.Accounts))
	}

	if len(l.Transactions) != 2 {
		t.Fatalf("transactions = %d", len(l.Transactions))
	}
}

func TestParsePortfolioFixture(t *testing.T) {
	testFile := fixturePath(t, "portfolio_001.qstock")

	l, err := ParseFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	if l.Format != FormatVersion {
		t.Fatalf("format = %q, want %q", l.Format, FormatVersion)
	}
	if l.Currency != "USD" {
		t.Fatalf("currency = %q, want USD", l.Currency)
	}

	if len(l.Portfolios) != 1 {
		t.Fatalf("portfolios = %d, want 1", len(l.Portfolios))
	}
	assertPortfolio(t, l.Portfolios[0], Portfolio{
		Name:        "personal",
		DisplayName: "Personal Portfolio",
		Pos:         Position{File: testFile, Line: 8},
	})

	if len(l.Accounts) != 3 {
		t.Fatalf("accounts = %d, want 3", len(l.Accounts))
	}
	assertAccount(t, l.Accounts[0], Account{
		Name:        "schwab",
		DisplayName: "Charles Schwab",
		Portfolio:   "personal",
		Currency:    "USD",
		Pos:         Position{File: testFile, Line: 10},
	})
	assertAccount(t, l.Accounts[1], Account{
		Name:        "fidelity",
		DisplayName: "Fidelity Brokerage",
		Portfolio:   "personal",
		Currency:    "USD",
		Pos:         Position{File: testFile, Line: 12},
	})
	assertAccount(t, l.Accounts[2], Account{
		Name:        "schwab-ira",
		DisplayName: "Schwab IRA",
		Portfolio:   "personal",
		Currency:    "USD",
		Pos:         Position{File: testFile, Line: 16},
	})

	if len(l.Transactions) != 6 {
		t.Fatalf("transactions = %d, want 6", len(l.Transactions))
	}

	assertTx(t, l.Transactions[0], Transaction{
		Date:     mustDate(t, "2025-01-15"),
		Type:     TxBuy,
		Symbol:   "AAPL",
		Quantity: dec(t, "10"),
		Price:    dec(t, "187.35"),
		Currency: "USD",
		Account:  "fidelity",
		Fee:      dec(t, "0.00"),
		Note:     "initial Apple position",
		Pos:      Position{File: testFile, Line: 20},
	})
	assertTx(t, l.Transactions[1], Transaction{
		Date:     mustDate(t, "2025-02-03"),
		Type:     TxBuy,
		Symbol:   "VOO",
		Quantity: dec(t, "3"),
		Price:    dec(t, "1433.7333333333333333"),
		Currency: "USD",
		Account:  "schwab-ira",
		Fee:      dec(t, "0.00"),
		Pos:      Position{File: testFile, Line: 25},
	})
	assertTx(t, l.Transactions[2], Transaction{
		Date:     mustDate(t, "2025-03-15"),
		Type:     TxDividend,
		Symbol:   "AAPL",
		Amount:   dec(t, "12.40"),
		Currency: "USD",
		Account:  "fidelity",
		Mode:     DividendCash,
		Pos:      Position{File: testFile, Line: 29},
	})
	assertTx(t, l.Transactions[3], Transaction{
		Date:           mustDate(t, "2025-03-15"),
		Type:           TxDividend,
		Symbol:         "VOO",
		Amount:         dec(t, "9.30"),
		Currency:       "USD",
		Account:        "schwab-ira",
		Mode:           DividendReinvest,
		ReinvestShares: dec(t, "0.0215"),
		ReinvestPrice:  dec(t, "432.56"),
		Pos:            Position{File: testFile, Line: 33},
	})
	assertTx(t, l.Transactions[4], Transaction{
		Date:     mustDate(t, "2025-04-01"),
		Type:     TxSell,
		Symbol:   "AAPL",
		Quantity: dec(t, "2"),
		Price:    dec(t, "205.10"),
		Currency: "USD",
		Account:  "fidelity",
		Fee:      dec(t, "1.25"),
		Pos:      Position{File: testFile, Line: 39},
	})
	assertTx(t, l.Transactions[5], Transaction{
		Date:     mustDate(t, "2025-06-10"),
		Type:     TxSplit,
		Symbol:   "NVDA",
		SplitNum: dec(t, "10"),
		SplitDen: dec(t, "1"),
		Account:  "fidelity",
		Pos:      Position{File: testFile, Line: 43},
	})
}

func TestParseIndentedMetadataWithoutRecord(t *testing.T) {
	_, err := ParseFile(fixturePath(t, "invalid_metadata_001.qstock"))
	if err == nil {
		t.Fatal("expected parse error")
	}

	var parseErr *ParseError
	if !errors.As(err, &parseErr) {
		t.Fatalf("err = %T, want *ParseError", err)
	}
	if parseErr.Line != 4 {
		t.Fatalf("line = %d, want 4", parseErr.Line)
	}
	if !strings.Contains(parseErr.Msg, "no preceding record") {
		t.Fatalf("message = %q, want no preceding record", parseErr.Msg)
	}
}

func fixturePath(t *testing.T, name string) string {
	t.Helper()
	path, err := filepath.Abs(filepath.Join("..", "..", "testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return path
}

func assertPortfolio(t *testing.T, got Portfolio, want Portfolio) {
	t.Helper()
	if got != want {
		t.Fatalf("portfolio = %#v, want %#v", got, want)
	}
}

func assertAccount(t *testing.T, got Account, want Account) {
	t.Helper()
	if got != want {
		t.Fatalf("account = %#v, want %#v", got, want)
	}
}

func assertTx(t *testing.T, got Transaction, want Transaction) {
	t.Helper()
	if !got.Date.Equal(want.Date) {
		t.Fatalf("date = %s, want %s", got.Date.Format(time.DateOnly), want.Date.Format(time.DateOnly))
	}
	if got.Type != want.Type ||
		got.Symbol != want.Symbol ||
		got.Currency != want.Currency ||
		got.Account != want.Account ||
		got.Mode != want.Mode ||
		got.Note != want.Note ||
		got.Pos != want.Pos {
		t.Fatalf("transaction = %#v, want %#v", got, want)
	}

	assertDecimal(t, "quantity", got.Quantity, want.Quantity)
	assertDecimal(t, "price", got.Price, want.Price)
	assertDecimal(t, "amount", got.Amount, want.Amount)
	assertDecimal(t, "fee", got.Fee, want.Fee)
	assertDecimal(t, "reinvest shares", got.ReinvestShares, want.ReinvestShares)
	assertDecimal(t, "reinvest price", got.ReinvestPrice, want.ReinvestPrice)
	assertDecimal(t, "split numerator", got.SplitNum, want.SplitNum)
	assertDecimal(t, "split denominator", got.SplitDen, want.SplitDen)
}

func assertDecimal(t *testing.T, name string, got decimal.Decimal, want decimal.Decimal) {
	t.Helper()
	if !got.Equal(want) {
		t.Fatalf("%s = %s, want %s", name, got, want)
	}
}

func dec(t *testing.T, s string) decimal.Decimal {
	t.Helper()
	d, err := decimal.NewFromString(s)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func mustDate(t *testing.T, s string) time.Time {
	t.Helper()
	d, err := time.Parse(time.DateOnly, s)
	if err != nil {
		t.Fatal(err)
	}
	return d
}
