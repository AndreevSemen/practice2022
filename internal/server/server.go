package server

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"summer/practice2022/internal/config"
	"summer/practice2022/internal/logic"
	"summer/practice2022/internal/structures"
	st "summer/practice2022/internal/structures"
)

type Logic interface {
	GetTokensByLoginAndPassword(login, password string) (st.Tokens, error)
	GetTokensByRefreshToken(refresh string) (st.Tokens, error)
}

type Server struct {
	cfg config.Config
	l   Logic
	mux *http.ServeMux
}

func NewServer(cfg config.Config, l Logic) *Server {
	s := &Server{
		cfg: cfg,
		l:   l,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/auth", s.handleAuth)
	mux.HandleFunc("/refresh", s.handleRefresh)

	return s
}

func (s *Server) Serve() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.ServerPort))
	if err != nil {
		return err
	}

	if s.cfg.HTTPSMode {
		crt := s.cfg.HTTPSCrtFile
		key := s.cfg.HTTPSKeyFile
		return http.ServeTLS(lis, s.mux, crt, key)
	} else {
		return http.Serve(lis, s.mux)
	}
}

func (s *Server) handleAuth(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	type response struct {
		Tokens structures.Tokens `json:"tokens"`
	}

	dec := json.NewDecoder(r.Body)

	var req request
	if err := dec.Decode(&req); err != nil {
		fmt.Fprintf(w, `{"error":%s}`, strconv.Quote(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokens, err := s.l.GetTokensByLoginAndPassword(req.Login, req.Password)
	if err == logic.ErrPasswordIsIncorrect {
		fmt.Fprintf(w, `{"error":%s}`, strconv.Quote(err.Error()))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		fmt.Fprintf(w, `{"error":%s}`, strconv.Quote(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := response{
		Tokens: tokens,
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		fmt.Fprintf(w, `{"error":%s}`, strconv.Quote(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleRefresh(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Tokens structures.Tokens `json:"tokens"`
	}

	type response struct {
		Tokens structures.Tokens `json:"tokens"`
	}

	dec := json.NewDecoder(r.Body)

	var req request
	if err := dec.Decode(&req); err != nil {
		fmt.Fprintf(w, `{"error":%s}`, strconv.Quote(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokens, err := s.l.GetTokensByRefreshToken(req.Tokens.Refresh)
	if err == logic.ErrRefreshTokenIsAlreadyUsed {
		fmt.Fprintf(w, `{"error":%s}`, strconv.Quote(err.Error()))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		fmt.Fprintf(w, `{"error":%s}`, strconv.Quote(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := response{
		Tokens: tokens,
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		fmt.Fprintf(w, `{"error":%s}`, strconv.Quote(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
