package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type logService struct {
	CommitLog *CommitLog
}

type WriteRequest struct {
	Record Record `json:"record"`
}

type WriteResponse struct {
	Offset uint64 `json:"offset"`
}

type ReadRequest struct {
	Offset uint64 `json:"offset"`
}

type ReadResponse struct {
	Record Record `json:"record"`
}

func NewHTTPServer(addr string) *http.Server {
	svc := &logService{
		CommitLog: NewCommitLog(),
	}

	r := mux.NewRouter()

	r.HandleFunc("/", svc.handleWrite).Methods("POST")
	r.HandleFunc("/{offset}", svc.handleRead).Methods("GET")

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

func (svc *logService) handleWrite(w http.ResponseWriter, r *http.Request) {
	req := &WriteRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to parse write request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	off, err := svc.CommitLog.Append(req.Record)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to append record to log: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	resp := &WriteResponse{
		Offset: off,
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to encode response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (svc *logService) handleRead(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	off, err := strconv.ParseUint(vars["offset"], 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to parse given offset: %s", err.Error()), http.StatusBadRequest)
		return
	}

	rec, err := svc.CommitLog.Read(off)
	if err != nil {
		if err == ErrOffsetNotFound {
			http.Error(w, fmt.Sprintf("unable to locate record: %s", err.Error()), http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("unable to read record from log: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	resp := &ReadResponse{
		Record: rec,
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to encode response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}
