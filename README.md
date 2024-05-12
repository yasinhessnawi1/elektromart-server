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

##### If running sql server from localhost:

1. To run sql database locally first install [xamp](https://www.apachefriends.org/)

2. Then run xamp manager (it may ask for root/admin permission) and start MYSQL Database and Apache Web Server.
3. Make a new database using the sql file provided in this wiki: ```https://gitlab.stud.idi.ntnu.no/yasinmh/e-commerce-website-database/-/wikis/home```
4. Note!!: you can clone the wiki using this command: 
   ``` 
   git clone git@gitlab.stud.idi.ntnu.no:yasinmh/e-commerce-website-database.wiki.git
   cd e-commerce-website-database.wiki
   ```
5. note!!: keep ports used as well as database user, password and name ready because you will use them in the .evn file


##### Common steps for running locally or using deployed server

Ensure you have Golang and MySQL installed on your system. Follow these steps to set up the API:

1. Clone the repository:
   ```
   git clone git@gitlab.stud.idi.ntnu.no:yasinmh/e-commerce-website-database.git
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
DATABASE_PORT=8000 (or your own database port)(found in the xamp in sql config)
PORT=8081 (or another port you wish the program to handle from)
DB_HOST=localhost:{port} (if you want to use the running local client)
DB_PORT=3306 (port of the database default is 3306)
DB_USER={user} (specify user used)
DB_PASSWORD={password} (specify password of the database)
DB_NAME={name} (name of the database)
```
Note that to run using the deployed server you need only configure 'PORT' all other values must remain unchanged.

5. Run the application:

```
go run ./cmd/main.go
```

## API Endpoints

Below are detailed descriptions of the API endpoints including request methods, paths, expected requests, and responses.
1. note!!: most of the id's are made up just to give an example,
   to be able to actually make a post, delete or put request, 
   foreign keys id's needs to be retried and used in the requests. 
2. Please consider using the postman file provided in the wiki page for testing.(see setup and installation part to 
for command to download the file from wiki)
3. All the requests are made based on the local running of the code, to test the backend deployment please see 
the link provided in the report and replace is with the "http//:localhost:8081"

### Products

**GET /products**: Retrieves all products.
```
http://localhost:8081/products
```
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
```
http://localhost:8081/products
```
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
```
http://localhost:8081/products/{id}
```
**Response**: Status: 200 OK

**Body**:

```
{
    "message": "Product deleted successfully."
}
```

**PUT /products/{id}**: Updates a product by ID.
```
http://localhost:8081/products/{id}
```
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
```
http://localhost:8081/users/{id}
http://localhost:8081/users
```
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
```
http://localhost:8081/users/{id}
```
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
```
http://localhost:8081/users/{id}
```
**Response**: Status: 200 OK

**Body**:

```
{
    "message": "User deleted successfully."
}
```

**PUT /users/{id}**: Updates a user by ID.
```
http://localhost:8081/users/{id}
```
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
```
http://localhost:8081/orders/{id}
http://localhost:8081/orders
```
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
```
http://localhost:8081/orders/{id}
```
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
```
http://localhost:8081/orders/{id}
```
**Response**: Status: 200 OK

**Body**:

```
{
    "message": "Order deleted successfully."
}
```

**PUT /orders/{id}**: Updates an order by ID.
```
http://localhost:8081/orders/{id}
```
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
```
http://localhost:8081/payments/{id}
http://localhost:8081/payments
```

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
```
http://localhost:8081/payments/{id}
```
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
```
http://localhost:8081/payments/{id}
```
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
```
http://localhost:8081/payments/{id}
```
**Response**: Status: 200 OK

**Body**:

```
{
    "message": "Payment deleted successfully."
}
```

**PUT /payments/{id}**: Updates a payment by ID.
```
http://localhost:8081/payments/{id}
```
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
```
http://localhost:8081/orderItems/{id}
http://localhost:8081/orderItems/
```
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
```
http://localhost:8081/orderItems/{id}
```
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
```
http://localhost:8081/orderItems/{id}
```
**Response**: Status: 200 OK

**Body**:

```
{
    "message": "OrderItem deleted successfully."
}
```

**PUT /orderItems/{id}**: Updates an orderItem by ID.
```
http://localhost:8081/orderItems/{id}
```
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

**GET /categories**: Retrieves all categories.
```
http://localhost:8081/categories/{id}
http://localhost:8081/categories
```
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

**POST /categories**: Adds a new category to the catalog.
```
http://localhost:8081/categories/{id}
```
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

**DELETE /categories/{id}**: Deletes a category by ID.
```
http://localhost:8081/categories/{id}
```

**Response**: Status: 200 OK

**Body**:

```
{
    "message": "Category deleted successfully."
}
```

**PUT /categories/{id}**: Updates a category by ID.
```
http://localhost:8081/categories/{id}
```

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
```
http://localhost:8081/brand/{id}
http://localhost:8081/brand
```

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
```
http://localhost:8081/brand/{id}
```
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
```
http://localhost:8081/brand/{id}
```
**Response**: Status: 200 OK

**Body**:

```
{
    "message": "Brand deleted successfully."
}
```

**PUT /brand/{id}**: Updates a brand by ID.
```
http://localhost:8081/brand/{id}
```
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
### review

**GET /reviews**: Retrieves all reviews.
```
http://localhost:8081/reviews/{id}
http://localhost:8081/reviews
```

**Response**: Status: 200 OK

**Body**:

```
[
    {
        "ID": 1516010785,
        "CreatedAt": "2024-04-27T07:57:40Z",
        "UpdatedAt": "2024-04-27T07:57:40Z",
        "DeletedAt": null,
        "product_id": 1192428622,
        "user_id": 4265253256,
        "rating": 0,
        "comment": "this is a bad product",
        "review_date": "2025-04-20T00:00:00Z"
    }
    ...
]
```

**POST /reviews**: Adds a new review to the catalog.
```
http://localhost:8081/reviews/{id}
```
**Request Body**:

```
{
        "product_id": 1192228322,
        "user_id": 4261254256,
        "rating": 2,
        "comment": "this is a new test review",
        "review_date": "2025-04-20"
    }
```

**Response**: Status: 201 Created

**Body**:

```
{
    "ID": 1402172417,
    "CreatedAt": "2024-05-12T10:36:14.081+02:00",
    "UpdatedAt": "2024-05-12T10:36:14.081+02:00",
    "DeletedAt": null,
    "product_id": 1192128622,
    "user_id": 4265255256,
    "rating": 2,
    "comment": "this is a new test review",
    "review_date": "2025-04-20"
}
```

**DELETE /reviews/{id}**: Deletes a reviews by ID.
```
http://localhost:8081/reviews/{id}
```
**Response**: Status: 200 OK

**Body**:

```
{
    "message": "reviews deleted successfully."
}
```

**PUT /reviews/{id}**: Updates a reviews by ID.
```
http://localhost:8081/reviews/{id}
```
**Request Body**:

```
{
    
    "ID": 1402172417,
    "CreatedAt": "2024-05-12T10:36:14.081+02:00",
    "UpdatedAt": "2024-05-12T10:36:14.081+02:00",
    "DeletedAt": null,
    "product_id": 1192128622,
    "user_id": 4265255256,
    "rating": 4, -> this value to be changed.
    "comment": "this is a new test review",
    "review_date": "2025-04-20"

}
```

**Response**: Status: 200 OK

**Body**:

```
{
    "message": "reviews updated successfully."
}
```
### shippingDetails

**GET /shippingDetails**: Retrieves all shippingDetails.
```
http://localhost:8081/shippingDetails/{id}
http://localhost:8081/shippingDetails
```

**Response**: Status: 200 OK

**Body**:

```
[
    {
        "ID": 1359557427,
        "CreatedAt": "2024-04-27T05:48:35Z",
        "UpdatedAt": "2024-04-27T05:48:35Z",
        "DeletedAt": null,
        "order_id": 1350357487,
        "address": "In any where",
        "shipping_date": "2025-04-20T00:00:00Z",
        "estimated_arrival": "2026-04-30T00:00:00Z",
        "status": "completed"
    }
    ...
]
```

**POST /shippingDetails**: Adds a new shippingDetails to the catalog.
```
http://localhost:8081/shippingDetails/{id}
```
**Request Body**:

```
{
        "order_id": 1352357487,
        "address": "test adress",
        "shipping_date": "2025-04-20",
        "estimated_arrival": "2026-04-30",
        "status": "completed"
    }
```

**Response**: Status: 201 Created

**Body**:

```
{
    "ID": 1828821868,
    "CreatedAt": "2024-05-12T10:43:49.148+02:00",
    "UpdatedAt": "2024-05-12T10:43:49.148+02:00",
    "DeletedAt": null,
    "order_id": 1359350487,
    "address": "test adress",
    "shipping_date": "2025-04-20",
    "estimated_arrival": "2026-04-30",
    "status": "completed"
}
```

**DELETE /shippingDetails/{id}**: Deletes a shippingDetails by ID.
```
http://localhost:8081/shippingDetails/{id}
```
**Response**: Status: 200 OK

**Body**:

```
{
    "message": "shippingDetails deleted successfully."
}
```

**PUT /shippingDetails/{id}**: Updates a shippingDetails by ID.
```
http://localhost:8081/shippingDetails/{id}
```
**Request Body**:

```
{
        "order_id": 1359357487,
        "address": "new test adress", -> this value to be updated
        "shipping_date": "2025-04-20",
        "estimated_arrival": "2026-04-30",
        "status": "completed"
    }
```

**Response**: Status: 200 OK

**Body**:

```
{
    "message": "shippingDetails updated successfully."
}
```

### searching
- This part explains how the searching mechanism works on the different endpoints:

| Resource          | Operation            | Endpoint                                                      |
|-------------------|----------------------|---------------------------------------------------------------|
| Users             | Search Users         | `GET http://localhost:8081/search-users/?username={name}`     |
|                   |                      | `or /?email={email}`                                          |
|                   |                      | `or /?firstname={firstname}&lastname={lastname}`              |
|                   |                      | `or /?address={address}`                                      |
| Shipping Details  | Search Shipping      | `GET http://localhost:8081/search-shippingDetails/?order_id={id}` |
|                   | Details             | `or /?address={address}`                                      |
|                   |                      | `or /?status={status}`                                        |
| Reviews           | Search Reviews       | `GET http://localhost:8081/search-reviews/?product_id={id}`   |
|                   |                      | `or /?comment={comment}`                                      |
|                   |                      | `or /?rating={rating}&user_id={user_id}&review_date={date}`   |
| Products          | Search Products      | `GET http://localhost:8081/search-products/?name={name}`      |
|                   |                      | `or /?price={price}`                                          |
|                   |                      | `or /?brand_name={brand}&category_name={category}`            |
| Brands            | Search Brands        | `GET http://localhost:8081/search-brands/?name={name}`        |
|                   |                      | `or /?description={description}`                              |
| Categories        | Search Categories    | `GET http://localhost:8081/search-categories/?name={name}`    |
|                   |                      | `or /?description={description}`                              |
| Orders            | Search Orders        | `GET http://localhost:8081/search-orders/?user_id={id}`       |
|                   |                      | `or /?total_amount={amount}`                                  |
|                   |                      | `or /?status={status}`                                        |
| Order Items       | Search Order Items   | `GET http://localhost:8081/search-orderItems/?order_id={id}`  |
|                   |                      | `or /?quantity={quantity}`                                    |
|                   |                      | `or /?product_id={product_id}`                                |
| Payments          | Search Payments      | `GET http://localhost:8081/search-payments/?payment_method={method}` |
|                   |                      | `or /?amount={amount}`                                        |
|                   |                      | `or /?order_id={order_id}`                                    |


