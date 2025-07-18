"use client";

import { useState } from "react";
import { useAuth } from "@/lib/auth-context";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { api } from "@/lib/api";

export function ProfileSecurity() {
  const { user } = useAuth();
  const [loading, setLoading] = useState(false);
  const [passwordForm, setPasswordForm] = useState({
    current_password: "",
    new_password: "",
    confirm_password: "",
  });
  const [success, setSuccess] = useState("");
  const [error, setError] = useState("");

  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setPasswordForm({
      ...passwordForm,
      [e.target.name]: e.target.value,
    });
    setError("");
    setSuccess("");
  };

  const handlePasswordSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError("");
    setSuccess("");

    if (passwordForm.new_password !== passwordForm.confirm_password) {
      setError("New passwords do not match");
      setLoading(false);
      return;
    }

    if (passwordForm.new_password.length < 8) {
      setError("Password must be at least 8 characters long");
      setLoading(false);
      return;
    }

    try {
      const response = await api.changePassword({
        current_password: passwordForm.current_password,
        new_password: passwordForm.new_password,
      });
      
      if (response.success) {
        setSuccess("Password updated successfully!");
        setPasswordForm({
          current_password: "",
          new_password: "",
          confirm_password: "",
        });
      } else {
        setError(response.error || "Failed to update password");
      }
    } catch (err: any) {
      setError(err.message || "Failed to update password");
    } finally {
      setLoading(false);
    }
  };

  const passwordStrength = (password: string) => {
    let strength = 0;
    if (password.length >= 8) strength++;
    if (/[a-z]/.test(password)) strength++;
    if (/[A-Z]/.test(password)) strength++;
    if (/[0-9]/.test(password)) strength++;
    if (/[^A-Za-z0-9]/.test(password)) strength++;
    return strength;
  };

  const getStrengthColor = (strength: number) => {
    if (strength <= 2) return "bg-red-500";
    if (strength <= 3) return "bg-yellow-500";
    return "bg-green-500";
  };

  const getStrengthText = (strength: number) => {
    if (strength <= 2) return "Weak";
    if (strength <= 3) return "Medium";
    return "Strong";
  };

  const strength = passwordStrength(passwordForm.new_password);

  return (
    <div className="space-y-8">
      <div>
        <h2 className="text-2xl font-bold text-white mb-2">Security Settings</h2>
        <p className="text-gray-400">
          Manage your password and account security preferences.
        </p>
      </div>

      {/* Account Status */}
      <div className="bg-white/5 border border-white/10 rounded-lg p-6">
        <h3 className="text-lg font-semibold text-white mb-4">Account Status</h3>
        <div className="space-y-3">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <i className="fas fa-envelope text-gray-400"></i>
              <span className="text-gray-300">Email Verification</span>
            </div>
            <span
              className={`px-3 py-1 rounded-full text-xs font-medium ${
                user?.email_verified
                  ? "bg-green-500/20 text-green-400"
                  : "bg-yellow-500/20 text-yellow-400"
              }`}
            >
              {user?.email_verified ? "Verified" : "Unverified"}
            </span>
          </div>
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <i className="fas fa-shield-alt text-gray-400"></i>
              <span className="text-gray-300">Two-Factor Authentication</span>
            </div>
            <span
              className={`px-3 py-1 rounded-full text-xs font-medium ${
                user?.two_factor_enabled
                  ? "bg-green-500/20 text-green-400"
                  : "bg-red-500/20 text-red-400"
              }`}
            >
              {user?.two_factor_enabled ? "Enabled" : "Disabled"}
            </span>
          </div>
        </div>
      </div>

      {/* Password Change */}
      <div className="bg-white/5 border border-white/10 rounded-lg p-6">
        <h3 className="text-lg font-semibold text-white mb-4">Change Password</h3>
        
        {success && (
          <div className="p-4 bg-green-500/20 border border-green-500/30 rounded-lg mb-4">
            <div className="flex items-center">
              <i className="fas fa-check-circle text-green-400 mr-2"></i>
              <span className="text-green-300">{success}</span>
            </div>
          </div>
        )}

        {error && (
          <div className="p-4 bg-red-500/20 border border-red-500/30 rounded-lg mb-4">
            <div className="flex items-center">
              <i className="fas fa-exclamation-circle text-red-400 mr-2"></i>
              <span className="text-red-300">{error}</span>
            </div>
          </div>
        )}

        <form onSubmit={handlePasswordSubmit} className="space-y-6">
          <div className="space-y-2">
            <label htmlFor="current_password" className="block text-sm font-medium text-gray-300">
              Current Password
            </label>
            <Input
              id="current_password"
              name="current_password"
              type="password"
              value={passwordForm.current_password}
              onChange={handlePasswordChange}
              placeholder="Enter your current password"
              required
            />
          </div>

          <div className="space-y-2">
            <label htmlFor="new_password" className="block text-sm font-medium text-gray-300">
              New Password
            </label>
            <Input
              id="new_password"
              name="new_password"
              type="password"
              value={passwordForm.new_password}
              onChange={handlePasswordChange}
              placeholder="Enter your new password"
              required
            />
            {passwordForm.new_password && (
              <div className="mt-2">
                <div className="flex items-center justify-between text-sm">
                  <span className="text-gray-400">Password strength:</span>
                  <span className={`font-medium ${
                    strength <= 2 ? "text-red-400" : 
                    strength <= 3 ? "text-yellow-400" : "text-green-400"
                  }`}>
                    {getStrengthText(strength)}
                  </span>
                </div>
                <div className="w-full bg-gray-700 rounded-full h-2 mt-1">
                  <div
                    className={`h-2 rounded-full transition-all duration-300 ${getStrengthColor(strength)}`}
                    style={{ width: `${(strength / 5) * 100}%` }}
                  ></div>
                </div>
              </div>
            )}
          </div>

          <div className="space-y-2">
            <label htmlFor="confirm_password" className="block text-sm font-medium text-gray-300">
              Confirm New Password
            </label>
            <Input
              id="confirm_password"
              name="confirm_password"
              type="password"
              value={passwordForm.confirm_password}
              onChange={handlePasswordChange}
              placeholder="Confirm your new password"
              required
            />
            {passwordForm.confirm_password && passwordForm.new_password !== passwordForm.confirm_password && (
              <p className="text-red-400 text-sm mt-1">
                <i className="fas fa-exclamation-triangle mr-1"></i>
                Passwords do not match
              </p>
            )}
          </div>

          <div className="flex justify-end">
            <Button
              type="submit"
              variant="primary"
              loading={loading}
              disabled={loading || passwordForm.new_password !== passwordForm.confirm_password}
            >
              Update Password
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}