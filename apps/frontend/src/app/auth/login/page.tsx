"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useAuth } from "@/lib/auth-context";
import { Input, InputLabel, Checkbox } from "@/components/ui/input/index";
import { Button } from "@/components/ui/button";
import AuthLayout from "@/components/auth-layout";
import { AuthGuard } from "@/components/auth/AuthGuard";
import { validateField } from "@/lib/validation";

function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});
  const [rememberMe, setRememberMe] = useState(false);

  const { login } = useAuth();
  const router = useRouter();

  // Validation functions
  const validateEmail = (value: string) => {
    const result = validateField("email", value);
    setFieldErrors(prev => ({ ...prev, email: result.error || "" }));
    return result.isValid;
  };

  const validatePassword = (value: string) => {
    const result = validateField("password", value);
    setFieldErrors(prev => ({ ...prev, password: result.error || "" }));
    return result.isValid;
  };

  const validateForm = () => {
    const emailValid = validateEmail(email);
    const passwordValid = validatePassword(password);
    return emailValid && passwordValid;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setFieldErrors({});

    // Validate form before submission
    if (!validateForm()) {
      return;
    }

    setLoading(true);

    try {
      await login(email, password);
      // Login successful - the auth context will handle routing based on 2FA status
      router.push("/dashboard");
    } catch (err: any) {
      setError(err.message || "Login failed");
    } finally {
      setLoading(false);
    }
  };

  const handleEmailChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setEmail(value);
    if (fieldErrors.email) {
      validateEmail(value);
    }
  };

  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setPassword(value);
    if (fieldErrors.password) {
      validatePassword(value);
    }
  };

  return (
    <AuthLayout>
      <div className="space-y-6">
        <div className="text-center mb-8">
          <h2 className="text-2xl md:text-3xl font-black mb-4 leading-tight">
            Sign in to your Account
          </h2>
          <p className="text-gray-300">Access your Orvex account</p>
        </div>

        {error && (
          <div className="px-4 py-3 rounded-xl flex items-center border bg-red-500/10 border-red-500/20 text-red-400">
            <i className="fas fa-exclamation-triangle mr-3"></i>
            <span className="font-medium">{error}</span>
          </div>
        )}

        <form className="space-y-6" onSubmit={handleSubmit}>
          <div>
            <InputLabel>Email address</InputLabel>
            <Input
              type="email"
              value={email}
              onChange={handleEmailChange}
              placeholder="Enter your email"
              required
              error={fieldErrors.email}
            />
          </div>

          <div>
            <InputLabel>Password</InputLabel>
            <Input
              type="password"
              value={password}
              onChange={handlePasswordChange}
              placeholder="Enter your password"
              required
              error={fieldErrors.password}
            />
          </div>


          <div className="flex items-center justify-between">
            <Checkbox
              label="Remember me"
              className="mr-3 w-4 h-4 rounded bg-white/5 border-white/10 text-blue-500 focus:ring-blue-500/50 focus:ring-2"
              checked={rememberMe}
              onChange={(e) => setRememberMe(e.target.checked)}
            />
            <Link
              href="#"
              className="text-sm text-blue-400 hover:text-blue-300 transition-colors font-medium"
            >
              Forgot password?
            </Link>
          </div>

          <Button
            type="submit"
            loading={loading}
            fullWidth
            size="lg"
            variant="primary"
            icon="fas fa-arrow-right"
            iconPosition="right"
            rounded="md"
          >
            Let's go!
          </Button>
        </form>

        <div className="text-center pt-6 border-t border-white/10">
          <p className="text-gray-400 text-sm">
            Don't have an account?{" "}
            <Link
              href="/auth/register"
              className="text-blue-400 hover:text-blue-300 transition-colors font-medium"
            >
              Create one now
            </Link>
          </p>
        </div>
      </div>
    </AuthLayout>
  );
}

// Wrap the component with AuthGuard
export default function LoginPageWithGuard() {
  return (
    <AuthGuard>
      <LoginPage />
    </AuthGuard>
  );
}
