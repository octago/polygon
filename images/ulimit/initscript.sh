#!/bin/bash

app &> /var/log/app.log || echo "ulimitabuser failed"

exec "$@"
