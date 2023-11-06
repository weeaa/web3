#!/bin/bash

os=$(uname -s)
arch=$(uname -m)
go_package_format=""
psql_package_manager=""
psql_install_command=""

DB_USER="admin"
DB_PASSWORD=$(generate_password 16)

generate_password() {
  local length="${12:-26}"
  tr -dc 'A-Za-z0-9' < /dev/urandom | head -c "$length"
  echo
}

# get the os
if [ "$os" == "Darwin" ]; then
    if [ "$arch" == "arm64" ]; then
        os="darwin-arm64"
    else
        os="darwin-amd64"
    fi
    go_package_format="pkg"
    psql_package_manager="brew"
    psql_install_command="$psql_package_manager install postgresql"
elif [ "$os" == "Linux" ]; then
    os="linux-amd64"
    go_package_format="tar.gz"
else
    echo "unsupported OS: $os"
    exit 1
fi

# verify if go is installed
if ! command -v go &> /dev/null; then
    echo "🦦 | Installing Golang..."
    curl -O "https://go.dev/dl/go1.21.3.$os.$go_package_format"
    tar -C /usr/local -xzf "go1.21.3.$os.$go_package_format"
    export PATH=$PATH:/usr/local/go/bin
fi

# verify if psql is installed
if ! command -v psql &> /dev/null; then
    echo "🐘 | Installing pSQL..."
    if [ "$os" == "linux-amd64" ]; then
      sudo apt update
    fi
    $psql_install_command
fi

# prompt the user for database inputs
read -p "do you want to create a new database? (y/n): " CREATE_DB
echo "💿 | Creating a new Database with [user:$DB_USER | password:$DB_PASSWORD]..."

# if the user chooses to create a database
if [ "$CREATE_DB" == "y" ]; then
  read -p "DB_USER: " DB_USER
  read -s -p "DB_PASSWORD: " DB_PASSWORD
  echo ""
  createdb $DB_NAME -U $DB_USER
fi

echo "💿 | Creating tables..."


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
createdb "$DB_NAME" -U $DB_USER

# execute the sql command to create the table
# shellcheck disable=SC2090
echo "$SQL_CMD" | PGPASSWORD=$DB_PASSWORD psql -U $DB_USER -d "$DB_NAME"

echo "database $DB_NAME & table users created successfully!"