#!/bin/bash

# lil bash script to create a database & table & install psql (psql install thru the script didn't work on macOS – unsure if it's related to my OS)

# check if psql is installed
if ! command -v psql &> /dev/null; then
    echo "psql could not be found, attempting to install..."
    sudo apt update
    sudo apt install -y postgresql-client
fi

# prompt the user for database inputs
read -p "do you want to create a new database? (yes/no): " CREATE_DB

# if the user chooses to create a database
if [ "$CREATE_DB" == "yes" ]; then
  read -p "DB_NAME: " DB_NAME
  read -p "DB_USER: " DB_USER
  read -s -p "DB_PASSWORD: " DB_PASSWORD
  echo ""
  createdb $DB_NAME -U $DB_USER
fi

# sql command to create the table
SQL_CMD="
DO \$$
BEGIN
   IF NOT EXISTS (
      SELECT FROM pg_catalog.pg_tables
      WHERE  schemaname = 'public'
      AND    tablename  = 'users'
      ) THEN

   CREATE TABLE users (
      base_address text,
      status text,
      twitter_username text,
      twitter_name text,
      twitter_url text,
      user_id integer
   );

   END IF;
END
\$$;
"

# create the database
createdb $DB_NAME -U $DB_USER

# execute the sql command to create the table
echo $SQL_CMD | PGPASSWORD=$DB_PASSWORD psql -U $DB_USER -d $DB_NAME

echo "database $DB_NAME & table users created successfully!"