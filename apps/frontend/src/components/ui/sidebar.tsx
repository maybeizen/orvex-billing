"use client";

import React, { ReactNode } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { ProfileCard } from "./profile-card";
import Image from "next/image";

export interface SidebarItemProps {
  icon?: string;
  label: string;
  href?: string;
  active?: boolean;
  disabled?: boolean;
  badge?: string | number;
  onClick?: () => void;
  className?: string;
}

export interface SidebarCategoryProps {
  title: string;
  children: ReactNode;
  collapsible?: boolean;
  defaultCollapsed?: boolean;
  className?: string;
}

export interface SidebarProps {
  children: ReactNode;
  className?: string;
  width?: "sm" | "md" | "lg";
  showProfile?: boolean;
}

export const SidebarItem: React.FC<SidebarItemProps> = ({
  icon,
  label,
  href,
  active: forcedActive,
  disabled = false,
  badge,
  onClick,
  className = "",
}) => {
  const pathname = usePathname();
  const isActive =
    forcedActive !== undefined ? forcedActive : href === pathname;

  const baseClasses = `
    flex items-center gap-3 px-3 py-2 rounded text-sm transition-all duration-200
    ${disabled ? "opacity-50 cursor-not-allowed" : "cursor-pointer"}
    ${
      isActive
        ? "bg-violet-500/15 text-white"
        : "text-white/80 hover:bg-white/10 hover:translate-x-0.5"
    }
    ${className}
  `;

  const content = (
    <>
      {icon && (
        <div className="flex-shrink-0 w-5 flex justify-center">
          <i className={`${icon} text-base`} />
        </div>
      )}

      <span className="flex-1 truncate">{label}</span>

      {badge && (
        <div className="flex-shrink-0">
          <span className="inline-flex items-center justify-center px-2 py-1 text-xs font-medium rounded-full bg-white/20 text-white/80">
            {badge}
          </span>
        </div>
      )}
    </>
  );

  if (href && !disabled) {
    return (
      <Link href={href} className={baseClasses}>
        {content}
      </Link>
    );
  }

  return (
    <div className={baseClasses} onClick={disabled ? undefined : onClick}>
      {content}
    </div>
  );
};

export const SidebarCategory: React.FC<SidebarCategoryProps> = ({
  title,
  children,
  collapsible = false,
  defaultCollapsed = false,
  className = "",
}) => {
  const [isCollapsed, setIsCollapsed] = React.useState(defaultCollapsed);
  const [isAnimating, setIsAnimating] = React.useState(false);

  const toggleCollapsed = () => {
    if (collapsible) {
      setIsAnimating(true);
      setIsCollapsed(!isCollapsed);
      setTimeout(() => setIsAnimating(false), 300);
    }
  };

  return (
    <div className={`space-y-2 ${className}`}>
      <div
        className={`
          flex items-center justify-between px-3 py-2 text-xs font-semibold text-white/50 uppercase tracking-wide transition-colors duration-200
          ${collapsible ? "cursor-pointer hover:text-white/70" : ""}
        `}
        onClick={toggleCollapsed}
      >
        <span>{title}</span>
        {collapsible && (
          <i
            className={`fas fa-chevron-${
              isCollapsed ? "right" : "down"
            } text-xs transition-transform duration-200 ${
              isAnimating ? "scale-110" : ""
            }`}
          />
        )}
      </div>

      <div
        className={`overflow-hidden transition-all duration-300 ease-in-out ${
          !collapsible || !isCollapsed
            ? "max-h-96 opacity-100"
            : "max-h-0 opacity-0"
        }`}
      >
        <div className="space-y-1 ml-2">{children}</div>
      </div>
    </div>
  );
};

export const Sidebar: React.FC<SidebarProps> = ({
  children,
  className = "",
  width = "md",
  showProfile = true,
}) => {
  const widthClasses = {
    sm: "w-56",
    md: "w-72",
    lg: "w-80",
  };

  return (
    <div
      className={`
      ${widthClasses[width]} h-screen fixed left-0 top-0 z-40
      bg-black/70 backdrop-blur-lg border-r border-white/20
      flex flex-col shadow-2xl
      ${className}
    `}
    >
      <div className="flex-shrink-0 p-4 border-b border-white/10">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-lg flex items-center justify-center">
            <Image src="/orvex.webp" alt="Orvex" width={40} height={40} />
          </div>
          <div>
            <h2 className="text-white font-semibold text-lg">Orvex</h2>
            <p className="text-white/40 text-xs">orvex.cc</p>
          </div>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto py-4 px-4 space-y-3 custom-scrollbar">
        {children}
      </div>

      {showProfile && (
        <div className="flex-shrink-0 p-4 border-t border-white/20 bg-black/30">
          <ProfileCard compact={width === "sm"} />
        </div>
      )}

      <style jsx global>{`
        .custom-scrollbar::-webkit-scrollbar {
          width: 4px;
        }
        .custom-scrollbar::-webkit-scrollbar-track {
          background: transparent;
        }
        .custom-scrollbar::-webkit-scrollbar-thumb {
          background: rgba(255, 255, 255, 0.1);
          border-radius: 2px;
        }
        .custom-scrollbar::-webkit-scrollbar-thumb:hover {
          background: rgba(255, 255, 255, 0.2);
        }
      `}</style>
    </div>
  );
};

export default Sidebar;
