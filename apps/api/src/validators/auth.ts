import { z } from "zod";

export const registerSchema = z
  .object({
    email: z
      .email("Invalid email address")
      .min(1, "Email is required")
      .max(254, "Email must be less than 254 characters")
      .toLowerCase()
      .trim(),
    password: z
      .string()
      .min(8, "Password must be at least 8 characters")
      .max(128, "Password must be less than 128 characters")
      .regex(
        /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/,
        "Password must contain at least one lowercase letter, one uppercase letter, and one number"
      ),
    confirmPassword: z.string().min(1, "Password confirmation is required"),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: "Passwords don't match",
    path: ["confirmPassword"],
  });

export const loginSchema = z.object({
  email: z
    .email("Invalid email address")
    .min(1, "Email is required")
    .max(254, "Email must be less than 254 characters")
    .toLowerCase()
    .trim(),
  password: z
    .string()
    .min(1, "Password is required")
    .max(128, "Password too long"),
});

export const forgotPasswordSchema = z.object({
  email: z
    .email("Invalid email address")
    .min(1, "Email is required")
    .max(254, "Email must be less than 254 characters")
    .toLowerCase()
    .trim(),
});

export const resetPasswordSchema = z
  .object({
    token: z
      .string()
      .min(1, "Reset token is required")
      .max(256, "Invalid token format"),
    password: z
      .string()
      .min(8, "Password must be at least 8 characters")
      .max(128, "Password must be less than 128 characters")
      .regex(
        /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/,
        "Password must contain at least one lowercase letter, one uppercase letter, and one number"
      ),
    confirmPassword: z.string().min(1, "Password confirmation is required"),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: "Passwords don't match",
    path: ["confirmPassword"],
  });

export type RegisterInput = z.infer<typeof registerSchema>;
export type LoginInput = z.infer<typeof loginSchema>;
export type ForgotPasswordInput = z.infer<typeof forgotPasswordSchema>;
export type ResetPasswordInput = z.infer<typeof resetPasswordSchema>;
