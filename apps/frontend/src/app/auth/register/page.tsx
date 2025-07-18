"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { api } from "@/lib/api";
import { Button } from "@/components/ui/button";
import AuthLayout from "@/components/auth-layout";
import { Input, InputLabel } from "@/components/ui/input/index";

export default function RegisterPage() {
  const [formData, setFormData] = useState({
    email: "",
    username: "",
    password: "",
    confirmPassword: "",
    firstName: "",
    lastName: "",
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);

  const router = useRouter();

  const handleChange = (field: string, value: string) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
    if (error) setError("");
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    if (formData.password !== formData.confirmPassword) {
      setError("Passwords do not match");
      setLoading(false);
      return;
    }

    if (formData.password.length < 8) {
      setError("Password must be at least 8 characters long");
      setLoading(false);
      return;
    }

    if (formData.username.length < 3 || formData.username.length > 20) {
      setError("Username must be between 3 and 20 characters");
      setLoading(false);
      return;
    }

    if (!/^[a-zA-Z0-9_]*$/.test(formData.username)) {
      setError("Username can only contain letters, numbers, and underscores");
      setLoading(false);
      return;
    }

    if (!/^[a-zA-Z]/.test(formData.username)) {
      setError("Username must start with a letter");
      setLoading(false);
      return;
    }

    try {
      const response = await api.register({
        email: formData.email,
        username: formData.username,
        password: formData.password,
        first_name: formData.firstName,
        last_name: formData.lastName,
      });

      if (response.success) {
        setSuccess(true);
        setTimeout(() => {
          router.push("/auth/login");
        }, 3000);
      } else {
        setError(response.error || response.message || "Registration failed");
      }
    } catch (err: any) {
      setError(err.message || "Registration failed");
    } finally {
      setLoading(false);
    }
  };

  if (success) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-black text-white">
        <div className="max-w-md mx-auto px-4 text-center">
          <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-xl p-10">
            <div className="text-green-400 text-6xl mb-6">
              <i className="fas fa-check-circle"></i>
            </div>
            <h2 className="text-3xl font-bold text-white mb-4">
              Welcome to Orvex!
            </h2>
            <p className="text-gray-300 mb-6">
              Your account has been created successfully. You can now access all
              features of our platform.
            </p>
            <div className="flex items-center justify-center text-blue-400">
              <i className="fas fa-spinner fa-spin mr-2"></i>
              <span className="text-sm">Redirecting to login...</span>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <AuthLayout>
      <div className="space-y-6">
        <div className="text-center mb-8">
          <h2 className="text-2xl md:text-3xl font-black mb-4 leading-tight">
            Create your Account
          </h2>
          <p className="text-gray-300">Join thousands of players on Orvex</p>
        </div>

        {error && (
          <div className="bg-red-500/10 border border-red-500/20 text-red-400 px-4 py-3 rounded-xl flex items-center">
            <i className="fas fa-exclamation-triangle mr-3"></i>
            <span className="font-medium">{error}</span>
          </div>
        )}

        <form className="space-y-6" onSubmit={handleSubmit}>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <InputLabel>First name</InputLabel>
              <Input
                value={formData.firstName}
                onChange={(e) => handleChange("firstName", e.target.value)}
                placeholder="John"
                required
              />
            </div>
            <div>
              <InputLabel>Last name</InputLabel>
              <Input
                value={formData.lastName}
                onChange={(e) => handleChange("lastName", e.target.value)}
                placeholder="Doe"
                required
              />
            </div>
          </div>

          <div>
            <InputLabel>Username</InputLabel>
            <Input
              value={formData.username}
              onChange={(e) => handleChange("username", e.target.value)}
              placeholder="johndoe123"
              required
            />
          </div>

          <div>
            <InputLabel>Email</InputLabel>
            <Input
              type="email"
              value={formData.email}
              onChange={(e) => handleChange("email", e.target.value)}
              placeholder="john@example.com"
              required
            />
          </div>

          <div>
            <InputLabel>Password</InputLabel>
            <Input
              type="password"
              value={formData.password}
              onChange={(e) => handleChange("password", e.target.value)}
              placeholder="Create a strong password"
              required
            />
          </div>

          <div>
            <InputLabel>Confirm password</InputLabel>
            <Input
              type="password"
              value={formData.confirmPassword}
              onChange={(e) => handleChange("confirmPassword", e.target.value)}
              placeholder="Confirm your password"
              error={
                formData.confirmPassword &&
                formData.password !== formData.confirmPassword
                  ? "Passwords do not match"
                  : ""
              }
              required
            />
          </div>

          <Button
            type="submit"
            loading={loading}
            fullWidth
            size="lg"
            variant="primary"
            icon="fas fa-rocket"
            iconPosition="left"
            rounded="md"
          >
            Create Account
          </Button>
        </form>

        <div className="text-center pt-6 border-t border-white/10">
          <p className="text-gray-400 text-sm">
            Already have an account?{" "}
            <Link
              href="/auth/login"
              className="text-blue-400 hover:text-blue-300 transition-colors font-medium"
            >
              Sign in
            </Link>
          </p>
        </div>
      </div>
    </AuthLayout>
  );
}
