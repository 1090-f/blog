# API Documentation

## Scope

This document reflects the API surface that is currently implemented in the repository.
If older planning notes conflict with code, treat the code and this file as the latest baseline.

## Base

- Health check path: `/health`
- API base path: `/api`
- Authenticated requests use `Authorization: Bearer <token>`
- `POST /api/upload` is available only when both `upload.dir` and `upload.url` are configured
- `jwt.secret` in `config/config.yaml` must be replaced before production use
- Optional first-admin provisioning is controlled by `admin_bootstrap`

## Response Envelope

Success:

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

Error:

```json
{
  "code": 4001,
  "message": "invalid request params",
  "data": null
}
```

## Health

### GET `/health`

Service health check.

## Public APIs

### GET `/api/categories`

List all categories.

### GET `/api/tags`

List all tags. Tags are independent from categories and can be attached to multiple articles.

### GET `/api/site-stats`

Get public site statistics. Response data fields are `articleCount`, `categoryCount`, `tagCount`, `totalWords`, `firstPublishedAt`, and `lastActivityAt`.

### GET `/api/site-activity`

Get article activity grouped by date for the selected month.

Query:

- `year`: optional; defaults to the current year when omitted
- `month`: optional; defaults to the current month when omitted; valid values are `1` to `12`

### GET `/api/articles`

List published articles.

Query:

- `page`: optional, default `1`
- `pageSize`: optional, default `10`
- `categoryId`: optional
- `tagId`: optional
- `keyword`: optional, fuzzy match against article title

### GET `/api/articles/latest`

List latest published articles.

Query:

- `limit`: optional, default `10`, max `20`

### GET `/api/articles/popular`

List popular published articles ordered by `viewCount DESC, createdAt DESC`.

Query:

- `limit`: optional, default `10`, max `20`

### GET `/api/articles/:id`

Get published article detail and increase `viewCount`.

### GET `/api/articles/:id/full`

Get article aggregate detail, including:

- article fields, category, tags and author
- `comments`
- `commentCount`

Does not increase `viewCount`.

### GET `/api/articles/:id/comments`

List public comments for one published article.

### POST `/api/comments`

Create a comment for a published article. Authentication is optional:

- when a valid JWT is supplied, the comment is associated with the current user
- without a JWT, `guestName` and `guestEmail` are required; `guestWebsite` is optional and must use `http` or `https`

Guest request example:

```json
{
  "articleId": 1,
  "replyToId": 2,
  "content": "Nice article",
  "guestName": "Alice",
  "guestEmail": "alice@example.com",
  "guestWebsite": "https://example.com"
}
```

`replyToId` is optional. Guest email addresses are used only for comment records and are not returned by public comment APIs.

## Auth APIs

### POST `/api/auth/register`

Register a normal user.

Request:

```json
{
  "username": "alice",
  "password": "123456",
  "nickname": "Alice"
}
```

### POST `/api/auth/login`

Login and get JWT token with current user info.

Request:

```json
{
  "username": "alice",
  "password": "123456"
}
```

Response:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "jwt-token",
    "user": {
      "id": 1,
      "username": "alice",
      "nickname": "Alice",
      "role": "user",
      "avatar": "",
      "status": 1
    }
  }
}
```

### GET `/api/user/session`

Get the current authenticated user session information.

## Logged-In User APIs

These APIs require authentication but do not require admin role.

### DELETE `/api/comments/:id`

Delete a comment created by the current user.

### POST `/api/articles`

Create an article as the current logged-in user.

Request:

```json
{
  "title": "First Post",
  "summary": "Short summary",
  "content": "Markdown or HTML content",
  "coverImage": "/uploads/articles/2026/06/cover.jpg",
  "status": "published",
  "categoryId": 1,
  "tagIds": [1, 3]
}
```

Current behavior:

- normal logged-in users can create articles
- current implementation also allows normal logged-in users to publish articles

### POST `/api/categories`

Create a category as the current logged-in user.

Request:

```json
{
  "name": "Backend",
  "description": "Server side posts"
}
```

## Upload API

### POST `/api/upload`

Upload one image file using form field `file`.

Requires:

- valid JWT token
- upload storage configured

Validation:

- allowed extensions: `.jpg`, `.jpeg`, `.png`, `.webp`
- content type must be image/jpeg, image/png, or image/webp
- file size must not exceed `upload.max_size_bytes`

Default configuration in repository:

- `upload.max_size_bytes = 5242880` (5 MB)

Response:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "url": "/uploads/articles/2026/06/example.jpg"
  }
}
```

## Admin APIs

All admin APIs require:

- valid JWT token
- current user role is `admin`

### GET `/api/admin/dashboard`

Get aggregate admin dashboard statistics.

Response data fields:

- `articleCount`
- `publishedCount`
- `draftCount`
- `categoryCount`
- `commentCount`
- `userCount`
- `totalViews`

### GET `/api/admin/articles`

List articles for admin management.

Query:

- `page`: optional
- `pageSize`: optional
- `categoryId`: optional
- `tagId`: optional
- `status`: optional, `draft` or `published`
- `keyword`: optional, fuzzy match against title

### POST `/api/admin/articles`

Create an article.

### PUT `/api/admin/articles/:id`

Update an article.

### DELETE `/api/admin/articles/:id`

Delete an article.

Notes:

- deletion is rejected when comments still reference the article

### POST `/api/admin/categories`

Create a category.

### PUT `/api/admin/categories/:id`

Update a category.

### DELETE `/api/admin/categories/:id`

Delete a category.

Notes:

- deletion is rejected when articles still use the category

### POST `/api/admin/tags`

Create a tag.

Request:

```json
{
  "name": "Go"
}
```

### PUT `/api/admin/tags/:id`

Update a tag.

### DELETE `/api/admin/tags/:id`

Delete a tag. Deletion is rejected when the tag is still attached to an article.

### GET `/api/admin/comments`

List comments for admin management. Supports `page`, `pageSize`, `keyword`, `articleId`, and `status` filters.

### PUT `/api/admin/comments/:id/status`

Show or hide a comment with `{ "status": 1 }` or `{ "status": 0 }`.

### DELETE `/api/admin/comments/:id`

Delete a comment.

### GET `/api/admin/users`

List users for admin management.

Query:

- `page`: optional
- `pageSize`: optional
- `keyword`: optional, match username or nickname
- `role`: optional, `admin` or `user`
- `status`: optional, `0` or `1`

### PUT `/api/admin/users/:id/status`

Update user enabled status.

Request:

```json
{
  "status": 0
}
```

Rules:

- `1` means enabled
- `0` means disabled
- disabled users are blocked by auth middleware from protected APIs

## Notes

- Public article APIs only return `published` articles
- Public comment APIs only return comments with `status = 1`
- The public and management HTTP services both expose `/health` and the base API routes; only the management service exposes `/api/admin/*`
- Admin article queries support title keyword filtering
- Logger middleware records `method`, `path`, `status`, and `latency`
- Recovery middleware logs panic details and still returns the unified error envelope
