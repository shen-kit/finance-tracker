package backend

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func SetupDb(path string) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// create tables
	sql := `
  CREATE TABLE IF NOT EXISTS category (
      cat_id   INTEGER     NOT NULL  PRIMARY KEY,
      cat_name VARCHAR(20) NOT NULL,
      cat_type CHAR(1)     NOT NULL  CHECK (cat_type IN ('I', 'E')),
      cat_desc VARCHAR(40)
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
	_, err = db.Exec(sql)
	if err != nil {
		log.Printf("%q: %s\n", err, sql)
	}
}
