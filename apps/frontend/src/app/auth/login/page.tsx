"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useAuth } from "@/lib/auth-context";
import { Input, InputLabel, Checkbox } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import AuthLayout from "@/components/auth-layout";

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [totpCode, setTotpCode] = useState("");
  const [showTotp, setShowTotp] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [rememberMe, setRememberMe] = useState(false);

  const { login } = useAuth();
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      await login(email, password, totpCode || undefined);
      router.push("/dashboard");
    } catch (err: any) {
      if (err.message.includes("2FA") || err.message.includes("TOTP")) {
        setShowTotp(true);
        setError("Please enter your 2FA code");
      } else {
        setError(err.message || "Login failed");
      }
    } finally {
      setLoading(false);
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
          <div
            className={`px-4 py-3 rounded-xl flex items-center border ${
              showTotp
                ? "bg-yellow-500/10 border-yellow-500/20 text-yellow-400"
                : "bg-red-500/10 border-red-500/20 text-red-400"
            }`}
          >
            <i
              className={`${
                showTotp ? "fas fa-shield-alt" : "fas fa-exclamation-triangle"
              } mr-3`}
            ></i>
            <span className="font-medium">{error}</span>
          </div>
        )}

        <form className="space-y-6" onSubmit={handleSubmit}>
          <div>
            <InputLabel>Email address</InputLabel>
            <Input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="Enter your email"
              required
            />
          </div>

          <div>
            <InputLabel>Password</InputLabel>
            <Input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Enter your password"
              required
            />
          </div>

          {showTotp && (
            <div>
              <InputLabel>2FA Code</InputLabel>
              <Input
                type="text"
                value={totpCode}
                onChange={(e) => setTotpCode(e.target.value)}
                placeholder="Enter 6-digit code"
                required
              />
            </div>
          )}

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
            {showTotp ? "Verify & Sign In" : "Let's go!"}
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
