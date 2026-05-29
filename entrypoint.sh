#!/bin/sh
echo "127.0.0.1 test-acme.example.com" >> /etc/hosts
exec "$@"
