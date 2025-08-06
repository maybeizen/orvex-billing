"use client";

import React from "react";
import { Sidebar, SidebarItem, SidebarCategory } from "@/components/ui/sidebar";
import { AuthProvider } from "@/contexts/auth-context";
import { ProtectedPage } from "@/components/protected-page";

export default function AdminDashboardLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <AuthProvider>
      <ProtectedPage>
        <div className="bg-neutral-950 min-h-screen">
          <Sidebar width="md" title="Orvex" subtitle="Admin Access">
            <SidebarItem
              icon="fas fa-home"
              label="Dashboard"
              href="/dashboard"
            />
            <SidebarItem
              icon="fas fa-objects-column"
              label="Services"
              href="/dashboard/services"
            />
            <SidebarItem
              icon="fas fa-signal"
              label="Status"
              href="/dashboard/status"
            />

            <SidebarCategory title="Account & Billing" collapsible>
              <SidebarItem
                icon="fas fa-user"
                label="My Account"
                href="/dashboard/account"
              />
              <SidebarItem
                icon="fas fa-store"
                label="Store"
                href="/dashboard/store"
              />
              <SidebarItem
                icon="fas fa-file-invoice"
                label="Invoices"
                href="/dashboard/invoices"
              />
              <SidebarItem
                icon="fas fa-credit-card"
                label="Payment Methods"
                href="/dashboard/payment-methods"
              />
              <SidebarItem
                icon="fas fa-history"
                label="History"
                href="/dashboard/history"
              />
            </SidebarCategory>

            <SidebarCategory title="Help & Legal" collapsible>
              <SidebarItem
                icon="fas fa-ticket-alt"
                label="Tickets"
                href="/dashboard/tickets"
              />
              <SidebarItem
                icon="fas fa-book"
                label="Knowledge Base"
                href="/dashboard/knowledge-base"
              />
              <SidebarItem
                icon="fas fa-envelope"
                label="Contact"
                href="/dashboard/contact"
              />
              <SidebarItem
                icon="fas fa-file-contract"
                label="Terms & Policy"
                href="/dashboard/terms"
              />
            </SidebarCategory>
          </Sidebar>

          <div className="ml-72 min-h-screen">
            <main className="p-8 min-h-screen">{children}</main>
          </div>
        </div>
      </ProtectedPage>
    </AuthProvider>
  );
}
