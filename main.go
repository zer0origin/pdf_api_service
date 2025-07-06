package main

import (
	v1 "pdf_service_api/v1"
)

func main() {
	router := v1.SetupRouter()
	_ = router.Run() // listen and serve on 0.0.0.0:8080
}
