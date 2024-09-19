package server

import (
	"5e-shop/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	r.HandlerFunc(http.MethodGet, "/get-campaign-shops", s.getCampaignShopsHandler)
	r.HandlerFunc(http.MethodGet, "/get-current-shop", s.getCurrentShopHandler)

	r.HandlerFunc(http.MethodGet, "/health", s.healthHandler)

	return r
}

func (s *Server) getCurrentShopHandler(w http.ResponseWriter, r *http.Request) {
	marshalledData, err := extractJsonBody(r)
	if err != nil {
		utils.HandleResponseError(w, fmt.Sprintf("error extracting json body from request. Err: %v", err), 504)
		return
	}

	campaignIdString, ok := marshalledData["campaignId"]
	if !ok {
		utils.HandleResponseError(w, "could not find campaignId value in request body", 404)
		return
	}
	campaignId, err := primitive.ObjectIDFromHex(campaignIdString)
	if err != nil {
		utils.HandleResponseError(w, fmt.Sprintf("error extracting ID from message. Err: %v", err), 504)
		return
	}

	campaign, err := s.db.GetCampaign(context.Background(), campaignId)
	if err != nil {
		utils.HandleResponseError(w, fmt.Sprintf("error getting campaign from DB. Err: %v", err), 504)
		return
	}
	currShop, err := s.db.GetShop(context.Background(), campaign.ActiveShop)
	if err != nil {
		utils.HandleResponseError(w, fmt.Sprintf("error getting shop from DB. Err: %v", err), 504)
		return
	}

	jsonResp, err := json.Marshal(currShop)
	if err != nil {
		log.Printf("error handling JSON marshal. Err: %v", err)
		w.WriteHeader(504)
		return
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) getCampaignShopsHandler(w http.ResponseWriter, r *http.Request) {
	marshalledData, err := extractJsonBody(r)
	if err != nil {
		utils.HandleResponseError(w, fmt.Sprintf("error extracting json body from request. Err: %v", err), 504)
		return
	}

	campaignIdString, ok := marshalledData["campaignId"]
	if !ok {
		utils.HandleResponseError(w, "could not find campaignId value in request body", 404)
		return
	}
	campaignId, err := primitive.ObjectIDFromHex(campaignIdString)
	if err != nil {
		utils.HandleResponseError(w, fmt.Sprintf("error extracting ID from message. Err: %v", err), 504)
		return
	}

	campaignShops, err := s.db.GetCampaignShops(r.Context(), campaignId)
	if err != nil {
		utils.HandleResponseError(w, fmt.Sprintf("error getting campaign shops from DB. Err: %v", err), 504)
		return
	}

	jsonResp, err := json.Marshal(campaignShops)
	if err != nil {
		log.Printf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(s.db.Health())

	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
		w.WriteHeader(504)
		return
	}

	_, _ = w.Write(jsonResp)
}
