package log

import (
	"testing"
	"time"
)

func TestTimeUnit_Format(t *testing.T) {
	tests := []struct {
		name string
		tr   TimeUnit
		want string
	}{
		{"Minute", Minute, ".%Y%m%d%H%M"},
		{"Hour", Hour, ".%Y%m%d%H"},
		{"Day", Day, ".%Y%m%d"},
		{"Month", Month, ".%Y%m"},
		{"Year", Year, ".%Y"},
		{"default", TimeUnit("xxx"), ".%Y%m%d"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Format(); got != tt.want {
				t.Errorf("TimeUnit.Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeUnit_RotationGap(t *testing.T) {
	tests := []struct {
		name string
		tr   TimeUnit
		want time.Duration
	}{
		{"Minute", Minute, time.Minute},
		{"Hour", Hour, time.Hour},
		{"Day", Day, time.Hour * 24},
		{"Month", Month, time.Hour * 24 * 30},
		{"Year", Year, time.Hour * 24 * 365},
		{"default", TimeUnit("xxx"), time.Hour * 24},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.RotationGap(); got != tt.want {
				t.Errorf("TimeUnit.RotationGap() = %v, want %v", got, tt.want)
			}
		})
	}
}
