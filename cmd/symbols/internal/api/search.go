package api

import (
	"context"
	"strings"
	"time"

	"github.com/cockroachdb/errors"

	"github.com/sourcegraph/sourcegraph/cmd/symbols/internal/database/store"
	"github.com/sourcegraph/sourcegraph/cmd/symbols/internal/types"
	obsv "github.com/sourcegraph/sourcegraph/internal/observation"
	"github.com/sourcegraph/sourcegraph/internal/search/result"
)

const searchTimeout = 60 * time.Second

func (h *apiHandler) handleSearchInternal(ctx context.Context, args types.SearchArgs) (_ *result.Symbols, err error) {
	ctx, traceLog, endObservation := h.operations.search.WithAndLogger(ctx, &err, obsv.Args{LogFields: []obsv.Field{
		obsv.String("repo", string(args.Repo)),
		obsv.String("commitID", string(args.CommitID)),
		obsv.String("query", args.Query),
		obsv.Bool("isRegExp", args.IsRegExp),
		obsv.Bool("isCaseSensitive", args.IsCaseSensitive),
		obsv.Int("numIncludePatterns", len(args.IncludePatterns)),
		obsv.String("includePatterns", strings.Join(args.IncludePatterns, ":")),
		obsv.String("excludePattern", args.ExcludePattern),
		obsv.Int("first", args.First),
	}})
	defer endObservation(1, obsv.Args{})

	ctx, cancel := context.WithTimeout(ctx, searchTimeout)
	defer cancel()

	dbFile, err := h.cachedDatabaseWriter.GetOrCreateDatabaseFile(ctx, args)
	if err != nil {
		return nil, errors.Wrap(err, "databaseWriter.GetOrCreateDatabaseFile")
	}
	traceLog(obsv.String("dbFile", dbFile))

	var results result.Symbols
	err = store.WithSQLiteStore(dbFile, func(db store.Store) (err error) {
		if results, err = db.Search(ctx, args); err != nil {
			return errors.Wrap(err, "store.Search")
		}

		return nil
	})

	return &results, err
}
