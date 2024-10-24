#!/bin/bash

WORK_DIR="$(
    cd $(dirname $0)
    pwd
)/.."
cd $WORK_DIR

python3 -m venv ".venv"
source .venv/bin/activate
pip install -r requirements.txt
