package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func Routes(db *Store) chi.Router {
	svc := &restService{
		db: db,
	}

	router := chi.NewRouter()

	router.Post("/save", svc.saveFoto)
	router.Get("/raw_foto/{id}", svc.getRawFoto)
	router.Get("/preview/{id}", svc.getPreview)
	router.Get("/fotos", svc.getFotos)
	router.Delete("/foto/{id}", svc.removeFoto)

	return router
}

type restService struct {
	db *Store
}

// POST /save
func (svc *restService) saveFoto(w http.ResponseWriter, r *http.Request) {
	rawFoto, err := ioutil.ReadAll(r.Body)
	if err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}

	const previewSize = 200
	preview, err := Preview(rawFoto, previewSize)
	if err != nil {
		preview = []byte{}
	}

	id, err := svc.db.SaveFoto(rawFoto, preview)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}

	render.JSON(w, r, id)
}

// GET /raw_foto/{id}
func (svc *restService) getRawFoto(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}

	raw_foto, err := svc.db.GetRawFoto(id)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}
	if len(raw_foto) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(raw_foto)))
	io.Copy(w, bytes.NewReader(raw_foto))
}

// GET /preview/{id}
func (svc *restService) getPreview(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}

	preview, err := svc.db.GetPreview(id)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}
	if len(preview) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(preview)))
	io.Copy(w, bytes.NewReader(preview))
}

// GET /fotos
func (svc *restService) getFotos(w http.ResponseWriter, r *http.Request) {
	ids, err := svc.db.GetFotos()
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}

	render.JSON(w, r, ids)
}

// DELETE /foto/{id}
func (svc *restService) removeFoto(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}

	err = svc.db.RemoveFoto(id)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// `statusCode` - HTTP status code
func SendError(w http.ResponseWriter, statusCode int, err error) {
	log.Printf("[ERROR] %v", err)

	w.WriteHeader(statusCode)
}
