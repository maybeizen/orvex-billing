"use client";

import React, { useState, useRef, useEffect } from "react";
import Image from "next/image";
import Link from "next/link";
import { useAuth } from "@/contexts/auth-context";

export interface ProfileCardProps {
  className?: string;
  compact?: boolean;
}

export const ProfileCard: React.FC<ProfileCardProps> = ({
  className = "",
  compact = false,
}) => {
  const { user, isAuthenticated, logout } = useAuth();
  const [showMenu, setShowMenu] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        setShowMenu(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  if (!isAuthenticated || !user) {
    return null;
  }

  const initials = user.username?.[0]?.toUpperCase() || "U";

  return (
    <div className={`relative ${className}`} ref={menuRef}>
      <div
        className="bg-white/5 border border-white/10 rounded-lg p-2.5 cursor-pointer hover:bg-white/8 transition-all duration-200"
        onClick={() => setShowMenu(!showMenu)}
      >
        <div className="flex items-center gap-3">
          <div className="flex-shrink-0">
            {user.avatar ? (
              <Image
                src={user.avatar}
                alt={user.username || "User"}
                width={28}
                height={28}
                className="rounded-full object-cover"
              />
            ) : (
              <div className="w-7 h-7 bg-blue-600 rounded-full flex items-center justify-center text-white font-medium text-xs">
                {initials}
              </div>
            )}
          </div>

          <div className="flex-1 min-w-0">
            <p className="font-medium text-white truncate text-sm">
              @{user.username}
            </p>
          </div>

          <div className="flex items-center justify-between mr-1">
            <i
              className={`fas fa-ellipsis-vertical text-white/40 text-lg transition-transform duration-200`}
            />
          </div>
        </div>
      </div>

      <div
        className={`absolute bottom-full left-0 right-0 mb-2 bg-black/90 backdrop-blur-md border border-white/20 rounded-lg shadow-2xl z-50 overflow-hidden transition-all duration-300 ease-out transform origin-bottom ${
          showMenu
            ? "opacity-100 scale-100 translate-y-0"
            : "opacity-0 scale-95 translate-y-2 pointer-events-none"
        }`}
      >
        <div className="py-1.5">
          <Link
            href="/dashboard/account"
            className="flex items-center gap-2.5 px-3 py-2 text-white/80 hover:bg-white/10 hover:text-white transition-colors duration-150 text-sm"
            onClick={() => setShowMenu(false)}
          >
            <i className="fas fa-user-circle w-3.5 text-center" />
            <span>Account Settings</span>
          </Link>

          {user.role === "admin" && (
            <Link
              href="/admin/dashboard"
              className="flex items-center gap-2.5 px-3 py-2 text-white/80 hover:bg-white/10 hover:text-white transition-colors duration-150 text-sm"
              onClick={() => setShowMenu(false)}
            >
              <i className="fas fa-crown w-3.5 text-center" />
              <span>Admin Panel</span>
            </Link>
          )}

          <div className="border-t border-white/10 my-1" />

          <button
            onClick={() => {
              setShowMenu(false);
              logout();
            }}
            className="w-full flex items-center gap-2.5 px-3 py-2 text-red-400 hover:bg-red-500/10 hover:text-red-300 transition-colors duration-150 text-sm"
          >
            <i className="fas fa-sign-out-alt w-3.5 text-center" />
            <span>Sign Out</span>
          </button>
        </div>
      </div>
    </div>
  );
};

export default ProfileCard;
