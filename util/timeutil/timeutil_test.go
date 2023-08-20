package timeutil

import (
	"testing"
	"time"
)

func TestGetShowTime(t *testing.T) {
	type args struct {
		ts time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "刚刚发布",
			args: args{
				ts: time.Now(),
			},
			want: "刚刚发布",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetShowTime(tt.args.ts); got != tt.want {
				t.Errorf("GetShowTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
