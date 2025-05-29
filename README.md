# ecommerce-order-normalizer-go

A Golang service to normalize and clean platform-specific e-commerce order items into a unified format.  
Built with Hexagonal Architecture. Includes full unit tests.

## 🧪 Unit Test

To run the unit tests for the order normalization logic:

```bash
# 1. Clone the project
git clone https://github.com/obobnilnil/ecommerce-order-normalizer-golang.git
cd ecommerce-order-normalizer-golang

# 2. Install dependencies
go mod tidy

# 3. (Optional) Run the service manually for testing
go run main.go

# 4. Run all unit tests: Ensure you're in the root directory of the project (where `go.mod` is located)
go test -v -count=1 ./internal/e-commerce/service
```
<p align="center"> <img src="https://github.com/user-attachments/assets/487a2082-c8f0-42a8-bba1-94dc6744baf5" alt="Terminal Output" /> <br /> Output from terminal showing successful test results</i> </p>

## 🧪 POSTMAN
📦 Sample POST Requests and Expected Responses<br>
## Case 1: Normal product (no prefix, no bundle)

Request
```json
[
  {
    "no": 1,
    "platformProductId": "FG0A-CLEAR-IPHONE16PROMAX",
    "qty": 2,
    "unitPrice": 50,
    "totalPrice": 100
  }
]
```
Response
```json
{
  "data": [
    {
      "no": 1,
      "productId": "FG0A-CLEAR-IPHONE16PROMAX",
      "materialId": "FG0A-CLEAR",
      "modelId": "IPHONE16PROMAX",
      "qty": 2,
      "unitPrice": 50,
      "totalPrice": 100
    },
    {
      "no": 2,
      "productId": "WIPING-CLOTH",
      "qty": 2,
      "unitPrice": 0,
      "totalPrice": 0
    },
    {
      "no": 3,
      "productId": "CLEAR-CLEANNER",
      "qty": 2,
      "unitPrice": 0,
      "totalPrice": 0
    }
  ],
  "message": "Order normalized successfully.",
  "status": "OK"
}
```
## Case 2: Product with prefix (x2-3&)
Request
```json
[
  {
    "no": 1,
    "platformProductId": "x2-3&FG0A-CLEAR-IPHONE16PROMAX",
    "qty": 2,
    "unitPrice": 50,
    "totalPrice": 100
  }
]
```
Response
```json
{
    "data": [
        {
            "no": 1,
            "productId": "FG0A-CLEAR-IPHONE16PROMAX",
            "materialId": "FG0A-CLEAR",
            "modelId": "IPHONE16PROMAX",
            "qty": 2,
            "unitPrice": 50,
            "totalPrice": 100
        },
        {
            "no": 2,
            "productId": "WIPING-CLOTH",
            "qty": 2,
            "unitPrice": 0,
            "totalPrice": 0
        },
        {
            "no": 3,
            "productId": "CLEAR-CLEANNER",
            "qty": 2,
            "unitPrice": 0,
            "totalPrice": 0
        }
    ],
    "message": "Order normalized successfully.",
    "status": "OK"
}
```
## Case 3: Prefixed + Multiplied Bundle *3
Request
```json
[
  {
    "no": 1,
    "platformProductId": "x2-3&FG0A-MATTE-IPHONE16PROMAX*3",
    "qty": 1,
    "unitPrice": 90,
    "totalPrice": 90
  }
]
```
Response
```json
{
  "data": [
    {
      "no": 1,
      "productId": "FG0A-MATTE-IPHONE16PROMAX",
      "materialId": "FG0A-MATTE",
      "modelId": "IPHONE16PROMAX",
      "qty": 3,
      "unitPrice": 30,
      "totalPrice": 90
    },
    {
      "no": 2,
      "productId": "WIPING-CLOTH",
      "qty": 3,
      "unitPrice": 0,
      "totalPrice": 0
    },
    {
      "no": 3,
      "productId": "MATTE-CLEANNER",
      "qty": 3,
      "unitPrice": 0,
      "totalPrice": 0
    }
  ],
  "message": "Order normalized successfully.",
  "status": "OK"
}
```
## Case 4: Bundle with 2 products using %20x
Request
```json
[
  {
    "no": 1,
    "platformProductId": "FG0A-CLEAR-OPPOA3/%20xFG0A-CLEAR-OPPOA3-B",
    "qty": 1,
    "unitPrice": 80,
    "totalPrice": 80
  }
]
```
Response
```json
{
  "data": [
    {
      "no": 1,
      "productId": "FG0A-CLEAR-OPPOA3",
      "materialId": "FG0A-CLEAR",
      "modelId": "OPPOA3",
      "qty": 1,
      "unitPrice": 40,
      "totalPrice": 40
    },
    {
      "no": 2,
      "productId": "FG0A-CLEAR-OPPOA3-B",
      "materialId": "FG0A-CLEAR",
      "modelId": "OPPOA3-B",
      "qty": 1,
      "unitPrice": 40,
      "totalPrice": 40
    },
    {
      "no": 3,
      "productId": "WIPING-CLOTH",
      "qty": 2,
      "unitPrice": 0,
      "totalPrice": 0
    },
    {
      "no": 4,
      "productId": "CLEAR-CLEANNER",
      "qty": 2,
      "unitPrice": 0,
      "totalPrice": 0
    }
  ],
  "message": "Order normalized successfully.",
  "status": "OK"
}
```
## Case 5: Triple bundle
Request
```json
[
  {
    "no": 1,
    "platformProductId": "FG0A-CLEAR-OPPOA3/%20xFG0A-CLEAR-OPPOA3-B/FG0A-MATTE-OPPOA3",
    "qty": 1,
    "unitPrice": 120,
    "totalPrice": 120
  }
]
```
Response
```json
{
  "data": [
    {
      "no": 1,
      "productId": "FG0A-CLEAR-OPPOA3",
      "materialId": "FG0A-CLEAR",
      "modelId": "OPPOA3",
      "qty": 1,
      "unitPrice": 40,
      "totalPrice": 40
    },
    {
      "no": 2,
      "productId": "FG0A-CLEAR-OPPOA3-B",
      "materialId": "FG0A-CLEAR",
      "modelId": "OPPOA3-B",
      "qty": 1,
      "unitPrice": 40,
      "totalPrice": 40
    },
    {
      "no": 3,
      "productId": "FG0A-MATTE-OPPOA3",
      "materialId": "FG0A-MATTE",
      "modelId": "OPPOA3",
      "qty": 1,
      "unitPrice": 40,
      "totalPrice": 40
    },
    {
      "no": 4,
      "productId": "WIPING-CLOTH",
      "qty": 3,
      "unitPrice": 0,
      "totalPrice": 0
    },
    {
      "no": 5,
      "productId": "CLEAR-CLEANNER",
      "qty": 2,
      "unitPrice": 0,
      "totalPrice": 0
    },
    {
      "no": 6,
      "productId": "MATTE-CLEANNER",
      "qty": 1,
      "unitPrice": 0,
      "totalPrice": 0
    }
  ],
  "message": "Order normalized successfully.",
  "status": "OK"
}
```
## Case 6: Bundle with prefix and multiple products
Request
```json
[
  {
    "no": 1,
    "platformProductId": "--FG0A-CLEAR-OPPOA3*2/FG0A-MATTE-OPPOA3",
    "qty": 1,
    "unitPrice": 120,
    "totalPrice": 120
  }
]
```
Response
```json
{
  "data": [
    {
      "no": 1,
      "productId": "FG0A-CLEAR-OPPOA3",
      "materialId": "FG0A-CLEAR",
      "modelId": "OPPOA3",
      "qty": 2,
      "unitPrice": 40,
      "totalPrice": 80
    },
    {
      "no": 2,
      "productId": "FG0A-MATTE-OPPOA3",
      "materialId": "FG0A-MATTE",
      "modelId": "OPPOA3",
      "qty": 1,
      "unitPrice": 40,
      "totalPrice": 40
    },
    {
      "no": 3,
      "productId": "WIPING-CLOTH",
      "qty": 3,
      "unitPrice": 0,
      "totalPrice": 0
    },
    {
      "no": 4,
      "productId": "CLEAR-CLEANNER",
      "qty": 2,
      "unitPrice": 0,
      "totalPrice": 0
    },
    {
      "no": 5,
      "productId": "MATTE-CLEANNER",
      "qty": 1,
      "unitPrice": 0,
      "totalPrice": 0
    }
  ],
  "message": "Order normalized successfully.",
  "status": "OK"
}
```
## Case 7: Mix of bundle and normal product
Request
```json
[
  {
    "no": 1,
    "platformProductId": "--FG0A-CLEAR-OPPOA3*2/FG0A-MATTE-OPPOA3*2",
    "qty": 1,
    "unitPrice": 160,
    "totalPrice": 160
  },
  {
    "no": 2,
    "platformProductId": "FG0A-PRIVACY-IPHONE16PROMAX",
    "qty": 1,
    "unitPrice": 50,
    "totalPrice": 50
  }
]
```
Response
```json
{
  "data": [
    {
      "no": 1,
      "productId": "FG0A-CLEAR-OPPOA3",
      "materialId": "FG0A-CLEAR",
      "modelId": "OPPOA3",
      "qty": 2,
      "unitPrice": 40,
      "totalPrice": 80
    },
    {
      "no": 2,
      "productId": "FG0A-MATTE-OPPOA3",
      "materialId": "FG0A-MATTE",
      "modelId": "OPPOA3",
      "qty": 2,
      "unitPrice": 40,
      "totalPrice": 80
    },
    {
      "no": 3,
      "productId": "FG0A-PRIVACY-IPHONE16PROMAX",
      "materialId": "FG0A-PRIVACY",
      "modelId": "IPHONE16PROMAX",
      "qty": 1,
      "unitPrice": 50,
      "totalPrice": 50
    },
    {
      "no": 4,
      "productId": "WIPING-CLOTH",
      "qty": 5,
      "unitPrice": 0,
      "totalPrice": 0
    },
    {
      "no": 5,
      "productId": "CLEAR-CLEANNER",
      "qty": 2,
      "unitPrice": 0,
      "totalPrice": 0
    },
    {
      "no": 6,
      "productId": "MATTE-CLEANNER",
      "qty": 2,
      "unitPrice": 0,
      "totalPrice": 0
    },
    {
      "no": 7,
      "productId": "PRIVACY-CLEANNER",
      "qty": 1,
      "unitPrice": 0,
      "totalPrice": 0
    }
  ],
  "message": "Order normalized successfully.",
  "status": "OK"
}
```
<p align="center"> <img src="https://github.com/user-attachments/assets/2613c2a6-61cb-4af1-80b8-e476d92c615d" alt="Postman Response Screenshot" /> <br />  Screenshot of Postman response showing normalized output</i> </p>



