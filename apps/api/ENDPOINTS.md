# API Endpoints Documentation

## Base URL
All endpoints are prefixed with `/api`

## Authentication
Most endpoints require session-based authentication. Include session cookies in requests.

---

## Health & System

### Health Check
- **Endpoint**: `GET /api/health`
- **Authentication**: None
- **Description**: Returns system health status, database connectivity, and performance metrics
- **Response**:
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T00:00:00Z",
  "version": "1.0.0",
  "uptime": "24h30m15s",
  "checks": {
    "database": {
      "status": "healthy",
      "message": "Database connection is working",
      "latency": "5ms"
    },
    "sessions": {
      "status": "healthy", 
      "message": "Session storage is working",
      "latency": "2ms"
    }
  },
  "system": {
    "go_version": "go1.21.0",
    "goroutines": 15,
    "memory_alloc_mb": 12,
    "memory_sys_mb": 25
  }
}
```

---

## Authentication Endpoints

### Register
- **Endpoint**: `POST /api/v1/auth/register`
- **Authentication**: None
- **Request Body**:
```json
{
  "email": "user@example.com",
  "username": "username123",
  "password": "securepassword123",
  "first_name": "John",
  "last_name": "Doe"
}
```
- **Validation**:
  - Email: Required, valid email format
  - Username: Required, 3-50 chars, alphanumeric only
  - Password: Required, minimum 8 characters
  - First/Last Name: Required, 1-50 characters

### Login
- **Endpoint**: `POST /api/v1/auth/login`
- **Authentication**: None
- **Request Body**:
```json
{
  "email": "user@example.com",
  "password": "securepassword123",
  "totp_code": "123456"
}
```
- **Note**: `totp_code` required only if 2FA is enabled

### Logout
- **Endpoint**: `POST /api/v1/auth/logout`
- **Authentication**: Required
- **Request Body**: None

### Logout All Sessions
- **Endpoint**: `POST /api/v1/auth/logout-all`
- **Authentication**: Required
- **Request Body**: None

---

## User Profile Management

### Get Profile
- **Endpoint**: `GET /api/v1/user/profile`
- **Authentication**: Required
- **Request Body**: None
- **Response**: Returns complete user profile (excludes sensitive fields)

### Update Profile
- **Endpoint**: `PUT /api/v1/user/profile`
- **Authentication**: Required
- **Request Body**:
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "bio": "Software developer passionate about security"
}
```
- **Validation**:
  - First/Last Name: Required, 1-50 characters
  - Bio: Optional, max 500 characters

### Update Username
- **Endpoint**: `PUT /api/v1/user/username`
- **Authentication**: Required
- **Request Body**:
```json
{
  "username": "newusername123"
}
```
- **Validation**: Required, 3-50 chars, alphanumeric only, must be unique

### Update Email
- **Endpoint**: `PUT /api/v1/user/email`
- **Authentication**: Required
- **Request Body**:
```json
{
  "email": "newemail@example.com"
}
```
- **Validation**: Required, valid email format, must be unique
- **Note**: Resets email verification status

### Change Password
- **Endpoint**: `PUT /api/v1/user/password`
- **Authentication**: Required
- **Request Body**:
```json
{
  "current_password": "oldpassword123",
  "new_password": "newpassword123"
}
```
- **Validation**: 
  - Current password: Required, must match existing
  - New password: Required, minimum 8 characters
- **Note**: Invalidates all other user sessions

---

## Avatar Management

### Upload Avatar
- **Endpoint**: `POST /api/v1/user/avatar`
- **Authentication**: Required
- **Content-Type**: `multipart/form-data`
- **Form Field**: `avatar` (file)
- **File Requirements**:
  - Formats: PNG, JPG, JPEG, WebP only
  - Max size: 5MB
  - Security: Content-type validation, image format verification
- **Response**:
```json
{
  "success": true,
  "message": "Avatar uploaded successfully",
  "avatar_url": "/uploads/avatars/123_abc123.png"
}
```

### Delete Avatar
- **Endpoint**: `DELETE /api/v1/user/avatar`
- **Authentication**: Required
- **Request Body**: None

---

## User Preferences

### Update Notification Preferences
- **Endpoint**: `PUT /api/v1/user/notifications`
- **Authentication**: Required
- **Request Body**:
```json
{
  "email_notifications": true,
  "marketing_emails": false,
  "security_notifications": true
}
```

### Update Privacy Settings
- **Endpoint**: `PUT /api/v1/user/privacy`
- **Authentication**: Required
- **Request Body**:
```json
{
  "profile_public": false,
  "show_email": false
}
```

---

## Two-Factor Authentication

### Enable 2FA (Step 1)
- **Endpoint**: `POST /api/v1/user/2fa/enable`
- **Authentication**: Required
- **Request Body**: None
- **Response**:
```json
{
  "success": true,
  "message": "Scan this QR code with your authenticator app",
  "qr_url": "otpauth://totp/YourApp:user@example.com?secret=ABC123&issuer=YourApp",
  "secret": "ABC123DEFG456",
  "note": "Save this secret safely. You'll need to verify with a TOTP code to complete setup."
}
```

### Confirm 2FA (Step 2)
- **Endpoint**: `POST /api/v1/user/2fa/confirm`
- **Authentication**: Required
- **Request Body**:
```json
{
  "secret": "ABC123DEFG456",
  "totp_code": "123456"
}
```
- **Validation**: TOTP code must be exactly 6 digits

### Disable 2FA
- **Endpoint**: `POST /api/v1/user/2fa/disable`
- **Authentication**: Required
- **Request Body**:
```json
{
  "totp_code": "123456"
}
```
- **Validation**: TOTP code must be exactly 6 digits

---

## Account Management

### Delete Account
- **Endpoint**: `DELETE /api/v1/user/account`
- **Authentication**: Required
- **Request Body**:
```json
{
  "password": "currentpassword123",
  "confirmation": "DELETE_MY_ACCOUNT"
}
```
- **Validation**:
  - Password: Required, must match current password
  - Confirmation: Must be exactly "DELETE_MY_ACCOUNT"
- **Note**: Permanently deletes account, avatar, and all sessions

---

## Static Files

### Avatar Files
- **Endpoint**: `GET /uploads/avatars/{filename}`
- **Authentication**: None
- **Description**: Serves uploaded avatar images

---

## Other Endpoints (Placeholder)

### Announcements
- **Endpoint**: `GET /api/announcements`
- **Authentication**: None
- **Description**: Public announcements list

- **Endpoint**: `GET /api/announcements/{id}`
- **Authentication**: None
- **Description**: Get specific announcement

### Cart (Protected)
- **Endpoint**: `GET /api/cart`
- **Authentication**: Required
- **Description**: Get user's cart

- **Endpoint**: `POST /api/cart/checkout`
- **Authentication**: Required
- **Description**: Process checkout

- **Endpoint**: `POST /api/cart/coupons`
- **Authentication**: Required
- **Description**: Apply coupon to cart

### Services (Protected)
- **Endpoint**: `GET /api/services`
- **Authentication**: Required
- **Description**: List user's services

- **Endpoint**: `POST /api/services`
- **Authentication**: Required
- **Description**: Create new service

- **Endpoint**: `GET /api/services/invoices`
- **Authentication**: Required
- **Description**: Get service invoices

### Admin (Protected)
- **Endpoint**: `GET /api/v1/admin/dashboard`
- **Authentication**: Required (Admin role)
- **Description**: Admin dashboard

---

## Error Responses

All endpoints return errors in this format:
```json
{
  "error": "Error message description",
  "success": false
}
```

Common HTTP status codes:
- `200` - Success
- `201` - Created
- `400` - Bad Request (validation errors)
- `401` - Unauthorized (authentication required)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found
- `409` - Conflict (duplicate data)
- `500` - Internal Server Error
- `503` - Service Unavailable (health check failure)

## Security Features

- **Session-based authentication** (not JWT)
- **Bcrypt password hashing** with cost factor 12
- **Account lockout** after 5 failed login attempts
- **Secure file upload validation** with content-type detection
- **CSRF protection ready** with secure cookies
- **Input validation** on all endpoints
- **Rate limiting support** (IP-based)
- **Soft deletes** with GORM
- **Session cleanup** on password changes and account deletion