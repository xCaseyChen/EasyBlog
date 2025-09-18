overview
```javascript
Request:
GET /api/admin/overview
BODY {}

Response:
BODY 
{
    "success": true,
    "message": "xxx",
    "data": {
        "published": 10,
        "tags": ["xxx", "yyy"],
        "categories": ["iii", "jjj"]
    }
}
```

posts list api
```javascript
Request
GET /api/admin/posts?page=1&status=published&tags=xxx,yyy&category=zzz

Response
BODY
{
    "success": true,
    "message": "find post briefs",
    "data": {
        "post_briefs": [
            {
                "id": 1,
                "title": "example",
                "slug": "example-1",
                "created_at": "time-1",
                "updated_at": "time-2",
                "category": "zzz",
                "tags": ["xxx", "yyy"],
                "status": "published",
                "pinned": true
            }
        ],
        "page": 1,
        "size": 10,
        "total": 42
    }
}
```

Post publish
```javascript
Request
POST /api/admin/post
BODY
{
    "title": "ttt",
    "category": "zzz",
    "tags": ["xxx", "yyy"],
    "content": "ccc",
    "status": "published",
    "pinned": true
}

Response
BODY
{
    "success": true,
    "message": "xxx"
}
```

post get api
```javascript
Request
GET /api/admin/post/:id

Response
BODY
{
    "success": true,
    "message": "xxx",
    "data": {
        "post": {
            "id": 1,
            "title": "example",
            "slug": "example-1",
            "created_at": "time-1",
            "updated_at": "time-2",
            "category": "zzz",
            "tags": ["xxx", "yyy"],
            "status": "published",
            "pinned": true,
            "content": "xxx"
        }
    }
}
```

Post modify api
```javascript
Request
PUT /api/admin/post/:id
BODY
{
    "title": "example",
    "category": "zzz",
    "tags": ["xxx", "yyy"],
    "status": "published",
    "pinned": true,
    "content": "xxx"
}

Response
BODY
{
    "success": true,
    "message": "xxx"
}
```

post delete api
```javascript
Request
DELETE /api/admin/post/:id
BODY {}

Response
BODY
{
    "success": true,
    "message": "xxx"
}
```

```javascript
GET PUT /api/admin/home
GET PUT /api/admin/about
```