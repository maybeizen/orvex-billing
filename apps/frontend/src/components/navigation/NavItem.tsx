"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

interface NavItemProps {
  label: string;
  href: string;
  icon: string;
  className?: string;
  variant?: "default" | "nested";
}

export default function NavItem({
  label,
  href,
  icon,
  className = "",
  variant = "default",
}: NavItemProps) {
  const pathname = usePathname();
  const isActive = pathname === href;

  const baseClasses = "flex items-center transition-all duration-200";
  const variantClasses = {
    default: `px-4 py-3 text-sm font-medium rounded-xl mb-2 ${
      isActive
        ? "bg-white/10 text-white border border-white/20"
        : "text-gray-300 hover:text-white hover:bg-white/5"
    }`,
    nested: `px-4 py-2 text-sm rounded-lg ${
      isActive
        ? "bg-white/10 text-white border border-white/20"
        : "text-gray-400 hover:text-white hover:bg-white/5"
    }`,
  };

  const iconClasses = variant === "default" ? "mr-3 w-5 text-center" : "mr-3 w-4 text-center text-xs";

  return (
    <Link
      href={href}
      className={`${baseClasses} ${variantClasses[variant]} ${className}`}
    >
      <i className={`${icon} ${iconClasses}`}></i>
      {label}
    </Link>
  );
}