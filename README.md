# Dracory Web Framework

Dracory is a lightweight and efficient Golang web framework designed to simplify the development of web applications. It emphasizes clarity, speed, and ease of use, providing a solid foundation for building scalable and maintainable applications.

## Key Features

* **Simple and Efficient Routing:** Dracory uses declarative routing as in "github.com/gouniverse/router"
* **Clean Architecture:** The framework adheres to the MVC (Model-View-Controller) pattern, promoting a clear separation of concerns and organized code structure.
* **Configuration Management:** Dracory uses flexible configuration loading from .env files and environment variables, simplifying deployment and environment-specific setups.
* **Database Integration:** Seamless integration with `database/sql` and `sqlx` facilitates efficient database interactions and data modeling.
* **Middleware Support:** Robust middleware capabilities enable easy implementation of cross-cutting concerns like authentication, logging, and request processing.
* **Testing Focus:** Dracory encourages test-driven development with easy-to-use testing utilities and configuration isolation for concurrent testing.
* **Web Server:** Built-in web server functionality simplifies setting up a web server.
* **Web Authentication:** Built-in web authentication middleware simplifies securing your web applications.
* **Minimalistic Templating:** Basic HTML templating using the pure Go HTML builder library `github.com/gouniverse/hb` provides a straightforward approach to server-side rendering.

## Getting Started

1. Clone this repository
2. Copy `.env.example` to `.env` and configure your environment variables
3. Run `go run cmd/web/main.go` to start the server

## Project Structure

```
dracory-project/
├── cmd/
│   └── web/            # Main application entry point
├── internal/
│   ├── app/
│   │   ├── controllers/ # HTTP handlers and business logic
│   │   ├── models/      # Data models and database interactions
│   │   ├── config/      # Configuration handling
│   │   └── routes/      # Routing definitions
│   └── platform/
│       └── database/   # Database connection and utilities
├── web/                # Static assets (HTML, CSS, JS)
│   ├── static/
│   └── templates/
├── test/               # Integration and end-to-end tests
├── .env                # Environment variables
├── go.mod              # Go module definition
└── README.md           # Project documentation
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
