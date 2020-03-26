package model

import (
	"reflect"
	"testing"
)

func TestReadAll(t *testing.T) {
	tests := []struct {
		name    string
		want    []GraduationData
		wantErr bool
	}{
		{
			name:    "normal",
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransformYear(t *testing.T) {
	type args struct {
		raw string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal",
			args: args{
				raw: "2019",
			},
			want: "二〇一九",
		},
		{
			name: "abnormal",
			args: args{
				raw: "20i0",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TransformYear(tt.args.raw); got != tt.want {
				t.Errorf("TransformYear() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransMonthOrDate(t *testing.T) {
	type args struct {
		raw string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "个位",
			args: args{
				raw: "1",
			},
			want: "一",
		}, {
			name: "十",
			args: args{
				raw: "10",
			},
			want: "十",
		}, {
			name: "一开头的十位",
			args: args{
				raw: "11",
			},
			want: "十一",
		},{
			name: "非一开头的十位",
			args : args {
				raw: "23",
			},
			want: "二十三",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TransformMonthOrDate(tt.args.raw); got != tt.want {
				t.Errorf("TransMonthOrDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
