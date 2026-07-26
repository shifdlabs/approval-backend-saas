package aidraft

import (
	"context"

	request "Microservice/data/request/AIDraft"
	response "Microservice/data/response/AIDraft"
	"Microservice/helper"
)

// AIDraftService generates letter drafts by delegating to SAIS. No prompt or
// provider logic lives in Approval — SAIS owns both.
type AIDraftService interface {
	Generate(ctx context.Context, bearerToken string, req request.GenerateDraftRequest) (*response.GenerateDraftResponse, *helper.ErrorModel)
}
