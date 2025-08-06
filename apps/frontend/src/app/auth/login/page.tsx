"use client";

import axios from "axios";
import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { useNotifications } from "@/contexts/notification-context";
import { useAuth } from "@/contexts/auth-context";
import { validateLoginForm, useFormValidation } from "@/utils/validation";

interface ErrorField {
  field: string;
  message: string;
}

interface ErrorResponse {
  success: boolean;
  message: string;
  errors?: ErrorField[];
}

export default function Login() {
  const [loading, setLoading] = useState(false);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});
  const { success, error } = useNotifications();
  const { validateField } = useFormValidation();
  const { login } = useAuth();
  const router = useRouter();

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
    try {
      e.preventDefault();

      const validation = validateLoginForm(email, password);
      if (!validation.isValid) {
        setFieldErrors(validation.errors);
        error("Please correct the form errors", {
          title: "Validation Error",
          duration: 4000,
        });
        return;
      }

      setFieldErrors({});
      setLoading(true);

      await login(email, password);

      setLoading(false);
      success("Successfully signed in!", { title: "Welcome back!" });
      router.push("/dashboard");
    } catch (err: any) {
      console.error(err);

      let errorMessage = "Something went wrong";

      if (err.message) {
        errorMessage = err.message;
      } else if (err.response?.data?.message) {
        errorMessage = err.response.data.message;
      }

      error(errorMessage, {
        title: "Sign In Failed",
        duration: 6000,
      });
      setLoading(false);
    }
  };

  return (
    <>
      <div className="text-center mb-6">
        <h1 className="text-2xl font-bold text-white mb-1">Welcome back</h1>
        <p className="text-sm text-white/60">Sign in to your account</p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="space-y-1">
          <label
            htmlFor="email"
            className="block text-xs font-medium text-white/80"
          >
            Email address
          </label>
          <input
            id="email"
            name="email"
            type="email"
            placeholder="you@example.com"
            value={email}
            onChange={(e) => {
              setEmail(e.target.value);
              handleFieldChange("email", e.target.value);
            }}
            className={`w-full px-3 py-2 bg-white/5 border rounded-lg text-sm text-white placeholder-white/40 focus:outline-none focus:ring-2 transition-all duration-200 ${
              fieldErrors.email
                ? "border-red-500/50 focus:ring-red-500/50 focus:border-red-500/50"
                : "border-white/10 focus:ring-violet-500/80"
            }`}
            autoComplete="email"
            required
          />
          {fieldErrors.email && (
            <p className="text-red-400 text-xs mt-1">{fieldErrors.email}</p>
          )}
        </div>

        <div className="space-y-1">
          <div className="flex justify-between items-center">
            <label
              htmlFor="password"
              className="block text-xs font-medium text-white/80"
            >
              Password
            </label>
            <a
              href="/auth/forgot-password"
              className="text-xs text-blue-400 hover:text-blue-300 transition-colors"
            >
              Forgot password?
            </a>
          </div>
          <div className="relative">
            <input
              id="password"
              name="password"
              placeholder="Enter your password"
              type={showPassword ? "text" : "password"}
              value={password}
              onChange={(e) => {
                setPassword(e.target.value);
                handleFieldChange("password", e.target.value);
              }}
              className={`w-full px-3 py-2 pr-10 bg-white/5 border rounded-lg text-sm text-white placeholder-white/40 focus:outline-none focus:ring-2 transition-all duration-200 ${
                fieldErrors.password
                  ? "border-red-500/50 focus:ring-red-500/50 focus:border-red-500/50"
                  : "border-white/10 focus:ring-violet-500/80"
              }`}
              autoComplete="current-password"
              required
            />
            <button
              type="button"
              onClick={() => setShowPassword(!showPassword)}
              className="absolute right-2 top-1/2 -translate-y-1/2 p-1 text-white/40 hover:text-white/60 transition-colors"
            >
              <i
                className={`fas ${
                  showPassword ? "fa-eye-slash" : "fa-eye"
                } text-sm`}
              />
            </button>
          </div>
          {fieldErrors.password && (
            <p className="text-red-400 text-xs mt-1">{fieldErrors.password}</p>
          )}
        </div>

        <Button
          variant="primary"
          type="submit"
          loading={loading}
          fullWidth
          size="md"
          className="mt-6"
        >
          {loading ? "Signing in..." : "Sign in"}
        </Button>
      </form>

      <div className="mt-6">
        <div className="relative flex justify-center text-xs">
          <span className="px-4 text-white/50">Or continue with</span>
        </div>

        <div className="mt-4 grid grid-cols-2 gap-3">
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

      <div className="mt-6 text-center">
        <p className="text-xs text-white/60">
          Don't have an account?{" "}
          <a
            href="/auth/signup"
            className="text-blue-400 hover:text-blue-300 font-medium transition-colors"
          >
            Sign up
          </a>
        </p>
      </div>
    </>
  );
}
