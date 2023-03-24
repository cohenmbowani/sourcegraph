package repoupdater

import (
	"context"
	"testing"

	"github.com/sourcegraph/sourcegraph/internal/repos"
	proto "github.com/sourcegraph/sourcegraph/internal/repoupdater/v1"
	"github.com/sourcegraph/sourcegraph/internal/types"
)

func TestSyncExternalServiceRateLimited(t *testing.T) {
	// Setup test server and context
	server := &RepoUpdaterServiceServer{}
	ctx := context.Background()

	// Setup request and expected response
	req := &proto.SyncExternalServiceRequest{ExternalServiceId: 1}
	want := &proto.SyncExternalServiceResponse{}

	// Stub out externalServiceValidate to return a rate limit exceeded error
	server.externalServiceValidate = func(ctx context.Context, es *types.ExternalService, genericSrc repos.Source) error {
		return github.ErrRateLimited
	}

	// Call method under test
	got, err := server.SyncExternalService(ctx, req)
	if err == nil || err.Error() != "rate limit exceeded" {
		t.Errorf("SyncExternalService() error = %v, want rate limit exceeded error", err)
	}

	// Check that the response is as expected (empty)
	if got != nil {
		t.Errorf("SyncExternalService() = %v, want nil", got)
	}
}
