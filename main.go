package main

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
)

var db *gorm.DB

// Відбувається ініціалізація з'єднання з базою даних PostgreSQL
func InitDB() {
	//Підключення до PostgreSQL
	dsn := "host=localhost user=postgres password=yourpassword dbname=postgres port=5432 sslmode=disable"
	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	// Автоматично створює або оновлює таблиці в середині бази данних
	if err := db.AutoMigrate(&Calculation{}); err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
}

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

//Основні методи ORM - Create, Find, Update, Delete

// GetCalculations повертає всі обчислення у форматі JSON
func GetCalculations(c echo.Context) error {
	var calculations []Calculation
	// Find використовується для знаходження всіх записів в базі данних
	if err := db.Find(&calculations).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get calculations"})
	}
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

	// Створюємо новий вираз для запису в базу данних
	calc := Calculation{
		uuid.NewString(),
		req.Expression,
		result,
	}
	// Додаємо в історію
	if err := db.Create(&calc).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not add calculation"})
	}
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

	//Найщли вираз по id
	var calc Calculation
	if err := db.First(&calc, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Could not find expression"})
	}

	//Новий вираз який нам передав сайт
	calc.Expression = req.Expression
	calc.Result = result

	//Метод для збереження змін
	if err := db.Save(&calc).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not update calculation"})
	}
	return c.JSON(http.StatusOK, calc) //Повертаємо новий оновлений вираз
}

func DeleteCalculations(c echo.Context) error {
	//Шукаємо id по якому будемо видаляти
	id := c.Param("id")
	//Видаляємо
	if err := db.Delete(&Calculation{}, "id", id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not delete calculation"})
	}

	return c.NoContent(http.StatusNoContent)
}

func main() {
	InitDB()

	e := echo.New()

	e.Use(middleware.CORS())   // запит -> middleware -> обробка CORS -> передає на сервер
	e.Use(middleware.Logger()) // запрос -> middleware -> логірує -> передає на сервер

	e.GET("/calculations", GetCalculations)
	e.POST("/calculations", PostCalculations)
	e.PATCH("/calculations/:id", PatchCalculations)
	e.DELETE("/calculations/:id", DeleteCalculations)

	e.Start("localhost:8080")
}
