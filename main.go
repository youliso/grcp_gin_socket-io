package main

import (
	"grpc/into"
	"grpc/utils/db"
)

func main() {
	db.Init()
	into.Run()
}
