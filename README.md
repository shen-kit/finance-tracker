# Finance Tracker

> In Development!!

This is a TUI finance tracker written in Go, using a sqlite database.

## Features

- [ ] income / expenditure:
	- [X] record income/expenditure
	- [ ] use custom categories to query , e.g. charts of income/expenditure over time for a given category
	- [ ] filterable and sortable table view
- [X] investments:
	- [X] record buying/selling, and the buy/sell price
	- [X] gets the current stock price to show current value, profit/loss, etc.
- [X] summary displays:
	- [X] monthly summary of all records, total income/expenditure, and net value change
	- [X] yearly summary with totals by month and category
	- [X] investment summary, total quantity, average buy + current price, P/L, %P/L
- [X] responsive to terminal size

## Instructions

### Running the App

- run `$ finance-tracker <path-to-database>`
    - `$ finance-tracker test.db`
    - `$ finance-tracker ~/folder1/folder2/test.db`
- creates a database if one doesn't exist at the path, if not opens the existing one

### Controls (arrows or vim motions)

- navigation:
    - `j`/`k`/`↑`/`↓`/`g`/`G`: navigate lists/tables
    - `l`/`enter`: select a list item (make the table selectable)
    - `q`/`<C-c>`/`<C-q>`: back
    - `tab`/`<S-tab>`: next/previous item (must use this for forms when adding/editing)
    - `<C-d>`: quit
- when a table is focused:
    - `a`: add new item
    - `e`: edit selected item
    - `d`: delete selected item
    - `H`/`L` (capital): navigate back/forward a page
        - previous/next year/month for the summary pages
        - previous/next page for records/investments/categories
- shortcuts:
    - `y`: year view
    - `m`: month view
    - `r`: records
    - `c`: categories
    - `i`: investments

## Other Notes

+ the database is normalised to 3NF to minimise data anomalies
