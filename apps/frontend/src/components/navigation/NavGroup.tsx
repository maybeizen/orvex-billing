"use client";

import { useState, ReactNode } from "react";

interface NavGroupProps {
  title: string;
  children: ReactNode;
  icon: string;
  defaultExpanded?: boolean;
  className?: string;
}

export default function NavGroup({
  title,
  children,
  icon,
  defaultExpanded = false,
  className = "",
}: NavGroupProps) {
  const [isExpanded, setIsExpanded] = useState(defaultExpanded);

  const toggleExpanded = () => {
    setIsExpanded(!isExpanded);
  };

  return (
    <div className={className}>
      <button
        onClick={toggleExpanded}
        className="flex items-center justify-between w-full px-4 py-3 text-sm font-medium text-gray-300 hover:text-white hover:bg-white/5 rounded-xl transition-all duration-200"
      >
        <div className="flex items-center">
          <i className={`${icon} mr-3 w-5 text-center`}></i>
          {title}
        </div>
        <i
          className={`fas fa-chevron-down transition-transform duration-200 ${
            isExpanded ? "rotate-180" : ""
          }`}
        ></i>
      </button>
      {isExpanded && (
        <div className="mt-2 ml-4 space-y-1">
          {children}
        </div>
      )}
    </div>
  );
}