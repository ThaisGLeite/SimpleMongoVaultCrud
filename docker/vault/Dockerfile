# Pull base image
FROM vault:1.13.3

# Add a non-root user and group 
RUN addgroup devvault && \
    adduser -S -G devvault devuser

# Set environment variables
ENV VAULT_ADDR=http://127.0.0.1:8200 \
    VAULT_CONFIG_DIR=/vault/config

# Create a directory for Vault configuration
RUN mkdir -p $VAULT_CONFIG_DIR

# Copy the config.json file into the image
COPY ./config/config.json $VAULT_CONFIG_DIR/config.json

# Change to non-root user
USER devuser

# Expose the vault port
EXPOSE 8200

# Start vault
CMD vault server -config=$VAULT_CONFIG_DIR/config.json
