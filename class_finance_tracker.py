import datetime
import os
import sqlite3

import yfinance as yf
from tabulate import tabulate


class FinanceTracker:

    def __init__(self) -> None:
        base_dir = os.path.dirname(__file__)
        self.conn, self.cur = self.initialise_db(base_dir)

        self.options, self.opt_str = self.define_options()

    def main_loop(self) -> None:
        while True:
            self.select_action()()

    """
    Initialisation
    """

    def initialise_db(self, base_dir: str) -> tuple[sqlite3.Connection, sqlite3.Cursor]:
        """
        Create the database and required tables if not yet exists.
        Return the database connection and cursor objects.
        """

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

        # automatically convert date object <-> string
        sqlite3.register_adapter(datetime.date, lambda d: d.isoformat())
        sqlite3.register_converter(
            "DATE", lambda s: datetime.date.fromisoformat(s.decode("utf-8"))
        )

        conn = sqlite3.connect(
            os.path.join(base_dir, "database.db"), detect_types=sqlite3.PARSE_DECLTYPES
        )
        cur = conn.cursor()

        create_tables(cur)
        create_default_category(cur)
        conn.commit()

        return (conn, cur)

    def define_options(self) -> tuple[dict[str, tuple], str]:
        "Returns the options, and the string to display what options are available"
        options = {
            # create
            "a":  ("Add Record",        self.add_record),
            "ac": ("Add Category",      self.add_category),
            "ai": ("Add Investment",    self.add_investment),
            # read
            "lr": ("List Records",      self.list_records_for_category),
            "lc": ("List Categories",   self.list_categories),
            "m":  ("Show Month Report", self.display_month_report),
            "y":  ("Show Year Report",  self.display_year_report),
            # update
            "u":  ("Edit Record",       self.edit_record),
            "uc": ("Edit Category",     self.edit_category),
            "ui": ("Edit Investment",   self.edit_investment),
            # delete
            "d":  ("Delete Record",     self.delete_record),
            "dc": ("Delete Category",   self.delete_category),
            "di": ("Delete Investment", self.delete_investment),
            # quit
            "q":  ("Quit",              quit),
        }
        opt_str = (
            "What would you like to do?\n"
            + "\n".join(list(map(lambda t: f"{t[0]:>2} -> {t[1][0]}", options.items())))
            + "\n: "
        )
        return (options, opt_str)

    def select_action(self):
        while True:
            sel = input(self.opt_str).lower()
            if sel not in self.options.keys():
                print("Invalid option. Please try again.")
                continue
            return self.options[sel][1]

    """
    User Actions (CRUD)
    """

    # Create

    def add_record(self) -> None:
        """
        Insert a record of income/expenditure to the database.
        """

        date: datetime.date = self.get_date_from_input()
        category: int = self.get_category_id()
        description: str = input("Description: ")
        amount: float = self.get_float("Amount: ")

        # save to db
        self.cur.execute(
            "INSERT INTO RECORD (rec_date, cat_id, rec_desc, rec_amt) VALUES (?,?,?,?)",
            (date, category, description, amount),
        )
        self.conn.commit()

    def add_category(self) -> None:
        """
        Create a new category.
        """
        cname = input("New category name: ")
        self.cur.execute("INSERT INTO CATEGORY (cat_name) VALUES (?)", (cname,))
        self.conn.commit()

    def add_investment(self) -> None:
        date: datetime.date = self.get_date_from_input()
        code: str           = input("Stock Code: ").upper()
        qty: float          = self.get_float("Quantity: ")
        price: float        = self.get_float("Unit Price: ")

        # save to db
        self.cur.execute(
            "INSERT INTO INVESTMENT (inv_date, inv_code, inv_qty, inv_price) VALUES (?,?,?,?)",
            (date, code, qty, price)
        )
        self.conn.commit()

    # Read

    def list_categories(self) -> None:
        """
        Display a list of all categories next to their ID.
        """
        res = self.cur.execute("SELECT * FROM CATEGORY;")
        fmt = map(lambda x: f"{x[0]}: {x[1]}", res)
        print("\n".join(fmt))

    def list_records_for_category(self) -> None:
        """
        Request a category id, then display records for the chosen category.
        Format: | Date | Description | Amount |
        """
        cat_id = self.get_category_id()
        self.cur.execute(
            "SELECT rec_date, rec_desc, rec_amt FROM RECORD WHERE cat_id = ? ORDER BY rec_date DESC;",
            (cat_id,),
        )
        res = self.cur.fetchall()
        print(tabulate(res, headers=["Date", "Description", "Amout"], tablefmt="grid"))

    def display_month_report(self) -> None:
        """
        Display a detailed report on the month:
            - summary: income, expenditure, net change
            - records in tabular format
        """
        raise NotImplementedError

    def display_year_report(self) -> None:
        """
        Display a detailed report on the year:
            - summary: income, expenditure, net change
            - table:   monthly income/expenditure/net change, categories
            - charts:  select categories
        """
        raise NotImplementedError

    def display_investment_summary(self) -> None:
        """
        Display all current holdings, average buy price, current price, current value, profit, and profit percentage
        """
        raise NotImplementedError

    # Update

    def edit_record(self) -> None:
        """
        Edit a record of income/expenditure in the database.
        Returns True if successful, False otherwise.
        """
        raise NotImplementedError

    def edit_category(self) -> None:
        raise NotImplementedError

    def edit_investment(self) -> None:
        raise NotImplementedError

    # Delete

    def delete_record(self) -> bool:
        """
        Delete a record of income/expenditure from the database.
        Returns True if successful, False otherwise.
        """
        raise NotImplementedError

    def delete_category(self) -> None:
        raise NotImplementedError

    def delete_investment(self) -> None:
        raise NotImplementedError

    """
    Helper Functions -> Database Queries
    """

    def get_categories(self) -> dict[int, str]:
        """
        Returns: a dictionary of {cat_id: cat_name}
        """
        res = self.cur.execute("SELECT * FROM CATEGORY;")
        d = {}
        for r in res:
            d[r[0]] = r[1]
        return d


    """
    Helper Functions -> User Input
    """

    def get_category_id(self) -> int:
        """
        Show all categories and their IDs, request a category ID from the user and validate that it upholds referential integrity.
        """
        self.list_categories()
        while True:
            try:
                return int(input(": "))
            except ValueError:
                print("Category ID must be an integer.")

    @staticmethod
    def get_float(prompt: str) -> float:
        while True:
            try:
                return float(input(prompt))
            except ValueError:
                print("Please enter a number.")

    @staticmethod
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

    @staticmethod
    def get_stock_price(stock_code: str) -> float:
        """
        Uses yfinance (yahoo finance backend) to get the most recent stock price.
        Calculates stock price as the average of the bid and ask prices.

        Inputs:
            - stock_code (str): must reflect the stock code in yahoo finance
        """
        info = yf.Ticker(stock_code).info
        return (info["bid"] + info["ask"]) / 2


if __name__ == "__main__":
    ft = FinanceTracker()
    ft.main_loop()
