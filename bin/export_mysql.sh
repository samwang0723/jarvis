#!/bin/bash

# Check if the table name is provided
if [ -z "$1" ]; then
  echo "Usage: $0 "
  exit 1
fi

# For migrating the legacy MySQL tables
# Usage: ./bin/export_mysql.sh {TABLE_NAME}
TABLE=$1

kubectl exec -it mysql-primary-0 -- mysqldump -u root -piori2008 jarvis --tables ${TABLE} --no-create-info --skip-comments --complete-insert --compatible=ansi --compact > ./${TABLE}_raw.sql
sed -E 's/\([0-9]+,/\(/g' ./${TABLE}_raw.sql > ./${TABLE}_remove_id.sql
sed -E 's/"id", //g' ./${TABLE}_remove_id.sql > ./${TABLE}.sql

rm ./${TABLE}_raw.sql
rm ./${TABLE}_remove_id.sql
