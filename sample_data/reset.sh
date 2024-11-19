#!/usr/bin/env bash

test -e sample_data/data/food.db && rm sample_data/data/food.db

go run ./cmd init
sqlite3 sample_data/data/food.db < sample_data/seed.sql
