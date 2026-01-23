package model

import "testing"

func TestBookLevel_IsZero(t *testing.T) {
	tests := []struct {
		name  string
		level BookLevel
		want  bool
	}{
		{"zero", BookLevel{}, true},
		{"has price", BookLevel{Price: 100, Quantity: 0}, false},
		{"has quantity", BookLevel{Price: 0, Quantity: 10}, false},
		{"has both", BookLevel{Price: 100, Quantity: 10}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.level.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}
