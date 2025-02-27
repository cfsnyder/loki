package storage

import (
	"io"
	"time"

	"github.com/go-kit/log"
	"github.com/grafana/dskit/flagext"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/grafana/loki/pkg/storage/chunk"
	"github.com/grafana/loki/pkg/storage/chunk/cache"
	"github.com/grafana/loki/pkg/storage/chunk/gcp"
	"github.com/grafana/loki/pkg/storage/chunk/testutils"
	"github.com/grafana/loki/pkg/util/validation"
)

type fixture struct {
	fixture testutils.Fixture
}

func (f fixture) Name() string { return "caching-store" }
func (f fixture) Clients() (chunk.IndexClient, chunk.Client, chunk.TableClient, chunk.SchemaConfig, io.Closer, error) {
	limits, err := defaultLimits()
	if err != nil {
		return nil, nil, nil, chunk.SchemaConfig{}, nil, err
	}
	indexClient, chunkClient, tableClient, schemaConfig, closer, err := f.fixture.Clients()
	reg := prometheus.NewRegistry()
	logger := log.NewNopLogger()
	indexClient = newCachingIndexClient(indexClient, cache.NewFifoCache("index-fifo", cache.FifoCacheConfig{
		MaxSizeItems: 500,
		TTL:          5 * time.Minute,
	}, reg, logger), 5*time.Minute, limits, logger, false)
	return indexClient, chunkClient, tableClient, schemaConfig, closer, err
}

// Fixtures for unit testing the caching storage.
var Fixtures = []testutils.Fixture{
	fixture{gcp.Fixtures[0]},
}

func defaultLimits() (*validation.Overrides, error) {
	var defaults validation.Limits
	flagext.DefaultValues(&defaults)
	defaults.CardinalityLimit = 5
	return validation.NewOverrides(defaults, nil)
}
