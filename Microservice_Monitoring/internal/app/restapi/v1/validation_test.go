package v1

import (
	"net/http/httptest"
	"testing"

	"monitoring/internal/domain/models"

	"github.com/gin-gonic/gin"
)

func newTestContext(rawQuery string) *gin.Context {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/?"+rawQuery, nil)
	c.Request = req
	return c
}

func TestParsePagination_Defaults(t *testing.T) {
	c := newTestContext("")

	pagination, err := ParsePagination(c)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if pagination.Limit != defaultPageLimit {
		t.Fatalf("expected default limit %d, got %d", defaultPageLimit, pagination.Limit)
	}
	if pagination.Offset != 0 {
		t.Fatalf("expected default offset 0, got %d", pagination.Offset)
	}
}

func TestParsePagination_WithValidValues(t *testing.T) {
	c := newTestContext("limit=25&offset=10")

	pagination, err := ParsePagination(c)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if pagination.Limit != 25 {
		t.Fatalf("expected limit 25, got %d", pagination.Limit)
	}
	if pagination.Offset != 10 {
		t.Fatalf("expected offset 10, got %d", pagination.Offset)
	}
}

func TestParsePagination_LimitZeroFallsBackToDefault(t *testing.T) {
	c := newTestContext("limit=0")

	pagination, err := ParsePagination(c)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if pagination.Limit != defaultPageLimit {
		t.Fatalf("expected default limit %d, got %d", defaultPageLimit, pagination.Limit)
	}
}

func TestParsePagination_InvalidLimit(t *testing.T) {
	c := newTestContext("limit=abc")

	_, err := ParsePagination(c)
	if err == nil {
		t.Fatal("expected error for invalid limit, got nil")
	}
	if err.Error() != "invalid limit" {
		t.Fatalf("expected error 'invalid limit', got '%v'", err)
	}
}

func TestParsePagination_TooLargeLimit(t *testing.T) {
	c := newTestContext("limit=501")

	_, err := ParsePagination(c)
	if err == nil {
		t.Fatal("expected error for too large limit, got nil")
	}
	if err.Error() != "limit must be 500 or less" {
		t.Fatalf("expected limit bound error, got '%v'", err)
	}
}

func TestParsePagination_InvalidOffset(t *testing.T) {
	c := newTestContext("offset=x")

	_, err := ParsePagination(c)
	if err == nil {
		t.Fatal("expected error for invalid offset, got nil")
	}
	if err.Error() != "invalid offset" {
		t.Fatalf("expected error 'invalid offset', got '%v'", err)
	}
}

func TestParsePagination_NegativeOffset(t *testing.T) {
	c := newTestContext("offset=-1")

	_, err := ParsePagination(c)
	if err == nil {
		t.Fatal("expected error for negative offset, got nil")
	}
	if err.Error() != "offset must be 0 or greater" {
		t.Fatalf("expected negative offset error, got '%v'", err)
	}
}

func TestValidateConfiguratieForSensorType_NegativeMargin(t *testing.T) {
	sensorType := models.SensorType{DrempelIsRange: true}
	margin := -1.0
	config := &models.UpdateSensorConfiguratieRequest{MargePercentage: &margin}

	errMsg := validateConfiguratieForSensorType(sensorType, config)
	if errMsg != "margePercentage mag niet negatief zijn" {
		t.Fatalf("unexpected error message: %s", errMsg)
	}
}

func TestValidateConfiguratieForSensorType_RangeMissingValues(t *testing.T) {
	sensorType := models.SensorType{DrempelIsRange: true}
	config := &models.UpdateSensorConfiguratieRequest{MinValue: new(1.0), MaxValue: nil}

	errMsg := validateConfiguratieForSensorType(sensorType, config)
	if errMsg != "minValue en maxValue zijn verplicht voor range-sensoren" {
		t.Fatalf("unexpected error message: %s", errMsg)
	}
}

func TestValidateConfiguratieForSensorType_RangeMinMustBeLowerThanMax(t *testing.T) {
	sensorType := models.SensorType{DrempelIsRange: true}
	config := &models.UpdateSensorConfiguratieRequest{
		MinValue: new(10.0),
		MaxValue: new(10.0),
	}

	errMsg := validateConfiguratieForSensorType(sensorType, config)
	if errMsg != "minValue moet lager zijn dan maxValue" {
		t.Fatalf("unexpected error message: %s", errMsg)
	}
}

func TestValidateConfiguratieForSensorType_RangeValid(t *testing.T) {
	sensorType := models.SensorType{DrempelIsRange: true}
	config := &models.UpdateSensorConfiguratieRequest{
		MinValue: new(10.0),
		MaxValue: new(20.0),
	}

	errMsg := validateConfiguratieForSensorType(sensorType, config)
	if errMsg != "" {
		t.Fatalf("expected no error, got: %s", errMsg)
	}
}

func TestValidateConfiguratieForSensorType_NonRangeMinRequired(t *testing.T) {
	sensorType := models.SensorType{DrempelIsRange: false}
	config := &models.UpdateSensorConfiguratieRequest{MinValue: nil, MaxValue: new(20.0)}

	errMsg := validateConfiguratieForSensorType(sensorType, config)
	if errMsg != "minValue is verplicht als normwaarde voor niet-range sensoren" {
		t.Fatalf("unexpected error message: %s", errMsg)
	}
}

func TestValidateConfiguratieForSensorType_NonRangeClearsMaxValue(t *testing.T) {
	sensorType := models.SensorType{DrempelIsRange: false}
	config := &models.UpdateSensorConfiguratieRequest{MinValue: new(12.5), MaxValue: new(20.0)}

	errMsg := validateConfiguratieForSensorType(sensorType, config)
	if errMsg != "" {
		t.Fatalf("expected no error, got: %s", errMsg)
	}
	if config.MaxValue != nil {
		t.Fatal("expected MaxValue to be cleared for non-range sensor")
	}
}
