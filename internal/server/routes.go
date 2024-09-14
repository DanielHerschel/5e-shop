package server

import (
	"5e-shop/internal/domain"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func extractJsonBody(r *http.Request) (map[string]string, error) {
	rawData, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read request body")
	}

	var marshalledData map[string]string

	json.Unmarshal(rawData, &marshalledData)
	return marshalledData, nil
}

func (s *Server) RegisterRoutes() http.Handler {
	r := httprouter.New()
	r.HandlerFunc(http.MethodGet, "/get-shops", s.getAllShopsHandler)

	r.HandlerFunc(http.MethodGet, "/get-current-shop", s.getCurrentShopHandler)

	r.HandlerFunc(http.MethodGet, "/health", s.healthHandler)

	return r
}

func (s *Server) getCurrentShopHandler(w http.ResponseWriter, r *http.Request) {
	marshalledData, err := extractJsonBody(r)
	if err != nil {
		log.Printf("error extracting json body from request. Err: %v", err)
		w.WriteHeader(504)
		return
	}

	campaignName, ok := marshalledData["campaign"]
	if !ok {
		log.Printf("could not find campaign value in request body")
		w.WriteHeader(404)
		return
	}

	log.Printf("campaignName: %s", campaignName)

	// TODO: search db
	shop := domain.Shop{
		Name: "The Small Stand",
		Items: []domain.Item{
			{Name: "Shorsword", Rarity: domain.Common, Gold: 5},
		},
	}
	resp := domain.Campaign{Name: "Curse of Strahd", Shops: []domain.Shop{shop}, CurrentShop: shop}
	// --------------------------------------------------------------------------------------------

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) getAllShopsHandler(w http.ResponseWriter, r *http.Request) {
	marshalledData, err := extractJsonBody(r)
	if err != nil {
		log.Printf("error extracting json body from request. Err: %v", err)
		w.WriteHeader(504)
		return
	}

	campaignName, ok := marshalledData["campaign"]
	if !ok {
		log.Printf("could not find campaign value in request body")
		w.WriteHeader(404)
		return
	}

	log.Printf("campaignName: %s", campaignName)

	// TODO: search db
	resp := []domain.Shop{
		{Name: "The Small Stand", Items: []domain.Item{
			{Name: "Shortsword", Rarity: domain.Common, Gold: 5},
		}},
	}
	// -----------------------------------------------------

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(s.db.Health())

	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
