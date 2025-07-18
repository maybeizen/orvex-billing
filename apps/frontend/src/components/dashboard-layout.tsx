"use client";

import { ReactNode } from "react";
import DashboardSidebar from "./dashboard-sidebar";
import { TwoFactorGuard } from "./auth/TwoFactorGuard";

interface DashboardLayoutProps {
  children: ReactNode;
}

export default function DashboardLayout({ children }: DashboardLayoutProps) {
  return (
    <TwoFactorGuard>
      <div className="min-h-screen bg-black">
        <div className="flex">
          {/* Sidebar */}
          <DashboardSidebar />

          {/* Main content */}
          <div className="flex-1 ml-64">
            <main className="p-6 lg:p-8">
              {children}
            </main>
          </div>
        </div>
      </div>
    </TwoFactorGuard>
  );
}