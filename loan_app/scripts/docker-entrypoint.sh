#!/usr/bin/env bash
set -e

echo "Applying migrations"

cd /repo/migrations

DUMP_SCHEMA_AFTER_MIGRATION=$DUMP_SCHEMA_AFTER_MIGRATION bundle exec rake db:migrate
if [[ $? != 0 ]]; then
  echo "Failed to apply migrations. Exiting."
  exit 1
fi
echo "Applied migrations successfully"

cd -

exec "$@"