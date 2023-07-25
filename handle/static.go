package handle

import (
	"log"
	"net/http"
	"os"
)

const staticDirPath = "./static"

func StaticFilesHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if static directory exists
		if _, err := os.Stat(staticDirPath); os.IsNotExist(err) {
			log.Printf("Static directory does not exist: %s\n", staticDirPath)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		fs := http.FileServer(http.Dir(staticDirPath))

		// Create a new request scope
		http.StripPrefix("/", fs).ServeHTTP(w, r)
	})
}
