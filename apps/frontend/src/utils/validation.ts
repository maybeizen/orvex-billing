export interface ValidationResult {
  isValid: boolean;
  error?: string;
}

export interface FormValidationResult {
  isValid: boolean;
  errors: Record<string, string>;
}

// Email validation
export const validateEmail = (email: string): ValidationResult => {
  if (!email) {
    return { isValid: false, error: "Email is required" };
  }

  if (email.length > 254) {
    return { isValid: false, error: "Email is too long" };
  }

  const emailRegex = /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/;
  
  if (!emailRegex.test(email)) {
    return { isValid: false, error: "Please enter a valid email address" };
  }

  return { isValid: true };
};

// Password validation
export const validatePassword = (password: string): ValidationResult => {
  if (!password) {
    return { isValid: false, error: "Password is required" };
  }

  if (password.length < 8) {
    return { isValid: false, error: "Password must be at least 8 characters long" };
  }

  if (password.length > 128) {
    return { isValid: false, error: "Password is too long" };
  }

  if (!/[a-z]/.test(password)) {
    return { isValid: false, error: "Password must contain at least one lowercase letter" };
  }

  if (!/[A-Z]/.test(password)) {
    return { isValid: false, error: "Password must contain at least one uppercase letter" };
  }

  if (!/\d/.test(password)) {
    return { isValid: false, error: "Password must contain at least one number" };
  }

  if (!/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(password)) {
    return { isValid: false, error: "Password must contain at least one special character" };
  }

  return { isValid: true };
};

// Username validation
export const validateUsername = (username: string): ValidationResult => {
  if (!username) {
    return { isValid: false, error: "Username is required" };
  }

  if (username.length < 3) {
    return { isValid: false, error: "Username must be at least 3 characters long" };
  }

  if (username.length > 30) {
    return { isValid: false, error: "Username must be less than 30 characters" };
  }

  if (!/^[a-zA-Z0-9_-]+$/.test(username)) {
    return { isValid: false, error: "Username can only contain letters, numbers, hyphens, and underscores" };
  }

  if (/^[-_]/.test(username) || /[-_]$/.test(username)) {
    return { isValid: false, error: "Username cannot start or end with hyphens or underscores" };
  }

  if (/[-_]{2,}/.test(username)) {
    return { isValid: false, error: "Username cannot contain consecutive hyphens or underscores" };
  }

  return { isValid: true };
};

// Name validation
export const validateName = (name: string, fieldName: string = "Name"): ValidationResult => {
  if (!name) {
    return { isValid: false, error: `${fieldName} is required` };
  }

  if (name.length < 1) {
    return { isValid: false, error: `${fieldName} is required` };
  }

  if (name.length > 50) {
    return { isValid: false, error: `${fieldName} must be less than 50 characters` };
  }

  if (!/^[a-zA-Z\s'-]+$/.test(name)) {
    return { isValid: false, error: `${fieldName} can only contain letters, spaces, hyphens, and apostrophes` };
  }

  if (/\s{2,}/.test(name)) {
    return { isValid: false, error: `${fieldName} cannot contain consecutive spaces` };
  }

  if (/^[\s'-]/.test(name) || /[\s'-]$/.test(name)) {
    return { isValid: false, error: `${fieldName} cannot start or end with spaces or special characters` };
  }

  return { isValid: true };
};

// Confirm password validation
export const validateConfirmPassword = (password: string, confirmPassword: string): ValidationResult => {
  if (!confirmPassword) {
    return { isValid: false, error: "Please confirm your password" };
  }

  if (password !== confirmPassword) {
    return { isValid: false, error: "Passwords do not match" };
  }

  return { isValid: true };
};

// Login form validation
export const validateLoginForm = (email: string, password: string): FormValidationResult => {
  const errors: Record<string, string> = {};

  const emailValidation = validateEmail(email);
  if (!emailValidation.isValid) {
    errors.email = emailValidation.error!;
  }

  if (!password) {
    errors.password = "Password is required";
  }

  return {
    isValid: Object.keys(errors).length === 0,
    errors,
  };
};

// Registration form validation
export const validateRegistrationForm = (data: {
  firstName: string;
  lastName: string;
  username: string;
  email: string;
  password: string;
  confirmPassword: string;
}): FormValidationResult => {
  const errors: Record<string, string> = {};

  const firstNameValidation = validateName(data.firstName, "First name");
  if (!firstNameValidation.isValid) {
    errors.firstName = firstNameValidation.error!;
  }

  const lastNameValidation = validateName(data.lastName, "Last name");
  if (!lastNameValidation.isValid) {
    errors.lastName = lastNameValidation.error!;
  }

  const usernameValidation = validateUsername(data.username);
  if (!usernameValidation.isValid) {
    errors.username = usernameValidation.error!;
  }

  const emailValidation = validateEmail(data.email);
  if (!emailValidation.isValid) {
    errors.email = emailValidation.error!;
  }

  const passwordValidation = validatePassword(data.password);
  if (!passwordValidation.isValid) {
    errors.password = passwordValidation.error!;
  }

  const confirmPasswordValidation = validateConfirmPassword(data.password, data.confirmPassword);
  if (!confirmPasswordValidation.isValid) {
    errors.confirmPassword = confirmPasswordValidation.error!;
  }

  return {
    isValid: Object.keys(errors).length === 0,
    errors,
  };
};

// Real-time validation hook
export const useFormValidation = () => {
  const validateField = (fieldName: string, value: string, additionalValue?: string): string | null => {
    switch (fieldName) {
      case "email":
        const emailResult = validateEmail(value);
        return emailResult.isValid ? null : emailResult.error!;
      
      case "password":
        const passwordResult = validatePassword(value);
        return passwordResult.isValid ? null : passwordResult.error!;
      
      case "confirmPassword":
        if (!additionalValue) return "Password is required to confirm";
        const confirmResult = validateConfirmPassword(additionalValue, value);
        return confirmResult.isValid ? null : confirmResult.error!;
      
      case "username":
        const usernameResult = validateUsername(value);
        return usernameResult.isValid ? null : usernameResult.error!;
      
      case "firstName":
        const firstNameResult = validateName(value, "First name");
        return firstNameResult.isValid ? null : firstNameResult.error!;
      
      case "lastName":
        const lastNameResult = validateName(value, "Last name");
        return lastNameResult.isValid ? null : lastNameResult.error!;
      
      default:
        return null;
    }
  };

  return { validateField };
};