# Changelog

All notable changes to StockQuery will be documented in this file.

### Added

- Added the initial Go module for `github.com/dcxforge/stkq`.
- Added the `stkq` CLI entry point.
- Added a Cobra-based root command with global flags for portfolio file, config file, portfolio filter, account filter, offline mode, cache refresh, and color output.
- Added a `version` command that prints the binary version, commit, and build date.
- Added the initial ledger AST for `qstock_0.1.0` portfolio files, including portfolios, accounts, transactions, source positions, and supported buy, sell, dividend, and split transaction types.
- Added a `.qstock` parser with file parsing, comment stripping, quoted field tokenization, account and transaction metadata, decimal money values, `@` per-share pricing, `@@` total-price pricing, dividend reinvestment metadata, split ratios, and structured parse errors.
- Added Makefile targets for building, installing, testing, vetting, formatting, tidying dependencies, and cleaning build artifacts.
- Added dependency lock data for Cobra, pflag, and related transitive dependencies.
- Added `shopspring/decimal` as the decimal arithmetic dependency for money and quantity parsing.
- Added parser test fixtures for simple portfolios, full portfolio coverage, and invalid indented metadata.
- Added unit tests for parsing portfolios, accounts, transactions, metadata, decimal fields, source positions, symbol normalization, and parse errors.
- Added ledger validation with aggregated semantic errors for format directives, duplicate portfolios and accounts, account-to-portfolio references, transaction account references, and transaction-specific numeric constraints.

### Documentation

- Expanded the README into a product overview for StockQuery.
- Documented the project identity, target audience, default data directory, portfolio file location, v0.1 scope, and excluded features.
- Added the product philosophy and target user story.
- Added example CLI usage for initialization, transactions, holdings, valuation, research, and SEC filing lookup.
- Documented the default file layout and initial `qstock_0.1.0` portfolio file format, including accounts, buys, sells, dividends, splits, comments, and format versioning.
- Added a sample `examples/portfolio.qstock` file covering accounts, buys, sells, cash dividends, reinvested dividends, fees, notes, and splits.
- Added an ADR documenting why StockQuery uses arbitrary-precision decimal arithmetic instead of `float64` for money.
- Added inline documentation for ledger validation errors, aggregated validation reporting, validation pass ordering, and transaction-specific validation rules.

### Changed

- Updated the MIT license copyright holder.
- Updated ignored files to exclude build artifacts, `.vscode/`, and `.DS_Store`.
- Updated the documented portfolio format version from `qstock/1` to `qstock_0.1.0`.
