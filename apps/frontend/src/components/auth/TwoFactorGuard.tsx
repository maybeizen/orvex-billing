"use client";

import { useEffect } from "react";
import { useRouter, usePathname } from "next/navigation";
import { useAuth } from "@/lib/auth-context";
import { LoadingScreen } from "@/components/ui/loading-spinner";

interface TwoFactorGuardProps {
  children: React.ReactNode;
}

export function TwoFactorGuard({ children }: TwoFactorGuardProps) {
  const { user, loading, requires2FA } = useAuth();
  const router = useRouter();
  const pathname = usePathname();

  useEffect(() => {
    // Don't redirect if still loading
    if (loading) return;

    // Don't redirect if user is not logged in (let other auth guards handle this)
    if (!user) return;

    // Don't redirect if we're already on the 2FA page
    if (pathname === "/auth/2fa") return;

    // Don't redirect for auth pages
    if (pathname.startsWith("/auth/")) return;

    // If user requires 2FA verification, redirect to 2FA page
    if (requires2FA) {
      router.push(`/auth/2fa?returnTo=${encodeURIComponent(pathname)}`);
      return;
    }
  }, [user, loading, requires2FA, router, pathname]);

  // Show loading screen while checking auth status
  if (loading) {
    return <LoadingScreen message="Verifying authentication..." />;
  }

  // If user requires 2FA and we're not on the 2FA page, don't render content
  if (requires2FA && pathname !== "/auth/2fa") {
    return <LoadingScreen message="Redirecting to two-factor authentication..." />;
  }

  return <>{children}</>;
}