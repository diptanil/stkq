# Changelog

All notable changes to StockQuery will be documented in this file.

### Added

- Added the initial Go module for `github.com/dcxforge/stkq`.
- Added the `stkq` CLI entry point.
- Added a Cobra-based root command with global flags for portfolio file, config file, portfolio filter, account filter, offline mode, cache refresh, and color output.
- Added a `version` command that prints the binary version, commit, and build date.
- Added Makefile targets for building, installing, testing, vetting, formatting, tidying dependencies, and cleaning build artifacts.
- Added dependency lock data for Cobra, pflag, and related transitive dependencies.

### Documentation

- Expanded the README into a product overview for StockQuery.
- Documented the project identity, target audience, default data directory, portfolio file location, v0.1 scope, and excluded features.
- Added the product philosophy and target user story.
- Added example CLI usage for initialization, transactions, holdings, valuation, research, and SEC filing lookup.
- Documented the default file layout and initial `qstock/1` portfolio file format, including accounts, buys, sells, dividends, splits, comments, and format versioning.

### Changed

- Updated the MIT license copyright holder.
- Updated ignored files to exclude build artifacts, `.vscode/`, and `.DS_Store`.
