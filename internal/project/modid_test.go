package project

import "testing"

func TestValidateModID(t *testing.T) {
	tests := []struct {
		name    string
		modID   string
		wantErr bool
	}{
		{name: "valid", modID: "my_mod", wantErr: false},
		{name: "too short", modID: "t", wantErr: true},
		{name: "starts with number", modID: "1mod", wantErr: true},
		{name: "invalid chars", modID: "My Mod", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateModID(tt.modID)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error for %q", tt.modID)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error for %q: %v", tt.modID, err)
			}
		})
	}
}

func TestNormalizeAutoModID(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{input: "My Cool Mod", want: "my_cool_mod"},
		{input: "t", want: "t_mod"},
		{input: "123", want: "mod_123"},
		{input: "!!!", want: "mymod"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := NormalizeAutoModID(tt.input); got != tt.want {
				t.Fatalf("NormalizeAutoModID(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
