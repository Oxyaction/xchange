#!/bin/bash
set -e

# create database account if not exists
psql -U postgres -tc "SELECT 1 FROM pg_database WHERE datname = 'account'" | grep -q 1 || psql -U postgres -c "CREATE DATABASE account"
psql -U postgres -tc "SELECT 1 FROM pg_database WHERE datname = 'account_test'" | grep -q 1 || psql -U postgres -c "CREATE DATABASE account_test"
psql -U postgres -tc "SELECT 1 FROM pg_database WHERE datname = 'order'" | grep -q 1 || psql -U postgres -c "CREATE DATABASE \"order\""
psql -U postgres -tc "SELECT 1 FROM pg_database WHERE datname = 'order_test'" | grep -q 1 || psql -U postgres -c "CREATE DATABASE order_test"

psql -v ON_ERROR_STOP=1 --username postgres --dbname account <<-EOSQL
    CREATE EXTENSION "uuid-ossp";
EOSQL

psql -v ON_ERROR_STOP=1 --username postgres --dbname account_test <<-EOSQL
    CREATE EXTENSION "uuid-ossp";
EOSQL

psql -v ON_ERROR_STOP=1 --username postgres --dbname order <<-EOSQL
    CREATE EXTENSION "uuid-ossp";
EOSQL

psql -v ON_ERROR_STOP=1 --username postgres --dbname order_test <<-EOSQL
    CREATE EXTENSION "uuid-ossp";
EOSQL