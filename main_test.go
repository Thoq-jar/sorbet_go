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
		errMsg  string
	}{
		{
			name:  "simple key value",
			input: "name => john",
			want:  map[string]string{"name": "john"},
		},
		{
			name: "multiple key values",
			input: `first => john
second => doe`,
			want: map[string]string{
				"first":  "john",
				"second": "doe",
			},
		},
		{
			name: "with continuation",
			input: `list => first
> second
> third`,
			want: map[string]string{
				"list": "first,second,third",
			},
		},
		{
			name: "multiple keys with continuation",
			input: `first => 1
> 2
second => a
> b
> c`,
			want: map[string]string{
				"first":  "1,2",
				"second": "a,b,c",
			},
		},
		{
			name:    "invalid syntax",
			input:   "invalid line =>",
			wantErr: true,
			errType: "Syntax",
			errMsg:  "Syntax error! Expected [key] => [value] at: invalid line =>",
		},
		{
			name:    "continuation without key",
			input:   "> continuation",
			wantErr: true,
			errType: "SyntaxException",
			errMsg:  "Continuation line without a key at: > continuation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.wantErr {
						t.Errorf("Parse() panicked unexpectedly: %v", r)
						return
					}
					errorMsg := r.(string)
					expectedError := tt.errType + ": " + tt.errMsg
					if errorMsg != expectedError {
						t.Errorf("Parse() error = %v, want %v", errorMsg, expectedError)
					}
				} else if tt.wantErr {
					t.Error("Parse() expected error but got none")
				}
			}()

			got := Parse(tt.input)
			if tt.wantErr {
				t.Error("Parse() expected error but got none")
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
				return
			}

			for k, v := range tt.want {
				if gotV, ok := got[k]; !ok || gotV != v {
					t.Errorf("Parse() got[%s] = %v, want[%s] = %v", k, gotV, k, v)
				}
			}
		})
	}
}
