package db

import (
	"testing"

	"github.com/warmans/rsk-search/pkg/models"
)

func Test_upsertSynopsis(t *testing.T) {
	tests := []struct {
		name     string
		existing []models.Synopsis
		incoming []models.Synopsis
		want     []models.Synopsis
	}{
		{
			name:     "incoming entry with a new position range is appended",
			existing: []models.Synopsis{{Description: "one", StartPos: 10, EndPos: 20}},
			incoming: []models.Synopsis{{Description: "two", StartPos: 30, EndPos: 40}},
			want: []models.Synopsis{
				{Description: "one", StartPos: 10, EndPos: 20},
				{Description: "two", StartPos: 30, EndPos: 40},
			},
		},
		{
			name:     "incoming entry with an existing position range replaces it",
			existing: []models.Synopsis{{Description: "one", StartPos: 10, EndPos: 20}, {Description: "two", StartPos: 30, EndPos: 40}},
			incoming: []models.Synopsis{{Description: "updated", StartPos: 10, EndPos: 20}},
			want: []models.Synopsis{
				{Description: "updated", StartPos: 10, EndPos: 20},
				{Description: "two", StartPos: 30, EndPos: 40},
			},
		},
		{
			name:     "existing entry with no incoming counterpart is untouched",
			existing: []models.Synopsis{{Description: "one", StartPos: 10, EndPos: 20}},
			incoming: []models.Synopsis{{Description: "updated", StartPos: 30, EndPos: 40}},
			want: []models.Synopsis{
				{Description: "one", StartPos: 10, EndPos: 20},
				{Description: "updated", StartPos: 30, EndPos: 40},
			},
		},
		{
			name:     "entry with matching start but different end does not match",
			existing: []models.Synopsis{{Description: "one", StartPos: 10, EndPos: 20}},
			incoming: []models.Synopsis{{Description: "updated", StartPos: 10, EndPos: 25}},
			want: []models.Synopsis{
				{Description: "one", StartPos: 10, EndPos: 20},
				{Description: "updated", StartPos: 10, EndPos: 25},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := upsertSynopsis(tt.existing, tt.incoming)
			if len(got) != len(tt.want) {
				t.Fatalf("got %d entries, want %d entries: %+v", len(got), len(tt.want), got)
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("entry %d: got %+v, want %+v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func Test_upsertTrivia(t *testing.T) {
	tests := []struct {
		name     string
		existing []models.Trivia
		incoming []models.Trivia
		want     []models.Trivia
	}{
		{
			name:     "incoming entry with an existing position range replaces it",
			existing: []models.Trivia{{Description: "one", StartPos: 10, EndPos: 20}, {Description: "two", StartPos: 30, EndPos: 40}},
			incoming: []models.Trivia{{Description: "updated", StartPos: 30, EndPos: 40}},
			want: []models.Trivia{
				{Description: "one", StartPos: 10, EndPos: 20},
				{Description: "updated", StartPos: 30, EndPos: 40},
			},
		},
		{
			name:     "existing entry with no incoming counterpart is untouched",
			existing: []models.Trivia{{Description: "one", StartPos: 10, EndPos: 20}},
			incoming: []models.Trivia{{Description: "new", StartPos: 30, EndPos: 40}},
			want: []models.Trivia{
				{Description: "one", StartPos: 10, EndPos: 20},
				{Description: "new", StartPos: 30, EndPos: 40},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := upsertTrivia(tt.existing, tt.incoming)
			if len(got) != len(tt.want) {
				t.Fatalf("got %d entries, want %d entries: %+v", len(got), len(tt.want), got)
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("entry %d: got %+v, want %+v", i, got[i], tt.want[i])
				}
			}
		})
	}
}
