package httpapi

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/inconshreveable/log15"

	"github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/stores/dbstore"
	obsv "github.com/sourcegraph/sourcegraph/internal/observation"
)

// handleEnqueueSinglePayload handles a non-multipart upload. This creates an upload record
// with state 'queued', proxies the data to the bundle manager, and returns the generated ID.
func (h *UploadHandler) handleEnqueueSinglePayload(ctx context.Context, uploadState uploadState, body io.Reader) (_ interface{}, statusCode int, err error) {
	ctx, traceLog, endObservation := h.operations.handleEnqueueSinglePayload.WithAndLogger(ctx, &err, obsv.Args{})
	defer func() {
		endObservation(1, obsv.Args{LogFields: []obsv.Field{
			obsv.Int("statusCode", statusCode),
		}})
	}()

	tx, err := h.dbStore.Transact(ctx)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer func() { err = tx.Done(err) }()

	id, err := tx.InsertUpload(ctx, dbstore.Upload{
		Commit:            uploadState.commit,
		Root:              uploadState.root,
		RepositoryID:      uploadState.repositoryID,
		Indexer:           uploadState.indexer,
		AssociatedIndexID: &uploadState.associatedIndexID,
		State:             "uploading",
		NumParts:          1,
		UploadedParts:     []int{0},
	})
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	traceLog(obsv.Int("uploadID", id))

	size, err := h.uploadStore.Upload(ctx, fmt.Sprintf("upload-%d.lsif.gz", id), body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	traceLog(obsv.Int("gzippedUploadSize", int(size)))

	if err := tx.MarkQueued(ctx, id, &size); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	log15.Info(
		"codeintel.httpapi: enqueued upload",
		"id", id,
		"repository_id", uploadState.repositoryID,
		"commit", uploadState.commit,
	)

	// older versions of src-cli expect a string
	return struct {
		ID string `json:"id"`
	}{ID: strconv.Itoa(id)}, 0, nil
}
