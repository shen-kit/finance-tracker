import datetime
import os
import sqlite3

from tabulate import tabulate

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
Info Display
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


def display_category_records(conn: sqlite3.Connection, cur: sqlite3.Cursor, cat_id) -> None:
    """
    Display income/expenditure for a given category.
    Format: | Date | Description | Amount |
    """
    sql = "SELECT rec_date, rec_desc, rec_amt FROM RECORD WHERE cat_id = ? ORDER BY rec_date DESC;"
    cur.execute(sql, (cat_id,));
    res = cur.fetchall()

    s = tabulate(res, headers=["Date", "Description", "Amout"], tablefmt="grid")
    print(s)
    
    conn.commit()


def display_investment_summary() -> None:
    """
    Display all current holdings, average buy price, current price, current value, profit, and profit percentage
    """
    raise NotImplementedError


"""
Editing the Database
"""


def add_record(conn: sqlite3.Connection, cur: sqlite3.Cursor) -> bool:
    """
    Insert a record of income/expenditure to the database.
    Returns True if successful, False otherwise.
    """

    def get_date_from_input() -> datetime.date:
        """
        Returns a datetime Date object.
        Available formats:
          - YYYY-(M)M-(D)D
          - (M)M-(D)D          : current year
          - (D)D               : current year and month
        """
        while True:
            d: list = input("Date: ").split("-")

            if d == [""]:
                return datetime.date.today()

            try:
                # convert entered YYYY-MM-DD to integers
                d = list(map(lambda s: int(s), d))
                if len(d) == 1:
                    d.insert(0, datetime.datetime.now().month)
                if len(d) == 2:
                    d.insert(0, datetime.datetime.now().year)
                return datetime.date(d[0], d[1], d[2])

            except ValueError:
                print(
                    (
                        "Date must be valid, and in one of the following formats (blank for today):\n"
                        " 1. YYYY-MM-DD\n"
                        " 2. MM-DD (current year)\n"
                        " 3. DD (current month and year)\n"
                    )
                )

    def get_category_id() -> int:
        """
        Show all categories and their IDs, receive input and validate that it upholds referential integrity
        """
        while True:
            try:
                return int(input("Category ID: "))
            except ValueError:
                print("Category ID must be an integer.")

    def get_amount() -> float:
        while True:
            try:
                return float(input("Amount: "))
            except ValueError:
                print("Please enter a number.")

    date: datetime.date = get_date_from_input()
    category: int = get_category_id()
    description: str = input("Description: ")
    amount: float = get_amount()

    # save to db
    sql = "INSERT INTO RECORD (rec_date, cat_id, rec_desc, rec_amt) VALUES (?,?,?,?)"
    cur.execute(sql, (date, category, description, amount))
    conn.commit()
    return cur.rowcount > 0  # was the insertion successful?


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


def select_option() -> str:
    return input(
        (
            "\nOptions:\n"
            "a: Add record\n"
            "r: Records (edit/delete)\n"
            "c: Categories (add/edit/delete)\n"
            "m: View monthly report\n"
            "y: View annual report\n"
            "q: quit\n"
            ": "
        )
    ).lower()


if __name__ == "__main__":
    base_dir = os.path.dirname(__file__)
    conn, cur = initialise_db(base_dir)

    while True:
        i = select_option()

        match i:
            case "a":
                add_record(conn, cur)
            case "v":
                cur.execute("SELECT * FROM RECORD;")
                rows = cur.fetchall()
                print(rows)
            case "r":
                pass
            case "c":
                display_category_records(conn, cur, 0)
            case "m":
                pass
            case "y":
                pass
            case "q":
                print("\nExiting...\n")
                break
            case _:
                print("Invalid option entered.\n")
