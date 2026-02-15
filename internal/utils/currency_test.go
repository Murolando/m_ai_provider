package utils

import (
	"testing"
)

func TestGetUSDToRUBRate(t *testing.T) {
	rate, err := GetUSDToRUBRate()
	if err != nil {
		t.Errorf("GetUSDToRUBRate() failed with error: %v", err)
		return
	}

	if rate <= 0 {
		t.Errorf("Expected positive rate, got %f", rate)
	}

	if rate < 50 || rate > 200 {
		t.Logf("Warning: USD rate seems unusual: %f", rate)
	}

	t.Logf("Current USD to RUB rate: %f", rate)
}

func TestGetUSDToRUBRateRealistic(t *testing.T) {
	rate, err := GetUSDToRUBRate()
	if err != nil {
		t.Errorf("GetUSDToRUBRate() failed with error: %v", err)
		return
	}

	if rate <= 0 {
		t.Errorf("Rate should be positive, got %f", rate)
	}

	if rate > 1000 {
		t.Errorf("Rate seems too high, got %f", rate)
	}

	if rate < 10 {
		t.Errorf("Rate seems too low, got %f", rate)
	}
}
