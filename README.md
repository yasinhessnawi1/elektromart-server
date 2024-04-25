# E-commerce API Documentation

## Overview

This RESTful API provides backend functionality for an e-commerce platform, supporting operations on products, users,
orders, and payments. It is built using Golang and leverages GORM for robust database interactions, designed for high
performance and scalability.

## External Dependencies

- **MySQL Database**: Used for persistent storage of all e-commerce data including users, products, and orders.
- **Gin-Gonic**: Utilized for efficient HTTP request routing and middleware handling.

## Third-Party Libraries

- **GORM**: The ORM library for Golang, used for database operations.
- **Gin-Gonic/Gin**: Simplifies the setup of HTTP routes and server.
- **GoDotEnv**: Manages configuration and environment variables from `.env` files.

## Setup and Installation

Ensure you have Golang and MySQL installed on your system. Follow these steps to set up the API:

1. Clone the repository:
   ```
   git clone https://yourrepositorylink.com/yourproject.git
    ```
2. Navigate into the project directory:
    ```
    cd "E-commerce Website Database"
     ```
3. Install the necessary Go packages:

```
go get .
```

4. Create a .env file in the root directory and populate it with your database credentials and other configurations:

```
DB_USER=root
DB_PASSWORD=password
DB_NAME=ecommerce
```

5. Run the application:

```
go run ./cmd/main.go
```

## API Endpoints

Below are detailed descriptions of the API endpoints including request methods, paths, expected requests, and responses.

### Products

**GET /products**: Retrieves all products.
**Response**: Status: 200 OK

**Body**:

```
[
    {
       "id":37948844, 
      "name": "Gaming Laptop",
      "description": "A high-end gaming laptop with the latest specs.",
      "price": 1500.00,
       "stock_quantity": 100,
       "brand_ID": 1,
       "category_ID": 1
    }
]
```

**POST /products**: Adds a new product to the catalog.
**Request Body**:

```
{
  "name": "Gaming Laptop",
  "description": "A high-end gaming laptop with the latest specs.",
  "price": 1500.00,
  "stock_quantity": 100,
  "brand_ID": 1486838590,
  "category_ID": 3170579900 
   }
```

**Response**: Status: 201 Created

**Body**:

```
{
    "id": 265789789,
    "message": "Product added successfully."
}
```

**DELETE /products/{id}**: Deletes a product by ID.
**Response**: Status: 200 OK

**Body**:

```
{
    "message": "Product deleted successfully."
}
```

**PUT /products/{id}**: Updates a product by ID.

**Request Body**:

```
{
  "name": "Gaming Laptop",
  "description": "A high-end gaming laptop with the latest specs.",
  "price": 1500.00,
  "stock_quantity": 100,
  "brand_ID": 1486838590,
  "category_ID": 3170579900 
   }
```

**Response**: Status: 200 OK

**Body**:

```
{
    "message": "Product updated successfully."
}
```

### Users

**GET /users**: Retrieves all registered users.
**Response**: Status: 200 OK

**Body**:

```
[
    {
        "id": 1764478339,
         "username": "newuser",
          "password": "securepassword",
          "email": "user@example.com",
          "first_name": "John",
          "last_name": "Doe",
           "address": "1234 Elm Street, Anytown, Anystate"
    }
]
```

**POST /users**: Registers a new user.
**Request Body**:

```
{
  "username": "newuserv",
  "password": "securepassword",
  "email": "user@exampsdle.com",
  "first_name": "John",
  "last_name": "Doe",
  "address": "1234 Elm Street, Anytown, Anystate"
}
```

**Response**: Status: 201 Created

**Body**:

```
{
    "id": 456672,
    "message": "User registered successfully."
}
```

**DELETE /users/{id}**: Deletes a user by ID.
**Response**: Status: 200 OK

**Body**:

```
{
    "message": "User deleted successfully."
}
```

**PUT /users/{id}**: Updates a user by ID.

**Request Body**:

```
{
  "username": "newuserv",
  "password": "securepassword",
  "email": "user@exampsdle.com",
  "first_name": "John",
  "last_name": "Doe",
  "address": "1234 Elm Street, Anytown, Anystate"
}
```

**Response**: Status: 200 OK

**Body**:

```
{
    "message": "User updated successfully."
}
```

### Orders

**GET /orders**: Retrieves all orders.
**Response**: Status: 200 OK

**Body**:

```
[
    {
       
        "id": 145678,
        "user_ID": 15678657,
        "order_date": "2024-04-20",
        "total_amount": 3000.00,
        "status": "Pending"
}

    ...
]
```

**POST /orders**: Creates a new order.

**Request Body**:

```
{
  "user_ID": 1352172511,
  "order_date": "2024-04-20",
  "total_amount": 3000.00,
  "status": "Pending"
}
```

**Response**: Status: 201 Created

**Body**:

```
{
    "id": 28746434,
    "message": "Order placed successfully."
}
```

**DELETE /orders/{id}**: Deletes an order by ID.
**Response**: Status: 200 OK

**Body**:

```
{
    "message": "Order deleted successfully."
}
```

**PUT /orders/{id}**: Updates an order by ID.

**Request Body**:

```
{
  "user_ID": 1352172511,
  "order_date": "2024-04-20",
  "total_amount": 3000.00,
  "status": "Pending"
}
```

**Response**: Status: 200 OK

**Body**:

```
{
    "message": "Order updated successfully."
}
```

### Payments

**GET /payments**: Retrieves all payments.
**Response**: Status: 200 OK

**Body**:

```
[
    {

        "id": 145678,
        "order_ID": 425451,
         "payment_method": "Credit Card",
         "amount": 3000.00,
         "payment_date": "2024-04-20",
         "status": "Completed"
}
    ...
]
```

**GET /payments/{id}**: Retrieves a payment by ID.
**Response**: Status: 200 OK

**Body**:

```
{
   "id": 145678,
  "order_ID": 435365361,
  "payment_method": "Credit Card",
  "amount": 3000.00,
  "payment_date": "2024-04-20",
  "status": "Completed"

}
```

**POST /payments**: Processes a payment for an order.

**Request Body**:

```
{
    "order_ID": 71599938,
  "payment_method": "Credit Card",
  "amount": 3000.00,
  "payment_date": "2024-04-20",
  "status": "Completed"
}
```

**Response**: Status: 200 OK

**Body**:

```
{
    "message": "Payment processed successfully."
}
```

**DELETE /payments/{id}**: Deletes a payment by ID.
**Response**: Status: 200 OK

**Body**:

```
{
    "message": "Payment deleted successfully."
}
```

**PUT /payments/{id}**: Updates a payment by ID.

**Request Body**:

```
{
    "order_ID": 71599938,
  "payment_method": "Credit Card",
  "amount": 3000.00,
  "payment_date": "2024-04-20",
  "status": "Completed"
}
```

**Response**: Status: 200 OK

**Body**:

```
{
    "message": "Payment updated successfully."
}
```

### orderItems

**GET /orderItems**: Retrieves all orderItems.
**Response**: Status: 200 OK

**Body**:

```
[
    {
        "id": 145678,
        "order_ID": 123454,
         "product_ID": 245452451,
         "quantity": 2,
         "subtotal": 3000.00
    }
    ...
]
```

**POST /orderItems**: Adds a new orderItem to the order.

**Request Body**:

```
{
  "order_ID": 71599938,
  "product_ID": 36259144,
  "quantity": 2,
  "subtotal": 3000.00
}
```

**Response**: Status: 201 Created

**Body**:

```
{
    "id": 2342432423,
    "message": "OrderItem added successfully."
}
```

**DELETE /orderItems/{id}**: Deletes an orderItem by ID.

**Response**: Status: 200 OK

**Body**:

```
{
    "message": "OrderItem deleted successfully."
}
```

**PUT /orderItems/{id}**: Updates an orderItem by ID.

**Request Body**:

```
{
  "order_ID": 71599938,
  "product_ID": 36259144,
  "quantity": 2,
  "subtotal": 3000.00
}
```

**Response**: Status: 200 OK

**Body**:

```
{
    "message": "OrderItem updated successfully."
}
```

### Category

**GET /category**: Retrieves all categories.
**Response**: Status: 200 OK

**Body**:

```
[
    {
        "id": 1,  
       "name": "Electronics",
       "description": "All electronic devices and gadgets."
    }
    ...
]
```

**POST /category**: Adds a new category to the catalog.

**Request Body**:

```
{
    "name": "Clothing",
    "description": "Clothing products"
}
```

**Response**: Status: 201 Created

**Body**:

```
{
    "id": 22342343432,
    "message": "Category added successfully."
}
```

**DELETE /category/{id}**: Deletes a category by ID.

**Response**: Status: 200 OK

**Body**:

```
{
    "message": "Category deleted successfully."
}
```

**PUT /category/{id}**: Updates a category by ID.

**Request Body**:

```
{
    "name": "Updated Electronics",
    "description": "Updated electronic products"
}
```

**Response**: Status: 200 OK

**Body**:

```
{
    "message": "Category updated successfully."
}
```

### brand

**GET /brand**: Retrieves all brands.

**Response**: Status: 200 OK

**Body**:

```
[
    {
        "id": 0,
        "brand_ID": 1,
        "name": "Apple",
        "description": "Apple products"
    }
    ...
]
```

**POST /brand**: Adds a new brand to the catalog.

**Request Body**:

```
{
    "name": "Samsung",
    "description": "Samsung products"
}
```

**Response**: Status: 201 Created

**Body**:

```
{
    "id": 2324343
    "message": "Brand added successfully."
}
```

**DELETE /brand/{id}**: Deletes a brand by ID.

**Response**: Status: 200 OK

**Body**:

```
{
    "message": "Brand deleted successfully."
}
```

**PUT /brand/{id}**: Updates a brand by ID.

**Request Body**:

```
{
    "name": "Updated Apple",
    "description": "Updated Apple products"
}
```

**Response**: Status: 200 OK

**Body**:

```
{
    "message": "Brand updated successfully."
}
```
