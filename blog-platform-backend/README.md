# 📝 Blog Platform Backend

This is the backend for a blog platform featuring user authentication, blog creation, and commenting. The API is built with **Express.js** and uses **MySQL** as the database.

## ⚙️ Tech Stack

- **Runtime:** Node.js
- **Framework:** Express.js
- **Database:** MySQL
- **Authentication:** JWT-based
- **Environment Management:** dotenv

## 🚀 Setup Instructions

1. **Clone the Repository**

   ```bash
   git clone https://github.com/u-s1ddhar7h/internship-backend-tasks.git

   cd internship-backend-tasks/blog-platform-backend
   ```
2. **Enter the Nix-Dev Shell**

    ```bash
    nix develop
    ```

3. **Install Dependencies**

    ```bash
    npm install 
    ```

4. **Configure the Environment Variables**

    Create a `.env` file in the root directory:

    ```text
    DB_HOST=localhost
    DB_USER=your_db_user
    DB_PASSWORD=your_db_password
    DB_NAME=your_db_name
    JWT_SECRET=your_super_secret_32_characters_long_JWT_Token
    ```

5. **Run the Server**

    ```bash
    npm start
    ```

## 📌 API Endpoints
**🔐 Auth**

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/register` | Register a new user |
| `POST` | `/login` | Login with credentials |
| `POST` | `/me` | Get current user details (auth required) |

**📝 Blog**

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/` | Create a new blog post (auth required) |
| `GET` | `/` | Retrieve all blog posts |
| `GET` | `/:id` | Retrieve a specific blog post by ID |
| `GET` | `/user/my-blogs` | Get blogs created by the authenticated user |
| `PUT` | `/:id` | Update a blog post (auth required) |
| `DELETE` | `/:id` | Delete a blog post (auth required) |

**💬 Comment**

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/:blogId` | Add a comment to a blog post (auth required) |
| `GET` | `/:blogId` | Get all comments for a blog |
| `DELETE` | `/:commentId` | Delete a comment (auth required) |

## 📬 Postman Collection

You can test the API using this Postman collection:  
👉 [View Collection on Postman](https://api.postman.com/collections/44023205-f1ea9c5d-f063-4704-bfa9-6a093e6a3da5?access_key=PMAT-01JW3W7WDQGYCE7W5GTJ73RG1E)

---
