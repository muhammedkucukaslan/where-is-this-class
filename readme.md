<h1 align="center">Where Is This Class? </h1>

<p align="center">
  <img src="./public/image.png" alt="Banner" width="600"/>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Golang-00ADD8?style=for-the-badge&logo=go&logoColor=white"/>
  <img src="https://img.shields.io/badge/Fiber-00A99D?style=for-the-badge&logo=go-fiber&logoColor=white"/>
  <img src="https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white"/>
  <img src="https://img.shields.io/badge/JWT-000000?style=for-the-badge&logo=JSON%20web%20tokens&logoColor=white"/>
</p>

A simple API that helps students and staff locate classrooms in a university by entering a class code. It supports multilingual descriptions and records the most searched rooms.

## Live Demo

[Here](https://muhammedkucukaslan.github.io/where-is-this-class-frontend/)

## Installation

```bash
git clone https://github.com/muhammedkucukaslan/where-is-this-class.git
cd where-is-this-class
```

Set up your `.env` file in the project root:

```env
CLIENT_URL=http://localhost:3000
JWT_SECRET=secret
ADMIN_PASSWORD=admin123
DATABASE_URL=postgres://postgres:postgres@localhost:5432/postgres
```

You may do not want to use sql. Only thing you should do is change `repository.go` with respect to the interface that `hander.go` wants.

Run the project:

```bash
go run .
```

## API Endpoints

### Login as Admin

Creates a JWT token and sets it as a cookie.

```http
POST /login
```

#### Request Body

```json
{
  "password": "string"
}
```

#### Example cURL

```bash
curl -X POST "http://localhost:8000/login"   -H "Content-Type: application/json"   -d '{"password": "admin123"}'
```

### Get a Classroom by Code

Returns classroom details and location instructions in the specified language.

```http
GET /classrooms/:code?language={language}
```

#### Example

```bash
curl -X GET "http://localhost:8000/classrooms/EG010?language=tr"
```

### Get Most Searched Classrooms

Returns a list of the most frequently searched classrooms.

```http
GET /classrooms/most-visited
```

#### Example

```bash
curl -X GET "http://localhost:8000/classrooms/most-visited"
```

### Create a Classroom

Requires authentication (JWT cookie). Adds a new classroom with optional image and translations.

```http
POST /classrooms
```

#### Request Body

```json
{
  "code": "string",
  "floor": "integer",
  "imageUrl": "string",
  "translations": [
    {
      "language": "string",
      "building": "string",
      "description": "string"
    }
  ]
}
```

> Note: I just configured the validation for three languages. You can change validation in the `handler.go:116`.

#### Example 

```bash
curl -X POST "http://localhost:8000/classrooms" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "EG010",
    "floor": 0,
    "imageUrl": "https://example.com/image.jpeg",
    "translations": [
      {
        "language": "en",
        "building": "Education Sciences",
        "description": "Classroom description in English."
      }
    ]
  }'
```

## Notes

- Make sure the database is running and accessible.
- Responses are in JSON format.
- You can integrate this backend with a frontend client at `CLIENT_URL`.
