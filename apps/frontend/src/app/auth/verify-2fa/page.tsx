"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/lib/auth-context";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { api } from "@/lib/api";

export default function Verify2FAPage() {
  const router = useRouter();
  const { user, refreshUser } = useAuth();
  const [code, setCode] = useState("");
  const [backupCode, setBackupCode] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [showBackupCode, setShowBackupCode] = useState(false);

  useEffect(() => {
    // Redirect if user is not authenticated or 2FA is not enabled
    if (!user) {
      router.push("/auth/login");
      return;
    }

    // Check if 2FA is already verified
    const checkVerificationStatus = async () => {
      try {
        const response = await api.getProfile();
        if (response.success && response.two_factor_verified) {
          router.push("/dashboard");
        }
      } catch (err) {
        console.error("Failed to check verification status:", err);
      }
    };

    checkVerificationStatus();
  }, [user, router]);

  const handleVerifyCode = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!code.trim()) return;

    setLoading(true);
    setError("");

    try {
      const response = await api.verify2FA(code.trim());
      if (response.success) {
        await refreshUser();
        router.push("/dashboard");
      } else {
        setError(response.error || "Invalid verification code");
      }
    } catch (err: any) {
      setError(err.message || "Failed to verify code");
    } finally {
      setLoading(false);
      setCode("");
    }
  };

  const handleVerifyBackupCode = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!backupCode.trim()) return;

    setLoading(true);
    setError("");

    try {
      const response = await api.verifyBackupCode(backupCode.trim());
      if (response.success) {
        await refreshUser();
        router.push("/dashboard");
      } else {
        setError(response.error || "Invalid backup code");
      }
    } catch (err: any) {
      setError(err.message || "Failed to verify backup code");
    } finally {
      setLoading(false);
      setBackupCode("");
    }
  };

  if (!user) {
    return (
      <div className="min-h-screen bg-black flex items-center justify-center">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-white"></div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-black flex items-center justify-center">
      <div className="w-full max-w-md mx-auto p-6">
        <div className="bg-white/5 border border-white/10 rounded-lg p-8">
          <div className="text-center mb-8">
            <div className="w-16 h-16 bg-gradient-to-r from-blue-400 via-purple-400 to-cyan-400 rounded-xl flex items-center justify-center mx-auto mb-4">
              <i className="fas fa-shield-alt text-2xl text-white"></i>
            </div>
            <h1 className="text-2xl font-bold text-white mb-2">
              Two-Factor Authentication
            </h1>
            <p className="text-gray-400">
              Please enter your authentication code to continue
            </p>
          </div>

          {error && (
            <div className="mb-6 p-4 bg-red-500/20 border border-red-500/30 rounded-lg">
              <div className="flex items-center">
                <i className="fas fa-exclamation-circle text-red-400 mr-2"></i>
                <span className="text-red-300">{error}</span>
              </div>
            </div>
          )}

          {!showBackupCode ? (
            <form onSubmit={handleVerifyCode} className="space-y-6">
              <div>
                <label className="block text-sm font-medium text-white mb-2">
                  Authentication Code
                </label>
                <Input
                  type="text"
                  value={code}
                  onChange={(e) => setCode(e.target.value.replace(/\D/g, "").slice(0, 6))}
                  placeholder="Enter 6-digit code"
                  maxLength={6}
                  className="text-center text-lg tracking-widest"
                  autoComplete="one-time-code"
                  autoFocus
                />
                <p className="text-xs text-gray-400 mt-2">
                  Enter the 6-digit code from your authenticator app
                </p>
              </div>

              <Button
                type="submit"
                variant="primary"
                className="w-full"
                loading={loading}
                disabled={loading || code.length !== 6}
              >
                Verify Code
              </Button>

              <div className="text-center">
                <button
                  type="button"
                  onClick={() => setShowBackupCode(true)}
                  className="text-sm text-blue-400 hover:text-blue-300 transition-colors"
                >
                  Use backup code instead
                </button>
              </div>
            </form>
          ) : (
            <form onSubmit={handleVerifyBackupCode} className="space-y-6">
              <div>
                <label className="block text-sm font-medium text-white mb-2">
                  Backup Code
                </label>
                <Input
                  type="text"
                  value={backupCode}
                  onChange={(e) => setBackupCode(e.target.value.toLowerCase().replace(/[^a-z0-9]/g, ""))}
                  placeholder="Enter backup code"
                  maxLength={8}
                  className="text-center text-lg tracking-wider"
                  autoComplete="one-time-code"
                  autoFocus
                />
                <p className="text-xs text-gray-400 mt-2">
                  Enter one of your 8-character backup codes
                </p>
              </div>

              <Button
                type="submit"
                variant="primary"
                className="w-full"
                loading={loading}
                disabled={loading || backupCode.length !== 8}
              >
                Verify Backup Code
              </Button>

              <div className="text-center">
                <button
                  type="button"
                  onClick={() => setShowBackupCode(false)}
                  className="text-sm text-blue-400 hover:text-blue-300 transition-colors"
                >
                  Use authenticator code instead
                </button>
              </div>
            </form>
          )}

          <div className="mt-8 pt-6 border-t border-white/10">
            <div className="text-center text-sm text-gray-400">
              <p>Having trouble? Contact support for assistance.</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}