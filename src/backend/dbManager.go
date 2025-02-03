package backend

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// prepared statements
var insInvStmt *sql.Stmt
var insRecStmt *sql.Stmt
var insCatStmt *sql.Stmt

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
		{date: time.Now().AddDate(0, -1, 0), code: "IVV", qty: 10, unitprice: 600},
		{date: time.Now().AddDate(0, -1, 0), code: "VGS.AX", qty: 5, unitprice: 600},
		{date: time.Now().AddDate(0, -1, 0), code: "IVV", qty: 10, unitprice: 600},
	}
	for _, inv := range investments {
		insertInvestment(inv)
	}

	categories := [...]Category{
		{name: "Work", isIncome: true, desc: "income from work"},
		{name: "Groceries", isIncome: false, desc: "grocery spending"},
	}
	for _, cat := range categories {
		insertCategory(cat)
	}

	records := [...]Record{
		{date: time.Now(), desc: "new record desc 1", amt: 100, catId: 1},
		{date: time.Now(), desc: "new record desc 2", amt: -200, catId: 2},
		{date: time.Now(), desc: "new record desc 3", amt: 300, catId: 1},
	}
	for _, rec := range records {
		insertRecord(rec)
	}

	fmt.Println("Inserted dummy data")

}

// helper functions

func insertRecord(rec Record) {
	_, date, desc, amt, cat_id := rec.spread()
	if _, err := insRecStmt.Exec(date, desc, amt, cat_id); err != nil {
		log.Fatal("Failed to insert into category: ", err.Error())
	}
}

func insertCategory(cat Category) {
	_, name, isIncome, desc := cat.spread()
	if _, err := insCatStmt.Exec(name, isIncome, desc); err != nil {
		log.Fatal("Failed to insert into category: ", err.Error())
	}
}

func insertInvestment(inv Investment) {
	_, date, code, qty, unitprice := inv.spread()
	if _, err := insInvStmt.Exec(date, code, qty, unitprice); err != nil {
		log.Fatal("Failed to insert into investment: ", err.Error())
	}
}
