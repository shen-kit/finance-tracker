# Finance Tracker

> In Development!!

This is a TUI finance tracker written in Go, using a sqlite database.

## Features

- income / expenditure:
	- record income/expenditure
	- use custom categories to query , e.g. charts of income/expenditure over time for a given category
	- filterable and sortable table view
- investments:
	- record buying/selling, and the buy/sell price
	- uses `yfinance` to get the current stock price to show current value, profit/loss, etc.
	- filterable and sortable table view
- displays:
	- monthly and yearly summaries
	- summaries by category
	- investment summaries / detailed views

## Other Notes

+ the database is normalised to 3NF to minimise data anomalies
