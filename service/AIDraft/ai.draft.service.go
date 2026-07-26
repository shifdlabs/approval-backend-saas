package aidraft

import (
	"context"

	request "Microservice/data/request/AIDraft"
	response "Microservice/data/response/AIDraft"
	"Microservice/helper"
)

// AIDraftService generates letter draft options by delegating to SAIS. No
// prompt or provider logic lives in Approval — SAIS owns both.
type AIDraftService interface {
	// Generate returns one draft per style in draftStyles, generated in
	// parallel. It only fails when every style failed; a partial result (some
	// styles missing) is still returned as a success.
	Generate(ctx context.Context, bearerToken string, req request.GenerateDraftRequest) (*response.GenerateDraftVariantsResponse, *helper.ErrorModel)
}
