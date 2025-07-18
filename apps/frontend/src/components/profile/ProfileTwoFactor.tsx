"use client";

import { useState } from "react";
import { useAuth } from "@/lib/auth-context";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { api } from "@/lib/api";

interface TwoFactorSetup {
  secret: string;
  qr_code: string;
}

type SetupStep = "idle" | "password" | "setup" | "verify" | "complete";
type DisableStep = "idle" | "verify";

export function ProfileTwoFactor() {
  const { user, refreshUser } = useAuth();
  const [loading, setLoading] = useState(false);
  const [setupStep, setSetupStep] = useState<SetupStep>("idle");
  const [disableStep, setDisableStep] = useState<DisableStep>("idle");
  const [success, setSuccess] = useState("");
  const [error, setError] = useState("");
  
  // Setup states
  const [password, setPassword] = useState("");
  const [setupData, setSetupData] = useState<TwoFactorSetup | null>(null);
  const [verificationCode, setVerificationCode] = useState("");
  const [backupCodes, setBackupCodes] = useState<string[]>([]);
  const [showBackupCodes, setShowBackupCodes] = useState(false);
  
  // Disable states
  const [disableCode, setDisableCode] = useState("");

  const resetStates = () => {
    setSetupStep("idle");
    setDisableStep("idle");
    setPassword("");
    setSetupData(null);
    setVerificationCode("");
    setBackupCodes([]);
    setDisableCode("");
    setError("");
    setSuccess("");
  };

  const handleStartSetup = () => {
    resetStates();
    setSetupStep("password");
  };

  const handlePasswordSubmit = async () => {
    if (!password.trim()) return;

    setLoading(true);
    setError("");

    try {
      const response = await api.setup2FA();
      if (response.success && response.data) {
        setSetupData(response.data);
        setSetupStep("setup");
      } else {
        setError(response.error || "Failed to setup 2FA");
      }
    } catch (err: any) {
      setError(err.message || "Failed to setup 2FA");
    } finally {
      setLoading(false);
    }
  };

  const handleVerifyAndEnable = async () => {
    if (!verificationCode || !setupData || !password) return;

    setLoading(true);
    setError("");

    try {
      const response = await api.enable2FA({
        secret: setupData.secret,
        token: verificationCode,
        password: password,
      });
      
      if (response.success) {
        setBackupCodes(response.data?.backup_codes || []);
        setSetupStep("complete");
        setSuccess("Two-factor authentication enabled successfully!");
        await refreshUser();
      } else {
        setError(response.error || "Invalid verification code");
      }
    } catch (err: any) {
      setError(err.message || "Failed to enable 2FA");
    } finally {
      setLoading(false);
    }
  };

  const handleStartDisable = () => {
    resetStates();
    setDisableStep("verify");
  };

  const handleDisable2FA = async () => {
    if (!disableCode.trim()) return;

    setLoading(true);
    setError("");

    try {
      const response = await api.disable2FA(disableCode);
      if (response.success) {
        setSuccess("Two-factor authentication disabled successfully!");
        setDisableStep("idle");
        await refreshUser();
      } else {
        setError(response.error || "Invalid verification code");
      }
    } catch (err: any) {
      setError(err.message || "Failed to disable 2FA");
    } finally {
      setLoading(false);
    }
  };

  const handleGenerateBackupCodes = async () => {
    setLoading(true);
    setError("");
    setSuccess("");

    try {
      const response = await api.generateBackupCodes();
      if (response.success && response.data) {
        setBackupCodes(response.data.backup_codes);
        setShowBackupCodes(true);
        setSuccess("New backup codes generated!");
      } else {
        setError(response.error || "Failed to generate backup codes");
      }
    } catch (err: any) {
      setError(err.message || "Failed to generate backup codes");
    } finally {
      setLoading(false);
    }
  };

  const downloadBackupCodes = () => {
    const text = backupCodes.join("\n");
    const blob = new Blob([text], { type: "text/plain" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = "orvex-backup-codes.txt";
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold text-white mb-2">Two-Factor Authentication</h2>
        <p className="text-gray-400">
          Add an extra layer of security to your account with 2FA.
        </p>
      </div>

      {success && (
        <div className="p-4 bg-green-500/20 border border-green-500/30 rounded-lg">
          <div className="flex items-center">
            <i className="fas fa-check-circle text-green-400 mr-2"></i>
            <span className="text-green-300">{success}</span>
          </div>
        </div>
      )}

      {error && (
        <div className="p-4 bg-red-500/20 border border-red-500/30 rounded-lg">
          <div className="flex items-center">
            <i className="fas fa-exclamation-circle text-red-400 mr-2"></i>
            <span className="text-red-300">{error}</span>
          </div>
        </div>
      )}

      {/* Current Status */}
      <div className="bg-white/5 border border-white/10 rounded-lg p-6">
        <div className="flex items-center justify-between">
          <div>
            <h3 className="text-lg font-semibold text-white mb-2">Current Status</h3>
            <div className="flex items-center space-x-3">
              <span
                className={`px-3 py-1 rounded-full text-sm font-medium ${
                  user?.two_factor_enabled
                    ? "bg-green-500/20 text-green-400"
                    : "bg-red-500/20 text-red-400"
                }`}
              >
                {user?.two_factor_enabled ? "Enabled" : "Disabled"}
              </span>
              {user?.two_factor_enabled && (
                <span className="text-green-400">
                  <i className="fas fa-shield-alt mr-1"></i>
                  Your account is protected
                </span>
              )}
            </div>
          </div>
          
          {!user?.two_factor_enabled && setupStep === "idle" && (
            <Button
              variant="primary"
              onClick={handleStartSetup}
              icon="fas fa-shield-alt"
            >
              Enable 2FA
            </Button>
          )}
        </div>
      </div>

      {/* Setup Process - Password Step */}
      {setupStep === "password" && (
        <div className="bg-white/5 border border-white/10 rounded-lg p-6">
          <h3 className="text-lg font-semibold text-white mb-4">Step 1: Confirm Your Password</h3>
          <p className="text-gray-400 mb-6">Please enter your current password to continue with 2FA setup.</p>
          
          <div className="max-w-md space-y-4">
            <Input
              type="password"
              placeholder="Enter your password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              onKeyPress={(e) => e.key === "Enter" && handlePasswordSubmit()}
              autoFocus
            />
            
            <div className="flex space-x-3">
              <Button
                variant="ghost"
                onClick={() => setSetupStep("idle")}
              >
                Cancel
              </Button>
              <Button
                variant="primary"
                onClick={handlePasswordSubmit}
                loading={loading}
                disabled={loading || !password.trim()}
              >
                Continue
              </Button>
            </div>
          </div>
        </div>
      )}

      {/* Setup Process - QR Code & Verification */}
      {setupStep === "setup" && setupData && (
        <div className="bg-white/5 border border-white/10 rounded-lg p-6">
          <h3 className="text-lg font-semibold text-white mb-4">Step 2: Scan QR Code & Verify</h3>
          
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
            {/* QR Code */}
            <div className="space-y-4">
              <div>
                <h4 className="font-medium text-white mb-2">Scan with your authenticator app</h4>
                <p className="text-sm text-gray-400 mb-4">
                  Use Google Authenticator, Authy, or similar apps to scan this code:
                </p>
              </div>
              
              <div className="flex justify-center">
                <div className="bg-white p-4 rounded-lg">
                  <img 
                    src={setupData.qr_code} 
                    alt="2FA QR Code" 
                    className="w-48 h-48"
                  />
                </div>
              </div>
              
              <div className="bg-gray-800 p-3 rounded-lg">
                <p className="text-xs text-gray-400 mb-1">Manual entry code:</p>
                <code className="text-sm text-gray-300 break-all">{setupData.secret}</code>
              </div>
            </div>

            {/* Verification */}
            <div className="space-y-4">
              <div>
                <h4 className="font-medium text-white mb-2">Enter verification code</h4>
                <p className="text-sm text-gray-400 mb-4">
                  Enter the 6-digit code from your authenticator app:
                </p>
              </div>
              
              <div className="space-y-4">
                <Input
                  type="text"
                  placeholder="000000"
                  value={verificationCode}
                  onChange={(e) => setVerificationCode(e.target.value.replace(/\D/g, "").slice(0, 6))}
                  onKeyPress={(e) => e.key === "Enter" && handleVerifyAndEnable()}
                  maxLength={6}
                  className="text-center text-2xl tracking-wider"
                  autoFocus
                />
                
                <div className="flex space-x-3">
                  <Button
                    variant="ghost"
                    onClick={() => setSetupStep("password")}
                  >
                    Back
                  </Button>
                  <Button
                    variant="primary"
                    onClick={handleVerifyAndEnable}
                    loading={loading}
                    disabled={loading || verificationCode.length !== 6}
                    className="flex-1"
                  >
                    Verify & Enable
                  </Button>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Setup Complete - Show Backup Codes */}
      {setupStep === "complete" && backupCodes.length > 0 && (
        <div className="bg-yellow-500/10 border border-yellow-500/30 rounded-lg p-6">
          <div className="flex items-start space-x-3 mb-4">
            <i className="fas fa-exclamation-triangle text-yellow-400 text-lg mt-1"></i>
            <div>
              <h3 className="text-lg font-semibold text-white mb-2">Save Your Backup Codes</h3>
              <p className="text-yellow-200 text-sm">
                These codes can be used to access your account if you lose your authenticator device.
                <strong> Each code can only be used once.</strong>
              </p>
            </div>
          </div>
          
          <div className="grid grid-cols-2 gap-2 font-mono text-sm bg-black/30 p-4 rounded-lg mb-4">
            {backupCodes.map((code, index) => (
              <div key={index} className="text-gray-300 py-1">
                {code}
              </div>
            ))}
          </div>
          
          <div className="flex space-x-3">
            <Button
              variant="glass"
              onClick={downloadBackupCodes}
              icon="fas fa-download"
            >
              Download Codes
            </Button>
            
            <Button
              variant="primary"
              onClick={() => {
                setSetupStep("idle");
                setBackupCodes([]);
              }}
            >
              Complete Setup
            </Button>
          </div>
        </div>
      )}

      {/* Backup Codes Display */}
      {showBackupCodes && backupCodes.length > 0 && setupStep === "idle" && (
        <div className="bg-yellow-500/10 border border-yellow-500/30 rounded-lg p-6">
          <div className="flex items-start space-x-3 mb-4">
            <i className="fas fa-exclamation-triangle text-yellow-400 text-lg mt-1"></i>
            <div>
              <h3 className="text-lg font-semibold text-white mb-2">Backup Codes</h3>
              <p className="text-yellow-200 text-sm">
                Save these backup codes in a safe place. You can use them to access your account if you lose your authenticator device.
                <strong> Each code can only be used once.</strong>
              </p>
            </div>
          </div>
          
          <div className="grid grid-cols-2 gap-2 font-mono text-sm bg-black/30 p-4 rounded-lg mb-4">
            {backupCodes.map((code, index) => (
              <div key={index} className="text-gray-300 py-1">
                {code}
              </div>
            ))}
          </div>
          
          <div className="flex space-x-3">
            <Button
              variant="glass"
              onClick={downloadBackupCodes}
              icon="fas fa-download"
            >
              Download Codes
            </Button>
            
            <Button
              variant="ghost"
              onClick={() => setShowBackupCodes(false)}
            >
              I've Saved Them
            </Button>
          </div>
        </div>
      )}

      {/* Disable 2FA */}
      {disableStep === "verify" && (
        <div className="bg-white/5 border border-white/10 rounded-lg p-6">
          <h3 className="text-lg font-semibold text-white mb-4">Disable Two-Factor Authentication</h3>
          
          <div className="bg-red-500/10 border border-red-500/30 rounded-lg p-4 mb-6">
            <div className="flex items-start space-x-3">
              <i className="fas fa-exclamation-triangle text-red-400 text-lg mt-1"></i>
              <div>
                <h4 className="font-medium text-white mb-1">Warning</h4>
                <p className="text-red-200 text-sm">
                  This will remove 2FA protection from your account. You'll only need your password to sign in.
                </p>
              </div>
            </div>
          </div>

          <div className="max-w-md space-y-4">
            <div>
              <label className="block text-sm font-medium text-white mb-2">
                Verification Code
              </label>
              <Input
                type="text"
                placeholder="000000"
                value={disableCode}
                onChange={(e) => setDisableCode(e.target.value.replace(/\D/g, "").slice(0, 6))}
                onKeyPress={(e) => e.key === "Enter" && handleDisable2FA()}
                maxLength={6}
                className="text-center text-2xl tracking-wider"
                autoFocus
              />
              <p className="text-xs text-gray-400 mt-2">
                Enter the 6-digit code from your authenticator app
              </p>
            </div>

            <div className="flex space-x-3">
              <Button 
                variant="ghost" 
                onClick={() => setDisableStep("idle")}
              >
                Cancel
              </Button>
              <Button
                variant="danger"
                onClick={handleDisable2FA}
                loading={loading}
                disabled={loading || disableCode.length !== 6}
              >
                Disable 2FA
              </Button>
            </div>
          </div>
        </div>
      )}

      {/* Manage 2FA (when enabled) */}
      {user?.two_factor_enabled && setupStep === "idle" && disableStep === "idle" && (
        <div className="bg-white/5 border border-white/10 rounded-lg p-6 space-y-4">
          <h3 className="text-lg font-semibold text-white">Manage 2FA</h3>
          
          <div className="space-y-3">
            <div className="flex items-center justify-between p-4 bg-white/5 rounded-lg">
              <div>
                <h4 className="font-medium text-white">Backup Codes</h4>
                <p className="text-sm text-gray-400">Generate new backup codes</p>
              </div>
              <Button
                variant="glass"
                onClick={handleGenerateBackupCodes}
                loading={loading}
                disabled={loading}
                icon="fas fa-refresh"
              >
                Regenerate
              </Button>
            </div>
            
            <div className="flex items-center justify-between p-4 bg-red-500/10 border border-red-500/30 rounded-lg">
              <div>
                <h4 className="font-medium text-white">Disable 2FA</h4>
                <p className="text-sm text-red-300">Remove two-factor authentication from your account</p>
              </div>
              <Button
                variant="danger"
                onClick={handleStartDisable}
                icon="fas fa-times"
              >
                Disable
              </Button>
            </div>
          </div>
        </div>
      )}

      {/* Information */}
      <div className="bg-blue-500/10 border border-blue-500/30 rounded-lg p-6">
        <div className="flex items-start space-x-3">
          <i className="fas fa-info-circle text-blue-400 text-lg mt-1"></i>
          <div>
            <h3 className="text-lg font-semibold text-white mb-2">What is 2FA?</h3>
            <div className="text-blue-200 text-sm space-y-2">
              <p>
                Two-factor authentication adds an extra layer of security to your account by requiring
                a second form of verification in addition to your password.
              </p>
              <p>
                You'll need an authenticator app like Google Authenticator, Authy, or 1Password to
                generate time-based codes.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}