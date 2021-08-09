package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteValidate(t *testing.T) {
	tests := []struct {
		name    string
		delete  CacheInvalidation
		wantErr error
	}{
		{
			name: "Validate empty delete",
			delete: CacheInvalidation{
				Paths: []CachePath{},
				Tags:  []string{},
			},
			wantErr: ErrInvalidDeleteCache,
		},
		{
			name: "Validate delete with paths and tags",
			delete: CacheInvalidation{
				Paths: []CachePath{
					{
						Service: "catalog",
						Path:    "/products",
					},
				},
				Tags: []string{"catalog"},
			},
			wantErr: ErrInvalidDeleteCache,
		},
		{
			name: "Validate valid delete",
			delete: CacheInvalidation{
				Tags: []string{"catalog"},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.delete.Validate()

			assert.Equal(t, err, tt.wantErr)
		})
	}
}

func TestErrCachePathFormat(t *testing.T) {
	err := &ErrCachePathNotFound{
		CachePath: CachePath{
			Service: "catalog",
			Path:    "/products",
		},
	}

	assert.EqualError(t, err, "Cache path not found: catalog/products")
}
