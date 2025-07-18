"use client";

import { useState, useCallback, useEffect } from "react";
import { useAuth } from "@/lib/auth-context";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { api } from "@/lib/api";

interface TwoFactorSetup {
  secret: string;
  qr_code: string;
}

interface TwoFactorState {
  // UI State
  isLoading: boolean;
  error: string;
  success: string;
  
  // Setup Flow
  setupStep: 'idle' | 'password' | 'qr' | 'complete';
  setupData: TwoFactorSetup | null;
  
  // Disable Flow
  disableStep: 'idle' | 'verify';
  
  // Backup Code Generation Flow
  generateStep: 'idle' | 'verify';
  
  // Form Data
  password: string;
  verificationCode: string;
  disableCode: string;
  generateCode: string;
  
  // Backup Codes
  backupCodes: string[];
  showBackupCodes: boolean;
}

const initialState: TwoFactorState = {
  isLoading: false,
  error: '',
  success: '',
  setupStep: 'idle',
  setupData: null,
  disableStep: 'idle',
  generateStep: 'idle',
  password: '',
  verificationCode: '',
  disableCode: '',
  generateCode: '',
  backupCodes: [],
  showBackupCodes: false,
};

export function ProfileTwoFactor() {
  const { user, refreshUser } = useAuth();
  const [state, setState] = useState<TwoFactorState>(initialState);

  // Helper function to update state
  const updateState = useCallback((updates: Partial<TwoFactorState>) => {
    setState(prev => ({ ...prev, ...updates }));
  }, []);

  // Reset to initial state
  const resetState = useCallback(() => {
    setState(initialState);
  }, []);

  // Clear messages after a delay
  useEffect(() => {
    if (state.success || state.error) {
      const timer = setTimeout(() => {
        updateState({ success: '', error: '' });
      }, 5000);
      return () => clearTimeout(timer);
    }
  }, [state.success, state.error, updateState]);

  // Handle API errors consistently
  const handleApiError = useCallback((error: any) => {
    const message = error?.response?.data?.error || error?.message || 'An unexpected error occurred';
    updateState({ error: message, isLoading: false });
  }, [updateState]);

  // API call wrapper with error handling
  const apiCall = useCallback(async (
    apiFunction: () => Promise<any>,
    onSuccess?: (data: any) => void,
    onError?: (error: any) => void
  ) => {
    try {
      updateState({ isLoading: true, error: '', success: '' });
      const response = await apiFunction();
      
      if (response.success) {
        onSuccess?.(response);
        updateState({ isLoading: false });
      } else {
        const error = response.error || 'Operation failed';
        updateState({ error, isLoading: false });
        onError?.(error);
      }
    } catch (error) {
      handleApiError(error);
      onError?.(error);
    }
  }, [updateState, handleApiError]);

  // Start 2FA setup
  const startSetup = useCallback(() => {
    resetState();
    updateState({ setupStep: 'password' });
  }, [resetState, updateState]);

  // Handle password submission
  const handlePasswordSubmit = useCallback(async () => {
    if (!state.password.trim()) {
      updateState({ error: 'Password is required' });
      return;
    }

    await apiCall(
      () => api.setup2FA(),
      (response) => {
        updateState({
          setupData: response.data,
          setupStep: 'qr',
          password: state.password // Keep password for final step
        });
      }
    );
  }, [state.password, apiCall, updateState]);

  // Handle verification and enable 2FA
  const handleVerifyAndEnable = useCallback(async () => {
    if (!state.verificationCode || !state.setupData || !state.password) {
      updateState({ error: 'All fields are required' });
      return;
    }

    if (state.verificationCode.length !== 6) {
      updateState({ error: 'Verification code must be 6 digits' });
      return;
    }

    await apiCall(
      () => api.enable2FA({
        secret: state.setupData!.secret,
        token: state.verificationCode,
        password: state.password,
      }),
      (response) => {
        updateState({
          backupCodes: response.data?.backup_codes || [],
          setupStep: 'complete',
          success: 'Two-factor authentication enabled successfully!'
        });
        refreshUser();
      }
    );
  }, [state.verificationCode, state.setupData, state.password, apiCall, updateState, refreshUser]);

  // Start disable 2FA
  const startDisable = useCallback(() => {
    resetState();
    updateState({ disableStep: 'verify' });
  }, [resetState, updateState]);

  // Handle disable 2FA
  const handleDisable = useCallback(async () => {
    if (!state.disableCode.trim()) {
      updateState({ error: 'Verification code is required' });
      return;
    }

    if (state.disableCode.length !== 6) {
      updateState({ error: 'Verification code must be 6 digits' });
      return;
    }

    await apiCall(
      () => api.disable2FA(state.disableCode),
      () => {
        updateState({
          disableStep: 'idle',
          success: 'Two-factor authentication disabled successfully!'
        });
        refreshUser();
      }
    );
  }, [state.disableCode, apiCall, updateState, refreshUser]);

  // Generate backup codes
  const generateBackupCodes = useCallback(async () => {
    updateState({ generateStep: 'verify' });
  }, [updateState]);

  // Handle backup code generation with TOTP
  const handleGenerateBackupCodes = useCallback(async () => {
    if (!state.generateCode.trim()) {
      updateState({ error: 'Verification code is required' });
      return;
    }

    if (state.generateCode.length !== 6) {
      updateState({ error: 'Verification code must be 6 digits' });
      return;
    }

    await apiCall(
      () => api.generateBackupCodes(state.generateCode),
      (response) => {
        updateState({
          backupCodes: response.data?.backup_codes || [],
          setupStep: 'complete', // Reuse the setup complete step to show codes
          success: 'New backup codes generated successfully!',
          generateStep: 'idle',
          generateCode: ''
        });
      }
    );
  }, [state.generateCode, apiCall, updateState]);

  // Download backup codes
  const downloadBackupCodes = useCallback(() => {
    if (!state.backupCodes.length) return;
    
    const text = state.backupCodes.join('\n');
    const blob = new Blob([text], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'orvex-backup-codes.txt';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  }, [state.backupCodes]);

  // Handle keyboard shortcuts
  const handleKeyPress = useCallback((e: React.KeyboardEvent, action: () => void) => {
    if (e.key === 'Enter' && !state.isLoading) {
      e.preventDefault();
      action();
    }
  }, [state.isLoading]);

  // Input change handlers
  const handlePasswordChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    updateState({ password: e.target.value });
  }, [updateState]);

  const handleVerificationCodeChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value.replace(/\D/g, '').slice(0, 6);
    updateState({ verificationCode: value });
  }, [updateState]);

  const handleDisableCodeChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value.replace(/\D/g, '').slice(0, 6);
    updateState({ disableCode: value });
  }, [updateState]);

  const handleGenerateCodeChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value.replace(/\D/g, '').slice(0, 6);
    updateState({ generateCode: value });
  }, [updateState]);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h2 className="text-2xl font-bold text-white mb-2">Two-Factor Authentication</h2>
        <p className="text-gray-400">
          Add an extra layer of security to your account with 2FA.
        </p>
      </div>

      {/* Success Message */}
      {state.success && (
        <div className="p-4 bg-green-500/20 border border-green-500/30 rounded-lg">
          <div className="flex items-center">
            <i className="fas fa-check-circle text-green-400 mr-2"></i>
            <span className="text-green-300">{state.success}</span>
          </div>
        </div>
      )}

      {/* Error Message */}
      {state.error && (
        <div className="p-4 bg-red-500/20 border border-red-500/30 rounded-lg">
          <div className="flex items-center">
            <i className="fas fa-exclamation-circle text-red-400 mr-2"></i>
            <span className="text-red-300">{state.error}</span>
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
          
          {!user?.two_factor_enabled && state.setupStep === 'idle' && (
            <Button
              variant="primary"
              onClick={startSetup}
              disabled={state.isLoading}
              icon="fas fa-shield-alt"
            >
              Enable 2FA
            </Button>
          )}
        </div>
      </div>

      {/* Setup Step 1: Password Confirmation */}
      {state.setupStep === 'password' && (
        <div className="bg-white/5 border border-white/10 rounded-lg p-6">
          <h3 className="text-lg font-semibold text-white mb-4">Step 1: Confirm Your Password</h3>
          <p className="text-gray-400 mb-6">Please enter your current password to continue with 2FA setup.</p>
          
          <div className="max-w-md space-y-4">
            <Input
              type="password"
              placeholder="Enter your password"
              value={state.password}
              onChange={handlePasswordChange}
              onKeyPress={(e) => handleKeyPress(e, handlePasswordSubmit)}
              disabled={state.isLoading}
              autoFocus
            />
            
            <div className="flex space-x-3">
              <Button
                variant="ghost"
                onClick={resetState}
                disabled={state.isLoading}
              >
                Cancel
              </Button>
              <Button
                variant="primary"
                onClick={handlePasswordSubmit}
                loading={state.isLoading}
                disabled={state.isLoading || !state.password.trim()}
              >
                Continue
              </Button>
            </div>
          </div>
        </div>
      )}

      {/* Setup Step 2: QR Code & Verification */}
      {state.setupStep === 'qr' && state.setupData && (
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
                    src={state.setupData.qr_code} 
                    alt="2FA QR Code" 
                    className="w-48 h-48"
                  />
                </div>
              </div>
              
              <div className="bg-gray-800 p-3 rounded-lg">
                <p className="text-xs text-gray-400 mb-1">Manual entry code:</p>
                <code className="text-sm text-gray-300 break-all">{state.setupData.secret}</code>
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
                  value={state.verificationCode}
                  onChange={handleVerificationCodeChange}
                  onKeyPress={(e) => handleKeyPress(e, handleVerifyAndEnable)}
                  maxLength={6}
                  className="text-center text-2xl tracking-wider"
                  disabled={state.isLoading}
                  autoFocus
                />
                
                <div className="flex space-x-3">
                  <Button
                    variant="ghost"
                    onClick={() => updateState({ setupStep: 'password' })}
                    disabled={state.isLoading}
                  >
                    Back
                  </Button>
                  <Button
                    variant="primary"
                    onClick={handleVerifyAndEnable}
                    loading={state.isLoading}
                    disabled={state.isLoading || state.verificationCode.length !== 6}
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

      {/* Setup Step 3: Complete - Show Backup Codes */}
      {state.setupStep === 'complete' && state.backupCodes.length > 0 && (
        <div className="bg-yellow-500/10 border border-yellow-500/30 rounded-lg p-6">
          <div className="flex items-start space-x-3 mb-4">
            <i className="fas fa-exclamation-triangle text-yellow-400 text-lg mt-1"></i>
            <div>
              <h3 className="text-lg font-semibold text-white mb-2">
                {user?.two_factor_enabled ? 'New Backup Codes Generated' : 'Save Your Backup Codes'}
              </h3>
              <p className="text-yellow-200 text-sm">
                These codes can be used to access your account if you lose your authenticator device.
                <strong> Each code can only be used once.</strong>
                {user?.two_factor_enabled && (
                  <span className="block mt-2">
                    <strong>Important:</strong> Your previous backup codes are no longer valid.
                  </span>
                )}
              </p>
            </div>
          </div>
          
          <div className="grid grid-cols-2 gap-2 font-mono text-sm bg-black/30 p-4 rounded-lg mb-4">
            {state.backupCodes.map((code, index) => (
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
              onClick={resetState}
            >
              {user?.two_factor_enabled ? 'Done' : 'Complete Setup'}
            </Button>
          </div>
        </div>
      )}


      {/* Generate Backup Codes */}
      {state.generateStep === 'verify' && (
        <div className="bg-white/5 border border-white/10 rounded-lg p-6">
          <h3 className="text-lg font-semibold text-white mb-4">Generate New Backup Codes</h3>
          
          <div className="bg-yellow-500/10 border border-yellow-500/30 rounded-lg p-4 mb-6">
            <div className="flex items-start space-x-3">
              <i className="fas fa-exclamation-triangle text-yellow-400 text-lg mt-1"></i>
              <div>
                <h4 className="font-medium text-white mb-1">Important</h4>
                <p className="text-yellow-200 text-sm">
                  Generating new backup codes will invalidate all existing backup codes. Save the new codes in a secure location.
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
                value={state.generateCode}
                onChange={handleGenerateCodeChange}
                onKeyPress={(e) => handleKeyPress(e, handleGenerateBackupCodes)}
                maxLength={6}
                className="text-center text-2xl tracking-wider"
                disabled={state.isLoading}
                autoFocus
              />
              <p className="text-xs text-gray-400 mt-2">
                Enter the 6-digit code from your authenticator app
              </p>
            </div>

            <div className="flex space-x-3">
              <Button 
                variant="ghost" 
                onClick={() => updateState({ generateStep: 'idle' })}
                disabled={state.isLoading}
              >
                Cancel
              </Button>
              <Button
                variant="primary"
                onClick={handleGenerateBackupCodes}
                loading={state.isLoading}
                disabled={state.isLoading || state.generateCode.length !== 6}
              >
                Generate New Codes
              </Button>
            </div>
          </div>
        </div>
      )}

      {/* Disable 2FA */}
      {state.disableStep === 'verify' && (
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
                value={state.disableCode}
                onChange={handleDisableCodeChange}
                onKeyPress={(e) => handleKeyPress(e, handleDisable)}
                maxLength={6}
                className="text-center text-2xl tracking-wider"
                disabled={state.isLoading}
                autoFocus
              />
              <p className="text-xs text-gray-400 mt-2">
                Enter the 6-digit code from your authenticator app
              </p>
            </div>

            <div className="flex space-x-3">
              <Button 
                variant="ghost" 
                onClick={() => updateState({ disableStep: 'idle' })}
                disabled={state.isLoading}
              >
                Cancel
              </Button>
              <Button
                variant="danger"
                onClick={handleDisable}
                loading={state.isLoading}
                disabled={state.isLoading || state.disableCode.length !== 6}
              >
                Disable 2FA
              </Button>
            </div>
          </div>
        </div>
      )}

      {/* Manage 2FA (when enabled) */}
      {user?.two_factor_enabled && state.setupStep === 'idle' && state.disableStep === 'idle' && state.generateStep === 'idle' && (
        <div className="bg-white/5 border border-white/10 rounded-lg p-6 space-y-4">
          <h3 className="text-lg font-semibold text-white">Manage 2FA</h3>
          
          <div className="space-y-3">
            <div className="flex items-center justify-between p-4 bg-white/5 rounded-lg">
              <div>
                <h4 className="font-medium text-white">Backup Codes</h4>
                <p className="text-sm text-gray-400">Generate new backup codes to replace existing ones</p>
              </div>
              <Button
                variant="glass"
                onClick={generateBackupCodes}
                loading={state.isLoading}
                disabled={state.isLoading}
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
                onClick={startDisable}
                disabled={state.isLoading}
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