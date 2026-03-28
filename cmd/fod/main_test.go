package main

import (
	"testing"

	"github.com/urfave/cli"
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

func TestFlags_ModeUsage(t *testing.T) {
	definedFlags, err := flags()
	if err != nil {
		t.Fatalf("flags() error = %v", err)
	}

	for _, f := range definedFlags {
		flag, ok := f.(cli.StringFlag)
		if !ok || flag.Name != "mode, m" {
			continue
		}

		want := "start mode: d|dir|directory for directory select mode, f|file for file select mode"
		if flag.Usage != want {
			t.Fatalf("mode usage = %q, want %q", flag.Usage, want)
		}
		return
	}

	t.Fatal("mode flag not found")
}
