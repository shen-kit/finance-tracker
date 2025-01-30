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
        self.id: int = id
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


class Investment(Entity):

    def __init__(self, id, date, code, qty, unit_price) -> None:
        super().__init__()
        self.id: int = id
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


class Category(Entity):

    def __init__(self, id, name, desc) -> None:
        super().__init__()
        self.id: int = id
        self.name: str = name
        self.desc: str = desc

    @override
    def to_tuple(self, include_id: bool) -> tuple:
        if include_id:
            return (self.id, self.name, self.desc)
        else:
            return (self.name, self.desc)

    @override
    @classmethod
    def from_list(cls, l: list) -> Category:
        return Category(*l)


class FinanceTracker:

    # constants
    INVESTMENTS_PER_PAGE = 2
    RECORDS_PER_PAGE = 2

    def __init__(self, db_path) -> None:
        self.conn, self.cur = self.connect_to_db(db_path)
        self.initialise_db()

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
                    cat_name VARCHAR(30) NOT NULL
                );

                CREATE TABLE IF NOT EXISTS record (
                    rec_id   INTEGER     NOT NULL  PRIMARY KEY,
                    rec_date DATE        NOT NULL,
                    rec_desc VARCHAR(50) NOT NULL,
                    cat_id   INTEGER     NOT NULL,  
                    rec_amt  NUMBER(7,2) NOT NULL,
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

    # getter helpers

    def get_investments_recent(self, page: int) -> list[Investment]:
        sql = "SELECT (inv_id, inv_date, inv_code, inv_qty, inv_unitprice) FROM investment ORDER BY inv_date DESC LIMIT ?, ?"
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
        min_cost: float,
        max_cost: float,
        start_date: dt.date,
        end_date: dt.date,
        code: str,
    ) -> list[Investment]:
        raise NotImplementedError

    def get_investments_summary(self):
        raise NotImplementedError

    def get_categories(self) -> list[Category]:
        raise NotImplementedError

    def get_records_recent(self, page: int) -> list[Record]:
        sql = "SELECT (rec_id, rec_date, rec_desc, cat_id, rec_amt) FROM record ORDER BY rec_date DESC LIMIT ?, ?"
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
        min_cost: float,
        max_cost: float,
        start: dt.date,
        end: dt.date,
        desc: str,
        cat_id: int,
    ) -> list[Record]:
        raise NotImplementedError

    def get_income_sum(self, start_date: dt.date, end_date: dt.date) -> float:
        raise NotImplementedError

    def get_expenditure_sum(self, start_date: dt.date, end_date: dt.date) -> float:
        raise NotImplementedError

    def get_category_sum(
        self, cat_id: int, start_date: dt.date, end_date: dt.date
    ) -> float:
        raise NotImplementedError


if __name__ == "__main__":
    base_dir = os.path.dirname(__file__)
    db_path = os.path.join(base_dir, "testing.db")
    ft = FinanceTracker(db_path)
