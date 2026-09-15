# SchoolHub API

Учебный backend для практики: регистрация, вход с токеном и CRUD учеников.

## Запуск

```powershell
cd practica
go run .
```

Сервер слушает `http://localhost:8080`. Данные хранятся в памяти и сбрасываются после перезапуска.

## Быстрая проверка в Postman

1. `POST /auth/register`

```json
{
  "email": "student@example.com",
  "password": "secret123"
}
```

2. `POST /auth/login` с теми же данными. Скопируйте `token` из ответа.
3. Для `/me` и всех `/students` добавьте заголовок `Authorization: Bearer <token>`.
4. `POST /students`:

```json
{
  "full_name": "Алия Садыкова",
  "class": "10А",
  "age": 16,
  "date_of_birth": "2010-03-15",
  "email": "aliya@example.com",
  "gender": "female"
}
```

Поддерживаются `GET /students`, `GET /students/{id}`, `PUT/PATCH /students/{id}` и `DELETE /students/{id}`.
