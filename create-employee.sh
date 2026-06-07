#!/bin/bash

curl -X POST http://localhost:8080/employees \
  -H "Content-Type: application/json" \
  -d '{
    "id": "emp_001",
    "name": "Jane Doe",
    "department": "Engineering",
    "role": "Software Engineer",
    "salary": 95000
  }'
