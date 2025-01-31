from __future__ import annotations

import datetime as dt
import os
import sqlite3
from abc import ABC, abstractmethod
from typing import override


class Entity(ABC):
    @abstractmethod
    def to_tuple(self, include_id: bool) -> tuple:
        """
        Returns a tuple of the attribute entities in the same order as defined in the schema.

        Inputs:
            - include_id (bool): whether to include the ID as part of the tuple
        """
        return ()

    @classmethod
    @abstractmethod
    def from_list(cls, l: list) -> Entity:
        pass


class Record(Entity):

    def __init__(self, id, date, desc, amt, cat_id) -> None:
        super().__init__()
        self.id: int | None = id
        self.date: dt.date = date
        self.desc: str = desc
        self.amt: float = amt
        self.cat_id: int = cat_id

    @override
    def to_tuple(self, include_id: bool) -> tuple:
        if include_id:
            return (self.id, self.date, self.desc, self.amt, self.cat_id)
        else:
            return (self.date, self.desc, self.amt, self.cat_id)

    @override
    @classmethod
    def from_list(cls, l: list) -> Record:
        return Record(*l)

    @override
    def __repr__(self) -> str:
        return f"Record ({self.id}, {self.date}, {self.desc}, ${self.amt:.2f}, {self.cat_id})"

class Investment(Entity):

    def __init__(self, id, date, code, qty, unit_price) -> None:
        super().__init__()
        self.id: int | None = id
        self.date: dt.date = date
        self.code: str = code
        self.qty: int = qty
        self.unit_price: float = unit_price

    @override
    def to_tuple(self, include_id: bool) -> tuple:
        if include_id:
            return (self.id, self.date, self.code, self.qty, self.unit_price)
        else:
            return (self.date, self.code, self.qty, self.unit_price)

    @override
    @classmethod
    def from_list(cls, l: list) -> Investment:
        return Investment(*l)

    @override
    def __repr__(self) -> str:
        return f"Investment ({self.id}, {self.date}, {self.code}, {self.qty} x ${self.unit_price:.2f})"


class Category(Entity):

    def __init__(self, id, name, desc, ctype) -> None:
        super().__init__()
        if ctype not in ("I", "E"):
            raise Exception(
                f"Category type must be 'I' or 'E'. '{ctype}' was supplied."
            )
        self.id: int | None = id
        self.name: str = name
        self.desc: str = desc
        self.ctype: str = ctype

    @override
    def to_tuple(self, include_id: bool) -> tuple:
        if include_id:
            return (self.id, self.name, self.desc, self.ctype)
        else:
            return (self.name, self.desc, self.ctype)

    @override
    @classmethod
    def from_list(cls, l: list) -> Category:
        return Category(*l)

    @override
    def __repr__(self) -> str:
        return f"Category ({self.id}, {self.name}, {self.desc}, {self.ctype})"


class FinanceTracker:

    # constants
    INVESTMENTS_PER_PAGE = 2
    RECORDS_PER_PAGE = 2

    def __init__(self, db_path) -> None:
        self.conn, self.cur = self.connect_to_db(db_path)
        self.initialise_db()

    # database connection management

    def connect_to_db(self, db_path) -> tuple[sqlite3.Connection, sqlite3.Cursor]:
        """
        Establish a database connection, return the connection and the cursor.
        Creates a new database if one at the specified path does not yet exist.
        """
        # auto convert dt.date <-> sql DATE type
        sqlite3.register_adapter(dt.date, lambda d: d.isoformat())
        sqlite3.register_converter(
            "DATE", lambda s: dt.date.fromisoformat(s.decode("utf-8"))
        )

        conn = sqlite3.connect(db_path, detect_types=sqlite3.PARSE_DECLTYPES)
        cur = conn.cursor()
        return conn, cur

    def initialise_db(self) -> None:
        """
        Create table structure if not yet exists
        """
        self.cur.executescript(
            """
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
                """
        )

    def close_db(self):
        self.conn.close()
        print("\nDatabase connection closed.\n")

    # insert data

    def add_record(self, record: Record) -> None:
        self.cur.execute(
            "INSERT INTO record (rec_date, rec_desc, rec_amt, cat_id) VALUES (?,?,?,?)",
            record.to_tuple(False),
        )
        self.conn.commit()

    def add_category(self, category: Category) -> None:
        self.cur.execute(
            "INSERT INTO category (cat_name, cat_desc, cat_type) VALUES (?,?,?)",
            category.to_tuple(False),
        )
        self.conn.commit()

    def add_investment(self, investment: Investment) -> None:
        self.cur.execute(
            "INSERT INTO investment (inv_date, inv_code, inv_qty, inv_unitprice) VALUES (?,?,?,?)",
            investment.to_tuple(False),
        )
        self.conn.commit()

    # getter helpers

    def get_investments_recent(self, page: int) -> list[Investment]:
        sql = "SELECT inv_id, inv_date, inv_code, inv_qty, inv_unitprice FROM investment ORDER BY inv_date DESC LIMIT ?, ?;"
        res = self.cur.execute(
            sql,
            (
                page * FinanceTracker.INVESTMENTS_PER_PAGE,
                FinanceTracker.INVESTMENTS_PER_PAGE,
            ),
        )
        return list(map(Investment.from_list, res))

    def get_investments_filter(
        self,
        *,
        min_cost: float = -999999999,
        max_cost: float = 999999999,
        start_date: dt.date = dt.date.fromordinal(1),
        end_date: dt.date = dt.date(9999, 12, 31),
        code: str = "",
    ) -> list[Investment]:
        sql = """SELECT *
                 FROM investment
                 WHERE inv_qty*inv_unitprice BETWEEN ? AND ?
                   AND inv_date BETWEEN ? AND ?
                   AND inv_code LIKE ?;
                 """
        res = self.cur.execute(
            sql, (min_cost, max_cost, start_date, end_date, "%" + code + "%")
        )
        return list(map(Investment.from_list, res))

    def get_investments_summary(self):
        raise NotImplementedError

    def get_categories(self, ctype: int = 0) -> list[Category]:
        """
        Return a list of categories.
        Input:
            type (int): 0 for all, -1 for expenditure, 1 for income
        """
        sql = "SELECT cat_id, cat_name, cat_desc, cat_type FROM category"
        match ctype:
            case -1:
                sql += " WHERE UPPER(cat_type) = 'E'"
            case 1:
                sql += " WHERE UPPER(cat_type) = 'I'"
        res = self.cur.execute(sql + ";")
        return list(map(Category.from_list, res))

    def get_records_recent(self, page: int) -> list[Record]:
        sql = "SELECT rec_id, rec_date, rec_desc, rec_amt, cat_id FROM record ORDER BY rec_date DESC LIMIT ?, ?;"
        res = self.cur.execute(
            sql,
            (
                page * FinanceTracker.RECORDS_PER_PAGE,
                FinanceTracker.RECORDS_PER_PAGE,
            ),
        )
        return list(map(Record.from_list, res))

    def get_records_filter(
        self,
        *,
        min_cost: float = -99999999,
        max_cost: float = 99999999,
        start_date: dt.date = dt.date.fromordinal(1),
        end_date: dt.date = dt.date(9999, 12, 31),
        desc: str = "",
        cat_id: int = -1,
    ) -> list[Record]:
        sql = """SELECT * FROM record
                 WHERE rec_amt BETWEEN ? AND ?
                   AND rec_date BETWEEN ? AND ?
                   AND rec_desc LIKE ?
                   AND cat_id """
        if cat_id > 0:
            sql += "= ?;"
        else:
            sql += "> ?;"
        res = self.cur.execute(
            sql, (min_cost, max_cost, start_date, end_date, "%" + desc + "%", cat_id)
        )

        return list(map(Record.from_list, res))

    def get_income_sum(self, start_date: dt.date, end_date: dt.date) -> float:
        sql = """SELECT SUM(rec_amt)
                    FROM record
                    WHERE cat_id IN (SELECT cat_id FROM category WHERE UPPER(cat_type) = 'I')
                      AND rec_date BETWEEN ? AND ?;
                """
        self.cur.execute(sql, (start_date, end_date))
        return self.cur.fetchone()[0]

    def get_expenditure_sum(self, start_date: dt.date, end_date: dt.date) -> float:
        sql = """SELECT SUM(rec_amt)
                    FROM record
                    WHERE cat_id IN (SELECT cat_id FROM category WHERE UPPER(cat_type) = 'E')
                      AND rec_date >= ? AND rec_date <= ?;
                """
        self.cur.execute(sql, (start_date, end_date))
        return self.cur.fetchone()[0]

    def get_category_sum(
        self, cat_id: int, start_date: dt.date, end_date: dt.date
    ) -> float:
        sql = "SELECT SUM(rec_amt) FROM record WHERE cat_id = ? AND rec_date >= ? AND rec_date <= ?;"
        self.cur.execute(sql, (cat_id, start_date, end_date))
        return self.cur.fetchone()[0]


def create_dummy_data(ft: FinanceTracker):
    # categories
    ft.add_category(Category(None, "cat1", "cat desc 1", "I"))
    ft.add_category(Category(None, "cat2", "cat desc 2", "E"))
    ft.add_category(Category(None, "cat3", "cat desc 3", "E"))

    # investments
    ft.add_investment(Investment(None, dt.date(2024, 1, 1), "IVV", 5, 600))
    ft.add_investment(Investment(None, dt.date(2024, 6, 1), "VGS", 8, 700))
    ft.add_investment(Investment(None, dt.date(2025, 1, 1), "IVV", 15, 400))

    # records
    ft.add_record(Record(None, dt.date(2024, 1, 1), "record 1", -50, 2))
    ft.add_record(Record(None, dt.date(2025, 1, 1), "record 2", -10, 2))
    ft.add_record(Record(None, dt.date(2025, 1, 10), "record 3", 500, 1))


if __name__ == "__main__":
    base_dir = os.path.dirname(__file__)
    db_path = os.path.join(base_dir, "testing.db")
    ft = FinanceTracker(db_path)

    # create_dummy_data(ft)

    print("get_investments_recent: " + str(ft.get_investments_recent(0)))
    print("get_investments_filter: " + str(ft.get_investments_filter()))
    print("get_categories        : " + str(ft.get_categories()))
    print("get_records_recent    : " + str(ft.get_records_recent(0)))
    print("get_records_filter    : " + str(ft.get_records_filter()))
    print("get_income_sum        : " + str(ft.get_income_sum(dt.date(2023, 1, 1), dt.date(2030, 1, 1))))
    print("get_expenditure_sum   : " + str(ft.get_expenditure_sum(dt.date(2023, 1, 1), dt.date(2030, 1, 1))))
    print("get_category_sum      : " + str(ft.get_category_sum(1, dt.date(2023, 1, 1), dt.date(2030, 1, 1))))
