package courses

import (
	"net/http"
)

// Router faz o dispatch de requisições para os handlers apropriados.
func Router(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		ListHandler(w, r)
	} else if r.Method == http.MethodPost {
		CreateHandler(w, r)
	} else {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}
