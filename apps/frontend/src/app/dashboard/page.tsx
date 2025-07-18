"use client";

import { useEffect, useState } from "react";
import { useAuth } from "@/lib/auth-context";
import { useRouter } from "next/navigation";
import DashboardLayout from "@/components/dashboard-layout";
import { LoadingScreen } from "@/components/ui/loading-spinner";
import { Button } from "@/components/ui/button";
import { api } from "@/lib/api";

interface DashboardStats {
  services: number;
  invoices: number;
  revenue: number;
}

interface Service {
  id: string;
  name: string;
  status: string;
  created_at: string;
  plan: string;
}

interface Invoice {
  id: string;
  amount: number;
  status: string;
  created_at: string;
  due_date: string;
}

export default function DashboardPage() {
  const { user, loading, requires2FA } = useAuth();
  const router = useRouter();
  const [activeTab, setActiveTab] = useState<"services" | "invoices">("services");
  const [stats, setStats] = useState<DashboardStats>({
    services: 0,
    invoices: 0,
    revenue: 0,
  });
  const [services, setServices] = useState<Service[]>([]);
  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [loadingStats, setLoadingStats] = useState(true);

  useEffect(() => {
    if (!loading) {
      if (!user) {
        router.push("/auth/login");
      } else if (requires2FA) {
        router.push("/auth/verify-2fa");
      }
    }
  }, [user, loading, requires2FA, router]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [servicesRes, invoicesRes] = await Promise.all([
          api.getServices().catch(() => ({ success: false, data: [] })),
          api.getInvoices().catch(() => ({ success: false, data: [] })),
        ]);

        const servicesData = servicesRes.success ? servicesRes.data || [] : [];
        const invoicesData = invoicesRes.success ? invoicesRes.data || [] : [];

        setServices(servicesData);
        setInvoices(invoicesData);
        setStats({
          services: servicesData.length,
          invoices: invoicesData.length,
          revenue: invoicesData.reduce(
            (sum: number, inv: any) => sum + (inv.amount || 0),
            0
          ),
        });
      } catch (error) {
        console.error("Failed to fetch dashboard data:", error);
      } finally {
        setLoadingStats(false);
      }
    };

    if (user) {
      fetchData();
    }
  }, [user]);

  if (loading) {
    return <LoadingScreen message="Loading dashboard..." />;
  }

  if (!user) {
    return null;
  }

  return (
    <DashboardLayout>
      <div className="space-y-8">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-4xl font-black text-white mb-2">
            Welcome back,{" "}
            <span className="text-transparent bg-clip-text bg-gradient-to-r from-blue-400 via-purple-400 to-cyan-400">
              {user.first_name}!
            </span>
          </h1>
          <p className="text-gray-300 text-lg">
            Here's an overview of your account and services.
          </p>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {[
            {
              title: "Active Services",
              value: stats.services,
              icon: "fas fa-server",
              color: "text-blue-400",
              bgColor: "bg-blue-500/10",
              borderColor: "border-blue-500/20",
            },
            {
              title: "Total Invoices",
              value: stats.invoices,
              icon: "fas fa-file-invoice",
              color: "text-purple-400",
              bgColor: "bg-purple-500/10",
              borderColor: "border-purple-500/20",
            },
            {
              title: "Credit Balance",
              value: `$${stats.revenue.toFixed(2)}`,
              icon: "fas fa-wallet",
              color: "text-green-400",
              bgColor: "bg-green-500/10",
              borderColor: "border-green-500/20",
            },
          ].map((stat, index) => (
            <div
              key={stat.title}
              className={`p-6 ${stat.bgColor} rounded-xl border ${stat.borderColor} backdrop-blur-sm`}
            >
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-gray-300 text-sm font-medium">{stat.title}</p>
                  <p className="text-2xl font-bold text-white mt-1">
                    {loadingStats ? (
                      <div className="animate-pulse bg-gray-600 h-8 w-20 rounded"></div>
                    ) : (
                      stat.value
                    )}
                  </p>
                </div>
                <div className="w-12 h-12 bg-white/10 rounded-lg flex items-center justify-center">
                  <i className={`${stat.icon} ${stat.color} text-xl`}></i>
                </div>
              </div>
            </div>
          ))}
        </div>

        {/* Main Content Container */}
        <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-xl p-6">
          {/* Tab Navigation */}
          <div className="flex space-x-1 bg-white/5 p-1 rounded-lg mb-6">
            <button
              onClick={() => setActiveTab("services")}
              className={`flex-1 px-4 py-2 rounded-md text-sm font-medium transition-all duration-200 ${
                activeTab === "services"
                  ? "bg-white/10 text-white shadow-sm"
                  : "text-gray-400 hover:text-white hover:bg-white/5"
              }`}
            >
              <i className="fas fa-server mr-2"></i>
              Services
            </button>
            <button
              onClick={() => setActiveTab("invoices")}
              className={`flex-1 px-4 py-2 rounded-md text-sm font-medium transition-all duration-200 ${
                activeTab === "invoices"
                  ? "bg-white/10 text-white shadow-sm"
                  : "text-gray-400 hover:text-white hover:bg-white/5"
              }`}
            >
              <i className="fas fa-file-invoice mr-2"></i>
              Invoices
            </button>
          </div>

          {/* Tab Content */}
          <div className="min-h-[400px]">
            {activeTab === "services" && (
              <div className="space-y-4">
                <div className="flex items-center justify-between mb-4">
                  <h3 className="text-xl font-bold text-white">Your Services</h3>
                  <Button
                    onClick={() => router.push("/dashboard/services/create")}
                    variant="glass"
                    size="sm"
                    icon="fas fa-plus"
                    iconPosition="left"
                  >
                    Purchase a Service
                  </Button>
                </div>
                {loadingStats ? (
                  <div className="space-y-4">
                    {[1, 2, 3].map((i) => (
                      <div key={i} className="animate-pulse bg-white/5 rounded-lg p-4">
                        <div className="h-4 bg-gray-600 rounded w-1/4 mb-2"></div>
                        <div className="h-3 bg-gray-600 rounded w-1/2"></div>
                      </div>
                    ))}
                  </div>
                ) : services.length === 0 ? (
                  <div className="text-center py-12">
                    <i className="fas fa-server text-4xl text-gray-400 mb-4"></i>
                    <p className="text-gray-400 text-lg">No services found</p>
                    <p className="text-gray-500 text-sm">Create your first service to get started</p>
                  </div>
                ) : (
                  <div className="space-y-4">
                    {services.map((service) => (
                      <div
                        key={service.id}
                        className="flex items-center justify-between p-4 bg-white/5 rounded-lg border border-white/10 hover:bg-white/10 transition-colors duration-200"
                      >
                        <div className="flex items-center space-x-4">
                          <div className="w-10 h-10 bg-blue-500/20 rounded-lg flex items-center justify-center">
                            <i className="fas fa-server text-blue-400"></i>
                          </div>
                          <div>
                            <h4 className="text-white font-medium">{service.name}</h4>
                            <p className="text-gray-400 text-sm">{service.plan}</p>
                          </div>
                        </div>
                        <div className="flex items-center space-x-4">
                          <span
                            className={`px-3 py-1 rounded-full text-xs font-medium ${
                              service.status === "active"
                                ? "bg-green-500/20 text-green-400"
                                : "bg-red-500/20 text-red-400"
                            }`}
                          >
                            {service.status}
                          </span>
                          <button className="text-gray-400 hover:text-white">
                            <i className="fas fa-chevron-right"></i>
                          </button>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            )}

            {activeTab === "invoices" && (
              <div className="space-y-4">
                <div className="flex items-center justify-between mb-4">
                  <h3 className="text-xl font-bold text-white">Recent Invoices</h3>
                  <Button
                    onClick={() => router.push("/dashboard/invoices")}
                    variant="glass"
                    size="sm"
                    icon="fas fa-eye"
                    iconPosition="left"
                  >
                    View All
                  </Button>
                </div>
                {loadingStats ? (
                  <div className="space-y-4">
                    {[1, 2, 3].map((i) => (
                      <div key={i} className="animate-pulse bg-white/5 rounded-lg p-4">
                        <div className="h-4 bg-gray-600 rounded w-1/4 mb-2"></div>
                        <div className="h-3 bg-gray-600 rounded w-1/2"></div>
                      </div>
                    ))}
                  </div>
                ) : invoices.length === 0 ? (
                  <div className="text-center py-12">
                    <i className="fas fa-file-invoice text-4xl text-gray-400 mb-4"></i>
                    <p className="text-gray-400 text-lg">No invoices found</p>
                    <p className="text-gray-500 text-sm">Your invoices will appear here</p>
                  </div>
                ) : (
                  <div className="space-y-4">
                    {invoices.slice(0, 5).map((invoice) => (
                      <div
                        key={invoice.id}
                        className="flex items-center justify-between p-4 bg-white/5 rounded-lg border border-white/10 hover:bg-white/10 transition-colors duration-200"
                      >
                        <div className="flex items-center space-x-4">
                          <div className="w-10 h-10 bg-purple-500/20 rounded-lg flex items-center justify-center">
                            <i className="fas fa-file-invoice text-purple-400"></i>
                          </div>
                          <div>
                            <h4 className="text-white font-medium">#{invoice.id}</h4>
                            <p className="text-gray-400 text-sm">
                              Due: {new Date(invoice.due_date).toLocaleDateString()}
                            </p>
                          </div>
                        </div>
                        <div className="flex items-center space-x-4">
                          <span className="text-white font-medium">
                            ${invoice.amount.toFixed(2)}
                          </span>
                          <span
                            className={`px-3 py-1 rounded-full text-xs font-medium ${
                              invoice.status === "paid"
                                ? "bg-green-500/20 text-green-400"
                                : invoice.status === "pending"
                                ? "bg-yellow-500/20 text-yellow-400"
                                : "bg-red-500/20 text-red-400"
                            }`}
                          >
                            {invoice.status}
                          </span>
                          <button className="text-gray-400 hover:text-white">
                            <i className="fas fa-chevron-right"></i>
                          </button>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            )}
          </div>
        </div>
      </div>
    </DashboardLayout>
  );
}
