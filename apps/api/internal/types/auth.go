package types

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
	TOTPCode string `json:"totp_code,omitempty"`
}

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,min=3,max=50,alphanum"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required,min=1,max=50"`
	LastName  string `json:"last_name" validate:"required,min=1,max=50"`
}

type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	User    *User  `json:"user,omitempty"`
}

type PasswordChangeRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	TOTPCode        string `json:"totp_code,omitempty"`
}

type ProfileUpdateRequest struct {
	FirstName string `json:"first_name" validate:"required,min=1,max=50"`
	LastName  string `json:"last_name" validate:"required,min=1,max=50"`
	Bio       string `json:"bio" validate:"max=500"`
}

type UsernameUpdateRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
	TOTPCode string `json:"totp_code,omitempty"`
}

type EmailUpdateRequest struct {
	Email    string `json:"email" validate:"required,email"`
	TOTPCode string `json:"totp_code,omitempty"`
}

type NotificationPreferencesRequest struct {
	EmailNotifications     bool `json:"email_notifications"`
	MarketingEmails       bool `json:"marketing_emails"`
	SecurityNotifications bool `json:"security_notifications"`
}

type PrivacySettingsRequest struct {
	ProfilePublic bool `json:"profile_public"`
	ShowEmail     bool `json:"show_email"`
}

type TwoFactorRequest struct {
	Enable   bool   `json:"enable"`
	TOTPCode string `json:"totp_code,omitempty"`
}

type AccountDeletionRequest struct {
	Password     string `json:"password" validate:"required"`
	Confirmation string `json:"confirmation" validate:"required,eq=DELETE_MY_ACCOUNT"`
}