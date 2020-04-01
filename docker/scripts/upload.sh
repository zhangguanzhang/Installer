#!/bin/bash
curl -X POST localhost:8080/api/v1/upload     -H "Content-Type: multipart/form-data" -F "file=@machines.xlsx"
