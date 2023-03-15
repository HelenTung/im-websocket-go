package main

import "main/router"

func main() {
	e := router.Router()
	e.Run(":8080")
}
