package main

import (
    "log"
	"os"
)

var Logger = log.New(os.Stdout, "logger: ", log.Lmicroseconds)