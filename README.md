# RISA_RESTFULL

Restfull Backend for Risa Aplication Pelindo III using Golang (Fiber) and MongoDB

## Dependency lokal

- [ErruUtils](https://github.com/muchlist/erru_utils_go)  
  Library ini digunakan untuk memformat response error dan logger sehingga response error memiliki format yang standart di setiap service (berguna jika akan mengimplementasikan microservice).

## Dependency pihak ketiga

- [Go Fiber Framework](https://github.com/gofiber/fiber/v2) : Web framework golang yang memiliki kemiripan dengan express js dan menggunakan fast-http (tidak berbeda jauh dengan gin dan echo).
- [Mongo go driver](https://go.mongodb.org/mongo-driver) : Saat ini service ini full menggunakan MongoDB.
- [JWT go](https://github.com/dgrijalva/jwt-go)
- [Ozzo validation](github.com/go-ozzo/ozzo-validation/**v4) : Library yang digunakan untuk validasi request body dari user. (Karena Go Fiber tidak memiliki input validasi seperti Binding di Gin)