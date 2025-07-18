"use client";

import { useAuth } from "@/lib/auth-context";
import { useRouter } from "next/navigation";
import {
  Navigation,
  NavItem,
  NavGroup,
  UserProfile,
} from "./navigation";

export default function DashboardSidebar() {
  const { user, logout } = useAuth();
  const router = useRouter();

  const handleLogout = async () => {
    await logout();
    router.push("/");
  };

  return (
    <Navigation>
      <div className="flex-1 space-y-4">
        <div>
          <NavItem
            label="Dashboard"
            href="/dashboard"
            icon="fas fa-tachometer-alt"
          />
          <NavItem
            label="Services"
            href="/dashboard/services"
            icon="fas fa-server"
          />
          {user?.role === "admin" && (
            <NavItem
              label="Admin Panel"
              href="/admin"
              icon="fas fa-crown"
            />
          )}
        </div>

        <div className="space-y-2">
          <NavGroup title="Billing" icon="fas fa-credit-card">
            <NavItem
              label="Invoices"
              href="/dashboard/invoices"
              icon="fas fa-file-invoice"
              variant="nested"
            />
            <NavItem
              label="Payment Methods"
              href="/dashboard/billing/payment"
              icon="fas fa-wallet"
              variant="nested"
            />
            <NavItem
              label="History"
              href="/dashboard/billing/history"
              icon="fas fa-history"
              variant="nested"
            />
          </NavGroup>

          <NavGroup title="Support" icon="fas fa-headset">
            <NavItem
              label="Tickets"
              href="/dashboard/support/tickets"
              icon="fas fa-ticket-alt"
              variant="nested"
            />
            <NavItem
              label="Knowledge Base"
              href="/dashboard/support/kb"
              icon="fas fa-book"
              variant="nested"
            />
            <NavItem
              label="Contact"
              href="/dashboard/support/contact"
              icon="fas fa-envelope"
              variant="nested"
            />
          </NavGroup>
        </div>
      </div>

      <div className="mt-auto">
        <div className="pt-4 border-t border-white/10 space-y-2">
          <NavItem
            label="Settings"
            href="/dashboard/settings"
            icon="fas fa-cog"
          />
          <NavItem
            label="Documentation"
            href="/dashboard/docs"
            icon="fas fa-book-open"
          />
        </div>

        <UserProfile user={user} onLogout={handleLogout} />
      </div>
    </Navigation>
  );
}
