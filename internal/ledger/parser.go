package ledger

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type ParseError struct {
	File string
	Line int
	Msg  string
}

func (e *ParseError) Error() string {
	if e.File == "" {
		return fmt.Sprintf("line %d: %s", e.Line, e.Msg)
	}
	return fmt.Sprintf("%s:%d %s", e.File, e.Line, e.Msg)
}

func ParseFile(path string) (*Ledger, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	return Parse(f, path)
}

type pending struct {
	kind string
	line int
	idx  int
}

func Parse(r io.Reader, file string) (*Ledger, error) {
	l := &Ledger{}
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	var cur *pending
	lineNo := 0

	for sc.Scan() {
		lineNo++
		raw := sc.Text()
		line := stripComment(raw)
		if strings.TrimSpace(line) == "" {
			continue
		}

		if isIndented(raw) {
			if cur == nil {
				return nil, &ParseError{
					File: file,
					Line: lineNo,
					Msg:  "indented metadeta with no preceding record",
				}
			}
			if err := applyMetadata(l, cur, line, lineNo, file); err != nil {
				return nil, err
			}
			continue
		}

		fields := splitFeilds(line)
		if len(fields) == 0 {
			continue
		}

		switch fields[0] {
		case "format":
			if len(fields) != 2 {
				return nil, &ParseError{
					File: file,
					Line: lineNo,
					Msg:  `expected: format ` + FormatVersion,
				}
			}
			l.Format = fields[1]
			cur = nil

		case "currency":
			if len(fields) != 2 {
				return nil, &ParseError{
					File: file,
					Line: lineNo,
					Msg:  `expected: currency USD`,
				}
			}
			l.Currency = fields[1]
			cur = nil

		case "portfolio":
			portfo, err := parsePortfolioLine(fields[1:], lineNo, file)
			if err != nil {
				return nil, err
			}
			portfo.Pos = Position{
				File: file,
				Line: lineNo,
			}
			l.Portfolios = append(l.Portfolios, portfo)
			cur = &pending{
				kind: "portfolio",
				line: lineNo,
				idx:  len(l.Portfolios) - 1,
			}

		case "account":
			acct, err := parseAccountLine(fields[1:], lineNo, file)
			if err != nil {
				return nil, err
			}
			acct.Pos = Position{
				File: file,
				Line: lineNo,
			}
			l.Accounts = append(l.Accounts, acct)
			cur = &pending{
				kind: "account",
				line: lineNo,
				idx:  len(l.Accounts) - 1,
			}

		default:
			tx, err := parseTransactionHeader(fields, lineNo, file)
			if err != nil {
				return nil, err
			}
			l.Transactions = append(l.Transactions, tx)
			cur = &pending{
				kind: "tx",
				line: lineNo,
				idx:  len(l.Transactions) - 1,
			}
		}
	}

	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("read %s: %w", file, err)
	}

	return l, nil
}

// stripComment removes a trailing `; ...` comment,
// it preserves any `;` that might appear inside a quoted string `"`
func stripComment(s string) string {
	inQuote := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '"' {
			inQuote = !inQuote
		}
		if c == ';' && !inQuote {
			return s[:i]
		}
	}
	return s
}

// isIndented reports whether the line begins with a space or tab.
func isIndented(s string) bool {
	return len(s) > 0 && (s[0] == ' ' || s[0] == '\t')
}

// splitFields tokenizes a line, respecting double quotes around
// values that contain spaces.
// Quotes are stripped from the output.
func splitFeilds(s string) []string {
	var out []string
	var b strings.Builder

	inQuote := false
	flush := func() {
		if b.Len() > 0 {
			out = append(out, b.String())
			b.Reset()
		}
	}

	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '"':
			inQuote = !inQuote
		case (c == ' ' || c == '\t') && !inQuote:
			flush()
		default:
			b.WriteByte(c)
		}
	}
	flush()
	return out
}

func parsePortfolioLine(fields []string, lineNo int, file string) (Portfolio, error) {
	if len(fields) < 1 {
		return Portfolio{}, &ParseError{
			File: file,
			Line: lineNo,
			Msg:  "portfolio: requires a name",
		}
	}
	p := Portfolio{
		Name:        fields[0],
		DisplayName: fields[0],
	}
	if len(fields) >= 2 {
		p.DisplayName = strings.Join(fields[1:], " ")
	}
	return p, nil
}

func parseAccountLine(fields []string, lineNo int, file string) (Account, error) {
	if len(fields) < 1 {
		return Account{}, &ParseError{
			File: file,
			Line: lineNo,
			Msg:  "account: requires a name",
		}
	}
	a := Account{
		Name:        fields[0],
		DisplayName: fields[0],
	}
	if len(fields) > 4 {
		return Account{}, &ParseError{
			File: file,
			Line: lineNo,
			Msg:  "account: too many fields",
		}
	}
	if len(fields) >= 2 {
		a.DisplayName = fields[1]
	}
	if len(fields) >= 3 {
		a.Portfolio = fields[2]
	}
	if len(fields) >= 4 {
		a.Currency = fields[3]
	}
	return a, nil
}

func parseTransactionHeader(fields []string, lineNo int, file string) (Transaction, error) {
	tx := Transaction{
		Pos: Position{
			File: file,
			Line: lineNo,
		},
	}
	if len(fields) < 3 {
		return tx, &ParseError{
			File: file,
			Line: lineNo,
			Msg:  "transaction: requires date type, and symbol",
		}
	}

	date, err := time.Parse("2006-01-02", fields[0])
	if err != nil {
		return tx, &ParseError{
			File: file,
			Line: lineNo,
			Msg:  fmt.Sprintf("transaction: invalid date %q: expected YYYY-MM-DD", fields[0]),
		}
	}
	tx.Date = date

	switch fields[1] {
	case "buy":
		tx.Type = TxBuy
	case "sell":
		tx.Type = TxSell
	case "dividend":
		tx.Type = TxDividend
	case "split":
		tx.Type = TxSplit
	default:
		return tx, &ParseError{
			File: file,
			Line: lineNo,
			Msg:  fmt.Sprintf("transaction: unsupported transaction type %q", fields[1])}
	}

	tx.Symbol = strings.ToUpper(fields[2])

	switch tx.Type {
	// 2024-01-15 buy AAPL 10 @ 187.35 USD
	// OR:
	// 2024-01-15 buy AAPL 10 @@ 1873.5 USD
	// Currency is optional
	case TxBuy, TxSell:
		if len(fields) < 6 || len(fields) > 7 {
			return tx, &ParseError{
				File: file,
				Line: lineNo,
				Msg: `transaction (buy|sell): 
				expected - DATE buy|sell SYMBOL QTY @|@@ PRICE [CURRENCY]
				[CURRENCY] - Optinal`,
			}
		}

		if !(fields[4] == "@" || fields[4] == "@@") {
			return tx, &ParseError{
				File: file,
				Line: lineNo,
				Msg: `transaction (buy|sell): 
				expected - DATE buy|sell SYMBOL QTY @|@@ PRICE [CURRENCY]
				[CURRENCY] - Optinal`,
			}
		}

		q, err := decimal.NewFromString(fields[3])
		if err != nil {
			return tx, &ParseError{
				File: file,
				Line: lineNo,
				Msg:  "transaction: invalid buy|sell quantity",
			}
		}

		p, err := decimal.NewFromString(fields[5])
		if err != nil {
			return tx, &ParseError{
				File: file,
				Line: lineNo,
				Msg:  "transaction: invalid buy|sell price",
			}
		}

		if fields[4] == "@@" && q.IsZero() {
			return tx, &ParseError{
				File: file,
				Line: lineNo,
				Msg:  "transaction: quantity cannot be 0",
			}
		}

		tx.Quantity = q
		if fields[4] == "@" {
			tx.Price = p
		} else {
			tx.Price = p.Div(q)
		}

		if len(fields) == 7 {
			tx.Currency = fields[6]
		}

	// 2024-08-15 dividend AAPL 6.00 USD
	// Currency is optional
	case TxDividend:
		if len(fields) < 4 || len(fields) > 5 {
			return tx, &ParseError{
				File: file,
				Line: lineNo,
				Msg: `transaction (dividend): 
				expected - DATE dividend SYMBOL AMOUNT [CURRENCY]
				[CURRENCY] - Optinal`,
			}
		}

		amount, err := decimal.NewFromString(fields[3])
		if err != nil {
			return tx, &ParseError{
				File: file,
				Line: lineNo,
				Msg:  "transaction: dividend invalid amount"}
		}
		tx.Amount = amount

		if len(fields) == 5 {
			tx.Currency = fields[4]
		}

	// 2024-08-29 split AAPL 4:1
	case TxSplit:
		if len(fields) != 4 {
			return tx, &ParseError{
				File: file,
				Line: lineNo,
				Msg: `transaction (split): 
				expected - DATE split SYMBOL N:M`,
			}
		}
		n, m, err := parseSplitRatio(fields[3])

		if err != nil {
			return tx, &ParseError{
				File: file,
				Line: lineNo,
				Msg:  err.Error()}
		}

		tx.SplitNum = n
		tx.SplitDen = m
	}

	return tx, nil
}

func parseSplitRatio(s string) (decimal.Decimal, decimal.Decimal, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return decimal.Zero, decimal.Zero,
			fmt.Errorf("transaction (split): invalid split ratio %q: expected N:M", s)
	}

	n, err := decimal.NewFromString(parts[0])
	if err != nil {
		return decimal.Zero, decimal.Zero,
			fmt.Errorf("transaction (split): invalid split numerator %q", parts[0])
	}

	m, err := decimal.NewFromString(parts[1])
	if err != nil {
		return decimal.Zero, decimal.Zero,
			fmt.Errorf("transaction (split): invalid split denominator %q", parts[1])
	}
	if m.IsZero() {
		return decimal.Zero, decimal.Zero,
			fmt.Errorf("transaction (split): denominator cannot be 0")
	}
	return n, m, nil
}

func applyMetadata(l *Ledger, cur *pending, line string, lineNo int, file string) error {
	key, value, err := splitKV(line, lineNo, file)
	if err != nil {
		return err
	}

	switch cur.kind {
	case "portfolio":
		_, _ = key, value

	case "account":
		acc := &l.Accounts[cur.idx]
		switch key {
		case "portfolio":
			acc.Portfolio = value
		case "currency":
			acc.Currency = value
		}
	case "tx":
		tx := &l.Transactions[cur.idx]

		switch key {
		case "account":
			tx.Account = value

		case "fee":
			// "0.00 USD" -> take the number part
			parts := strings.Fields(value)

			if len(parts) > 2 {
				return &ParseError{
					File: file,
					Line: lineNo,
					Msg: `transaction (fee):
					expected - AMOUNT [CURRENCY]
					[CURRENCY] - Optinal`,
				}
			}

			if len(parts) == 0 {
				return &ParseError{
					File: file,
					Line: lineNo,
					Msg:  "transaction: empty fee",
				}
			}

			d, err := decimal.NewFromString(parts[0])
			if err != nil {
				return &ParseError{
					File: file,
					Line: lineNo,
					Msg:  "transaction: invalid fee",
				}
			}
			tx.Fee = d
			if len(parts) == 2 {
				tx.Currency = parts[1]
			}

		case "mode":
			switch value {
			case "cash":
				tx.Mode = DividendCash
			case "reinvest":
				tx.Mode = DividendReinvest
			default:
				return &ParseError{
					File: file,
					Line: lineNo,
					Msg:  fmt.Sprintf(`mode must be "cash" or "reinvest", got %q`, value),
				}
			}

		case "shares":
			d, err := decimal.NewFromString(value)
			if err != nil {
				return &ParseError{
					File: file,
					Line: lineNo,
					Msg:  "invalid shares",
				}
			}

			tx.ReinvestShares = d

		case "price":
			parts := strings.Fields(value)

			if len(parts) > 2 {
				return &ParseError{
					File: file,
					Line: lineNo,
					Msg: `transaction (price):
					expected - AMOUNT [CURRENCY]
					[CURRENCY] - Optinal`,
				}
			}

			if len(parts) == 0 {
				return &ParseError{
					File: file,
					Line: lineNo,
					Msg:  "transaction: empty fee",
				}
			}

			d, err := decimal.NewFromString(parts[0])
			if err != nil {
				return &ParseError{
					File: file,
					Line: lineNo,
					Msg:  "transaction: invalid price",
				}
			}

			tx.ReinvestPrice = d

		case "note":
			tx.Note = value

		}
	}

	return nil
}

func splitKV(line string, lineNo int, file string) (string, string, error) {
	line = strings.TrimSpace(line)
	idx := strings.Index(line, ":")
	if idx < 0 {
		return "", "", &ParseError{
			File: file,
			Line: lineNo,
			Msg:  `expected: "key: value"`,
		}
	}

	return strings.TrimSpace(line[:idx]),
		strings.TrimSpace(line[idx+1:]),
		nil
}
