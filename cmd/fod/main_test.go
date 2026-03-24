package main

import (
	"testing"

	fod "github.com/ykanda/fod/pkg"
)

func TestValidateDialogResult(t *testing.T) {
	tests := []struct {
		name    string
		result  fod.ResultCode
		wantErr bool
	}{
		{
			name:    "ok",
			result:  fod.ResultOK,
			wantErr: false,
		},
		{
			name:    "cancel",
			result:  fod.ResultCancel,
			wantErr: false,
		},
		{
			name:    "none",
			result:  fod.ResultNone,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDialogResult(tt.result)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateDialogResult(%v) error = %v, wantErr %v", tt.result, err, tt.wantErr)
			}
		})
	}
}
