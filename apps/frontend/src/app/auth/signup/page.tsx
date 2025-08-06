"use client";

import axios from "axios";
import React, { useState } from "react";
import { Button } from "@/components/ui/button";
import { useNotifications } from "@/contexts/notification-context";
import {
  validateRegistrationForm,
  useFormValidation,
} from "@/utils/validation";
import { PasswordStrength } from "@/components/ui/password-strength";

interface ErrorField {
  field: string;
  message: string;
}

interface ErrorResponse {
  success: boolean;
  message: string;
  errors?: ErrorField[];
}

export default function SignUp() {
  const [loading, setLoading] = useState(false);
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});
  const [passwordFocused, setPasswordFocused] = useState(false);
  const { success, error } = useNotifications();
  const { validateField } = useFormValidation();

  const handleFieldChange = (fieldName: string, value: string) => {
    if (fieldErrors[fieldName]) {
      setFieldErrors((prev) => {
        const newErrors = { ...prev };
        delete newErrors[fieldName];
        return newErrors;
      });
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const validation = validateRegistrationForm({
      firstName,
      lastName,
      username,
      email,
      password,
      confirmPassword,
    });

    if (!validation.isValid) {
      setFieldErrors(validation.errors);
      error("Please correct the form errors", {
        title: "Validation Error",
        duration: 4000,
      });
      return;
    }

    try {
      setFieldErrors({});
      setLoading(true);

      const res = await axios.post(
        "http://localhost:3001/api/v1/auth/register",
        { firstName, lastName, username, email, password, confirmPassword },
        { withCredentials: true }
      );

      setLoading(false);
      success("Account created successfully!", {
        title: "Welcome to Orvex!",
        duration: 6000,
      });
      console.log(`Registered successfully! ${res.data}`);
    } catch (err: any) {
      console.error(err);
      const errorData: ErrorResponse = err.response?.data;

      if (errorData?.errors && errorData.errors.length > 0) {
        errorData.errors.forEach((fieldError) => {
          error(`${fieldError.field}: ${fieldError.message}`, {
            title: "Validation Error",
            duration: 6000,
          });
        });
      } else {
        error(errorData?.message || "Something went wrong", {
          title: "Registration Failed",
          duration: 6000,
        });
      }
      setLoading(false);
    }
  };

  return (
    <>
      <div className="text-center mb-8">
        <h1 className="text-3xl font-bold text-white mb-2">Create account</h1>
        <p className="text-white/60">Get started with your free account</p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-6">
        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-2">
            <label
              htmlFor="firstName"
              className="block text-sm font-medium text-white/80"
            >
              First name
            </label>
            <input
              id="firstName"
              name="firstName"
              type="text"
              placeholder="John"
              value={firstName}
              onChange={(e) => {
                setFirstName(e.target.value);
                handleFieldChange("firstName", e.target.value);
              }}
              className={`w-full px-4 py-3 bg-white/5 border rounded-lg text-white placeholder-white/40 focus:outline-none focus:ring-2 transition-all duration-200 ${
                fieldErrors.firstName
                  ? "border-red-500/50 focus:ring-red-500/50 focus:border-red-500/50"
                  : "border-white/10 focus:ring-blue-500/50 focus:border-blue-500/50"
              }`}
              autoComplete="given-name"
              required
            />
            {fieldErrors.firstName && (
              <p className="text-red-400 text-sm mt-1">
                {fieldErrors.firstName}
              </p>
            )}
          </div>
          <div className="space-y-2">
            <label
              htmlFor="lastName"
              className="block text-sm font-medium text-white/80"
            >
              Last name
            </label>
            <input
              id="lastName"
              name="lastName"
              type="text"
              placeholder="Doe"
              value={lastName}
              onChange={(e) => {
                setLastName(e.target.value);
                handleFieldChange("lastName", e.target.value);
              }}
              className={`w-full px-4 py-3 bg-white/5 border rounded-lg text-white placeholder-white/40 focus:outline-none focus:ring-2 transition-all duration-200 ${
                fieldErrors.lastName
                  ? "border-red-500/50 focus:ring-red-500/50 focus:border-red-500/50"
                  : "border-white/10 focus:ring-blue-500/50 focus:border-blue-500/50"
              }`}
              autoComplete="family-name"
              required
            />
            {fieldErrors.lastName && (
              <p className="text-red-400 text-sm mt-1">
                {fieldErrors.lastName}
              </p>
            )}
          </div>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-2">
            <label
              htmlFor="username"
              className="block text-sm font-medium text-white/80"
            >
              Username
            </label>
            <input
              id="username"
              name="username"
              type="text"
              placeholder="johndoe"
              value={username}
              onChange={(e) => {
                setUsername(e.target.value);
                handleFieldChange("username", e.target.value);
              }}
              className={`w-full px-4 py-3 bg-white/5 border rounded-lg text-white placeholder-white/40 focus:outline-none focus:ring-2 transition-all duration-200 ${
                fieldErrors.username
                  ? "border-red-500/50 focus:ring-red-500/50 focus:border-red-500/50"
                  : "border-white/10 focus:ring-blue-500/50 focus:border-blue-500/50"
              }`}
              autoComplete="username"
              required
            />
            {fieldErrors.username && (
              <p className="text-red-400 text-sm mt-1">
                {fieldErrors.username}
              </p>
            )}
          </div>
          <div className="space-y-2">
            <label
              htmlFor="email"
              className="block text-sm font-medium text-white/80"
            >
              Email address
            </label>
            <input
              id="email"
              name="email"
              type="email"
              placeholder="john@example.com"
              value={email}
              onChange={(e) => {
                setEmail(e.target.value);
                handleFieldChange("email", e.target.value);
              }}
              className={`w-full px-4 py-3 bg-white/5 border rounded-lg text-white placeholder-white/40 focus:outline-none focus:ring-2 transition-all duration-200 ${
                fieldErrors.email
                  ? "border-red-500/50 focus:ring-red-500/50 focus:border-red-500/50"
                  : "border-white/10 focus:ring-blue-500/50 focus:border-blue-500/50"
              }`}
              autoComplete="email"
              required
            />
            {fieldErrors.email && (
              <p className="text-red-400 text-sm mt-1">{fieldErrors.email}</p>
            )}
          </div>
        </div>

        <div className="space-y-2">
          <label
            htmlFor="password"
            className="block text-sm font-medium text-white/80"
          >
            Password
          </label>
          <div className="relative">
            <input
              id="password"
              name="password"
              placeholder="Create a strong password"
              type={showPassword ? "text" : "password"}
              value={password}
              onChange={(e) => {
                setPassword(e.target.value);
                handleFieldChange("password", e.target.value);
              }}
              onFocus={() => setPasswordFocused(true)}
              onBlur={() => setPasswordFocused(false)}
              className={`w-full px-4 py-3 pr-12 bg-white/5 border rounded-lg text-white placeholder-white/40 focus:outline-none focus:ring-2 transition-all duration-200 ${
                fieldErrors.password
                  ? "border-red-500/50 focus:ring-red-500/50 focus:border-red-500/50"
                  : "border-white/10 focus:ring-blue-500/50 focus:border-blue-500/50"
              }`}
              autoComplete="new-password"
              required
            />
            <button
              type="button"
              onClick={() => setShowPassword(!showPassword)}
              className="absolute right-3 top-1/2 -translate-y-1/2 p-1 text-white/40 hover:text-white/60 transition-colors"
              aria-label={showPassword ? "Hide password" : "Show password"}
            >
              <i
                className={`fas ${
                  showPassword ? "fa-eye-slash" : "fa-eye"
                } text-lg`}
              />
            </button>
          </div>
          {fieldErrors.password && (
            <p className="text-red-400 text-sm mt-1">{fieldErrors.password}</p>
          )}
          <PasswordStrength
            password={password}
            showPopup={true}
            isFocused={passwordFocused}
          />
        </div>

        <div className="space-y-2">
          <label
            htmlFor="confirmPassword"
            className="block text-sm font-medium text-white/80"
          >
            Confirm password
          </label>
          <div className="relative">
            <input
              id="confirmPassword"
              name="confirmPassword"
              placeholder="Confirm your password"
              type={showConfirmPassword ? "text" : "password"}
              value={confirmPassword}
              onChange={(e) => {
                setConfirmPassword(e.target.value);
                handleFieldChange("confirmPassword", e.target.value);
              }}
              className={`w-full px-4 py-3 pr-12 bg-white/5 border rounded-lg text-white placeholder-white/40 focus:outline-none focus:ring-2 transition-all duration-200 ${
                fieldErrors.confirmPassword
                  ? "border-red-500/50 focus:ring-red-500/50 focus:border-red-500/50"
                  : "border-white/10 focus:ring-blue-500/50 focus:border-blue-500/50"
              }`}
              autoComplete="new-password"
              required
            />
            <button
              type="button"
              onClick={() => setShowConfirmPassword(!showConfirmPassword)}
              className="absolute right-3 top-1/2 -translate-y-1/2 p-1 text-white/40 hover:text-white/60 transition-colors"
              aria-label={
                showConfirmPassword ? "Hide password" : "Show password"
              }
            >
              <i
                className={`fas ${
                  showConfirmPassword ? "fa-eye-slash" : "fa-eye"
                } text-lg`}
              />
            </button>
          </div>
          {fieldErrors.confirmPassword && (
            <p className="text-red-400 text-sm mt-1">
              {fieldErrors.confirmPassword}
            </p>
          )}
        </div>

        <Button
          variant="primary"
          type="submit"
          loading={loading}
          fullWidth
          size="lg"
          className="mt-8"
        >
          {loading ? "Creating account..." : "Create account"}
        </Button>
      </form>

      <div className="mt-8">
        <div className="relative flex justify-center text-sm">
          <span className="px-4 text-white/50">Or continue with</span>
        </div>

        <div className="mt-6 grid grid-cols-2 gap-3">
          <Button
            variant="glass"
            icon="fab fa-google"
            iconPosition="left"
            fullWidth
            className="justify-center"
          >
            Google
          </Button>
          <Button
            variant="glass"
            icon="fab fa-discord"
            iconPosition="left"
            fullWidth
            className="justify-center"
          >
            Discord
          </Button>
        </div>
      </div>

      <div className="mt-8 text-center">
        <p className="text-sm text-white/60">
          Already have an account?{" "}
          <a
            href="/auth/login"
            className="text-blue-400 hover:text-blue-300 font-medium transition-colors"
          >
            Sign in
          </a>
        </p>
      </div>
    </>
  );
}
