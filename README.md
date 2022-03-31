# Mock ATM 

A simple and lightwight Restful API using Golang with Fiber Framework and JWT authentication. 

## Description

This project simulate some common atm transactions suchs:
- customer logging in to an ATM by providing pin
- customer can view current balance
- customer can deposit money
- customer can widthraw money

## Getting Started

### Installing
- Make sure you have Go installed [download](https://go.dev/dl/). Version 1.14 or higher is required.

- Clone this repo-
- #### Add sent .env in same directory folder

### Executing program locally
-  Open a terminal
-  #### Make sure sent .env is added to the same directory folder (should be same level as other files like main.go and not inside a child folder)
-  Cd to project directory and run comand ```go run main.go```
-  After you run main, wait until compile (you'll know its done once you see the fiber table that display fiber and the host url will be visible)
-  Once its running, you can use the listed api below:
```
[GET] /health 
- returns 200 if server is running. A handy endpoint to test if main.go was run properly
```
```
[POST] /login 
[content-type] application/json
[Body] Sample:
{
    "card_number": 9999777755552222,
    "pin": 1212
}
```

```
[Get] /account/:type/balance 
[content-type] application/json
- Must be login to use this
```
```
[PUT] /account/transaction
[content-type] application/json
[Body] Sample:
{
    "card_number": 9999777755552222,
    "pin": 1212
}
- Must be login to use this
- sample acceptable account type(checking/saving)
```

# Test Notes
Test ATM account sample: 
```
card_number: 9999777755552222,
pin: 1212
has a `checking` and `saving`  
there's only 2 allowed balance update transaction:`widthrawal` and `deposit`
```
- Success auth is only valid for 10minutes and has no current functionality to refresh expired token
- this does have API for creating new account/closing account/ creating new user/logout/etc
Tip: if your using VSCODE have [plugin](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) for http rest use the http.rest for quick

