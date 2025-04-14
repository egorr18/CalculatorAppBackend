package calculationService

// Calculation містить дані про обчислення
type Calculation struct {
	ID         string `json:"id"`
	Expression string `json:"expression"`
	Result     string `json:"result"`
}

// CalculationRequest представляє запит для обчислення
type CalculationRequest struct {
	Expression string `json:"expression"`
}
