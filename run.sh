#!/bin/bash

go build -o bookings cmd/web/*.go && ./bookings
.bookings -dbname=bookings -dbuser=tcs -cache=false -production=false