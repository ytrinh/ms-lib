package errors

import "testing"

func TestE(t *testing.T) {
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Empty", args{args: nil}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := E(tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("E() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
