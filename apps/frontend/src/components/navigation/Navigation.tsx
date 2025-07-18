"use client";

import { ReactNode } from "react";
import Link from "next/link";

interface NavigationProps {
  children: ReactNode;
  className?: string;
  brand?: {
    text: string;
    subtitle?: string;
    href?: string;
    icon?: string;
  };
}

export default function Navigation({
  children,
  className = "",
  brand = {
    text: "Orvex",
    subtitle: "Dashboard",
    href: "/dashboard",
    icon: "fas fa-bolt",
  },
}: NavigationProps) {
  return (
    <aside className={`fixed left-0 top-0 h-screen w-64 bg-black border-r border-white/10 ${className}`}>
      <div className="flex flex-col h-full">
        <header className="p-6 border-b border-white/10">
          <Link href={brand.href || "/dashboard"} className="flex items-center space-x-3">
            <div className="w-10 h-10 bg-gradient-to-r from-blue-400 via-purple-400 to-cyan-400 rounded-xl flex items-center justify-center">
              <i className={`${brand.icon} text-white text-xl`}></i>
            </div>
            <div>
              <h1 className="text-xl font-black text-transparent bg-clip-text bg-gradient-to-r from-blue-400 via-purple-400 to-cyan-400">
                {brand.text}
              </h1>
              {brand.subtitle && (
                <p className="text-xs text-gray-400">{brand.subtitle}</p>
              )}
            </div>
          </Link>
        </header>

        <nav className="flex-1 overflow-y-auto py-6 px-4 flex flex-col">
          {children}
        </nav>
      </div>
    </aside>
  );
}