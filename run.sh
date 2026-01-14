#!/bin/bash

go build -o bookings cmd/web/*.go 
./bookings -dbname=bookings -dbuser=tcs -dbhost=localhost -cache=false -production=false