FROM golang:1.14-buster

# Install required tools
RUN apt-get update -y && apt-get upgrade -y && apt-get install -y wait-for-it

# Add application user
RUN useradd -m -s /bin/bash app

# Set working directory
WORKDIR /app

# Copy source over
COPY --chown=app:app . .

# Switch to the app user
USER app
