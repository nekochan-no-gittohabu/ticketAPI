package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"ticketAPI/commonerrors"
	"ticketAPI/ticket"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	router chi.Mux
	svc    *ticket.Service
}

func New(s *ticket.Service) *Server {
	server := &Server{
		router: *chi.NewRouter(),
		svc:    s,
	}
	server.routes()
	return server
}

func (s Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

func (s *Server) routes() {
	s.router.Get("/ticket/{id}", s.getTicket)

	s.router.Route("/ticket_options", func(r chi.Router) {
		r.Post("/", s.createTicket)
		r.Post("/{id}/purchases", s.purchaseTicket)
	})
}

func (s *Server) getTicket(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(req, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t, err := s.svc.GetTicket(req.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(t); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) createTicket(w http.ResponseWriter, req *http.Request) {
	t := ticket.Ticket{}
	err := json.NewDecoder(req.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	t, err = s.svc.CreateTicketOption(req.Context(), t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(t); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) purchaseTicket(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(req, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	r := ticket.PurchaseRequest{}
	err = json.NewDecoder(req.Body).Decode(&r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = s.svc.PurchaseTicket(req.Context(), id, r)
	if err != nil {
		if errors.Is(errors.New(commonerrors.ErrTooManyPurchases), err) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
}
