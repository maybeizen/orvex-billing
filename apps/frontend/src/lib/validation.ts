// Email validation
export const isValidEmail = (email: string): boolean => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
};

// Password validation
export const isValidPassword = (password: string): boolean => {
  return password.length >= 8;
};

// Username validation
export const isValidUsername = (username: string): boolean => {
  const usernameRegex = /^[a-zA-Z0-9_-]+$/;
  return username.length >= 3 && username.length <= 50 && usernameRegex.test(username);
};

// Name validation
export const isValidName = (name: string): boolean => {
  return name.trim().length >= 1 && name.trim().length <= 50;
};

// TOTP code validation
export const isValidTOTPCode = (code: string): boolean => {
  return /^\d{6}$/.test(code);
};

// Backup code validation
export const isValidBackupCode = (code: string): boolean => {
  return /^[a-f0-9]{8}$/.test(code.toLowerCase());
};

// Validation error messages
export const validationErrors = {
  email: {
    required: "Email is required",
    invalid: "Please enter a valid email address",
  },
  password: {
    required: "Password is required",
    minLength: "Password must be at least 8 characters",
    mismatch: "Passwords do not match",
  },
  username: {
    required: "Username is required",
    invalid: "Username can only contain letters, numbers, dashes, and underscores",
    minLength: "Username must be at least 3 characters",
    maxLength: "Username must be no more than 50 characters",
  },
  name: {
    required: "Name is required",
    minLength: "Name must be at least 1 character",
    maxLength: "Name must be no more than 50 characters",
  },
  totp: {
    required: "Verification code is required",
    invalid: "Verification code must be 6 digits",
  },
  backupCode: {
    required: "Backup code is required",
    invalid: "Backup code must be 8 characters",
  },
};

// Form field validation
export interface ValidationResult {
  isValid: boolean;
  error?: string;
}

export const validateField = (field: string, value: string, options?: any): ValidationResult => {
  if (!value || value.trim() === "") {
    return {
      isValid: false,
      error: validationErrors[field as keyof typeof validationErrors]?.required || "This field is required",
    };
  }

  switch (field) {
    case "email":
      return {
        isValid: isValidEmail(value),
        error: !isValidEmail(value) ? validationErrors.email.invalid : undefined,
      };
    
    case "password":
      return {
        isValid: isValidPassword(value),
        error: !isValidPassword(value) ? validationErrors.password.minLength : undefined,
      };
    
    case "confirmPassword":
      return {
        isValid: value === options?.password,
        error: value !== options?.password ? validationErrors.password.mismatch : undefined,
      };
    
    case "username":
      if (!isValidUsername(value)) {
        if (value.length < 3) {
          return { isValid: false, error: validationErrors.username.minLength };
        }
        if (value.length > 50) {
          return { isValid: false, error: validationErrors.username.maxLength };
        }
        return { isValid: false, error: validationErrors.username.invalid };
      }
      return { isValid: true };
    
    case "firstName":
    case "lastName":
      return {
        isValid: isValidName(value),
        error: !isValidName(value) ? 
          (value.length < 1 ? validationErrors.name.minLength : validationErrors.name.maxLength) : undefined,
      };
    
    case "totp":
      return {
        isValid: isValidTOTPCode(value),
        error: !isValidTOTPCode(value) ? validationErrors.totp.invalid : undefined,
      };
    
    case "backupCode":
      return {
        isValid: isValidBackupCode(value),
        error: !isValidBackupCode(value) ? validationErrors.backupCode.invalid : undefined,
      };
    
    default:
      return { isValid: true };
  }
};

// Validate entire form
export const validateForm = (formData: Record<string, any>, rules: Record<string, any>): Record<string, string> => {
  const errors: Record<string, string> = {};
  
  for (const [field, value] of Object.entries(formData)) {
    if (rules[field]) {
      const result = validateField(field, value, formData);
      if (!result.isValid && result.error) {
        errors[field] = result.error;
      }
    }
  }
  
  return errors;
};