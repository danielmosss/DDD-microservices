package asset

import "time"

type LifecycleState string

const (
	LifecycleCreated     LifecycleState = "CREATED"
	LifecycleActive      LifecycleState = "ACTIVE"
	LifecycleMaintenance LifecycleState = "MAINTENANCE"
	LifecycleRetired     LifecycleState = "RETIRED"
)

type CriticalityLevel string

const (
	CriticalityLow      CriticalityLevel = "LOW"
	CriticalityMedium   CriticalityLevel = "MEDIUM"
	CriticalityHigh     CriticalityLevel = "HIGH"
	CriticalityCritical CriticalityLevel = "CRITICAL"
)

type RiskLevel string

const (
	RiskLow      RiskLevel = "LOW"
	RiskMedium   RiskLevel = "MEDIUM"
	RiskHigh     RiskLevel = "HIGH"
	RiskCritical RiskLevel = "CRITICAL"
)

type StrategyType string

const (
	StrategyPreventive     StrategyType = "PREVENTIVE"
	StrategyCorrective     StrategyType = "CORRECTIVE"
	StrategyConditionBased StrategyType = "CONDITION_BASED"
)

type Location struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}

type Asset struct {
	AssetID        string           `json:"assetId"`
	Type           string           `json:"type"`
	Location       Location         `json:"location"`
	Criticality    CriticalityLevel `json:"criticality"`
	LifecycleState LifecycleState   `json:"lifecycleState"`
	Value          string           `json:"value"`

	Requirements          []AssetRequirement    `json:"requirements"`
	Performances          []AssetPerformance    `json:"performances"`
	RiskAssessments       []RiskAssessment      `json:"riskAssessments"`
	MaintenanceStrategies []MaintenanceStrategy `json:"maintenanceStrategies"`
}

type AssetRequirement struct {
	RequirementID      string  `json:"requirementId"`
	AvailabilityTarget float64 `json:"availabilityTarget"`
	PerformanceTarget  string  `json:"performanceTarget"`
	MaxDowntimeHours   int     `json:"maxDowntimeHours"`
}

type AssetPerformance struct {
	PerformanceID  string    `json:"performanceId"`
	Availability   float64   `json:"availability"`
	ConditionScore int       `json:"conditionScore"`
	AssessmentDate time.Time `json:"assessmentDate"`
	DataSource     string    `json:"dataSource"`
}

type RiskAssessment struct {
	RiskID           string    `json:"riskId"`
	RiskLevel        RiskLevel `json:"riskLevel"`
	ImpactScore      int       `json:"impactScore"`
	ProbabilityScore int       `json:"probabilityScore"`
	AssessmentDate   time.Time `json:"assessmentDate"`
}

type MaintenanceStrategy struct {
	StrategyID                string       `json:"strategyId"`
	StrategyType              StrategyType `json:"strategyType"`
	MaintenanceIntervalMonths int          `json:"maintenanceIntervalMonths"`
	ValidFrom                 time.Time    `json:"validFrom"`
}
