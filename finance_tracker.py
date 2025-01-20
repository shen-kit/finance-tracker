import datetime
import os
import sqlite3

import yfinance as yf

"""
Helper Functions
"""


def get_stock_price(stock_code: str) -> float:
    """
    Uses yfinance (yahoo finance backend) to get the most recent stock price.
    Calculates stock price as the average of the bid and ask prices.
    =======================
    Inputs:
        - stock_code (str): must reflect the stock code in yahoo finance
    """
    info = yf.Ticker(stock_code).info
    return (info["bid"] + info["ask"]) / 2


"""
Database Initialisation
"""


def initialise_db(base_dir: str) -> tuple[sqlite3.Connection, sqlite3.Cursor]:

    def create_tables(cur: sqlite3.Cursor):
        cur.execute(
            """CREATE TABLE IF NOT EXISTS CATEGORY (
            cat_id     INTEGER        PRIMARY KEY,
            cat_name   VARCHAR(30)    NOT NULL);
            """
        )
        cur.execute(
            """CREATE TABLE IF NOT EXISTS RECORD (
            rec_id   INTEGER      PRIMARY KEY,
            cat_id   INTEGER      NOT NULL         DEFAULT "1",
            rec_date DATE         NOT NULL,
            rec_desc VARCHAR(100) NOT NULL,
            rec_amt  NUMBER(9,2)  NOT NULL,
            FOREIGN KEY (cat_id) REFERENCES CATEGORY (cat_id) ON UPDATE SET DEFAULT);"""
        )
        cur.execute(
            """CREATE TABLE IF NOT EXISTS INVESTMENT (
            inv_id     INTEGER    PRIMARY KEY,
            inv_code   VARCHAR(7) NOT NULL,
            inv_date   DATE       NOT NULL,
            inv_price  INTEGER    NOT NULL,
            inv_qty    SMALLINT   NOT NULL);
            """
        )

    def create_default_category(cur: sqlite3.Cursor):
        """create default fallback category"""
        if cur.execute("SELECT * FROM CATEGORY;").fetchone() is None:
            cur.execute("INSERT INTO CATEGORY (cat_id, cat_name) VALUES (1, '-/-')")

    # convert datetime.date into a string (YYYY-MM-DD) when storing
    sqlite3.register_adapter(datetime.date, lambda d: d.isoformat())
    # convert string into datetime.date when retrieving
    sqlite3.register_converter(
        "DATE", lambda s: datetime.date.fromisoformat(s.decode("utf-8"))
    )

    # PARSE_DECLTYPES matches declared types with their converter
    # used for DATE -> datetime.date object conversion
    conn = sqlite3.connect(
        os.path.join(base_dir, "database.db"), detect_types=sqlite3.PARSE_DECLTYPES
    )
    cur = conn.cursor()

    create_tables(cur)
    create_default_category(cur)
    conn.commit()

    return (conn, cur)


"""
Querying the Database
"""

def display_month_report() -> None:
    """
    Display a detailed report on the month:
        - summary: income, expenditure, net change
        - records in tabular format
    """
    raise NotImplementedError

def display_year_report() -> None:
    """
    Display a detailed report on the year:
        - summary: income, expenditure, net change
        - table:   monthly income/expenditure/net change, categories
        - charts:  select categories
    """
    raise NotImplementedError

def display_category_records() -> None:
    """
    Display income/expenditure for a given category
    """
    raise NotImplementedError

def get_investment_summary() -> None:
    """
    Display all current holdings, average buy price, current price, current value, profit, and profit percentage
    """
    raise NotImplementedError


"""
Editing the Database
"""

def insert_record() -> bool:
    """
    Insert a record of income/expenditure to the database.
    Returns True if successful, False otherwise.
    """
    raise NotImplementedError

def delete_record() -> bool:
    """
    Delete a record of income/expenditure from the database.
    Returns True if successful, False otherwise.
    """
    raise NotImplementedError

def edit_record() -> bool:
    """
    Edit a record of income/expenditure in the database.
    Returns True if successful, False otherwise.
    """
    raise NotImplementedError

def create_category() -> bool:
    """
    Create a new category.
    Returns True if successful, False otherwise.
    """
    raise NotImplementedError

def edit_category() -> bool:
    """
    Edit the name of a category.
    Returns True if successful, False otherwise.
    """
    raise NotImplementedError

def delete_category() -> bool:
    """
    Delete a category.
    The category of records with this category will be changed to '-/-'.

    Returns True if successful, False otherwise.
    """
    raise NotImplementedError

"""
Main
"""

if __name__ == "__main__":
    base_dir = os.path.dirname(__file__)
    conn, cur = initialise_db(base_dir)
