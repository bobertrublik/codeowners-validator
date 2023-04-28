package check_test

import (
	"context"
	"errors"
	"testing"

	"go.szostok.io/codeowners-validator/internal/api"
	"go.szostok.io/codeowners-validator/internal/check"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRespectingCanceledContext(t *testing.T) {
	must := func(checker api.Checker, err error) api.Checker {
		require.NoError(t, err)
		return checker
	}

	checkers := []api.Checker{
		check.NewDuplicatedPattern(),
		check.NewFileExist(),
		check.NewValidSyntax(),
		check.NewNotOwnedFile(check.NotOwnedFileConfig{}),
		must(check.NewValidOwner(check.ValidOwnerConfig{Repository: "org/repo"}, nil, true)),
	}

	for _, checker := range checkers {
		sut := checker
		t.Run(checker.Name(), func(t *testing.T) {
			// given: canceled context
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			// when
			out, err := sut.Check(ctx, LoadInput(FixtureValidCODEOWNERS))

			// then
			assert.True(t, errors.Is(err, context.Canceled))
			assert.Empty(t, out)
		})
	}
}
