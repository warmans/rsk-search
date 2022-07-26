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
		// TODO: Add test cases.
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
