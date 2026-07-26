package aidraft

import (
	"context"
	"log"

	request "Microservice/data/request/AIDraft"
	response "Microservice/data/response/AIDraft"
	"Microservice/helper"
	"Microservice/pkg/sais"
)

type AIDraftServiceImpl struct {
	client *sais.Client
}

func NewAIDraftServiceImpl(client *sais.Client) AIDraftService {
	return &AIDraftServiceImpl{client: client}
}

func (s *AIDraftServiceImpl) Generate(ctx context.Context, bearerToken string, req request.GenerateDraftRequest) (*response.GenerateDraftResponse, *helper.ErrorModel) {
	if !s.client.Enabled() {
		return nil, &helper.ErrorModel{Code: 503, Message: "AI draft feature is not enabled"}
	}

	result, err := s.client.GenerateDraft(ctx, bearerToken, sais.DraftRequest{
		LetterType: req.LetterType,
		Purpose:    req.Purpose,
		KeyPoints:  req.KeyPoints,
		Recipient:  req.Recipient,
		Tone:       req.Tone,
	})
	if err != nil {
		// The SAIS error can name the provider and carry response fragments —
		// keep it server-side and hand the caller a generic message.
		log.Printf("aidraft: draft generation failed: %v", err)
		return nil, &helper.ErrorModel{Code: 502, Message: "AI service is unavailable. Please try again."}
	}

	return &response.GenerateDraftResponse{
		Draft:      result.Draft,
		TokensUsed: result.TokensUsed,
	}, nil
}
