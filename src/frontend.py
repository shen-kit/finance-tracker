from textual.app import App, ComposeResult
from textual.containers import Container, Horizontal, Vertical
from textual.events import Key
from textual.screen import ModalScreen, Screen
from textual.widget import Widget
from textual.widgets import (Button, Digits, Footer, Header, Input, Label,
                             ListItem, ListView, Select, Switch)

from backend import Category


class FocusableLabel(Label, Widget):

    def __init__(self, text: str, action=lambda: None):
        super().__init__()
        self.text = text
        self.action = action

    def compose(self) -> ComposeResult:
        yield Label(self.text)

    def on_mount(self):
        self.can_focus = True

    def on_focus(self) -> None:
        self.styles.color = "lime"

    def on_blur(self) -> None:
        self.styles.color = None

    def on_key(self, event: Key) -> None:
        if event.key == "enter":
            self.action()


class MyListView(ListView):
    BINDINGS = [
        ("j,down", "cursor_down"),
        ("k,up", "cursor_up"),
    ]


class AddCategoryScreen(Screen):

    def compose(self):
        yield Container(
            Label("New Category", id="title"),
            Input(placeholder="Name", id="name_input"),
            Input(placeholder="Description", id="desc_input"),
            Horizontal(Label("Type:"), Select([("Income", "I"), ("Expenditure", "E")])),
            Horizontal(
                FocusableLabel("Add"),
                FocusableLabel("Cancel"),
            ),
        )


class CategoriesScreen(Screen):

    def __init__(self):
        super().__init__()
        self.categories = [
            Category(1, "cat1", "category description 1", "I"),
            Category(2, "cat2", "category description 2", "E"),
            Category(3, "cat3", "category description 3", "E"),
        ]

    def compose(self) -> ComposeResult:
        cats = []
        for cat in self.categories:
            cats.append(ListItem(FocusableLabel(cat.name)))
        yield MyListView(*cats)


class FTApp(App):

    CSS_PATH = "app.css"
    SCREENS = {
        "categories": CategoriesScreen,
        "add_category": AddCategoryScreen,
    }
    BINDINGS = [
        ("j,down", "focus_next", "Focus next"),
        ("k,up", "focus_previous", "Focus next"),
        ("escape", "app.back", "Back"),
    ]

    def on_mount(self) -> None:
        self.theme = "catppuccin-mocha"

    def compose(self) -> ComposeResult:
        yield Header()
        yield Footer()
        self.listview = MyListView(
            ListItem(Label("View Categories"), id="categories"),
            ListItem(Label("Add Category"), id="add_category"),
        )
        yield self.listview

    def on_list_view_selected(self, event: ListView.Selected) -> None:
        if event.item.id is not None:
            self.push_screen(event.item.id)


if __name__ == "__main__":
    app = FTApp()
    app.run()
