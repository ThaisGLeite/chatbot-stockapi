package handle

import (
	"net/http"
)

func StaticFilesHandler() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
}
