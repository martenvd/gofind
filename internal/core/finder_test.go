package core

import "testing"

func TestClampIndex(t *testing.T) {
	tests := []struct {
		name     string
		current  int
		delta    int
		maxIndex int
		expected int
	}{
		{
			name:     "move down within bounds",
			current:  0,
			delta:    10,
			maxIndex: 20,
			expected: 10,
		},
		{
			name:     "move down clamped to max",
			current:  15,
			delta:    10,
			maxIndex: 20,
			expected: 20,
		},
		{
			name:     "move up within bounds",
			current:  15,
			delta:    -10,
			maxIndex: 20,
			expected: 5,
		},
		{
			name:     "move up clamped to zero",
			current:  3,
			delta:    -10,
			maxIndex: 20,
			expected: 0,
		},
		{
			name:     "zero delta stays put",
			current:  5,
			delta:    0,
			maxIndex: 20,
			expected: 5,
		},
		{
			name:     "empty list returns zero",
			current:  0,
			delta:    10,
			maxIndex: -1,
			expected: 0,
		},
		{
			name:     "single item list move down stays at zero",
			current:  0,
			delta:    10,
			maxIndex: 0,
			expected: 0,
		},
		{
			name:     "single item list move up stays at zero",
			current:  0,
			delta:    -10,
			maxIndex: 0,
			expected: 0,
		},
		{
			name:     "exact landing on max",
			current:  10,
			delta:    10,
			maxIndex: 20,
			expected: 20,
		},
		{
			name:     "exact landing on zero",
			current:  10,
			delta:    -10,
			maxIndex: 20,
			expected: 0,
		},
		{
			name:     "large delta down",
			current:  0,
			delta:    1000,
			maxIndex: 50,
			expected: 50,
		},
		{
			name:     "large delta up",
			current:  50,
			delta:    -1000,
			maxIndex: 50,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClampIndex(tt.current, tt.delta, tt.maxIndex)
			if got != tt.expected {
				t.Errorf("ClampIndex(%d, %d, %d) = %d, want %d",
					tt.current, tt.delta, tt.maxIndex, got, tt.expected)
			}
		})
	}
}
