Query post list
```javascript
Rquest:
GET /api/guest/posts?tags=xxx,yyy&category=zzz&before_id=iii&limit=jjj
BODY {}

Response:
BODY
{
    "success": true,
    "message": "find post briefs",
    "data": {
        "post_briefs": [
            {
                "id": 1,
                "title": "example",
                "slug": "example-1"
                "created_at": "time-1",
                "updated_at": "time-2",
                "category": "zzz",
                "tags": ["xxx", "yyy"]
            }
        ],
        "next_before_id": 20
    }
}
```