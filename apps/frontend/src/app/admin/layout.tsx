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
          <Sidebar width="md" title="Orvex" subtitle="Admin Panel">
            <SidebarItem
              icon="fas fa-objects-column"
              label="Dashboard"
              href="/admin/dashboard"
            />

            <SidebarCategory title="Management" collapsible>
              <SidebarItem
                icon="fas fa-users"
                label="Users"
                href="/admin/users"
              />
              <SidebarItem
                icon="fas fa-server"
                label="Services"
                href="/admin/services"
              />
              <SidebarItem
                icon="fas fa-file-invoice"
                label="Invoices"
                href="/admin/invoices"
              />
              <SidebarItem
                icon="fas fa-tags"
                label="Coupons"
                href="/admin/coupons"
              />
              <SidebarItem
                icon="fas fa-bullhorn"
                label="News"
                href="/admin/news"
              />
            </SidebarCategory>

            <SidebarCategory title="Infrastructure" collapsible>
              <SidebarItem
                icon="fas fa-network-wired"
                label="Nodes"
                href="/admin/nodes"
              />
              <SidebarItem
                icon="fas fa-tachometer-alt"
                label="Resource Monitor"
                href="/admin/resources"
              />
              <SidebarItem
                icon="fas fa-dragon"
                label="Pterodactyl"
                href="/admin/pterodactyl"
              />
            </SidebarCategory>

            <SidebarCategory title="Support" collapsible>
              <SidebarItem
                icon="fas fa-ticket-alt"
                label="Tickets"
                href="/admin/tickets"
              />
              <SidebarItem
                icon="fas fa-comment"
                label="AI Chat"
                href="/admin/ai-chat"
              />
            </SidebarCategory>

            <SidebarCategory title="Configuration" collapsible>
              <SidebarItem
                icon="fas fa-cogs"
                label="Settings"
                href="/admin/settings"
              />
              <SidebarItem
                icon="fas fa-credit-card"
                label="Gateways"
                href="/admin/gateways"
              />
              <SidebarItem
                icon="fas fa-key"
                label="API Keys"
                href="/admin/api-keys"
              />
              <SidebarItem
                icon="fas fa-clipboard-list"
                label="Logs"
                href="/admin/logs"
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
