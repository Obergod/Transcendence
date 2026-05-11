#!/usr/bin/env bash
set -e

mkdir -p /run/mysqld /var/lib/mysql
chown -R mysql:mysql /run/mysqld /var/lib/mysql

if [ ! -d /var/lib/mysql/mysql ]; then
  echo "Initializing MariaDB data directory..."
  mariadb-install-db --user=mysql --datadir=/var/lib/mysql
fi

echo "Starting temporary MariaDB server..."
mysqld_safe --datadir=/var/lib/mysql --user=mysql --skip-networking &
TMP_PID=$!

echo "Waiting for MariaDB to become ready..."
until mysqladmin ping -uroot --silent; do
  sleep 1
done

if ! mysql -uroot -p"${SQL_ROOT_PASSWORD:-}" -e "SELECT 1" >/dev/null 2>&1; then
  echo "Setting root password..."
  mysql -uroot -e "ALTER USER 'root'@'localhost' IDENTIFIED BY '${SQL_ROOT_PASSWORD}'; FLUSH PRIVILEGES;"
fi

if ! mysql -uroot -p"${SQL_ROOT_PASSWORD}" -e "USE ${SQL_DATABASE}" 2>/dev/null; then
  echo "Running initial application SQL setup..."
  mysql -uroot -p"${SQL_ROOT_PASSWORD}" <<-EOSQL
    CREATE DATABASE IF NOT EXISTS \`${SQL_DATABASE}\`;
    CREATE USER IF NOT EXISTS \`${SQL_USER}\`@'%' IDENTIFIED BY '${SQL_PASSWORD}';
    GRANT ALL PRIVILEGES ON \`${SQL_DATABASE}\`.* TO \`${SQL_USER}\`@'%';
    FLUSH PRIVILEGES;
EOSQL
fi

echo "Stopping temporary MariaDB server..."
mysqladmin -uroot -p"${SQL_ROOT_PASSWORD}" shutdown
wait "$TMP_PID" || true

echo "Starting MariaDB in foreground..."
exec mysqld_safe --datadir=/var/lib/mysql --user=mysql