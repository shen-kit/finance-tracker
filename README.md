# Finance Tracker

> In Development!!

This is a finance tracker written in Python using the `sqlite3` package.

## Features

- income / expenditure:
	- record income/expenditure
	- use custom categories to query later, e.g. charts of income/expenditure over time for a given category
- investments:
	- record buying/selling, and the buy/sell price
	- uses `yfinance` to get the current stock price to show current value, profit/loss, etc.
- displays:
	- monthly and yearly summaries
	- summaries by category
	- investment summaries / detailed views

## Other Notes

+ as everything is open-source and all data is stored in the `database.db` file, more queries can easily be written to support whatever functionality you require.
+ the database is normalised to 3NF to minimise data anomalies
