# 🛒 E-commerce Backend

This is the backend system for an e-commerce platform built using **Go (Gin)** and **MongoDB**. It includes features such as user authentication, product management, and order processing, all exposed via a RESTful API.

## ⚙️ Tech Stack

- **Language:** Go
- **Framework:** Gin
- **Database:** MongoDB
- **Authentication:** JWT-based
- **Environment Management:** GoDotEnv

## 🚀 Setup Instructions

1. **Clone the Repository**

    ```bash
    git clone https://github.com/u-s1ddhar7h/internship-backend-tasks.git

    cd internship-backend-tasks/ecommerce-backend
    ```

2. **Enter the Nix Dev Shell**

   ```bash
   nix develop
   ```

3. **Install Dependencies**

   ```bash
   go mod tidy
   ```

4. **Configure Environment Variables**

   Create a `.env` file in the root directory:
   ```env
   PORT=8080
   MONGO_URI=mongodb://localhost:27017
   DB_NAME=ecommerce
   JWT_SECRET=your_super_secret_32_characters_long_JWT_Token
   ```

5. **Run the Server**

   ```bash
   go run main.go
   ```

   The server will start on `http://localhost:8080`

## 📌 API Endpoints

### 🔐 Auth

| Method | Endpoint    | Description                  |
| ------ | ----------- | ---------------------------- |
| POST   | `/register` | Register a new user          |
| POST   | `/login`    | Authenticate and get a token |
| GET    | `/me`       | Get current user profile     |

### 📦 Products

| Method | Endpoint        | Description                          |
| ------ | --------------- | ------------------------------------ |
| POST   | `/products`     | Create a new product (auth required) |
| GET    | `/products`     | Get all products                     |
| GET    | `/products/:id` | Get product by ID                    |
| PUT    | `/products/:id` | Update product (auth required)       |
| DELETE | `/products/:id` | Delete product (auth required)       |

### 🧾 Orders

| Method | Endpoint      | Description                          |
| ------ | ------------- | ------------------------------------ |
| POST   | `/orders`     | Place a new order (auth required)    |
| GET    | `/orders`     | Get all orders of the logged-in user |
| GET    | `/orders/:id` | Get order details by ID              |

> **Note:** All endpoints that require authentication must include a JWT token in the `Authorization` header:
> `Authorization: Bearer <token>`

## 📬 Postman Collection

You can test the API using this Postman collection:

👉 [View Collection on Postman](https://api.postman.com/collections/44023205-f9943b02-9210-4468-a9f1-eb6cd8a0b78d?access_key=PMAT-01JW99BW3XN4HYZ6H5H6MTYXEM)

---
