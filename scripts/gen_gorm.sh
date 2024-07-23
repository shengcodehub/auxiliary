#!/bin/bash
# Generate GORM code from dev database schema
echo "start generating gorm files..."
go run ./scripts/gorm/*.go
echo "done!!!"