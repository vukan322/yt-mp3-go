package handler

import "net/http"

func (h *AppHandler) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	staticFs := http.FileServer(http.Dir("./web/static"))
	mux.Handle(h.BasePath+"/static/", http.StripPrefix(h.BasePath+"/static/", staticFs))

	mux.HandleFunc(h.BasePath+"/", h.HandleIndex)
	mux.HandleFunc(h.BasePath+"/info", h.HandleInfo)
	mux.HandleFunc(h.BasePath+"/download", h.HandleDownload)
	mux.HandleFunc(h.BasePath+"/downloads/", h.HandleServeDownload)
	mux.HandleFunc(h.BasePath+"/events", h.HandleStatusEvents)

	return mux
}
