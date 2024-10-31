import os

_DATA_DIR = os.path.join(os.path.dirname(__file__), "..", "resource")


def save_to_file(file_name, data):
    path = os.path.join(_DATA_DIR, file_name)
    dir = os.path.abspath(os.path.split(path)[0])
    os.makedirs(dir, exist_ok=True)
    with open(path, "+w") as file:
        file.write(data)


def load_from_file(file_name):
    path = os.path.join(_DATA_DIR, file_name)
    with open(path, "r") as file:
        return file.read()
