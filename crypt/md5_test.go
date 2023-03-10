package crypt

import "testing"

func TestMd5String(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				s: "abc",
			},
			want: "900150983cd24fb0d6963f7d28e17f72",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Md5String(tt.args.s); got != tt.want {
				t.Errorf("Md5String() = %v, want %v", got, tt.want)
			}
		})
	}
}
