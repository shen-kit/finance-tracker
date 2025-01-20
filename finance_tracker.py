import os
import sqlite3


def initialise_db(base_dir):

    def create_tables(cur):
        cur.execute(
            """CREATE TABLE IF NOT EXISTS CATEGORY (
            cat_id INT PRIMARY KEY,
            cat_name   VARCHAR(30) NOT NULL);"""
        )
        cur.execute(
            """CREATE TABLE IF NOT EXISTS RECORD (
            rec_id   INT          PRIMARY KEY,
            cat_id   INT          NOT NULL         DEFAULT "1",
            rec_date DATE         NOT NULL,
            rec_desc VARCHAR(100) NOT NULL,
            rec_amt  INT          NOT NULL,
            FOREIGN KEY (cat_id) REFERENCES CATEGORY (cat_id) ON UPDATE SET DEFAULT);"""
        )
        cur.execute(
            """CREATE TABLE IF NOT EXISTS INVESTMENT (
            inv_id     INT        PRIMARY KEY,
            inv_code   VARCHAR(7) NOT NULL,
            inv_date   DATE       NOT NULL,
            inv_price  INT        NOT NULL,
            inv_qty    SMALLINT   NOT NULL);"""
        )

    def create_default_category(cur):
        """create default fallback category"""
        if cur.execute("SELECT * FROM CATEGORY;").fetchone() is None:
            cur.execute("INSERT INTO CATEGORY (cat_id, cat_name) VALUES (1, '-/-')")

    conn = sqlite3.connect(os.path.join(base_dir, "database.db"))
    cur = conn.cursor()

    create_tables(cur)
    create_default_category(cur)
    conn.commit()

    return conn, cur


if __name__ == "__main__":

    base_dir = os.path.dirname(__file__)
    conn, cur = initialise_db(base_dir)
