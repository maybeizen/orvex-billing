"use client";

import { ReactNode } from "react";
import DashboardSidebar from "./dashboard-sidebar";

interface DashboardLayoutProps {
  children: ReactNode;
}

export default function DashboardLayout({ children }: DashboardLayoutProps) {
  return (
    <div className="min-h-screen bg-black">
      <div className="flex">
        <DashboardSidebar />
        <div className="flex-1 ml-64">
          <main className="p-6 lg:p-8">{children}</main>
        </div>
      </div>
    </div>
  );
}
