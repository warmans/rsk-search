package models

import "testing"

func Test_parsePositionSpec(t *testing.T) {
	tests := []struct {
		name    string
		pos     string
		want    int64
		want1   int64
		wantErr bool
	}{
		{
			name:    "empty string returns error",
			pos:     "",
			want:    0,
			want1:   0,
			wantErr: true,
		},
		{
			name:    "valid start position",
			pos:     "101",
			want:    101,
			want1:   102,
			wantErr: false,
		},
		{
			name:    "valid start and end position",
			pos:     "101-110",
			want:    101,
			want1:   110,
			wantErr: false,
		},
		{
			name:    "invalid stat position",
			pos:     "foo",
			want:    0,
			want1:   0,
			wantErr: true,
		},
		{
			name:    "invalid end position",
			pos:     "1-bar",
			want:    0,
			want1:   0,
			wantErr: true,
		},
		{
			name:    "more than one dash is invalid",
			pos:     "1-10-30",
			want:    0,
			want1:   0,
			wantErr: true,
		},
		{
			name:    "if the start and end pos are equal increment the end pos by 1",
			pos:     "2-2",
			want:    2,
			want1:   2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parsePositionRange(tt.pos)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePositionRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parsePositionRange() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("parsePositionRange() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
