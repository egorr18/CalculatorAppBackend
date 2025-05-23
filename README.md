# 📟 Go Calculator API

Це простий REST API калькулятор, написаний на Go з використанням фреймворку Echo і бібліотеки govaluate для парсингу математичних виразів.

---

## 🚀 Запуск

go run main.go

API буде доступне на http://localhost:8080

---

## Функціонал

- Обчислення математичних виразів (наприклад: 55+20*(3-1))
- Збереження історії обчислень у памʼяті
- Віддача результатів у форматі JSON
- Простий REST API з методами GET та POST PATCH та DELETE
- Можливість редагувати (PATCH) попередні обчислення
- Можливість видаляти (DELETE) окремі записи або всю історію
---

## Routes

| Метод | Шлях             | Опис                           |
|-------|------------------|--------------------------------|
| GET   | /calculations    | Отримати всі обчислення        |
| POST  | /calculations    | Надіслати вираз для обчислення |
| Patch  | /calculations/:id    | Оновити існуючий вираз за ID   |
| Delete  | /calculations/:id    | Видалити обчислення            |


---

## Залежності

- github.com/labstack/echo/v4
- github.com/labstack/echo/v4/middleware
- github.com/Knetic/govaluate
- github.com/google/uuid

---

## Приклад тіла запиту

{
  "expression": "10 + 20 * (3 - 1)"
}

---

## Особливості

-Повноцінна підтримка математичних виразів (+, -, *, /, %, дужки, логічні оператори)
-Генерація унікального UUID для кожного обчислення
-Вбудовані middleware для CORS безпеки та логування запитів
-Можливість не лише додавати, а й оновлювати/видаляти обчислення (CRUD API)
-Зрозумілі JSON-відповіді з деталями обчислень
-Автоматична валідація вхідних даних

