"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { motion } from "framer-motion";
import { useAuth } from "@/lib/auth-context";
import Input from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import Card, { CardContent, CardHeader } from "@/components/ui/card";

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [totpCode, setTotpCode] = useState("");
  const [showTotp, setShowTotp] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

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
    <div className="min-h-screen bg-gray-900 relative overflow-hidden">
      {/* Background effects */}
      <div className="absolute inset-0">
        <div className="absolute top-20 -left-4 w-72 h-72 bg-violet-500 rounded-full mix-blend-multiply filter blur-xl opacity-10 animate-pulse"></div>
        <div className="absolute -bottom-8 -right-4 w-72 h-72 bg-purple-500 rounded-full mix-blend-multiply filter blur-xl opacity-10 animate-pulse delay-1000"></div>
      </div>

      <div className="relative z-10 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
        <div className="sm:mx-auto sm:w-full sm:max-w-md">
          <motion.div
            initial={{ opacity: 0, y: -20 }}
            animate={{ opacity: 1, y: 0 }}
            className="text-center mb-8"
          >
            <Link href="/" className="inline-block">
              <h2 className="text-4xl font-bold bg-gradient-to-r from-violet-400 to-purple-400 bg-clip-text text-transparent">
                Orvex
              </h2>
            </Link>
            <h1 className="mt-6 text-3xl font-bold text-white">Welcome back</h1>
            <p className="mt-2 text-gray-400">
              Don't have an account?{" "}
              <Link
                href="/auth/register"
                className="text-violet-400 hover:text-violet-300 transition-colors"
              >
                Create one now
              </Link>
            </p>
          </motion.div>
        </div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2 }}
          className="sm:mx-auto sm:w-full sm:max-w-md"
        >
          <Card variant="glass">
            <CardHeader>
              <h3 className="text-lg font-semibold text-white">
                Sign in to your account
              </h3>
              <p className="text-gray-400 text-sm">
                Enter your credentials to continue
              </p>
            </CardHeader>
            <CardContent>
              <form className="space-y-6" onSubmit={handleSubmit}>
                {error && (
                  <motion.div
                    initial={{ opacity: 0, x: -10 }}
                    animate={{ opacity: 1, x: 0 }}
                    className={`px-4 py-3 rounded-lg flex items-center ${
                      showTotp
                        ? "bg-yellow-500/10 border border-yellow-500/20 text-yellow-400"
                        : "bg-red-500/10 border border-red-500/20 text-red-400"
                    }`}
                  >
                    <i
                      className={`${
                        showTotp
                          ? "fas fa-shield-alt"
                          : "fas fa-exclamation-triangle"
                      } mr-2`}
                    ></i>
                    {error}
                  </motion.div>
                )}

                <Input
                  label="Email address"
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="Enter your email"
                  icon="fas fa-envelope"
                  required
                />

                <Input
                  label="Password"
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="Enter your password"
                  icon="fas fa-lock"
                  required
                />

                {showTotp && (
                  <motion.div
                    initial={{ opacity: 0, height: 0 }}
                    animate={{ opacity: 1, height: "auto" }}
                    transition={{ duration: 0.3 }}
                  >
                    <Input
                      label="2FA Code"
                      type="text"
                      value={totpCode}
                      onChange={(e) => setTotpCode(e.target.value)}
                      placeholder="Enter 6-digit code"
                      icon="fas fa-shield-alt"
                      maxLength={6}
                      required
                    />
                  </motion.div>
                )}

                <div className="flex items-center justify-between text-sm">
                  <label className="flex items-center text-gray-400">
                    <input
                      type="checkbox"
                      className="mr-2 rounded bg-gray-700 border-gray-600"
                    />
                    Remember me
                  </label>
                  <Link
                    href="#"
                    className="text-violet-400 hover:text-violet-300 transition-colors"
                  >
                    Forgot password?
                  </Link>
                </div>

                <Button
                  type="submit"
                  loading={loading}
                  className="w-full"
                  size="lg"
                  variant="primary"
                  icon="fas fa-sign-in-alt"
                  iconPosition="left"
                  rounded="md"
                >
                  {showTotp ? "Verify & Sign In" : "Sign In"}
                </Button>

                <div className="relative">
                  <div className="absolute inset-0 flex items-center">
                    <div className="w-full border-t border-gray-700"></div>
                  </div>
                  <div className="relative flex justify-center text-sm">
                    <span className="px-2 bg-gray-800 text-gray-400">
                      New to Orvex?
                    </span>
                  </div>
                </div>

                <Button
                  variant="secondary"
                  className="w-full"
                  icon="fas fa-user-plus"
                  iconPosition="left"
                  rounded="md"
                  onClick={() => router.push("/auth/register")}
                >
                  Create Account
                </Button>
              </form>
            </CardContent>
          </Card>
        </motion.div>
      </div>
    </div>
  );
}
