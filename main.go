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

// CalculationExpression обробляє математичний вираз і повертає результат як рядок
func CalculationExpression(expression string) (string, error) {
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

// GetCalculations повертає всі обчислення у форматі JSON
func GetCalculations(c echo.Context) error {
	return c.JSON(http.StatusOK, calculations)
}

// PostCalculations Передає JSON з expression
func PostCalculations(c echo.Context) error {
	//Передаємо запрос
	var req CalculationRequest
	//«Декодуємо запит» — це означає «читаємо дані, які прийшли від клієнта, і перетворюємо їх у зручну структуру, з якою можна працювати в коді».
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Рахуємо новий результат
	result, err := CalculationExpression(req.Expression)
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
	//Дістає id з силки яку ми хочемо обновити
	id := c.Param("id")

	//Передаємо запрос
	var req CalculationRequest
	//Декодуємо запит
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	// Рахуємо новий результат
	result, err := CalculationExpression(req.Expression)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid expression"})
	}

	//Оновлюємо
	for i, calculation := range calculations {
		if calculation.ID == id {
			calculations[i].Expression = req.Expression
			calculations[i].Result = result
			return c.JSON(http.StatusOK, calculations[i])
		}
	}
	return c.JSON(http.StatusBadRequest, map[string]string{"error": "Calculation not found"})
}

func DeleteCalculations(c echo.Context) error {
	id := c.Param("id")

	for i, calculation := range calculations {
		if calculation.ID == id {
			calculations = append(calculations[:i], calculations[i+1:]...)
			//Повертає відповідь без тіла 'код 204'
			return c.NoContent(http.StatusNoContent)
		}
	}
	return c.JSON(http.StatusBadRequest, map[string]string{"error": "Calculation not found"})
}

func main() {
	e := echo.New()

	e.Use(middleware.CORS())   // запит -> middleware -> обробка CORS -> передає на сервер
	e.Use(middleware.Logger()) // запрос -> middleware -> логірує -> передає на сервер

	e.GET("/calculations", GetCalculations)
	e.POST("/calculations", PostCalculations)
	e.PATCH("/calculations/:id", PatchCalculations)
	e.DELETE("/calculations/:id", DeleteCalculations)

	e.Start("localhost:8080")
}
