#!/bin/bash
set -e

cat >> "$PGDATA/pg_hba.conf" <<EOF
# Replication entries
host replication postgres all md5
host replication postgres 0.0.0.0/0 md5
host replication postgres 172.19.0.0/16 md5
EOF