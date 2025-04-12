# Go Calculator API

Це простий REST API для обчислення математичних виразів, створений на основі **Go**, **Echo**, **GORM** та **PostgreSQL**. API підтримує створення, оновлення, видалення та отримання результатів обчислень.

---

## 🔧 Технології

- [Go](https://golang.org/)
- [Echo](https://echo.labstack.com/) – веб-фреймворк
- [GORM](https://gorm.io/) – ORM для роботи з базами даних
- [PostgreSQL](https://www.postgresql.org/)
- [govaluate](https://github.com/Knetic/govaluate) – парсинг і обчислення математичних виразів

---

## 🚀 Запуск проєкту

### Встановлення залежностей

```bash
go get github.com/Knetic/govaluate
go get github.com/google/uuid
go get github.com/labstack/echo/v4
go get github.com/labstack/echo/v4/middleware
go get gorm.io/driver/postgres
go get gorm.io/gorm
```

## 📬 API Маршрути

| Метод  | Шлях                | Опис                       |
| ------ | ------------------- | -------------------------- |
| GET    | `/calculations`     | Отримати всі записи у базі |
| POST   | `/calculations`     | Створити новий вираз       |
| PATCH  | `/calculations/:id` | Оновити існуючий вираз     |
| DELETE | `/calculations/:id` | Видалити вираз з бази      |

---

## 📤 Формат запиту та відповіді

### POST `/calculations`

**Запит:**

```json
{
"expression": "5 + 5 * 2"
}
```

**Відповідь:**

```json
{
"id": "f1c2a6b0-92b9-4d8c-8716-03fcaeb12aaf",
"expression": "5 + 5 * 2",
"result": "15"
}
```

---

## 📃 Структура БД (GORM)

```go
type Calculation struct {
ID         string `json:"id"`
Expression string `json:"expression"`
Result     string `json:"result"`
}
```

---

## ✅ Функціонал

- ✅ Обчислення виразів (наприклад: `55+20*(3-1)`)
- ✅ Збереження історії в PostgreSQL
- ✅ Повертає результат у JSON
- ✅ REST API: GET, POST, PATCH, DELETE

___

