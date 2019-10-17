package providers

import (
	"context"
	"github.com/nbedos/citop/cache"
	"os"
	"testing"
	"time"
)

func TestCircleCIClientFetchRepositoryBuilds(t *testing.T) {
	token := os.Getenv("CIRCLECI_API_TOKEN")
	if token == "" {
		t.Fatal("environment variable CIRCLECI_API_TOKEN not set")
	}

	client := NewCircleCIClient(CircleCIURL, "", token, 100*time.Millisecond)
	ctx := context.Background()
	buildc := make(chan cache.Build)
	errc := make(chan error, 1)

	repository, err := client.Repository(ctx, "https://github.com/nbedos/citop")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		if err := client.fetchRepositoryBuilds(ctx, repository, 20, buildc); err != nil {
			errc <- err
		}
		close(buildc)
		close(errc)
	}()

	for range buildc {
	}

	if err := <-errc; err != nil {
		t.Fatal(err)
	}
}
