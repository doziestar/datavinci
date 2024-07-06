#!/bin/bash

# Function to print usage information
print_usage() {
    echo "Usage: $0 <action> <db_type> <db_name> [migration_name]"
    echo "Actions:"
    echo "  diff    - Generate a new migration"
    echo "  apply   - Apply migrations"
    echo "DB Types:"
    echo "  sqlite"
    echo "  postgres"
    echo "  mariadb"
    echo "Example:"
    echo "  $0 diff sqlite mydb new_migration"
    echo "  $0 apply postgres mydb"
}

# Check if required arguments are provided
if [ $# -lt 3 ]; then
    print_usage
    exit 1
fi

ACTION=$1
DB_TYPE=$2
DB_NAME=$3
MIGRATION_NAME=$4

# Set the migrations directory
MIGRATIONS_DIR="file://ent/migrate/migrations"

# Function to get the appropriate --dev-url for diff action
get_dev_url() {
    case $DB_TYPE in
        sqlite)
            echo "sqlite://file?mode=memory&_fk=1"
            ;;
        postgres)
            echo "docker://postgres/15/$DB_NAME?search_path=public"
            ;;
        mariadb)
            echo "docker://mariadb/latest/$DB_NAME"
            ;;
        *)
            echo "Unsupported database type: $DB_TYPE"
            exit 1
            ;;
    esac
}

# Function to get the appropriate --url for apply action
get_url() {
    case $DB_TYPE in
        sqlite)
            echo "sqlite://$DB_NAME.db?_fk=1"
            ;;
        postgres)
            echo "postgres://postgres:pass@localhost:5432/$DB_NAME?search_path=public&sslmode=disable"
            ;;
        mariadb)
            echo "maria://root:pass@localhost:3306/$DB_NAME"
            ;;
        *)
            echo "Unsupported database type: $DB_TYPE"
            exit 1
            ;;
    esac
}

# Perform the requested action
case $ACTION in
    diff)
        if [ -z "$MIGRATION_NAME" ]; then
            echo "Error: Migration name is required for 'diff' action"
            print_usage
            exit 1
        fi
        DEV_URL=$(get_dev_url)
        atlas migrate diff $MIGRATION_NAME \
            --dir "$MIGRATIONS_DIR" \
            --to "ent://ent/schema" \
            --dev-url "$DEV_URL"
        ;;
    apply)
        URL=$(get_url)
        atlas migrate apply \
            --dir "$MIGRATIONS_DIR" \
            --url "$URL"
        ;;
    *)
        echo "Error: Invalid action. Use 'diff' or 'apply'"
        print_usage
        exit 1
        ;;
esac

echo "Operation completed successfully."
