#!/bin/bash

# Start Redis in the background
redis-server --daemonize yes

# Sleep for a while to ensure Redis has started
sleep 10

# Set up user and password
redis-cli ACL SETUSER thais on > /dev/null
redis-cli ACL SETUSER thais ~* +get +set +list +lpush > /dev/null
redis-cli ACL SETUSER thais >algodao11 < /dev/null

# Stop the background Redis server
redis-cli -a algodao11 shutdown

# Start Redis in the foreground
redis-server --appendonly yes --requirepass algodao11
