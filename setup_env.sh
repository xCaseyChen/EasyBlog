#!/bin/bash

# PostgreSQL service configurations
POSTGRES_USER="bloguser"
POSTGRES_PASSWORD=$(openssl rand -hex 16)
POSTGRES_DB="blogdb"
POSTGRES_PORT=5432

POSTGRES_HEALTHCHECK_RETRIES=10
POSTGRES_HEALTHCHECK_INTERVAL="3s"

# pgAdmin service configurations
PGADMIN_EMAIL="pgadmin@example.com"
PGADMIN_PASSWORD=$(openssl rand -hex 16)
PGADMIN_PORT=8080

# Server configurations
SERVER_PORT=80
JWT_SECRET=$(openssl rand -hex 32)

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
check_port_availability ${POSTGRES_PORT}
check_port_availability ${PGADMIN_PORT}
check_port_availability ${SERVER_PORT}

# Generate the .env configuration file
if test -e ".env.example"
then
    echo ".env is already exists"
    exit 1
fi

cat > .env.example <<EOF
# PostgreSQL username
POSTGRES_USER=${POSTGRES_USER}

# PostgreSQL password (randomly generated, 16 bytes characters)
POSTGRES_PASSWORD=${POSTGRES_PASSWORD}

# PostgreSQL database name
POSTGRES_DB=${POSTGRES_DB}

# PostgreSQL host port (used for host access to the container)
POSTGRES_PORT=${POSTGRES_PORT}

# PostgreSQL health check retry times
POSTGRES_HEALTHCHECK_RETRIES=${POSTGRES_HEALTHCHECK_RETRIES}

# PostgreSQL health check retry interval
POSTGRES_HEALTHCHECK_INTERVAL=${POSTGRES_HEALTHCHECK_INTERVAL}

# pgAdmin default login email (used to log in to pgAdmin)
PGADMIN_EMAIL=${PGADMIN_EMAIL}

# pgAdmin default login password (randomly generated, 16 bytes characters)
PGADMIN_PASSWORD=${PGADMIN_PASSWORD}

# pgAdmin host port (used for host access to the container)
PGADMIN_PORT=${PGADMIN_PORT}

# Server host port (used for host access to the server)
SERVER_PORT=${SERVER_PORT}

# Server JWT secret
JWT_SECRET=${JWT_SECRET}
EOF

echo "Successfully generated .env"