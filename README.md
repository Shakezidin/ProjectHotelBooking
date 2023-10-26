# Hotel Booking Application

This project is a hotel booking application written in Go Language. It utilizes the Gin HTTP web framework and PostgreSQL database for data management. The application offers features for hotel owners, administrators, and users, allowing hotel owners to add hotels and rooms, administrators to manage users and owners, and users to search for and book hotels.

## Frameworks and Libraries Used

- **Gin-Gonic:** The project is built on the Gin web framework, a popular choice for building web applications in Go.
```
go get -u github.com/gin-gonic/gin
```

- **PostgreSQL:** PostgreSQL, a powerful open-source relational database, is used for data management in the project.

- **GORM:** The GORM ORM tool is used to simplify database queries for better understanding. To install GORM and the PostgreSQL driver, use the following commands:
```
go get -u gorm.io/gorm
```
```
go get -u gorm.io/driver/postgres
```

- **Razorpay:** For payment processing, the project uses the Razorpay test case. To install the Razorpay library, use the following import statement:
```
github.com/razorpay/razorpay-go
```

- **Validator:** The package "validator" is used to implement value validations for structs and individual fields based on tags. To install Validator, use the following import statement:
```
github.com/go-playground/validator/v10
```

- **SMTP:** SMTP (Simple Mail Transfer Protocol) is used for sending emails. The SMTP settings are configured using environment variables (EMAIL and PASSWORD).

- **JWT:** JSON Web Tokens (JWT) are used for user authentication. To work with JWT, use the following import statement:
```
github.com/golang-jwt/jwt/v4
```

## How to Run the Project

To run the project, use the following command:

```bash
go run main.go

### API Documentation
```
https://documenter.getpostman.com/view/28796903/2s9YRFUpbF
```

## 👉 Signup as user 
  ### Endpoint :
  ```
  http://43.207.185.37:5000/user/signup
  ```  
  ### Method:
  `POST`
  
  ### Request Body:
  | Parameter     | Type    | Description              |
  |---------------|---------|--------------------------|
  | `username`    | string  | user  name of the user   |
  | `name`        | string  | name of the user         |
  | `email`       | string  | Email ID of the user     |
  | `password`    | string  | Password of the user     |
  | `phonenumber` | string  | Phone number of the user |
  
  ### Example Request:
  ```
   POST https://icrode-booking.online/user/signup
  -H "Content-Type: application/json" 
  -d '{
      "username" : "Tony_Stark",
      "name": "Stark",
      "email" : "tony@yopmail.com",
      "password" : "12345",
      "phonenumber": 1234567890
  }'
  ```
  
  ### Success Response:
  HTTP Code: `200 OK`
  
  ```
  {
    "status": "true", "messsage": "Go to user/signup/verification"
  }
  ```
  
## 👉 To verify the otp
  ### Endpoint :
  ```
  https://icrode-booking.online/user/signup/verification
  ```  
  ### Method:
  `POST`

   
  ### Request Body:
 ### Request Body:
  | Parameter   | Type     | Description       |
  |-------------|----------|-------------------|
  | `email`     | string   | Otp of the user   |
  | `otp`       | string   | Email of the user |
  
 
  ### Example Request:
  ```
   POST  https://icrode-booking.online/user/signup/verification
  -H "Content-Type: application/json" \
  -d '{
         "emai" : "1904",
         "otp"  : "tony@yopmail.com"
      }'
  ``` 
  ### Success Response:
  HTTP Code: `200 OK`
  
  ```
  {
    "status": "true", "message": "Otp Verification success. User creation done"
  }
  ```
  
## 👉 To login as a user
  ### Endpoint :
  ```
    https://icrode-booking.online/user/login
  ```  
  ### Method:
  `POST`
 
   ### Request Body:
  | Parameter     | Type    | Description              |
  |---------------|---------|--------------------------|
  | `username`    | string  | user name of the user    |
  | `password`    | string  | Password of the user     |
  
 ### Example Request:
  ```
   POST  https://icrode-booking.online/user/login
  -H "Content-Type: application/json" \
  -d '{
        "username" : "tony_Stark",
        "password" :"12345"
      }'
  ``` 

  ### Success Response:
  HTTP Code: `200 OK`
  
  ```
  {
    "loginStatus": "Success"
  }
  ```
 
