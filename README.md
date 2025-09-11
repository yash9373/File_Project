# File Vault Backend

This repository contains the backend for a secure file vault application. It provides a robust and scalable API for user authentication, encrypted file storage, and secure file sharing.

The backend is built with a focus on a clean architecture, separating concerns into distinct layers to promote maintainability, testability, and scalability.

## Technologies Used
- **Go (Golang):** The core programming language, chosen for its performance, concurrency features, and strong standard library.
- **Fiber:** A high-performance HTTP web framework for building the API endpoints.
- **MySQL:** The relational database used to store user information, file metadata, and sharing link details.
- **JWT (JSON Web Tokens):** Used for secure, stateless authentication and authorization.
- **AES-GCM:** An industry-standard authenticated encryption algorithm used to encrypt files.
- **Bcrypt:** A strong password-hashing algorithm to securely store user credentials.

## Backend Architecture
The backend follows a layered, service-oriented architecture, which ensures a clear separation of concerns.

- **main.go:** The application's entry point, responsible for initializing the database connection, setting up the router, and starting the server.
- **routes/:** Defines all the API endpoints and maps them to their respective controller functions.
- **middleware/:** Contains middleware functions, such as the auth middleware, which validates JWT tokens to protect routes.
- **controllers/:** Handles incoming HTTP requests. A controller's primary responsibility is to parse request data, call the appropriate service, and format the response.
- **services/:** Contains the core business logic of the application. This includes file encryption/decryption, user creation, and file sharing logic. Services act as a bridge between controllers and repositories.
- **repositories/:** An abstraction layer for database interactions. Each repository is responsible for all database operations for a specific model (e.g., user_repository handles all database queries related to the User model).
- **models/:** Defines the data structures (structs) that represent the data in the database, such as User, File, and ShareLink.

## How Encryption Works
For enhanced security, all files uploaded to the vault are encrypted at rest.

1. **File Upload:** When a user uploads a file, the backend receives the file data.
2. **Encryption:** The `crypto_service.go` component uses Advanced Encryption Standard (AES) in Galois/Counter Mode (GCM) to encrypt the file's binary data.
3. **Key:** A secure encryption key is used, which is configured in the `.env` file and must be kept secret.
4. **Initialization Vector (IV):** A unique, randomly generated IV is created for each file. This ensures that even if two identical files are uploaded, their encrypted contents will be different. The IV is stored along with the file metadata, but not the encryption key.
5. **Storage:** The encrypted file, along with its unique ID and the IV, is saved to the designated file storage location on the server. The original filename and other metadata are stored in the database.

During a file download, the encrypted data is retrieved and decrypted using the same encryption key and the stored IV, restoring the file to its original state before it's sent to the user.

## API Endpoints
The API is designed with RESTful principles, using standard HTTP methods for common actions.

| Method | Endpoint              | Description |
|--------|-----------------------|-------------|
| POST   | /auth/register        | Registers a new user. |
| POST   | /auth/login           | Authenticates a user and returns a JWT token. |
| POST   | /files/upload         | Uploads and encrypts a file. Requires authentication. |
| GET    | /files                | Retrieves a list of all files for the authenticated user. |
| GET    | /files/:id/download   | Downloads an encrypted file by its ID. Requires authentication. |
| DELETE | /files/:id            | Deletes a file. Requires authentication. |
| POST   | /share                | Creates a secure, shareable link for a file. |
| GET    | /share/:linkId        | Downloads a file using a public shareable link. No authentication required. |

## Getting Started

### Prerequisites
- Go (>= 1.18)
- MySQL database instance
- Git

### Setup Instructions
1. **Clone the Repository**
   ```bash
   git clone https://github.com/yash9373/file_project.git
   cd file_project
   ```

2. **Configure the Environment**  
   Create a `.env` file in the root directory and add the following variables:

   ```env
   DB_URI="user:password@tcp(127.0.0.1:3306)/database_name?charset=utf8mb4&parseTime=True&loc=Local"
   JWT_SECRET="your_secret_key"
   ENCRYPTION_KEY="a_32_byte_string_for_AES"
   ```

   Replace the placeholders with your actual database credentials and secrets.

3. **Run the Server**
   ```bash
   go run main.go
   ```

   The server will start on port `8080` (or the port defined in your configuration).

## Screenshots

### Login Page
![alt text](<Screenshot 2025-09-11 211051.png>)


### Dashboard
![alt text](<Screenshot 2025-09-11 211057.png>)

### Upload File
![alt text](<Screenshot 2025-09-11 212022.png>)

### Share Link Modal
![alt text](<Screenshot 2025-09-11 211115.png>)
![alt text](<Screenshot 2025-09-11 211104.png>)

---

Feel free to contact me with any questions or collaboration opportunities.
