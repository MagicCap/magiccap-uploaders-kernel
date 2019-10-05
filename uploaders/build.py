#!/usr/local/bin/python3

import json
import os
import base64


def v1_load_dir(d):
    """Loads V1 in directory form."""
    scandir = os.scandir(d)
    items = {}
    for file in scandir:
        if file.is_file():
            loaded = json.load(open(file.path))
            loaded['icon'] = base64.encodebytes(open(os.path.join(d, loaded['icon']), "rb").read()).decode()
            items[file.name.split(".")[0]] = loaded
    return items


def build_v1():
    """Builds API V1."""
    v1_imports = json.load(open("./v1_imports.json"))
    items = {}
    for value in v1_imports:
        items[value[0]] = v1_load_dir(value[1])
    json.dump(items, open("./v1.json", "w+"))


build_v1()
