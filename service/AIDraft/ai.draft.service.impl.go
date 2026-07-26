package aidraft

import (
	"context"
	"log"
	"sync"

	request "Microservice/data/request/AIDraft"
	response "Microservice/data/response/AIDraft"
	"Microservice/helper"
	"Microservice/pkg/sais"
)

// draftStyles are the fixed style options generated on every request, so the
// user can compare and pick rather than committing to one style up front.
// Each value is a tone SAIS already understands (see shifd-ai-service
// CLAUDE.md's draft-generation tone field) — adding a style here requires SAIS
// to support that tone value first.
var draftStyles = []string{"formal", "ringkas", "semi-formal"}

type AIDraftServiceImpl struct {
	client *sais.Client
}

func NewAIDraftServiceImpl(client *sais.Client) AIDraftService {
	return &AIDraftServiceImpl{client: client}
}

// variantResult is one goroutine's outcome, kept private so a failed style
// carries its error without needing a sentinel zero value.
type variantResult struct {
	variant response.DraftVariant
	err     error
}

func (s *AIDraftServiceImpl) Generate(ctx context.Context, bearerToken string, req request.GenerateDraftRequest) (*response.GenerateDraftVariantsResponse, *helper.ErrorModel) {
	if !s.client.Enabled() {
		return nil, &helper.ErrorModel{Code: 503, Message: "AI draft feature is not enabled"}
	}

	results := make([]variantResult, len(draftStyles))
	var wg sync.WaitGroup
	for i, style := range draftStyles {
		wg.Add(1)
		go func(i int, style string) {
			defer wg.Done()
			result, err := s.client.GenerateDraft(ctx, bearerToken, sais.DraftRequest{
				LetterType: req.LetterType,
				Purpose:    req.Purpose,
				KeyPoints:  req.KeyPoints,
				Recipient:  req.Recipient,
				Tone:       style,
			})
			if err != nil {
				results[i] = variantResult{err: err}
				return
			}
			results[i] = variantResult{variant: response.DraftVariant{
				Style:      style,
				Draft:      result.Draft,
				TokensUsed: result.TokensUsed,
			}}
		}(i, style)
	}
	wg.Wait()

	variants := make([]response.DraftVariant, 0, len(draftStyles))
	for _, r := range results {
		if r.err != nil {
			// One style failing (e.g. this provider's known intermittent
			// English-leakage or a timeout) must not sink the other styles —
			// the user still gets to choose among whatever came back.
			log.Printf("aidraft: style generation failed: %v", r.err)
			continue
		}
		variants = append(variants, r.variant)
	}

	if len(variants) == 0 {
		return nil, &helper.ErrorModel{Code: 502, Message: "AI service is unavailable. Please try again."}
	}

	return &response.GenerateDraftVariantsResponse{Variants: variants}, nil
}
