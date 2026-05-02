# StockQuery Overview

## Tagline

> Plain-text portfolio tracker and stock research for cli nerds who invest.

---

## Summary

**StockQuery** lets users keep owned stock positions in a readable text file, update that file from the terminal, and query portfolio holdings, estimated value, quotes, structured company facts, and SEC filing details without requiring brokerage integration.

---

## Project identity

| Item | Decision |
|---|---|
| Project name | StockQuery |
| Binary name | `stkq` |
| License | MIT |
| Primary audience | CLI Nerds who invest |
| Primary interface | Terminal CLI |
| Primary storage | Plain text file |
| Default data directory | `~/.stkq/` |
| Default portfolio file | `~/.stkq/portfolio.qstock` |
| Initial market | US stocks and ETFs |
| v0.1 portfolio scope | Owned positions only |
| v0.1 research scope | Structured facts and SEC filing details only |
| News support | :x: Later version, not v0.1 |

---

## Product philosophy

StockQuery should feel like a mix of:

- `hledger` for plain-text ownership and auditability
- `git` for predictable command structure
- `jq`/Unix tools for scriptable output later
- a personal research terminal for stock ownership context

Core principles:

1. **The portfolio file is the source of truth.**
2. **Fetched market data is cacheable, replaceable, and never silently authoritative.**
3. **The tool should work offline for portfolio accounting.**
4. **The default experience should not require an API key.**
5. **Users can later configure paid/free API providers if they want better data.**

---

## v0.1 scope

### Included

- Manual plain-text portfolio file
- CLI commands that append transactions to the text file
- Multiple named portfolios
- Multiple accounts
- Owned US stocks and ETFs
- Buy transactions
- Sell transactions
- Cash dividends
- Reinvested dividends
- Manual stock splits
- Weighted average cost basis
- Holdings table
- Estimated portfolio value
- Quote lookup using a no-key provider where possible
- Local market data cache
- SEC company metadata and filing history
- Structured SEC-backed fundamentals where practical
- Pretty terminal tables
- GitHub tagged release
- Homebrew install on macOS and Linux through a tap

### Excluded

- Options
- Crypto
- Mutual funds beyond simple ETF-like symbols
- Tax reports
- FIFO/LIFO/specific lot selection
- News fetching
- Advanced performance analytics
- Rebalancing recommendations

---

## Target user

Stock Query is for an anyone who:

- Owns stocks or ETFs
- Likes plain text and terminal workflows
- Does not want to depend on a brokerage UI for basic portfolio visibility
- Wants scriptable local data
- Wants SEC-backed structured research from the terminal
- Is comfortable editing a text file but still appreciates CLI helpers

---

## User story

> As an cli nerd who invests, I want to keep my portfolio in a readable text file, add buys/sells from the terminal, and query my holdings and SEC filings without logging into a brokerage account.

Example flow:

```bash
stkq init
stkq add buy AAPL 10 @ 187.35 --account fidelity --date 2025-01-15
stkq add dividend AAPL 12.40 --account fidelity --mode cash --date 2025-03-15
stkq holdings
stkq value
stkq research AAPL
stkq filings AAPL
```

---


## Default file layout

```text
~/.stkq/
  portfolio.qstock
  config.toml
  cache/
    quotes/
    sec/
```

---

## Portfolio file format


Example:

```text
format qstock/1
currency USD    

portfolio personal "Personal Portfolio"

account fidelity "Fidelity Brokerage"
  portfolio: personal
  currency: USD

account schwab-ira "Schwab IRA"
  portfolio: personal
  currency: USD

2025-01-15 buy AAPL 10 @ 187.35 USD
  account: fidelity
  fee: 0.00 USD
  note: initial Apple position

2025-02-03 buy VOO 3 @ 430.12 USD
  account: schwab-ira
  fee: 0.00 USD

2025-03-15 dividend AAPL 12.40 USD
  account: fidelity
  mode: cash

2025-03-15 dividend VOO 9.30 USD
  account: schwab-ira
  mode: reinvest
  shares: 0.0215
  price: 432.56 USD

2025-04-01 sell AAPL 2 @ 205.10 USD
  account: fidelity
  fee: 0.00 USD

2025-06-10 split NVDA 10:1
  account: fidelity
```

### Comment syntax

Use semicolon comments:

```text
; This is a comment.
```

### Format versioning

Always include:

```text
format qstock/1
```

This allows future grammar migrations.

---