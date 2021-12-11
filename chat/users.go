package chat

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

type usersResource struct {
	db     *gorm.DB
	logger *log.Logger
}

func (rs usersResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", rs.Create)

	return r
}

func (rs usersResource) Create(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Name     string
		Username string
	}

	reqBody := requestBody{}

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Couldn't decode request body", 400)
		return
	}

	if reqBody.Name == "" || reqBody.Username == "" {
		http.Error(w, "Username and Name can't be blank", 400)
		return
	}

	user := &User{Name: reqBody.Name, Username: reqBody.Username}
	result := rs.db.Create(&user)
	if result.Error != nil {
		log.Printf("error while creating user: %v", result.Error)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	marshaled, err := json.Marshal(map[string]*User{"user": user})
	if err != nil {
		log.Printf("error while marshaling user: %v", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Write(marshaled)
	w.Write([]byte("\r\n"))
	return
}
