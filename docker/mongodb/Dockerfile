# Use the official MongoDB Docker image as a base
FROM mongo:5.0.19

# MongoDB data will be stored in /data/db, make mongodb the owner
RUN mkdir -p /data/db && chown -R mongodb:mongodb /data/db

# Run the MongoDB server as the mongodb user
USER mongodb

# MongoDB listens on the default port 27017
EXPOSE 27017

# Set the command to run with the config file
CMD ["mongod", "-f", "/etc/mongod.conf"]

# Make /data/db a volume, which makes it writable even with a read-only filesystem
VOLUME ["/data/db"]
