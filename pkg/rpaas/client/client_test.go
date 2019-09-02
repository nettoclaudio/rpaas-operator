package client

import (
	"testing"
)

func Test_New(t *testing.T) {
	tests := []struct {
		name string
		cfg  Configuration
	}{
		{
			name: "when neither target or ",
			cfg:  Configuration{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}
