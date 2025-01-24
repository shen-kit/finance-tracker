import datetime
import os
import sqlite3
from sys import exit

import yfinance as yf
from pyfzf.pyfzf import FzfPrompt
from tabulate import tabulate


class Investment:
    def __init__(self, id, code, date, price, qty) -> None:
        self.id: int = id
        self.code: str = code
        self.date: datetime.date = date
        self.price: float = price
        self.qty: float = qty

    def to_tuple(self) -> tuple[int, datetime.date, str, float, float]:
        return (self.id, self.date, self.code, self.price, self.qty)

    def __str__(self):
        return f"{self.id:2}:   {self.date.isoformat()}   {self.code:7}   {self.price:8.2f}   {self.qty}"


class FinanceTracker:

    def __init__(self) -> None:
        base_dir = os.path.dirname(__file__)
        self.conn, self.cur = self.initialise_db(base_dir)
        self.options, self.opt_str = self.define_options()

    def main_loop(self) -> None:
        try:
            while True:
                self.select_action()()
        except KeyboardInterrupt:
            print()
            self.main_loop()

    def quit(self) -> None:
        self.conn.close()
        print("\nDatabase connection closed.\n")
        exit(0)

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
                    cat_name   VARCHAR(30)    NOT NULL
                ); """
            )
            cur.execute(
                """CREATE TABLE IF NOT EXISTS RECORD (
                    rec_id   INTEGER      PRIMARY KEY,
                    cat_id   INTEGER      NOT NULL         DEFAULT "1",
                    rec_date DATE         NOT NULL,
                    rec_desc VARCHAR(100) NOT NULL,
                    rec_amt  NUMBER(9,2)  NOT NULL,
                    FOREIGN KEY (cat_id)  REFERENCES CATEGORY (cat_id) ON UPDATE SET DEFAULT
                );"""
            )
            cur.execute(
                """CREATE TABLE IF NOT EXISTS INVESTMENT (
                    inv_id     INTEGER    PRIMARY KEY,
                    inv_code   VARCHAR(7) NOT NULL,
                    inv_date   DATE       NOT NULL,
                    inv_price  INTEGER    NOT NULL,
                    inv_qty    SMALLINT   NOT NULL
                );"""
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
            "a": ("Add Record", self.add_record),
            "ac": ("Add Category", self.add_category),
            "ai": ("Add Investment", self.add_investment),
            # read
            "lr": ("List Records", self.list_records_for_category),
            "lc": ("List Categories", self.list_categories),
            "li": ("List Investments", self.list_investments),
            "is": ("Show Investment Summary", self.display_investment_summary),
            "m": ("Show Month Report", self.display_month_report),
            "y": ("Show Year Report", self.display_year_report),
            # update
            "u": ("Edit Record", self.edit_record),
            "uc": ("Edit Category", self.edit_category),
            "ui": ("Edit Investment", self.edit_investment),
            # delete
            "d": ("Delete Record", self.delete_record),
            "dc": ("Delete Category", self.delete_category),
            "di": ("Delete Investment", self.delete_investment),
            # quit
            "q": ("Quit", self.quit),
        }
        opt_str = (
            "\nWhat would you like to do?\n"
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
            print()
            return self.options[sel][1]

    """
    User Actions (CRUD)
    """

    # Create

    def add_record(self) -> None:
        """
        Insert a record of income/expenditure to the database.
        """
        print("New Record:")
        date: datetime.date = self.get_date_from_input()
        cat_id, cat_name = self.get_category()
        print(f"Category: {cat_name}")
        description: str = input("Description: ")
        amount: float = self.get_float("Amount: ")

        # save to db
        self.cur.execute(
            "INSERT INTO RECORD (rec_date, cat_id, rec_desc, rec_amt) VALUES (?,?,?,?)",
            (date, cat_id, description, amount),
        )
        self.conn.commit()

        print("Record saved.")

    def add_category(self) -> None:
        """
        Create a new category.
        """
        print("New Category:")
        categories = self.get_categories().values()
        while True:
            # store lowercase to prevent duplicates, can capitalise when displaying
            cname = input("New category name: ").lower()
            if cname not in categories:
                break
            print("That category already exists!")
        self.cur.execute("INSERT INTO CATEGORY (cat_name) VALUES (?)", (cname,))
        self.conn.commit()
        print("Category saved.")

    def add_investment(self) -> None:
        print("New Investment:")
        date: datetime.date = self.get_date_from_input()
        code: str = input("Stock Code: ").upper()
        qty: float = self.get_float("Quantity: ")
        price: float = self.get_float("Unit Price: ")

        # save to db
        self.cur.execute(
            "INSERT INTO INVESTMENT (inv_date, inv_code, inv_qty, inv_price) VALUES (?,?,?,?)",
            (date, code, qty, price),
        )
        self.conn.commit()
        print("Investment saved.")

    # Read

    def list_categories(self) -> None:
        """
        Display a list of all categories next to their ID.
        """
        c = self.get_categories()
        fmt = map(lambda kv: f"{kv[1]}: {kv[0]}", c.items())
        print("Categories:")
        print("\n".join(fmt))

    def list_records_for_category(self) -> None:
        """
        Request a category id, then display records for the chosen category.
        Format: | Date | Description | Amount |
        """
        cat_id, _ = self.get_category()
        limit = self.get_int("Number of records to show: ")
        res = self.cur.execute(
            "SELECT rec_date, rec_desc, rec_amt FROM RECORD WHERE cat_id = ? ORDER BY rec_date DESC LIMIT ?;",
            (cat_id, limit),
        )
        print("\nRecords:")
        print(tabulate(res, headers=["Date", "Description", "Amout"]))

    def list_investments(self) -> None:
        """
        List all investments that have been made in table format
        """
        investments: list[Investment] = self.get_investments()
        print("Investments:")
        print(
            tabulate(
                [i.to_tuple() for i in investments],
                headers=["ID", "Date", "Code", "Unit Price", "Qty"],
            )
        )

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
        # get [ code, qty, total_purchase_cost ]
        res = self.cur.execute(
            "SELECT inv_code, SUM(inv_qty), SUM(inv_qty * inv_price) FROM INVESTMENT GROUP BY inv_code;"
        ).fetchall()
        table_data = []
        for r in res:
            unitprice = self.get_stock_price(r[0])
            table_data.append(
                [
                    r[0],  # code
                    r[1],  # qty
                    r[2] / r[1],  # avg buy price
                    r[2],  # total buy price
                    unitprice,  # current unit price
                    unitprice * r[1],  # current value
                    unitprice * r[1] - r[2],  # profit
                ]
            )
        print(
            tabulate(
                table_data,
                headers=[
                    "Code",
                    "Qty",
                    "Avg Buy",
                    "Buy Value",
                    "Curr Price",
                    "Curr Value",
                    "Profit/Loss",
                ],
                floatfmt=".2f",
            )
        )

    # Update

    def edit_record(self) -> None:
        """
        Edit a record of income/expenditure in the database.
        Returns True if successful, False otherwise.
        """
        raise NotImplementedError

    def edit_category(self) -> None:
        cid, oldname = self.get_category()
        newname = input("New category name: ").lower()
        self.cur.execute(
            "UPDATE CATEGORY SET cat_name = ? WHERE cat_id = ?", (newname, cid)
        )
        self.conn.commit()
        print(f"Category '{oldname}' renamed to '{newname}'.")

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
        cid, cname = self.get_category()
        self.cur.execute("DELETE FROM CATEGORY WHERE cat_id = ?;", (cid,))
        self.conn.commit()
        print(f"Category '{cname}' deleted.")

    def delete_investment(self) -> None:
        raise NotImplementedError

    """
    Helper Functions -> Database Queries
    """

    def get_categories(self) -> dict[str, int]:
        """
        Returns: a dictionary of {cat_name: cat_id}
        """
        res = self.cur.execute("SELECT * FROM CATEGORY;")
        return {r[1]: r[0] for r in res}

    def get_investments(self) -> list[Investment]:
        res = self.cur.execute("SELECT * FROM INVESTMENT;")
        return [Investment(*r) for r in res]

    """
    Helper Functions -> User Input
    """

    def get_category(self) -> tuple[int, str]:
        """
        Use FZF to let the user select a category.
        Return (id, name)
        """
        categories = self.get_categories()
        names = categories.keys()
        try:
            sel = FzfPrompt().prompt(names, "--cycle")[0]
            return (categories[sel], sel)
        except IndexError:  # Ctrl+C pressed during selection
            raise KeyboardInterrupt

    def get_investment(self) -> Investment:
        self.list_investments()
        while True:
            in_id = input("Investment ID: ")
            res = self.cur.execute(
                "SELECT * FROM INVESTMENT WHERE inv_id = ?", (in_id,)
            ).fetchone()
            if res is not None:
                return Investment(*res)
            print("Invalid selection.")

    @staticmethod
    def get_int(prompt: str, allow_blank=False) -> int | None:
        """
        Get an integer as input from the user.
        If `allow_blank`, return `None` when no input is given.
        """
        while True:
            try:
                inp = input(prompt)
                if allow_blank and inp == "":
                    return None
                return int(inp)
            except ValueError:
                print("Please enter an integer.")

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
