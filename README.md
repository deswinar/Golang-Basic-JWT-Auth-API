# Golang Basic JWT Authentication API

This is a basic Go (Golang) project that demonstrates JWT authentication with secret rotation and MySQL database integration. The project includes user registration, login, token refresh, and protected routes.

## Prerequisites

- Go (Golang) 1.18+ installed
- MySQL 5.7+ installed
- Git installed
- Postman or any API testing tool (optional)

## Project Setup

Follow the steps below to get the project up and running on your local machine.

### 1. Clone the Repository

First, clone the project repository to your local machine:

```bash
git clone https://github.com/deswinar/Golang-Basic-JWT-Auth-API.git
cd Golang-Basic-JWT-Auth-API
```
2. Install Dependencies
Install the required Go modules:

```bash
go mod tidy
```
3. Create .env File
In the root directory of the project, create a .env file with the following content:

```bash
# JWT Secrets (for rotating secrets)
JWT_ACTIVE_SECRET=your-secret
JWT_OLD_SECRET_1=your-secret-1
JWT_OLD_SECRET_2=your-secret-2

# Database configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your-password
DB_NAME=your_db
DB_CHARSET=utf8mb4
DB_LOC=Local
```
Replace the placeholders with your actual values:

JWT Secrets: These are the secrets used for validating JWT tokens. The active secret (JWT_ACTIVE_SECRET) will be used for token generation and verification, while the old secrets are used for rotating tokens.
Database Configuration: These are the credentials and settings for connecting to your MySQL database. Ensure that the DB_NAME exists in your MySQL server.
4. Set Up MySQL Database
Make sure your MySQL database is up and running. You can create a new database using the following command:

```bash
CREATE DATABASE your_db;
```
5. Configure Database Connection
The project uses GORM (Go ORM) to interact with the MySQL database. The connection parameters are configured in the .env file.

Make sure to update the following values in the .env file:

DB_HOST: The hostname or IP address of your MySQL server.
DB_PORT: The port MySQL is running on (default is 3306).
DB_USER: The MySQL username (default is root).
DB_PASSWORD: The password for your MySQL user.
DB_NAME: The name of your database.
DB_CHARSET: The charset for your database (e.g., utf8mb4).
DB_LOC: The time zone setting for the database connection (e.g., Local).
6. Run the Application
Once you've set up the .env file and your database, you can run the application using the following command:

```bash
go run main.go
```
This will start the application on http://localhost:8080.

7. Testing the API
You can now test the API using Postman or any other API client.

Available Endpoints:
POST /register: Register a new user.
POST /login: Log in a user and receive a JWT token.
POST /refresh-token: Refresh an expired JWT token.
POST /logout: Log out the user.
GET /protected: A protected route that requires a valid JWT token.
GET /test: Test endpoint to verify the API is working.
Example Request to Login:
```bash
POST /login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "your-password"
}
```
Example Request to Access Protected Route:
```bash
GET /protected
Authorization: Bearer <your-jwt-token>
```
8. Secret Rotation
The project includes JWT secret rotation. By default, the active secret (JWT_ACTIVE_SECRET) is used for new tokens, and old secrets (JWT_OLD_SECRET_1, JWT_OLD_SECRET_2) are used to verify older tokens.

The active secret rotates every 30 days, ensuring that JWT tokens remain secure over time. You can modify the secret rotation duration in the auth.StartSecretRotation function.

Additional Notes
The project uses GORM for ORM and MySQL for database storage.
The JWT tokens are signed and verified using HS256 algorithm.
The project includes basic rate limiting for the login and refresh token routes.
License
This project is open source and available under the MIT License.

For any questions or issues, feel free to open an issue on the repository or contact the project maintainer.

### Explanation:

- **Step-by-step instructions** for setting up and running the project.
- **Environment configuration** in the `.env` file, with clear placeholders and descriptions for each key.
- Instructions for setting up the MySQL database.
- **API testing instructions**, including example API requests.
- **Secret rotation** explanation with default JWT configuration.

### Next Steps:

1. Copy this content into a `README.md` file.
2. Commit and push the file to your GitHub repository.

This README will provide clear instructions for setting up your project, especially for developers who are not familiar with it.