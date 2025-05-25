# 🤖 Product Recommendation Engine Backend

This is a backend service for recommending beauty products based on user preferences and product ratings. It is built using **Python** and exposes a RESTful API that delivers recommendations using a machine learning model trained on Amazon's **beauty_products_rating.csv** dataset.

## ⚙️ Tech Stack

- **Language:** Python
- **Framework:** Flask
- **Data Source:** `beauty_products_rating.csv` (Amazon)
- **Machine Learning:** Scikit-learn

## 🚀 Setup Instructions

1. **Clone the Repository**
    ```bash
    git clone https://github.com/u-s1ddhar7h/internship-backend-tasks.git

    cd internship-backend-tasks/recommendation-engine
    ```

2. **Enter the Nix Dev Shell**

   ```bash
   nix develop
   ```

3. **Install Dependencies**

   ```bash
   pip install -r requirements.txt
   ```

4. **Run the Server**

   ```bash
   python app.py
   ```

    The server will start on `http://localhost:8000`

## 📌 API Endpoints

**🎯 Recommendations**

| Method | Endpoint                       | Description                                       |
| ------ | ------------------------------ | ------------------------------------------------- |
| GET    | `/recommend/<string:user_id>`  | Returns a list of recommended products for a user |

> **Note:** Ensure the dataset is preprocessed and the model is trained before calling the endpoints.

## 📂 Dataset

The dataset used is:
**`beauty_products_rating.csv`** — provided by Amazon, containing product ratings and review information for beauty-related products.

## 📬 Postman Collection

You can test the API using this Postman collection:

👉 [View Collection on Postman](https://www.postman.com/your-workspace/your-recommendation-collection-link)

---
