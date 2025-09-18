Setup password
```javascript
Rquest:
POST /api/setup/password
BODY
{
    "password": "xxx"
}

Response:
BODY
{
    "success": true,
    "message": "admin password set up"
}
```

Login as admin
```javascript
Request:
POST /api/login
BODY {
    "password": "xxx"
}

Response:
set cookies jwt token
BODY
{
    "success": true,
    "message": "xxx",
}
```