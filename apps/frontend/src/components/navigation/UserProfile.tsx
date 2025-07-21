"use client";

import { useState } from "react";
import Link from "next/link";

interface UserProfileProps {
  user?: {
    first_name?: string;
    last_name?: string;
    email?: string;
    avatar_url?: string;
    role?: string;
  };
  onLogout?: () => void;
  className?: string;
}

export default function UserProfile({
  user,
  onLogout,
  className = "",
}: UserProfileProps) {
  const [showUserMenu, setShowUserMenu] = useState(false);

  const toggleUserMenu = () => {
    setShowUserMenu(!showUserMenu);
  };

  const handleLogout = () => {
    setShowUserMenu(false);
    onLogout?.();
  };

  const getInitials = () => {
    return user?.first_name?.[0]?.toUpperCase() || "U";
  };

  const getDisplayName = () => {
    return (
      `${user?.first_name || ""} ${user?.last_name || ""}`.trim() || "User"
    );
  };

  const getAvatarUrl = (avatarUrl?: string) => {
    if (!avatarUrl) return null;
    if (avatarUrl.startsWith("http")) return avatarUrl;
    // Construct full URL for avatar - remove /api from base URL and add avatar path
    const baseUrl =
      process.env.NEXT_PUBLIC_API_URL || "http://localhost:3000/api";
    return baseUrl.replace("/api", "") + avatarUrl;
  };

  return (
    <div className={`p-4 border-t border-white/10 ${className}`}>
      <div className="relative">
        <button
          onClick={toggleUserMenu}
          className="flex items-center w-full p-3 bg-white/5 rounded-xl hover:bg-white/10 transition-all duration-200"
        >
          <div className="flex items-center flex-1">
            {user?.avatar_url ? (
              <img
                src={getAvatarUrl(user.avatar_url) || undefined}
                alt={getDisplayName()}
                className="w-10 h-10 rounded-lg border border-white/20"
              />
            ) : (
              <div className="w-10 h-10 bg-gradient-to-r from-blue-400 via-purple-400 to-cyan-400 rounded-lg flex items-center justify-center">
                <span className="text-white font-bold text-sm">
                  {getInitials()}
                </span>
              </div>
            )}
            <div className="ml-3 flex-1 text-left">
              <p className="text-sm font-medium text-white truncate max-w-[90px]">
                {getDisplayName()}
              </p>
              <p className="text-xs text-gray-400 truncate max-w-[90px]">
                {user?.email}
              </p>
            </div>
          </div>
          <i
            className={`fas fa-chevron-up transition-transform duration-200 text-gray-400 ${
              showUserMenu ? "rotate-180" : ""
            }`}
          ></i>
        </button>

        {showUserMenu && (
          <div className="absolute bottom-full left-0 right-0 mb-2 bg-white/10 backdrop-blur-sm border border-white/20 rounded-xl p-2 shadow-xl">
            <Link
              href="/dashboard/profile"
              className="flex items-center w-full px-3 py-2 text-sm text-gray-300 hover:text-white hover:bg-white/5 rounded-lg transition-all duration-200"
              onClick={() => setShowUserMenu(false)}
            >
              <i className="fas fa-user mr-3 w-4 text-center"></i>
              Profile Settings
            </Link>
            {user?.role === "admin" && (
              <Link
                href="/admin"
                className="flex items-center w-full px-3 py-2 text-sm text-gray-300 hover:text-white hover:bg-white/5 rounded-lg transition-all duration-200"
                onClick={() => setShowUserMenu(false)}
              >
                <i className="fas fa-crown mr-3 w-4 text-center"></i>
                Admin Panel
              </Link>
            )}
            {onLogout && (
              <button
                onClick={handleLogout}
                className="flex items-center w-full px-3 py-2 text-sm text-red-400 hover:text-red-300 hover:bg-white/5 rounded-lg transition-all duration-200"
              >
                <i className="fas fa-sign-out-alt mr-3 w-4 text-center"></i>
                Sign Out
              </button>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
