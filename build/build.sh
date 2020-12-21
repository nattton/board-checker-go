#!/bin/bash
mkdir -p bin
go build -o bin/web ./cmd/web/
go build -o bin/admin ./cmd/admin/