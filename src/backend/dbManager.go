package backend

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"math"
	"math/rand"
	"strings"
	"time"
)

var db *sql.DB

var PAGE_ROWS = 15

// prepared statements
var insInvStmt *sql.Stmt
var insRecStmt *sql.Stmt
var insCatStmt *sql.Stmt

var getInvRecStmt *sql.Stmt
var getInvFilStmt *sql.Stmt
var getRecRecStmt *sql.Stmt
var getCategoriesStmt *sql.Stmt

var getIncomeSumStmt *sql.Stmt
var getExpenditureSumStmt *sql.Stmt
var getCategorySumStmt *sql.Stmt

/* Connects to the database, then creates tables and prepared statements */
func SetupDb(path string) {

	createTables := func() {
		var err error
		sql := `
    PRAGMA foreign_keys = ON; -- enforce referential integrity

    CREATE TABLE IF NOT EXISTS category (
      cat_id       INTEGER     NOT NULL PRIMARY KEY,
      cat_name     VARCHAR(20) NOT NULL,
      cat_isincome BOOL        NOT NULL,
      cat_desc     VARCHAR(40)
    );

    CREATE TABLE IF NOT EXISTS record (
      rec_id   INTEGER     NOT NULL  PRIMARY KEY,
      rec_date DATE        NOT NULL,
      rec_desc VARCHAR(50) NOT NULL,
      rec_amt  NUMBER(9)   NOT NULL,
      cat_id   INTEGER,
      CONSTRAINT category_record_fk FOREIGN KEY (cat_id) REFERENCES category (cat_id) ON UPDATE CASCADE ON DELETE SET NULL
    );

    CREATE TABLE IF NOT EXISTS investment (
      inv_id        INTEGER     NOT NULL  PRIMARY KEY,
      inv_date      DATE        NOT NULL,
      inv_code      VARCHAR(10) NOT NULL,
      inv_qty       NUMBER(7,2) NOT NULL,
      inv_unitprice NUMBER(8)   NOT NULL
    );

    CREATE TABLE IF NOT EXISTS stock (
      st_code         VARCHAR(10) NOT NULL PRIMARY KEY,
      st_unitprice    NUMBER(8,2) NOT NULL,
      st_last_updated DATE        NOT NULL
    );
    `
		if _, err = db.Exec(sql); err != nil {
			log.Printf("%q: %s\n", err, sql)
		}
	}

	createPreparedStmts := func() {
		var err error
		// insertion statements
		insInvStmt, err = db.Prepare("INSERT INTO investment (inv_date, inv_code, inv_unitprice, inv_qty) VALUES (?,?,?,?)")
		if err != nil {
			log.Println("Failed initialising insInvStmt: ", err)
		}
		insRecStmt, err = db.Prepare("INSERT INTO record (rec_date, rec_desc, rec_amt, cat_id) VALUES (?,?,?,?)")
		if err != nil {
			log.Println("Failed initialising insRecStmt: ", err)
		}
		insCatStmt, err = db.Prepare("INSERT INTO category (cat_name, cat_isincome, cat_desc) VALUES (?,?,?)")
		if err != nil {
			log.Println("Failed initialising insCatStmt: ", err)
		}

		// query statements
		getInvRecStmt, err = db.Prepare(`SELECT inv_id, inv_date, inv_code, inv_unitprice, inv_qty
                                     FROM investment
                                     ORDER BY inv_date DESC
                                     LIMIT ?, ?`)
		if err != nil {
			log.Println("Failed initialising getInvRecStmt: ", err)
		}
		getInvFilStmt, err = db.Prepare(`SELECT inv_id, inv_date, inv_code, inv_unitprice, inv_qty
                                     FROM investment
                                     WHERE inv_qty*inv_unitprice BETWEEN ? AND ?
                                       AND inv_date BETWEEN ? AND ?
                                       AND inv_code LIKE ?
                                     ORDER BY inv_date DESC`)
		if err != nil {
			log.Println("Failed initialising getInvFilStmt: ", err)
		}
		getRecRecStmt, err = db.Prepare(`SELECT rec_id, rec_date, rec_desc, rec_amt, cat_id
                                     FROM record
                                     ORDER BY rec_date DESC
                                     LIMIT ?, ?`)
		if err != nil {
			log.Println("Failed initialising getRecRecStmt: ", err)
		}
		getCategoriesStmt, err = db.Prepare(`SELECT cat_id, cat_name, cat_desc, cat_isincome FROM category`)
		if err != nil {
			log.Println("Failed initialising getRecRecStmt: ", err)
		}

		getIncomeSumStmt, err = db.Prepare(`SELECT IFNULL(SUM(rec_amt), 0)
                                        FROM record
                                        WHERE cat_id IN (SELECT cat_id FROM category WHERE cat_isincome = true)
                                          AND rec_date BETWEEN ? AND ?`)
		if err != nil {
			log.Println("Failed initialising getRecRecStmt: ", err)
		}
		getExpenditureSumStmt, err = db.Prepare(`SELECT IFNULL(SUM(rec_amt), 0)
                                             FROM record
                                             WHERE cat_id IN (SELECT cat_id FROM category WHERE cat_isincome = false)
                                               AND rec_date BETWEEN ? AND ?`)
		if err != nil {
			log.Println("Failed initialising getRecRecStmt: ", err)
		}
		getCategorySumStmt, err = db.Prepare(`SELECT IFNULL(SUM(rec_amt), 0)
                                          FROM record
                                          WHERE cat_id = ? AND rec_date BETWEEN ? AND ?`)
		if err != nil {
			log.Println("Failed initialising getRecRecStmt: ", err)
		}
	}

	// open connection to db
	var err error
	db, err = sql.Open("sqlite3", path)
	db.SetMaxOpenConns(1)
	if err != nil {
		log.Fatal(err)
	}

	createTables()
	createPreparedStmts()
}

func CreateDummyData() {
	startDate, _ := makeDate(2024, 11, 1)

	// investments
	for range 15 {
		InsertInvestment(Investment{
			Date:      startDate.AddDate(0, 0, rand.Intn(100)),
			Code:      "IVV",
			Qty:       float32(rand.Intn(100)),
			Unitprice: rand.Intn(60000),
		})
	}
	for range 15 {
		InsertInvestment(Investment{
			Date:      startDate.AddDate(0, 0, rand.Intn(100)),
			Code:      "VGS.AX",
			Qty:       float32(rand.Intn(500)),
			Unitprice: rand.Intn(20000),
		})
	}

	// categories
	categories := [...]Category{
		{Name: "Work", IsIncome: true, Desc: "income from work"},
		{Name: "Allowance", IsIncome: true, Desc: "allowance from parents"},
		{Name: "Groceries", IsIncome: false, Desc: "groceries"},
		{Name: "Entertainment", IsIncome: false, Desc: "eating out, experiences, spending for fun"},
		{Name: "Gifts", IsIncome: false, Desc: "buying presents for others"},
	}
	for _, cat := range categories {
		InsertCategory(cat)
	}

	// records
	for i := range 20 { // income
		InsertRecord(Record{
			Date:  startDate.AddDate(0, 0, rand.Intn(100)),
			Desc:  "dummy income record " + fmt.Sprint(i),
			Amt:   600 + rand.Intn(70000),
			CatId: rand.Intn(2) + 1,
		})
	}
	for i := range 80 { // expenditure
		InsertRecord(Record{
			Date:  startDate.AddDate(0, 0, rand.Intn(100)),
			Desc:  "dummy expenditure record " + fmt.Sprint(i),
			Amt:   -rand.Intn(20000),
			CatId: rand.Intn(2) + 3,
		})
	}

	fmt.Println("Inserted dummy data")
}

// Inserting Rows

func InsertRecord(rec Record) {
	_, date, desc, amt, cat_id := rec.Spread()
	if _, err := insRecStmt.Exec(date, desc, amt, cat_id); err != nil {
		log.Fatal("Failed to insert into category: ", err.Error())
	}
}

func InsertCategory(cat Category) {
	_, name, isIncome, desc := cat.Spread()
	if _, err := insCatStmt.Exec(name, isIncome, desc); err != nil {
		log.Fatal("Failed to insert into category: ", err.Error())
	}
}

func InsertInvestment(inv Investment) {
	_, date, code, unitprice, qty := inv.Spread()
	if _, err := insInvStmt.Exec(date, code, unitprice, qty); err != nil {
		log.Fatal("Failed to insert into investment: ", err.Error())
	}
}

// Reading Rows

/* Returns investments made during within a date range */
func GetInvestmentsRecent(page int) []DataRow {
	rows, err := getInvRecStmt.Query(page*PAGE_ROWS, PAGE_ROWS)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	return dbRowsToInvestments(rows)
}

/* Returns investments matching a specified filter */
func GetInvestmentsFilter(opts FilterOpts) []DataRow {
	rows, err := getInvFilStmt.Query(opts.minCost, opts.maxCost, opts.startDate, opts.endDate, "%"+opts.code+"%")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	return dbRowsToInvestments(rows)
}

/* Returns records from within a date range */
func GetRecordsRecent(page int) []DataRow {
	rows, err := getRecRecStmt.Query(page*PAGE_ROWS, PAGE_ROWS)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	return dbRowsToRecords(rows)
}

/* Returns records matching a specified filter */
func GetRecordsFilter(opts FilterOpts) []DataRow {
	cmd := `SELECT rec_id, rec_date, rec_desc, rec_amt, cat_id
          FROM record
          WHERE rec_amt BETWEEN ? AND ?
            AND rec_date >= ? AND rec_date < ?
          ORDER BY rec_date ASC`
	args := []any{opts.minCost, opts.maxCost, opts.startDate, opts.endDate}

	// filter by category if some are selected
	if len(opts.catIds) > 0 {
		cmd += "AND cat_id IN (?" + strings.Repeat(", ?", len(opts.catIds)-1) + ")"
		for _, c := range opts.catIds {
			args = append(args, c)
		}
	}

	rows, err := db.Query(cmd, args...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	return dbRowsToRecords(rows)
}

/* Returns a list of records, the total income and total expenditure */
func GetMonthInfo(date time.Time) ([]DataRow, float32, float32) {
	mStart, mEnd := getMonthStartAndEnd(date)
	recs := GetRecordsFilter(
		NewFilterOpts().
			WithStartDate(mStart).
			WithEndDate(mEnd))
	income := GetIncomeSum(mStart, mEnd)
	expenditure := GetExpenditureSum(mStart, mEnd)
	return recs, income, expenditure
}

/* Returns a slice containing all of the categories */
func GetCategories(page int) []DataRow {
	rows, err := getCategoriesStmt.Query()
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	return dbRowsToCategories(rows)
}

/* Returns the total income over a date range (inclusive) */
func GetIncomeSum(startDate, endDate time.Time) float32 {
	var sum float32
	if err := getIncomeSumStmt.QueryRow(startDate, endDate).Scan(&sum); err != nil {
		panic(err)
	}
	return sum
}

/* Returns the total expenditure over a date range (inclusive), flips sign (expenditure > 0) */
func GetExpenditureSum(startDate, endDate time.Time) float32 {
	var sum float32
	if err := getExpenditureSumStmt.QueryRow(startDate, endDate).Scan(&sum); err != nil {
		panic(err)
	}
	return -sum
}

/* Returns the total money in/out for a given category over a date range */
func GetCategorySum(catId int, startDate, endDate time.Time) (float32, error) {
	var sum float32
	if err := getCategorySumStmt.QueryRow(catId, startDate, endDate).Scan(&sum); err != nil {
		return 0, err
	}
	return sum, nil
}

/* Returns rows of [catName, [sum(Month)]] */
func GetYearSummary(year int) []DataRow {

	dataRowFromRows := func(rows *sql.Rows) []DataRow {
		var res = []DataRow{}
		var c *CategoryYear
		for rows.Next() {
			var cid, month, amt int
			if err := rows.Scan(&cid, &month, &amt); err != nil {
				panic(err)
			}
			if c == nil || cid != c.CatId {
				c = &CategoryYear{CatId: cid}
				res = append(res, c)
			}
			c.MonthSums[month-1] = amt
		}
		return res
	}

	// income categories
	sql := `SELECT cat_id, SUBSTR(rec_date, 6, 2), SUM(rec_amt)
          FROM record NATURAL JOIN category
          WHERE SUBSTR(rec_date, 1, 4) = ? AND cat_isincome
          GROUP BY cat_id, SUBSTR(rec_date, 6, 2)
          ORDER BY cat_id, SUBSTR(rec_date, 6, 2) ASC;`
	rows, err := db.Query(sql, fmt.Sprint(year))
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var res = dataRowFromRows(rows)
	res = append(res, &CategoryYear{CatId: -2}) // divider row

	// total income
	sql_totals := `SELECT -3, SUBSTR(rec_date, 6, 2), SUM(rec_amt)
         FROM record NATURAL JOIN category
         WHERE SUBSTR(rec_date, 1, 4) = ? AND cat_isincome
         GROUP BY SUBSTR(rec_date, 6, 2)
         ORDER BY SUBSTR(rec_date, 6, 2) ASC;`
	rows, err = db.Query(sql_totals, fmt.Sprint(year))
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	res = append(res, dataRowFromRows(rows)...)
	res = append(res, &CategoryYear{CatId: -2}) // divider row

	// expenditure categories
	sql = strings.Replace(sql, "cat_isincome", "NOT cat_isincome", 1)
	rows, err = db.Query(sql, fmt.Sprint(year))
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	res = append(res, dataRowFromRows(rows)...)
	res = append(res, &CategoryYear{CatId: -2}) // divider row

	// total expenditure
	sql_totals = strings.Replace(sql_totals, "-3", "-4", 1)
	sql_totals = strings.Replace(sql_totals, "cat_isincome", "NOT cat_isincome", 1)
	rows, err = db.Query(sql_totals, fmt.Sprint(year))
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	res = append(res, dataRowFromRows(rows)...)

	res = append(res, &CategoryYear{CatId: -2}) // divider row

	// net change
	sql_totals = strings.Replace(sql_totals, "-4", "0", 1)
	sql_totals = strings.Replace(sql_totals, "AND NOT cat_isincome", "", 1)
	rows, err = db.Query(sql_totals, fmt.Sprint(year))
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	res = append(res, dataRowFromRows(rows)...)

	return res
}

func GetInvestmentSummary() []DataRow {
	sql := `SELECT inv_code, SUM(inv_qty), SUM(inv_qty * inv_unitprice) / SUM(inv_qty)
          FROM investment
          GROUP BY inv_code
          ORDER BY inv_code`

	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
	}

	var invRows []InvSummaryRow
	for rows.Next() {
		var row InvSummaryRow
		var avgBuyF float64
		if err := rows.Scan(&row.code, &row.qty, &avgBuyF); err != nil {
			panic(err)
		}
		row.avgBuy = int(avgBuyF)
		invRows = append(invRows, row)
	}
	rows.Close()

	for i := range len(invRows) {

		// only update once per day maximum
		r := db.QueryRow("SELECT st_unitprice FROM stock WHERE st_code = ? AND st_last_updated >= ?", invRows[i].code, time.Now().AddDate(0, 0, -1))

		// try to use unitprice from db, if not exists get from yahoo finance and update db
		if err := r.Scan(&invRows[i].curPrice); err != nil {
			invRows[i].curPrice, err = GetCurrentStockPrice(invRows[i].code)
			if err != nil {
				panic(err)
			}

			_, err := db.Exec("INSERT OR REPLACE INTO stock (st_code, st_unitprice, st_last_updated) VALUES (?,?,?)", invRows[i].code, invRows[i].curPrice, time.Now())
			if err != nil {
				panic(err)
			}
		}
	}

	dRows := make([]DataRow, len(invRows))
	for i, v := range invRows {
		dRows[i] = v
	}

	return dRows
}

// Frontend Helper Functions

func GetInvestmentsMaxPage() int {
	var res float64
	db.QueryRow("SELECT (COUNT(*) / ?) - 1 FROM investment", float32(PAGE_ROWS)).Scan(&res)
	return int(math.Ceil(res))
}

func GetRecordsMaxPage() int {
	var res float64
	db.QueryRow("SELECT (COUNT(*) / ?) - 1 FROM record", float32(PAGE_ROWS)).Scan(&res)
	return int(math.Ceil(res))
}

func GetCategoryNameFromId(catId int) string {
	switch catId {
	case 0:
		return "[orange::b:]Net Change"
	case -1:
		return "(deleted)"
	case -2:
		return ""
	case -3:
		return "[orange]Total Income"
	case -4:
		return "[orange]Total Expenditure"
	default:
		var res string
		db.QueryRow("SELECT cat_name FROM category WHERE cat_id = ?", catId).Scan(&res)
		return res
	}
}

func GetCategoryIdFromName(catName string) int {
	var res int
	db.QueryRow("SELECT cat_id FROM category WHERE cat_name = ?", catName).Scan(&res)
	return res
}

// Updating Rows

func UpdateRecord(id int, rec Record) {
	_, date, desc, amt, catId := rec.Spread()
	_, err := db.Exec("UPDATE record SET rec_date = ?, rec_desc = ?, rec_amt = ?, cat_id = ? WHERE rec_id = ?", date, desc, amt, catId, id)
	if err != nil {
		log.Fatal("Failed to insert into investment: ", err.Error())
	}
}

func UpdateCategory(id int, cat Category) {
	_, name, isIncome, desc := cat.Spread()
	_, err := db.Exec("UPDATE category SET cat_name = ?, cat_isincome = ?, cat_desc = ? WHERE cat_id = ?", name, isIncome, desc, id)
	if err != nil {
		log.Fatal("Failed to insert into investment: ", err.Error())
	}
}

func UpdateInvestment(id int, inv Investment) {
	_, date, code, unitprice, qty := inv.Spread()
	_, err := db.Exec("UPDATE investment SET inv_date = ?, inv_code = ?, inv_qty = ?, inv_unitprice = ? WHERE inv_id = ?", date, code, qty, unitprice, id)
	if err != nil {
		log.Fatal("Failed to insert into investment: ", err.Error())
	}
}

// Deleting Rows

func DeleteRecord(id int) error {
	_, err := db.Exec("DELETE FROM record WHERE rec_id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func DeleteCategory(id int) error {
	_, err := db.Exec("DELETE FROM category WHERE cat_id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func DeleteInvestment(id int) error {
	_, err := db.Exec("DELETE FROM investment WHERE inv_id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
