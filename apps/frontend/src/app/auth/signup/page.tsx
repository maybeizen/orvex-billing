"use client";

import axios from "axios";
import React, { useState } from "react";
import { Button } from "@/components/ui/button";

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
  const [errors, setErrors] = useState<ErrorField[]>([]);
  const [generalError, setGeneralError] = useState<string | null>(null);
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrors([]);
    setGeneralError(null);

    if (password !== confirmPassword) {
      setGeneralError("Passwords do not match");
      return;
    }

    try {
      setLoading(true);
      const res = await axios.post(
        "http://localhost:3001/api/v1/auth/register",
        { firstName, lastName, username, email, password, confirmPassword },
        { withCredentials: true }
      );

      setLoading(false);
      console.log(`Registered successfully! ${res.data}`);
    } catch (error: any) {
      console.error(error);
      const errorData: ErrorResponse = error.response?.data;
      
      if (errorData?.errors && errorData.errors.length > 0) {
        setErrors(errorData.errors);
      } else {
        setGeneralError(errorData?.message || "Something went wrong");
      }
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-zinc-950 p-4 max-w-3xl w-full">
      <div className="w-full max-w-md">
        <div className="bg-neutral-900/80 backdrop-blur-xl border border-neutral-800/50 rounded-2xl p-8 shadow-2xl">
          <div className="text-center mb-8">
            <h1 className="text-3xl font-bold text-white mb-2">
              Create account
            </h1>
            <p className="text-neutral-400">
              Get started with your free account
            </p>
          </div>

          {(generalError || errors.length > 0) && (
            <div className="mb-6 p-3 bg-red-500/10 border border-red-500/20 rounded-lg space-y-2">
              {generalError && (
                <p className="text-red-400 text-sm text-center">{generalError}</p>
              )}
              {errors.map((error, index) => (
                <div key={index} className="text-red-400 text-sm">
                  <span className="font-medium capitalize">{error.field}:</span> {error.message}
                </div>
              ))}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-6">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <label
                  htmlFor="firstName"
                  className="block text-sm font-medium text-neutral-300"
                >
                  First name
                </label>
                <input
                  id="firstName"
                  name="firstName"
                  type="text"
                  placeholder="John"
                  value={firstName}
                  onChange={(e) => setFirstName(e.target.value)}
                  className="w-full px-4 py-3 bg-neutral-800/50 border border-neutral-700/50 rounded-lg text-white placeholder-neutral-500 focus:outline-none focus:ring-2 focus:ring-violet-500/50 focus:border-violet-500/50 transition-all duration-200"
                  autoComplete="given-name"
                  required
                />
              </div>
              <div className="space-y-2">
                <label
                  htmlFor="lastName"
                  className="block text-sm font-medium text-neutral-300"
                >
                  Last name
                </label>
                <input
                  id="lastName"
                  name="lastName"
                  type="text"
                  placeholder="Doe"
                  value={lastName}
                  onChange={(e) => setLastName(e.target.value)}
                  className="w-full px-4 py-3 bg-neutral-800/50 border border-neutral-700/50 rounded-lg text-white placeholder-neutral-500 focus:outline-none focus:ring-2 focus:ring-violet-500/50 focus:border-violet-500/50 transition-all duration-200"
                  autoComplete="family-name"
                  required
                />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <label
                  htmlFor="username"
                  className="block text-sm font-medium text-neutral-300"
                >
                  Username
                </label>
                <input
                  id="username"
                  name="username"
                  type="text"
                  placeholder="johndoe"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  className="w-full px-4 py-3 bg-neutral-800/50 border border-neutral-700/50 rounded-lg text-white placeholder-neutral-500 focus:outline-none focus:ring-2 focus:ring-violet-500/50 focus:border-violet-500/50 transition-all duration-200"
                  autoComplete="username"
                  required
                />
              </div>
              <div className="space-y-2">
                <label
                  htmlFor="email"
                  className="block text-sm font-medium text-neutral-300"
                >
                  Email address
                </label>
                <input
                  id="email"
                  name="email"
                  type="email"
                  placeholder="john@example.com"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className="w-full px-4 py-3 bg-neutral-800/50 border border-neutral-700/50 rounded-lg text-white placeholder-neutral-500 focus:outline-none focus:ring-2 focus:ring-violet-500/50 focus:border-violet-500/50 transition-all duration-200"
                  autoComplete="email"
                  required
                />
              </div>
            </div>

            <div className="space-y-2">
              <label
                htmlFor="password"
                className="block text-sm font-medium text-neutral-300"
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
                  onChange={(e) => setPassword(e.target.value)}
                  className="w-full px-4 py-3 pr-12 bg-neutral-800/50 border border-neutral-700/50 rounded-lg text-white placeholder-neutral-500 focus:outline-none focus:ring-2 focus:ring-violet-500/50 focus:border-violet-500/50 transition-all duration-200"
                  autoComplete="new-password"
                  required
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 p-1 text-neutral-400 hover:text-neutral-300 transition-colors"
                  aria-label={showPassword ? "Hide password" : "Show password"}
                >
                  <i
                    className={`fas ${
                      showPassword ? "fa-eye-slash" : "fa-eye"
                    } text-lg`}
                  />
                </button>
              </div>
            </div>

            <div className="space-y-2">
              <label
                htmlFor="confirmPassword"
                className="block text-sm font-medium text-neutral-300"
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
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  className="w-full px-4 py-3 pr-12 bg-neutral-800/50 border border-neutral-700/50 rounded-lg text-white placeholder-neutral-500 focus:outline-none focus:ring-2 focus:ring-violet-500/50 focus:border-violet-500/50 transition-all duration-200"
                  autoComplete="new-password"
                  required
                />
                <button
                  type="button"
                  onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 p-1 text-neutral-400 hover:text-neutral-300 transition-colors"
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
            <div className="relative">
              <div className="absolute inset-0 flex items-center">
                <div className="w-full border-t border-neutral-700/50" />
              </div>
              <div className="relative flex justify-center text-sm">
                <span className="px-4 bg-neutral-900/80 text-neutral-500">
                  Or continue with
                </span>
              </div>
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
            <p className="text-sm text-neutral-400">
              Already have an account?{" "}
              <a
                href="/auth/login"
                className="text-violet-400 hover:text-violet-300 font-medium transition-colors"
              >
                Sign in
              </a>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
