# Swagger UI Guide

## Accessing the Documentation

Once your server is running, navigate to:

```
http://localhost:8080/swagger/index.html
```

## What You'll See

### Main Page Structure

```
┌─────────────────────────────────────────────────────────────┐
│  BriefBot API v1.0                                          │
│  BriefBot backend API for managing content items,          │
│  podcasts, and daily digests                                │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ▼ users - User management operations                      │
│     POST   /users                Create a new user         │
│     GET    /users                List all users            │
│     GET    /users/{id}           Get a user by ID          │
│     GET    /users/email/{email}  Get a user by email       │
│     PUT    /users/{id}           Update a user             │
│     DELETE /users/{id}           Delete a user             │
│                                                             │
│  ▼ items - Content item operations with async processing   │
│     POST   /items                Create a new content item │
│     GET    /items/{id}           Get an item by ID         │
│     GET    /items/user/{userID}  Get items by user         │
│     ... (11 total endpoints)                               │
│                                                             │
│  ▼ podcasts - Podcast generation and management            │
│     POST   /podcasts             Create a new podcast      │
│     GET    /podcasts/{id}        Get a podcast by ID       │
│     ... (15 total endpoints)                               │
│                                                             │
│  ▼ digest - Daily digest email triggers                    │
│     POST   /digest/trigger       Trigger for all users     │
│     ... (4 total endpoints)                                │
│                                                             │
│  Schemas ▼                                                  │
│     View all request/response models                       │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

## Using the Interactive Features

### 1. Viewing an Endpoint

Click on any endpoint to expand it:

```
POST /items - Create a new content item
├─ Description: Create a new content item from URL with async processing
├─ Request Body (application/json)
│  └─ CreateItemRequest
│     {
│       "user_id": 1,
│       "url": "https://example.com/article"
│     }
├─ Responses
│  ├─ 201: Created
│  │  └─ CreateItemResponse
│  ├─ 400: Bad Request
│  │  └─ ErrorResponse
│  └─ 500: Internal Server Error
│     └─ ErrorResponse
└─ [Try it out] button
```

### 2. Testing an Endpoint

1. Click **"Try it out"** button
2. Edit the request body/parameters
3. Click **"Execute"**
4. See the response:
   - Status code
   - Response body
   - Response headers
   - cURL command equivalent

Example response:
```json
{
  "item": {
    "id": 123,
    "user_id": 1,
    "url": "https://example.com/article",
    "title": "Example Article",
    "processing_status": "pending",
    "created_at": "2025-10-03T16:00:00Z"
  },
  "message": "Item created successfully and will be processed in the background",
  "processing_status": "pending"
}
```

### 3. Viewing Schemas

Click on **"Schemas"** at the bottom to see all data models:

```
Schemas
├─ CreateItemRequest
│  ├─ user_id (integer, required)
│  └─ url (string, required)
├─ CreateItemResponse
│  ├─ item (Item object)
│  ├─ message (string)
│  └─ processing_status (string)
├─ Item
│  ├─ id (integer)
│  ├─ user_id (integer)
│  ├─ url (string)
│  ├─ title (string)
│  ├─ summary (string)
│  ├─ processing_status (string)
│  └─ ... (more fields)
└─ ... (all other models)
```

## Tips

### Quick Testing Workflow

1. **Create a user**: `POST /users`
2. **Create an item**: `POST /items` with the user_id
3. **Check status**: `GET /items/{id}/status`
4. **Get user's items**: `GET /items/user/{userID}`
5. **Create podcast**: `POST /podcasts` with item_ids
6. **Get podcast**: `GET /podcasts/{id}`

### Using Query Parameters

For endpoints with query parameters (e.g., `/items/status?status=pending`):
- The parameter fields appear in the UI
- Fill them in before clicking Execute
- They're automatically added to the URL

### Understanding Response Codes

- **200 OK**: Successful GET/PATCH
- **201 Created**: Successful POST
- **202 Accepted**: Async operation started
- **400 Bad Request**: Invalid input
- **404 Not Found**: Resource doesn't exist
- **500 Internal Server Error**: Server error
- **503 Service Unavailable**: Service not configured

### SSE Endpoints

The SSE streaming endpoints are documented but can't be tested in Swagger UI:
- `/items/user/{userID}/stream`
- `/podcasts/user/{userID}/stream`

Use a tool like `curl` or EventSource API to test these:
```bash
curl http://localhost:8080/items/user/1/stream
```

## Exporting API Spec

### For Client Generation

Download the OpenAPI spec:
- JSON: `http://localhost:8080/swagger/doc.json`
- YAML: `docs/swagger.yaml`

Use with tools like:
- **openapi-generator** - Generate client SDKs
- **Postman** - Import and test
- **Insomnia** - Import and test
- **Stoplight** - API design and documentation

### Example: Generate TypeScript Client

```bash
# Download the spec
curl http://localhost:8080/swagger/doc.json > api-spec.json

# Generate TypeScript client
openapi-generator-cli generate \
  -i api-spec.json \
  -g typescript-axios \
  -o ./client
```

## Troubleshooting

### Swagger UI Not Loading

1. Check server is running: `curl http://localhost:8080/`
2. Check Swagger route: `curl http://localhost:8080/swagger/doc.json`
3. Check browser console for errors
4. Verify `docs/` directory exists with files

### Endpoints Not Showing

1. Regenerate docs: `make swagger`
2. Restart server
3. Clear browser cache
4. Check for compilation errors

### Can't Execute Requests

1. Check CORS settings in server
2. Verify server is accessible from browser
3. Check network tab for errors
4. Try with curl first to verify endpoint works

## Advanced Features

### Authorization

If your API uses authentication (Bearer tokens, API keys):
1. Click **"Authorize"** button at top
2. Enter your credentials
3. All subsequent requests include auth headers

### Saving Requests

Swagger UI doesn't save requests, but you can:
1. Copy the generated cURL command
2. Save to a file or collection
3. Use Postman/Insomnia for persistent collections

### Customizing Examples

To change example values, edit the annotations in handler files:
```go
// @Param user_id query int true "User ID" example(123)
```

Then regenerate: `make swagger`

---

Happy API testing! 🚀
