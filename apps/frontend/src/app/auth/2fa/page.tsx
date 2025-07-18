"use client";

import { useState, useEffect } from "react";
import { useAuth } from "@/lib/auth-context";
import { useRouter, useSearchParams } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { api } from "@/lib/api";

export default function TwoFactorPage() {
  const { user, refreshUser } = useAuth();
  const router = useRouter();
  const searchParams = useSearchParams();
  const [code, setCode] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [useBackupCode, setUseBackupCode] = useState(false);

  const returnTo = searchParams.get("returnTo") || "/dashboard";

  useEffect(() => {
    // If user is not logged in, redirect to login
    if (!user) {
      router.push("/auth/login");
      return;
    }

    // If 2FA is not enabled, redirect to dashboard
    if (!user.two_factor_enabled) {
      router.push(returnTo);
      return;
    }
  }, [user, router, returnTo]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!code.trim()) return;

    setLoading(true);
    setError("");

    try {
      const response = useBackupCode 
        ? await api.verifyBackupCode(code.trim())
        : await api.verify2FA(code.trim());
      
      if (response.success) {
        // 2FA verification successful, refresh user and redirect
        await refreshUser();
        router.push(returnTo);
      } else {
        setError(response.error || "Invalid verification code");
      }
    } catch (err: any) {
      setError(err.message || "Verification failed");
    } finally {
      setLoading(false);
    }
  };

  const handleCodeChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value.replace(/\D/g, "");
    setCode(value.slice(0, useBackupCode ? 8 : 6));
    setError("");
  };

  // Don't render anything if user is not available yet
  if (!user) {
    return null;
  }

  // Don't render if 2FA is not enabled
  if (!user.two_factor_enabled) {
    return null;
  }

  return (
    <div className="min-h-screen bg-black flex items-center justify-center p-4">
      <div className="max-w-md w-full space-y-8">
        {/* Header */}
        <div className="text-center">
          <div className="w-16 h-16 bg-gradient-to-r from-blue-400 via-purple-400 to-cyan-400 rounded-2xl flex items-center justify-center mx-auto mb-6">
            <i className="fas fa-shield-alt text-white text-2xl"></i>
          </div>
          <h1 className="text-3xl font-black text-white mb-2">
            Two-Factor Authentication
          </h1>
          <p className="text-gray-400">
            Enter the verification code from your authenticator app
          </p>
        </div>

        {/* Form */}
        <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-xl p-8">
          {error && (
            <div className="p-4 bg-red-500/20 border border-red-500/30 rounded-lg mb-6">
              <div className="flex items-center">
                <i className="fas fa-exclamation-circle text-red-400 mr-2"></i>
                <span className="text-red-300">{error}</span>
              </div>
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-6">
            <div className="space-y-2">
              <label htmlFor="code" className="block text-sm font-medium text-gray-300">
                {useBackupCode ? "Backup Code" : "Verification Code"}
              </label>
              <Input
                id="code"
                type="text"
                value={code}
                onChange={handleCodeChange}
                placeholder={useBackupCode ? "12345678" : "000000"}
                maxLength={useBackupCode ? 8 : 6}
                className="text-center text-2xl tracking-wider"
                autoComplete="one-time-code"
                autoFocus
                required
              />
              <p className="text-xs text-gray-400 text-center">
                {useBackupCode 
                  ? "Enter one of your 8-digit backup codes"
                  : "Enter the 6-digit code from your authenticator app"
                }
              </p>
            </div>

            <Button
              type="submit"
              variant="primary"
              fullWidth
              loading={loading}
              disabled={loading || code.length < (useBackupCode ? 8 : 6)}
            >
              Verify
            </Button>
          </form>

          {/* Backup Code Toggle */}
          <div className="mt-6 pt-6 border-t border-white/10">
            <button
              type="button"
              onClick={() => {
                setUseBackupCode(!useBackupCode);
                setCode("");
                setError("");
              }}
              className="w-full text-center text-sm text-gray-400 hover:text-white transition-colors duration-200"
            >
              {useBackupCode ? (
                <>
                  <i className="fas fa-mobile-alt mr-2"></i>
                  Use authenticator app instead
                </>
              ) : (
                <>
                  <i className="fas fa-key mr-2"></i>
                  Use backup code instead
                </>
              )}
            </button>
          </div>

          {/* Help */}
          <div className="mt-6 pt-6 border-t border-white/10">
            <div className="text-center space-y-2">
              <p className="text-xs text-gray-400">
                Having trouble? Contact support for assistance.
              </p>
              <button
                type="button"
                onClick={() => {
                  // Add logout logic here if needed
                  router.push("/auth/login");
                }}
                className="text-xs text-gray-400 hover:text-white transition-colors duration-200"
              >
                Sign out and try again
              </button>
            </div>
          </div>
        </div>

        {/* Security Notice */}
        <div className="bg-blue-500/10 border border-blue-500/30 rounded-lg p-4">
          <div className="flex items-start space-x-3">
            <i className="fas fa-info-circle text-blue-400 text-sm mt-0.5"></i>
            <div>
              <p className="text-blue-200 text-xs">
                For your security, you need to verify your identity with two-factor authentication
                before accessing your account.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}