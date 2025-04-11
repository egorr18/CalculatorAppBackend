package main

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

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

// Глобальна перемінна історія наших виразів
var calculations = []Calculation{}

// calculationExpression обробляє математичний вираз і повертає результат як рядок
func calculationExpression(expression string) (string, error) {
	expr, err := govaluate.NewEvaluableExpression(expression) // Створили вираз (55+55) // Дозволяє працювати з більш гнучкими та складними виразами, наприклад, з кастомними функціями.
	if err != nil {
		return "", err // Передали 55++55
	}
	// Простий та швидкий спосіб для обчислення виразів.
	result, err := expr.Evaluate(nil)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", result), nil
}

// getCalculations повертає всі обчислення у форматі JSON
func getCalculations(c echo.Context) error {
	return c.JSON(http.StatusOK, calculations)
}

// Передає JSON з expression
func postCalculations(c echo.Context) error {
	var req CalculationRequest
	//«Декодуємо запит» — це означає «читаємо дані, які прийшли від клієнта, і перетворюємо їх у зручну структуру, з якою можна працювати в коді».
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// рахуємо
	result, err := calculationExpression(req.Expression)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid expression"}) // map — структура для зберігання даних у форматі "ключ → значення"
	} // Якщо все ок

	// Створюємо новий вираз
	calc := Calculation{
		uuid.NewString(),
		req.Expression,
		result,
	}
	// Додаємо в історію
	calculations = append(calculations, calc)
	// і успішно повертаємо
	return c.JSON(http.StatusCreated, calc)
}

func PatchCalculations(c echo.Context) error {
	//достає id з силки яку ми хочемо обновити
	id := c.Param("id")

	//Передаємо запрос
	var req CalculationRequest
	//Декодуємо запит
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	// рахуємо новий результат
	result, err := CalculationExpression(req.Expression)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid expression"})
	}

	for i, calculation := range calculations {
		if calculation.ID == id {
			calculations[i].Expression = req.Expression
			calculations[i].Result = result
			return c.JSON(http.StatusOK, calculations[i])
		}
	}
	return c.JSON(http.StatusBadRequest, map[string]string{"error": "Calculation not found"})
}

func main() {
	e := echo.New()

	e.Use(middleware.CORS())   // запит -> middleware -> обробка CORS -> передає на сервер
	e.Use(middleware.Logger()) // запрос -> middleware -> логірує -> передає на сервер

	e.GET("/calculations", getCalculations)
	e.POST("/calculations", postCalculations)

	e.Start("localhost:8080")
}
