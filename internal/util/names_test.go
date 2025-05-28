package util_test

import (
	"testing"

	"deedles.dev/gir/internal/util"
	"github.com/stretchr/testify/require"
)

func TestToCamelCase(t *testing.T) {
	tests := []struct{ before, after string }{
		{"test_to_camel_case", "TestToCamelCase"},
		{"gi_repository_get_n_infos", "GiRepositoryGetNInfos"},
	}

	for _, test := range tests {
		t.Run(test.before, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, test.after, util.ToCamelCase(test.before))
		})
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct{ before, after string }{
		{"TestToCamelCase", "test_to_camel_case"},
		{"GiRepositoryGetNInfos", "gi_repository_get_n_infos"},
	}

	for _, test := range tests {
		t.Run(test.before, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, test.after, util.ToSnakeCase(test.before))
		})
	}
}
