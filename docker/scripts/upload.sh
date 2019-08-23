#!/bin/bash
curl -sX POST localhost:8080/api/v1/upload   -F "file=@machines.xlsx"  -H "Content-Type: multipart/form-data"
