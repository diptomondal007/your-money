#! /bin/bash

export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=test
export DB_PASSWORD=password
export DB_NAME=your-money

echo "building application ...."
go build -v .
./your-money serve