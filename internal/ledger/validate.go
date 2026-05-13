package ledger

import (
	"fmt"
	"strings"
)

// ValidationError describes a single semantic problem found in a ledger file.
type ValidationError struct {
	File string
	Line int
	Msg  string
	Hint string // "did you mean?" suggestions
}

func (e *ValidationError) Error() string {
	loc := ""
	if e.File != "" {
		loc = fmt.Sprintf("%s:%d: ", e.File, e.Line)
	} else {
		loc = fmt.Sprintf("line %d: ", e.Line)
	}

	if e.Hint != "" {
		return loc + e.Msg + "(" + e.Hint + ")"
	}

	return loc + e.Msg
}

// MultiError preserves all validation failures so users can fix several
// ledger issues in one edit cycle.
type MultiError struct {
	Errs []error
}

func (m *MultiError) Error() string {
	if len(m.Errs) == 0 {
		return ""
	}

	var b strings.Builder
	for i, e := range m.Errs {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(e.Error())
	}
	return b.String()
}

// Validate checks ledger-level invariants that the parser cannot know
// while reading the file, such as duplicate names and cross-record references.
func Validate(l *Ledger) error {

	var errs []error

	// The format directive gates the rest of the file grammar. Keep this
	// explicit so newer or older ledger files fail before downstream use.
	if l.Format == "" {
		errs = append(errs, &ValidationError{
			Msg: `missing format directive`,
		})
	} else if l.Format != FormatVersion {
		errs = append(errs, &ValidationError{
			Msg: fmt.Sprintf("unsupported format %q (expected %q)", l.Format, FormatVersion),
		})
	}

	// Collect portfolios first so accounts can validate their optional
	// portfolio reference in a single pass.
	portfolios := map[string]bool{}
	for _, p := range l.Portfolios {

		if portfolios[p.Name] {
			errs = append(errs, &ValidationError{
				File: p.Pos.File,
				Line: p.Pos.Line,
				Msg:  fmt.Sprintf("duplicate portfolio %q", p.Name),
			})
		}

		portfolios[p.Name] = true
	}

	// Accounts are stored by name because transactions reference accounts
	// by their stable ledger identifier, not by display name.
	accounts := map[string]Account{}
	for _, a := range l.Accounts {

		if _, ok := accounts[a.Name]; ok {
			errs = append(errs, &ValidationError{
				File: a.Pos.File,
				Line: a.Pos.Line,
				Msg:  fmt.Sprintf("duplicate account %q", a.Name),
			})
		}

		if a.Portfolio != "" && !portfolios[a.Portfolio] {
			errs = append(errs, &ValidationError{
				File: a.Pos.File,
				Line: a.Pos.Line,
				Msg: fmt.Sprintf("account %q references unknown portfolio %q",
					a.Name, a.Portfolio),
				// TODO: Add Hint
			})
		}

		accounts[a.Name] = a

	}

	// Transaction validation is type-specific because the same struct holds
	// all transaction variants and irrelevant fields remain at zero values.
	for _, tx := range l.Transactions {
		if tx.Account == "" {
			errs = append(errs, &ValidationError{
				File: tx.Pos.File,
				Line: tx.Pos.Line,
				Msg:  "transaction missing account",
			})
		} else if _, ok := accounts[tx.Account]; !ok {
			errs = append(errs, &ValidationError{
				File: tx.Pos.File,
				Line: tx.Pos.Line,
				Msg:  fmt.Sprintf("unkown account %q", tx.Account),
				// TODO: Add Hint
			})
		}

		switch tx.Type {
		case TxBuy, TxSell:
			if !tx.Quantity.IsPositive() {
				errs = append(errs, &ValidationError{
					File: tx.Pos.File,
					Line: tx.Pos.Line,
					Msg:  "transcation quantity must be positive",
				})
			}

			// Price may be zero for no-cost acquisitions, but it should never
			// be negative.
			if tx.Price.IsNegative() {
				errs = append(errs, &ValidationError{
					File: tx.Pos.File,
					Line: tx.Pos.Line,
					Msg:  "transcation price must be 0.00 or more",
				})
			}

		case TxDividend:
			if tx.Amount.IsNegative() {
				errs = append(errs, &ValidationError{
					File: tx.Pos.File,
					Line: tx.Pos.Line,
					Msg:  "dividend amount must be 0.00 or more",
				})
			}
			switch tx.Mode {
			case "":
				errs = append(errs, &ValidationError{
					File: tx.Pos.File,
					Line: tx.Pos.Line,
					Msg:  "dividend missing mode: cash|reinvest",
				})
			case DividendReinvest:
				if !tx.ReinvestShares.IsPositive() {
					errs = append(errs, &ValidationError{
						File: tx.Pos.File,
						Line: tx.Pos.Line,
						Msg:  "dividend reinvest shares must be positive",
					})
				}
				if !tx.ReinvestPrice.IsPositive() {
					errs = append(errs, &ValidationError{
						File: tx.Pos.File,
						Line: tx.Pos.Line,
						Msg:  "dividend reinvest price must be positive",
					})
				}
			}

		case TxSplit:
			if !tx.SplitNum.IsPositive() || !tx.SplitDen.IsPositive() {
				errs = append(errs, &ValidationError{
					File: tx.Pos.File,
					Line: tx.Pos.Line,
					Msg:  "split transaction numerator and denominator must be positive",
				})
			}

		default:
			errs = append(errs, &ValidationError{
				File: tx.Pos.File,
				Line: tx.Pos.Line,
				Msg:  fmt.Sprintf("unsupported transaction type %q", tx.Type),
			})
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return &MultiError{Errs: errs}
}
