#!/bin/bash

# PostgreSQL service configurations
POSTGRES_USER="bloguser"
POSTGRES_PASSWORD=$(openssl rand -hex 16)
POSTGRES_DB="blogdb"
POSTGRES_HOST_PORT=5432

# pgAdmin service configurations
PGADMIN_EMAIL="pgadmin@example.com"
PGADMIN_PASSWORD=$(openssl rand -hex 16)
PGADMIN_PORT=8080

#####################################

# Function to check port availability
function check_port_availability() {
    local port=$1
    if lsof -i :$port >/dev/null 2>&1
    then
        echo "Error: port $port is already in use."
        exit 1
    fi
}

# check ports availability
check_port_availability ${POSTGRES_HOST_PORT}
check_port_availability ${PGADMIN_PORT}

# Generate the .env configuration file
cat > .env <<EOF
# PostgreSQL username
POSTGRES_USER=${POSTGRES_USER}

# PostgreSQL password (randomly generated, 16 bytes characters)
POSTGRES_PASSWORD=${POSTGRES_PASSWORD}

# PostgreSQL database name
POSTGRES_DB=${POSTGRES_DB}

# PostgreSQL host port (used for host access to the container)
POSTGRES_HOST_PORT=${POSTGRES_HOST_PORT}

# pgAdmin default login email (used to log in to pgAdmin)
PGADMIN_EMAIL=${PGADMIN_EMAIL}

# pgAdmin default login password (randomly generated, 16 bytes characters)
PGADMIN_PASSWORD=${PGADMIN_PASSWORD}

# pgAdmin host port (used for host access to the container)
PGADMIN_PORT=${PGADMIN_PORT}
EOF

echo "Successfully generated .env"