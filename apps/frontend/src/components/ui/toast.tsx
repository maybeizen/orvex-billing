"use client";

import React, { useEffect, useState } from "react";

export type ToastType = "success" | "error" | "warning" | "info";
export type ToastPosition =
  | "top-left"
  | "top-right"
  | "bottom-left"
  | "bottom-right"
  | "top-center"
  | "bottom-center";

export interface ToastProps {
  id: string;
  type: ToastType;
  title?: string;
  message: string;
  duration?: number;
  position?: ToastPosition;
  disabled?: boolean;
  onClose: (id: string) => void;
}

const toastStyles = {
  success: {
    bg: "bg-green-500/10 border-green-500/20",
    icon: "fas fa-check-circle text-green-400",
    text: "text-green-400",
    progress: "bg-green-400",
  },
  error: {
    bg: "bg-red-500/10 border-red-500/20",
    icon: "fas fa-exclamation-circle text-red-400",
    text: "text-red-400",
    progress: "bg-red-400",
  },
  warning: {
    bg: "bg-yellow-500/10 border-yellow-500/20",
    icon: "fas fa-exclamation-triangle text-yellow-400",
    text: "text-yellow-400",
    progress: "bg-yellow-400",
  },
  info: {
    bg: "bg-blue-500/10 border-blue-500/20",
    icon: "fas fa-info-circle text-blue-400",
    text: "text-blue-400",
    progress: "bg-blue-400",
  },
};

export const Toast: React.FC<ToastProps> = ({
  id,
  type,
  title,
  message,
  duration = 5000,
  disabled = false,
  onClose,
}) => {
  const [progress, setProgress] = useState(100);
  const [isVisible, setIsVisible] = useState(false);
  const [isExiting, setIsExiting] = useState(false);

  const styles = toastStyles[type];

  useEffect(() => {
    const showTimer = setTimeout(() => setIsVisible(true), 10);

    if (!disabled && duration > 0) {
      const progressTimer = setInterval(() => {
        setProgress((prev) => {
          const newProgress = prev - 100 / (duration / 100);
          if (newProgress <= 0) {
            clearInterval(progressTimer);
            handleClose();
            return 0;
          }
          return newProgress;
        });
      }, 100);

      return () => {
        clearTimeout(showTimer);
        clearInterval(progressTimer);
      };
    }

    return () => clearTimeout(showTimer);
  }, [duration, disabled]);

  const handleClose = () => {
    setIsExiting(true);
    setTimeout(() => onClose(id), 300);
  };

  return (
    <div
      className={`
        transform transition-all duration-300 ease-in-out
        ${
          isVisible && !isExiting
            ? "translate-x-0 opacity-100 scale-100"
            : "translate-x-full opacity-0 scale-95"
        }
        relative max-w-sm w-full
        ${styles.bg} backdrop-blur-sm border rounded-lg shadow-lg
        p-4 mb-3
        ${disabled ? "opacity-60" : ""}
      `}
    >
      <button
        onClick={handleClose}
        className="absolute top-2 right-2 p-1 text-white/40 hover:text-white/60 transition-colors"
        disabled={disabled}
      >
        <i className="fas fa-times text-sm" />
      </button>

      <div className="flex items-start gap-3 pr-6">
        <div className="flex-shrink-0 mt-0.5">
          <i className={`${styles.icon} text-lg`} />
        </div>

        <div className="flex-1 min-w-0">
          {title && (
            <h4 className="text-sm font-semibold text-white mb-1 leading-tight">
              {title}
            </h4>
          )}
          <p className="text-sm text-white/80 leading-relaxed break-words">
            {message}
          </p>
        </div>
      </div>

      {!disabled && duration > 0 && (
        <div className="absolute bottom-0 left-0 right-0 h-1 bg-white/10 rounded-b-lg overflow-hidden">
          <div
            className={`h-full ${styles.progress} transition-all duration-100 ease-linear rounded-b-lg`}
            style={{ width: `${progress}%` }}
          />
        </div>
      )}
    </div>
  );
};
