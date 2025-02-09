package backend

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

const PAGE_ROWS = 15

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
      rec_amt  NUMBER(7,2) NOT NULL,
      cat_id   INTEGER     NOT NULL,
      CONSTRAINT category_record_fk FOREIGN KEY (cat_id) REFERENCES category (cat_id) ON UPDATE CASCADE ON DELETE SET NULL
    );
    CREATE TABLE IF NOT EXISTS investment (
      inv_id        INTEGER     NOT NULL  PRIMARY KEY,
      inv_date      DATE        NOT NULL,
      inv_code      VARCHAR(10) NOT NULL,
      inv_qty       NUMBER(7,2) NOT NULL,
      inv_unitprice NUMBER(8,2) NOT NULL
    );
    `
		if _, err = db.Exec(sql); err != nil {
			log.Printf("%q: %s\n", err, sql)
		}
	}

	createPreparedStmts := func() {
		var err error
		// insertion statements
		insInvStmt, err = db.Prepare("INSERT INTO investment (inv_date, inv_code, inv_qty, inv_unitprice) VALUES (?,?,?,?)")
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
		getInvRecStmt, err = db.Prepare(`SELECT inv_id, inv_date, inv_code, inv_qty, inv_unitprice
                                     FROM investment
                                     ORDER BY inv_date DESC
                                     LIMIT ?, ?`)
		if err != nil {
			log.Println("Failed initialising getInvRecStmt: ", err)
		}
		getInvFilStmt, err = db.Prepare(`SELECT inv_id, inv_date, inv_code, inv_qty, inv_unitprice
                                     FROM investment
                                     WHERE inv_qty*inv_unitprice BETWEEN ? AND ?
                                       AND inv_date BETWEEN ? AND ?
                                       AND inv_code LIKE ?`)
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

		getIncomeSumStmt, err = db.Prepare(`SELECT SUM(rec_amt)
                                        FROM record
                                        WHERE cat_id IN (SELECT cat_id FROM category WHERE cat_isincome = true)
                                          AND rec_date BETWEEN ? AND ?`)
		if err != nil {
			log.Println("Failed initialising getRecRecStmt: ", err)
		}
		getExpenditureSumStmt, err = db.Prepare(`SELECT SUM(rec_amt)
                                             FROM record
                                             WHERE cat_id IN (SELECT cat_id FROM category WHERE cat_isincome = false)
                                               AND rec_date BETWEEN ? AND ?`)
		if err != nil {
			log.Println("Failed initialising getRecRecStmt: ", err)
		}
		getCategorySumStmt, err = db.Prepare(`SELECT SUM(rec_amt)
                                          FROM record
                                          WHERE cat_id = ? AND rec_date BETWEEN ? AND ?`)
		if err != nil {
			log.Println("Failed initialising getRecRecStmt: ", err)
		}
	}

	// open connection to db
	var err error
	db, err = sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}

	createTables()
	createPreparedStmts()
}

func CreateDummyData() {
	investments := [...]Investment{
		{Date: time.Now().AddDate(0, -1, 0), Code: "IVV", Qty: 10, Unitprice: 600},
		{Date: time.Now().AddDate(0, -1, 0), Code: "VGS.AX", Qty: 5, Unitprice: 600},
		{Date: time.Now().AddDate(0, -1, 0), Code: "IVV", Qty: 10, Unitprice: 600},
	}
	for range 10 {
		for _, inv := range investments {
			InsertInvestment(inv)
		}
	}

	// categories := [...]Category{
	// 	{Name: "Work", IsIncome: true, Desc: "income from work"},
	// 	{Name: "Groceries", IsIncome: false, Desc: "grocery spending"},
	// }
	// for _, cat := range categories {
	// 	InsertCategory(cat)
	// }
	//
	// records := [...]Record{
	// 	{Date: time.Now(), Desc: "new record desc 1", Amt: 100, CatId: 1},
	// 	{Date: time.Now(), Desc: "new record desc 2", Amt: -200, CatId: 2},
	// 	{Date: time.Now(), Desc: "new record desc 3", Amt: 300, CatId: 1},
	// }
	// for _, rec := range records {
	// 	InsertRecord(rec)
	// }
	//
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
	_, date, code, qty, unitprice := inv.Spread()
	if _, err := insInvStmt.Exec(date, code, qty, unitprice); err != nil {
		log.Fatal("Failed to insert into investment: ", err.Error())
	}
}

// Reading Rows

/* Returns investments made during within a date range */
func GetInvestmentsRecent(page int) ([]Investment, error) {
	rows, err := getInvRecStmt.Query(page*PAGE_ROWS, PAGE_ROWS)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return dbRowsToInvestments(rows)
}

/* Returns investments matching a specified filter */
func GetInvestmentsFilter(opts FilterOpts) ([]Investment, error) {
	rows, err := getInvFilStmt.Query(opts.minCost, opts.maxCost, opts.startDate, opts.endDate, "%"+opts.code+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return dbRowsToInvestments(rows)
}

/* Returns records from within a date range */
func GetRecordsRecent(page int8) []Record {
	rows, err := getRecRecStmt.Query(page*PAGE_ROWS, PAGE_ROWS)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	return dbRowsToRecords(rows)
}

/* Returns records matching a specified filter */
func GetRecordsFilter(opts FilterOpts) []Record {
	cmd := `SELECT rec_id, rec_date, rec_desc, rec_amt, cat_id
          FROM record
          WHERE rec_amt BETWEEN ? AND ?
            AND rec_date BETWEEN ? AND ?`
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

/* Returns a slice containing all of the categories */
func GetCategories() []Category {
	rows, err := getCategoriesStmt.Query()
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	return dbRowsToCategories(rows)
}

/* Returns the total income over a date range (inclusive) */
func GetIncomeSum(startDate, endDate time.Time) (float32, error) {
	var sum float32
	if err := getIncomeSumStmt.QueryRow(startDate, endDate).Scan(&sum); err != nil {
		return 0, err
	}
	return sum, nil
}

/* Returns the total expenditure over a date range (inclusive), flips sign (expenditure > 0) */
func GetExpenditureSum(startDate, endDate time.Time) (float32, error) {
	var sum float32
	if err := getExpenditureSumStmt.QueryRow(startDate, endDate).Scan(&sum); err != nil {
		return 0, err
	}
	return -sum, nil
}

/* Returns the total money in/out for a given category over a date range */
func GetCategorySum(catId int, startDate, endDate time.Time) (float32, error) {
	var sum float32
	if err := getCategorySumStmt.QueryRow(catId, startDate, endDate).Scan(&sum); err != nil {
		return 0, err
	}
	return sum, nil
}

// Frontend Helper Functions

func GetInvestmentsPages() int8 {
	var res float64
	db.QueryRow("SELECT COUNT(*) / ? FROM investment", float32(PAGE_ROWS)).Scan(&res)
	return int8(math.Ceil(res))
}

func GetRecordsPages() int8 {
	var res float64
	db.QueryRow("SELECT COUNT(*) / ? FROM record", float32(PAGE_ROWS)).Scan(&res)
	return int8(math.Ceil(res))
}

func GetCategoryNameFromId(catId int) string {
	var res string
	db.QueryRow("SELECT cat_name FROM category WHERE cat_id = ?", catId).Scan(&res)
	return res
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
	_, date, code, qty, unitprice := inv.Spread()
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
