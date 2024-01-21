package custom

import (
	"reflect"
	"testing"
)

func TestReq_ConvertToPage(t *testing.T) {
	type fields struct {
		Page     int32
		PageSize int32
		Order    []*OrderParam
		Search   []*SearchParam
	}
	type args struct {
		total int32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *PaginatorReply
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				Page:     0,
				PageSize: -1,
				Order:    nil,
				Search:   nil,
			},
			args: args{
				total: 100,
			},
			want: &PaginatorReply{
				Page:      0,
				PageSize:  0,
				Total:     100,
				PrevPage:  0,
				NextPage:  0,
				TotalPage: 0,
			},
			wantErr: true,
		},
		{
			name: "test2",
			fields: fields{
				Page:     0,
				PageSize: 0,
				Order:    nil,
				Search:   nil,
			},
			args: args{
				total: 100,
			},
			want: &PaginatorReply{
				Page:      0,
				PageSize:  0,
				Total:     100,
				PrevPage:  0,
				NextPage:  0,
				TotalPage: 0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PaginatorReq{
				Page:     tt.fields.Page,
				PageSize: tt.fields.PageSize,
				Order:    tt.fields.Order,
				Search:   tt.fields.Search,
			}
			got, err := p.ConvertToPage(tt.args.total)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToPage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToPage() got = %v, want %v", got, tt.want)
			}
		})
	}
}
