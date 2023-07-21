package handle

import "net/http"

func Handle() {
	fs := http.FileServer(http.Dir("../static"))
	http.Handle("/", fs)
}
