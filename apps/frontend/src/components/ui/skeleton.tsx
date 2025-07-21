import React, { useEffect, useState } from "react";

interface SkeletonProps {
  className?: string;
  fakeDelay?: number;
  children?: React.ReactNode;
}

export const Skeleton: React.FC<SkeletonProps> = ({
  className = "",
  fakeDelay = 0,
  children,
}) => {
  const [isDelayed, setIsDelayed] = useState(fakeDelay > 0);

  useEffect(() => {
    if (fakeDelay > 0) {
      const timer = setTimeout(() => {
        setIsDelayed(false);
      }, fakeDelay);

      return () => clearTimeout(timer);
    }
  }, [fakeDelay]);

  if (isDelayed) {
    return (
      <div
        className={`animate-pulse bg-white/10 rounded ${className}`}
      />
    );
  }

  return <>{children}</>;
};

interface SkeletonTextProps {
  lines?: number;
  className?: string;
  fakeDelay?: number;
  children?: React.ReactNode;
}

export const SkeletonText: React.FC<SkeletonTextProps> = ({
  lines = 1,
  className = "",
  fakeDelay = 0,
  children,
}) => {
  const [isDelayed, setIsDelayed] = useState(fakeDelay > 0);

  useEffect(() => {
    if (fakeDelay > 0) {
      const timer = setTimeout(() => {
        setIsDelayed(false);
      }, fakeDelay);

      return () => clearTimeout(timer);
    }
  }, [fakeDelay]);

  if (isDelayed) {
    return (
      <div className={`space-y-2 ${className}`}>
        {Array.from({ length: lines }).map((_, i) => (
          <div
            key={i}
            className={`animate-pulse bg-white/10 rounded h-4 ${
              i === lines - 1 && lines > 1 ? "w-3/4" : "w-full"
            }`}
          />
        ))}
      </div>
    );
  }

  return <>{children}</>;
};

interface SkeletonCardProps {
  className?: string;
  fakeDelay?: number;
  children?: React.ReactNode;
}

export const SkeletonCard: React.FC<SkeletonCardProps> = ({
  className = "",
  fakeDelay = 0,
  children,
}) => {
  const [isDelayed, setIsDelayed] = useState(fakeDelay > 0);

  useEffect(() => {
    if (fakeDelay > 0) {
      const timer = setTimeout(() => {
        setIsDelayed(false);
      }, fakeDelay);

      return () => clearTimeout(timer);
    }
  }, [fakeDelay]);

  if (isDelayed) {
    return (
      <div className={`animate-pulse bg-white/5 rounded-lg p-4 border border-white/10 ${className}`}>
        <div className="flex items-center space-x-4">
          <div className="w-10 h-10 bg-white/10 rounded-lg" />
          <div className="flex-1 space-y-2">
            <div className="h-4 bg-white/10 rounded w-1/2" />
            <div className="h-3 bg-white/10 rounded w-1/3" />
          </div>
          <div className="w-16 h-6 bg-white/10 rounded" />
        </div>
      </div>
    );
  }

  return <>{children}</>;
};