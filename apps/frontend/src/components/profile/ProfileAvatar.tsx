"use client";

import { useState, useRef } from "react";
import { useAuth } from "@/lib/auth-context";
import { Button } from "@/components/ui/button";
import { api } from "@/lib/api";

export function ProfileAvatar() {
  const { user, refreshUser } = useAuth();
  const [loading, setLoading] = useState(false);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [success, setSuccess] = useState("");
  const [error, setError] = useState("");
  const fileInputRef = useRef<HTMLInputElement>(null);

  const getInitials = () => {
    const first = user?.first_name?.[0]?.toUpperCase() || "";
    const last = user?.last_name?.[0]?.toUpperCase() || "";
    return first + last || "U";
  };

  const getDisplayName = () => {
    return `${user?.first_name || ""} ${user?.last_name || ""}`.trim() || "User";
  };

  const getAvatarUrl = (avatarUrl?: string) => {
    if (!avatarUrl) return null;
    if (avatarUrl.startsWith('http')) return avatarUrl;
    // Construct full URL for avatar - remove /api from base URL and add avatar path
    const baseUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:3000/api";
    return baseUrl.replace('/api', '') + avatarUrl;
  };

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    // Validate file type
    if (!file.type.startsWith("image/")) {
      setError("Please select a valid image file");
      return;
    }

    // Validate file size (5MB max)
    if (file.size > 5 * 1024 * 1024) {
      setError("Image must be smaller than 5MB");
      return;
    }

    setSelectedFile(file);
    setError("");
    setSuccess("");

    // Create preview
    const reader = new FileReader();
    reader.onload = () => {
      setPreviewUrl(reader.result as string);
    };
    reader.readAsDataURL(file);
  };

  const handleUpload = async () => {
    if (!selectedFile) return;

    setLoading(true);
    setError("");
    setSuccess("");

    try {
      const formData = new FormData();
      formData.append("avatar", selectedFile);

      const response = await api.uploadAvatar(formData);
      if (response.success) {
        setSuccess("Avatar updated successfully!");
        setPreviewUrl(null);
        setSelectedFile(null);
        await refreshUser();
      } else {
        setError(response.error || "Failed to upload avatar");
      }
    } catch (err: any) {
      setError(err.message || "Failed to upload avatar");
    } finally {
      setLoading(false);
    }
  };

  const handleRemoveAvatar = async () => {
    setLoading(true);
    setError("");
    setSuccess("");

    try {
      const response = await api.removeAvatar();
      if (response.success) {
        setSuccess("Avatar removed successfully!");
        await refreshUser();
      } else {
        setError(response.error || "Failed to remove avatar");
      }
    } catch (err: any) {
      setError(err.message || "Failed to remove avatar");
    } finally {
      setLoading(false);
    }
  };

  const handleCancel = () => {
    setSelectedFile(null);
    setPreviewUrl(null);
    setError("");
    setSuccess("");
    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold text-white mb-2">Profile Avatar</h2>
        <p className="text-gray-400">
          Upload a profile picture to personalize your account.
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

      <div className="bg-white/5 border border-white/10 rounded-lg p-6">
        <div className="flex flex-col md:flex-row items-start md:items-center space-y-6 md:space-y-0 md:space-x-8">
          {/* Current Avatar */}
          <div className="flex flex-col items-center space-y-4">
            <div className="relative">
              {previewUrl ? (
                <img
                  src={previewUrl}
                  alt="Preview"
                  className="w-24 h-24 rounded-xl object-cover border-2 border-white/20"
                />
              ) : user?.avatar_url ? (
                <img
                  src={getAvatarUrl(user.avatar_url) || undefined}
                  alt={getDisplayName()}
                  className="w-24 h-24 rounded-xl object-cover border-2 border-white/20"
                />
              ) : (
                <div className="w-24 h-24 bg-gradient-to-r from-blue-400 via-purple-400 to-cyan-400 rounded-xl flex items-center justify-center">
                  <span className="text-white font-bold text-2xl">
                    {getInitials()}
                  </span>
                </div>
              )}
            </div>
            <p className="text-sm text-gray-400 text-center">Current Avatar</p>
          </div>

          {/* Upload Controls */}
          <div className="flex-1 space-y-4">
            <div>
              <input
                ref={fileInputRef}
                type="file"
                accept="image/*"
                onChange={handleFileSelect}
                className="hidden"
              />
              
              <div className="space-y-3">
                <div className="flex flex-wrap gap-3">
                  <Button
                    type="button"
                    variant="glass"
                    onClick={() => fileInputRef.current?.click()}
                    disabled={loading}
                    icon="fas fa-upload"
                  >
                    Choose Image
                  </Button>
                  
                  {selectedFile && (
                    <>
                      <Button
                        type="button"
                        variant="primary"
                        onClick={handleUpload}
                        loading={loading}
                        disabled={loading}
                        icon="fas fa-save"
                      >
                        Upload
                      </Button>
                      
                      <Button
                        type="button"
                        variant="ghost"
                        onClick={handleCancel}
                        disabled={loading}
                      >
                        Cancel
                      </Button>
                    </>
                  )}
                  
                  {user?.avatar_url && !selectedFile && (
                    <Button
                      type="button"
                      variant="danger"
                      onClick={handleRemoveAvatar}
                      loading={loading}
                      disabled={loading}
                      icon="fas fa-trash"
                    >
                      Remove
                    </Button>
                  )}
                </div>

                {selectedFile && (
                  <div className="text-sm text-gray-400">
                    <p>Selected: {selectedFile.name}</p>
                    <p>Size: {(selectedFile.size / 1024 / 1024).toFixed(2)} MB</p>
                  </div>
                )}
              </div>
            </div>

            <div className="text-sm text-gray-400 space-y-1">
              <p><strong>Requirements:</strong></p>
              <ul className="list-disc list-inside space-y-1 ml-4">
                <li>Maximum file size: 5MB</li>
                <li>Supported formats: JPG, PNG, WebP</li>
                <li>Recommended size: 400x400 pixels</li>
                <li>Square images work best</li>
              </ul>
            </div>
          </div>
        </div>
      </div>

      {/* Avatar Preview Gallery */}
      <div className="bg-white/5 border border-white/10 rounded-lg p-6">
        <h3 className="text-lg font-semibold text-white mb-4">Preview</h3>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="text-center">
            <div className="flex items-center justify-center mb-2">
              {user?.avatar_url ? (
                <img
                  src={getAvatarUrl(user.avatar_url) || undefined}
                  alt="Small"
                  className="w-8 h-8 rounded-lg object-cover"
                />
              ) : (
                <div className="w-8 h-8 bg-gradient-to-r from-blue-400 via-purple-400 to-cyan-400 rounded-lg flex items-center justify-center">
                  <span className="text-white font-bold text-xs">
                    {getInitials()}
                  </span>
                </div>
              )}
            </div>
            <p className="text-xs text-gray-400">Small (32px)</p>
          </div>

          <div className="text-center">
            <div className="flex items-center justify-center mb-2">
              {user?.avatar_url ? (
                <img
                  src={getAvatarUrl(user.avatar_url) || undefined}
                  alt="Medium"
                  className="w-12 h-12 rounded-lg object-cover"
                />
              ) : (
                <div className="w-12 h-12 bg-gradient-to-r from-blue-400 via-purple-400 to-cyan-400 rounded-lg flex items-center justify-center">
                  <span className="text-white font-bold text-sm">
                    {getInitials()}
                  </span>
                </div>
              )}
            </div>
            <p className="text-xs text-gray-400">Medium (48px)</p>
          </div>

          <div className="text-center">
            <div className="flex items-center justify-center mb-2">
              {user?.avatar_url ? (
                <img
                  src={getAvatarUrl(user.avatar_url) || undefined}
                  alt="Large"
                  className="w-16 h-16 rounded-lg object-cover"
                />
              ) : (
                <div className="w-16 h-16 bg-gradient-to-r from-blue-400 via-purple-400 to-cyan-400 rounded-lg flex items-center justify-center">
                  <span className="text-white font-bold text-lg">
                    {getInitials()}
                  </span>
                </div>
              )}
            </div>
            <p className="text-xs text-gray-400">Large (64px)</p>
          </div>
        </div>
      </div>
    </div>
  );
}