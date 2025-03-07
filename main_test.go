package sorbet

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    map[string]string
		wantErr bool
		errType string
	}{
		{
			name:  "simple key-value",
			input: "key => value",
			want:  map[string]string{"key": "value"},
		},
		{
			name:    "invalid syntax",
			input:   "key =>",
			wantErr: true,
			errType: "Syntax",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				for k, v := range tt.want {
					if got[k] != v {
						t.Errorf("Parse() got[%s] = %v, want[%s] = %v", k, got[k], k, v)
					}
				}
			}
			if tt.wantErr && err != nil {
				if serr, ok := err.(*SorbetError); !ok || serr.ErrorType != tt.errType {
					t.Errorf("Parse() error type = %v, want %v", serr.ErrorType, tt.errType)
				}
			}
		})
	}
}
