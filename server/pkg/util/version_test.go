package util

import "testing"

func TestNextVersion(t *testing.T) {
	type args struct {
		currentVersion string
		change         VersionChange
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "empty string returns error",
			args: args{
				currentVersion: "",
				change:         MinorVersion,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "wrong number of parts returns an error",
			args: args{
				currentVersion: "1.0.0.0",
				change:         MinorVersion,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "invalid parts returns an error",
			args: args{
				currentVersion: "1.0.foo",
				change:         MinorVersion,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "increment major version",
			args: args{
				currentVersion: "1.0.0",
				change:         MajorVersion,
			},
			want:    "2.0.0",
			wantErr: false,
		},
		{
			name: "increment minor version",
			args: args{
				currentVersion: "1.0.0",
				change:         MinorVersion,
			},
			want:    "1.1.0",
			wantErr: false,
		},
		{
			name: "increment patch version",
			args: args{
				currentVersion: "1.0.0",
				change:         PatchVersion,
			},
			want:    "1.0.1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NextVersion(tt.args.currentVersion, tt.args.change)
			if (err != nil) != tt.wantErr {
				t.Errorf("NextVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NextVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}
