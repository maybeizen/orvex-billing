"use client";

import React, { useState } from "react";

export interface PasswordStrengthProps {
  password: string;
  showDetails?: boolean;
  className?: string;
  showPopup?: boolean;
  isFocused?: boolean;
}

export type PasswordStrengthLevel = "weak" | "fair" | "good" | "strong";

interface PasswordRequirement {
  label: string;
  regex: RegExp;
  met: boolean;
}

export const calculatePasswordStrength = (password: string) => {
  if (!password) {
    return {
      score: 0,
      level: "weak" as PasswordStrengthLevel,
      percentage: 0,
      requirements: [],
    };
  }

  const requirements: PasswordRequirement[] = [
    {
      label: "At least 8 characters",
      regex: /.{8,}/,
      met: /.{8,}/.test(password),
    },
    {
      label: "At least one lowercase letter",
      regex: /[a-z]/,
      met: /[a-z]/.test(password),
    },
    {
      label: "At least one uppercase letter",
      regex: /[A-Z]/,
      met: /[A-Z]/.test(password),
    },
    {
      label: "At least one number",
      regex: /\d/,
      met: /\d/.test(password),
    },
    {
      label: "At least one special character",
      regex: /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/,
      met: /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(password),
    },
  ];

  const metRequirements = requirements.filter((req) => req.met).length;
  const score = metRequirements;
  const percentage = (score / requirements.length) * 100;

  let level: PasswordStrengthLevel;
  if (score <= 1) level = "weak";
  else if (score <= 2) level = "fair";
  else if (score <= 4) level = "good";
  else level = "strong";

  return {
    score,
    level,
    percentage,
    requirements,
  };
};

const strengthConfig = {
  weak: {
    color: "bg-red-500",
    textColor: "text-red-400",
    label: "Weak",
  },
  fair: {
    color: "bg-orange-500",
    textColor: "text-orange-400",
    label: "Fair",
  },
  good: {
    color: "bg-yellow-500",
    textColor: "text-yellow-400",
    label: "Good",
  },
  strong: {
    color: "bg-green-500",
    textColor: "text-green-400",
    label: "Strong",
  },
};

export const PasswordStrength: React.FC<PasswordStrengthProps> = ({
  password,
  showDetails = false,
  showPopup = false,
  isFocused = false,
  className = "",
}) => {
  const strength = calculatePasswordStrength(password);

  const showStrengthBar = password.length > 0;

  const config = showStrengthBar
    ? strengthConfig[strength.level]
    : { color: "bg-white/10", textColor: "text-white/40", label: "" };

  return (
    <div className={`mt-2 relative ${className}`}>
      {showStrengthBar && (
        <div className="flex items-center gap-2 mb-2">
          <div className="flex-1 h-2 bg-white/10 rounded-full overflow-hidden">
            <div
              className={`h-full ${config.color} transition-all duration-300 ease-out rounded-full`}
              style={{ width: `${strength.percentage}%` }}
            />
          </div>
          <div className="flex items-center gap-2">
            <span className={`text-xs font-medium ${config.textColor}`}>
              {config.label}
            </span>
          </div>
        </div>
      )}

      {showPopup && isFocused && showStrengthBar && (
        <div className="absolute top-full left-0 mt-2 z-50 w-72 p-4 bg-neutral-800 backdrop-blur-sm border border-neutral-700 rounded-lg shadow-xl">
          <div className="space-y-2">
            <p className="text-xs font-medium text-white/80 mb-3">
              Password requirements:
            </p>
            {strength.requirements.map((req, index) => (
              <div key={index} className="flex items-center gap-2 text-xs">
                <i
                  className={`fas ${
                    req.met
                      ? "fa-check text-green-400"
                      : "fa-times text-red-400"
                  } text-xs flex-shrink-0`}
                />
                <span
                  className={`${
                    req.met ? "text-green-400" : "text-white/60"
                  } transition-colors duration-200`}
                >
                  {req.label}
                </span>
              </div>
            ))}
          </div>
          <div className="absolute -top-1 left-6 w-2 h-2 bg-neutral-800 border-l border-t border-neutral-700 transform rotate-45"></div>
        </div>
      )}

      {showDetails && !showPopup && showStrengthBar && (
        <div className="space-y-1">
          <p className="text-xs text-white/60 mb-1">Password requirements:</p>
          {strength.requirements.map((req, index) => (
            <div key={index} className="flex items-center gap-2 text-xs">
              <i
                className={`fas ${
                  req.met ? "fa-check text-green-400" : "fa-times text-red-400"
                } text-xs`}
              />
              <span
                className={`${
                  req.met ? "text-green-400" : "text-white/60"
                } transition-colors duration-200`}
              >
                {req.label}
              </span>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};
