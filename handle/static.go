package handle

import (
	"net/http"
)

func StaticFilesHandler() {
	fs := http.FileServer(http.Dir("../static"))
	http.Handle("/", fs)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/createChatroom", CreateChatroomHandler)
	http.HandleFunc("/sendMessage", SendMessageHandler)
	http.HandleFunc("/retrieveMessages", RetrieveMessagesHandler)
}
