"use client";

import { useEffect, useState } from "react";
import { useAuth } from "@/lib/auth-context";
import { useRouter } from "next/navigation";
import { motion } from "framer-motion";
import DashboardLayout from "@/components/dashboard-layout";
import { LoadingScreen } from "@/components/ui/loading-spinner";
import { api } from "@/lib/api";

interface DashboardStats {
  services: number;
  invoices: number;
  revenue: number;
}

export default function DashboardPage() {
  const { user, loading } = useAuth();
  const router = useRouter();
  const [stats, setStats] = useState<DashboardStats>({
    services: 0,
    invoices: 0,
    revenue: 0,
  });
  const [loadingStats, setLoadingStats] = useState(true);

  useEffect(() => {
    if (!loading && !user) {
      router.push("/auth/login");
    }
  }, [user, loading, router]);

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const [servicesRes, invoicesRes] = await Promise.all([
          api.getServices().catch(() => ({ success: false, data: [] })),
          api.getInvoices().catch(() => ({ success: false, data: [] })),
        ]);

        const servicesData = servicesRes.success ? servicesRes.data || [] : [];
        const invoicesData = invoicesRes.success ? invoicesRes.data || [] : [];

        setStats({
          services: servicesData.length,
          invoices: invoicesData.length,
          revenue: invoicesData.reduce(
            (sum: number, inv: any) => sum + (inv.amount || 0),
            0
          ),
        });
      } catch (error) {
        console.error("Failed to fetch dashboard stats:", error);
      } finally {
        setLoadingStats(false);
      }
    };

    if (user) {
      fetchStats();
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
      <motion.div
        initial={{ opacity: 0, y: 30 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.8 }}
        className="space-y-10"
      >
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2, duration: 0.8 }}
          className="relative"
        >
          <div className="absolute inset-0 bg-gradient-to-r from-violet-600/10 to-cyan-600/10 rounded-3xl blur-xl"></div>
          <div className="relative bg-black/40 backdrop-blur-3xl border border-white/10 rounded-3xl p-10 shadow-2xl">
            <div className="flex items-center justify-between">
              <div>
                <motion.h1
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: 0.4, duration: 0.6 }}
                  className="text-5xl font-black text-white mb-4 font-['Sora']"
                >
                  <span className="text-transparent bg-clip-text bg-gradient-to-r from-violet-400 to-cyan-400">
                    Welcome back,
                  </span>
                  <span className="block text-white">{user.first_name}!</span>
                </motion.h1>
                <motion.p
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: 0.6, duration: 0.6 }}
                  className="text-white/70 text-xl"
                >
                  Here's an overview of your account and services.
                </motion.p>
              </div>
              <motion.div
                initial={{ opacity: 0, scale: 0.8 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ delay: 0.8, duration: 0.6 }}
                className="hidden lg:block"
              >
                <div className="w-32 h-32 bg-gradient-to-br from-violet-500/20 to-cyan-500/20 rounded-full flex items-center justify-center">
                  <i className="fas fa-user-circle text-6xl text-white/30"></i>
                </div>
              </motion.div>
            </div>
          </div>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.4, duration: 0.8 }}
          className="grid grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-3"
        >
          {[
            {
              title: "Active Services",
              value: stats.services,
              icon: "fas fa-server",
              gradient: "from-violet-500 to-purple-500",
              bgGradient: "from-violet-500/10 to-purple-500/10",
            },
            {
              title: "Total Invoices",
              value: stats.invoices,
              icon: "fas fa-file-invoice",
              gradient: "from-blue-500 to-cyan-500",
              bgGradient: "from-blue-500/10 to-cyan-500/10",
            },
            {
              title: "Total Revenue",
              value: `$${stats.revenue.toFixed(2)}`,
              icon: "fas fa-dollar-sign",
              gradient: "from-green-500 to-emerald-500",
              bgGradient: "from-green-500/10 to-emerald-500/10",
            },
          ].map((stat, index) => (
            <motion.div
              key={stat.title}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.6 + index * 0.2, duration: 0.6 }}
              whileHover={{ y: -8, scale: 1.02 }}
              className="group cursor-pointer"
            >
              <div className="relative">
                <div className="absolute inset-0 bg-gradient-to-br from-violet-600/5 to-cyan-600/5 rounded-2xl blur-xl"></div>
                <div
                  className={`relative bg-gradient-to-br ${stat.bgGradient} backdrop-blur-xl border border-white/10 rounded-2xl p-8 shadow-xl hover:shadow-2xl transition-all duration-300`}
                >
                  <div className="absolute inset-0 bg-gradient-to-br from-white/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300 rounded-2xl" />
                  <div className="relative z-10">
                    <div className="flex items-center justify-between mb-6">
                      <div
                        className={`w-16 h-16 rounded-xl bg-gradient-to-r ${stat.gradient} flex items-center justify-center group-hover:scale-110 transition-transform duration-300 shadow-lg`}
                      >
                        <i className={`${stat.icon} text-white text-2xl`}></i>
                      </div>
                      <div className="text-right">
                        <div className="text-3xl font-black text-white font-['Sora']">
                          {loadingStats ? (
                            <div className="animate-pulse bg-gray-600 h-8 w-20 rounded"></div>
                          ) : (
                            stat.value
                          )}
                        </div>
                        <div className="text-white/60 text-sm font-medium mt-1">
                          {stat.title}
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </motion.div>
          ))}
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.8, duration: 0.8 }}
          className="relative"
        >
          <div className="absolute inset-0 bg-gradient-to-r from-violet-600/10 to-cyan-600/10 rounded-3xl blur-xl"></div>
          <div className="relative bg-black/40 backdrop-blur-3xl border border-white/10 rounded-3xl p-10 shadow-2xl">
            <div className="mb-8">
              <h3 className="text-3xl font-bold text-white mb-2 font-['Sora']">
                Quick Actions
              </h3>
              <p className="text-white/60 text-lg">Shortcuts to common tasks</p>
            </div>
            <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
              {[
                {
                  title: "New Service",
                  description: "Create a new service for your account",
                  icon: "fas fa-plus",
                  gradient: "from-violet-500 to-purple-500",
                  onClick: () => router.push("/dashboard/services"),
                },
                {
                  title: "View Invoices",
                  description: "Check your billing history and invoices",
                  icon: "fas fa-file-invoice",
                  gradient: "from-blue-500 to-cyan-500",
                  onClick: () => router.push("/dashboard/invoices"),
                },
                {
                  title: "Manage Profile",
                  description: "Update your account settings and preferences",
                  icon: "fas fa-user-cog",
                  gradient: "from-green-500 to-emerald-500",
                  onClick: () => router.push("/dashboard/profile"),
                },
                {
                  title: "Support",
                  description: "Get help from our support team",
                  icon: "fas fa-headset",
                  gradient: "from-orange-500 to-red-500",
                  onClick: () => {},
                },
              ].map((action, index) => (
                <motion.div
                  key={action.title}
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: 1 + index * 0.1, duration: 0.6 }}
                  whileHover={{ y: -5, scale: 1.05 }}
                  className="cursor-pointer group"
                  onClick={action.onClick}
                >
                  <div className="bg-gray-800/30 backdrop-blur-sm p-6 rounded-xl border border-white/10 hover:border-white/20 transition-all duration-300 group-hover:shadow-xl">
                    <div
                      className={`w-14 h-14 rounded-xl bg-gradient-to-r ${action.gradient} flex items-center justify-center mb-4 group-hover:scale-110 transition-transform duration-300 shadow-lg`}
                    >
                      <i className={`${action.icon} text-white text-xl`}></i>
                    </div>
                    <h4 className="text-xl font-bold text-white mb-2 font-['Sora']">
                      {action.title}
                    </h4>
                    <p className="text-sm text-white/60 leading-relaxed">
                      {action.description}
                    </p>
                  </div>
                </motion.div>
              ))}
            </div>
          </div>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 1.2, duration: 0.8 }}
          className="relative"
        >
          <div className="absolute inset-0 bg-gradient-to-r from-violet-600/10 to-cyan-600/10 rounded-3xl blur-xl"></div>
          <div className="relative bg-black/40 backdrop-blur-3xl border border-white/10 rounded-3xl p-10 shadow-2xl">
            <div className="mb-8">
              <h3 className="text-3xl font-bold text-white mb-2 font-['Sora']">
                Account Status
              </h3>
              <p className="text-white/60 text-lg">
                Security and verification status
              </p>
            </div>
            <div className="space-y-6">
              <div className="flex items-center justify-between p-6 bg-gray-800/30 rounded-xl border border-white/10">
                <div className="flex items-center">
                  <div className="w-12 h-12 bg-gradient-to-r from-blue-500 to-cyan-500 rounded-lg flex items-center justify-center mr-4">
                    <i className="fas fa-envelope text-white text-xl"></i>
                  </div>
                  <div>
                    <span className="text-white font-medium text-lg">
                      Email Verification
                    </span>
                    <p className="text-white/60 text-sm">
                      Verify your email address
                    </p>
                  </div>
                </div>
                <span
                  className={`inline-flex items-center px-4 py-2 rounded-full text-sm font-medium ${
                    user.email_verified
                      ? "bg-green-500/20 text-green-400 border border-green-500/30"
                      : "bg-red-500/20 text-red-400 border border-red-500/30"
                  }`}
                >
                  {user.email_verified ? (
                    <>
                      <i className="fas fa-check mr-2"></i> Verified
                    </>
                  ) : (
                    <>
                      <i className="fas fa-exclamation-triangle mr-2"></i> Not
                      Verified
                    </>
                  )}
                </span>
              </div>
              <div className="flex items-center justify-between p-6 bg-gray-800/30 rounded-xl border border-white/10">
                <div className="flex items-center">
                  <div className="w-12 h-12 bg-gradient-to-r from-violet-500 to-purple-500 rounded-lg flex items-center justify-center mr-4">
                    <i className="fas fa-shield-alt text-white text-xl"></i>
                  </div>
                  <div>
                    <span className="text-white font-medium text-lg">
                      Two-Factor Authentication
                    </span>
                    <p className="text-white/60 text-sm">
                      Secure your account with 2FA
                    </p>
                  </div>
                </div>
                <span
                  className={`inline-flex items-center px-4 py-2 rounded-full text-sm font-medium ${
                    user.two_factor_enabled
                      ? "bg-green-500/20 text-green-400 border border-green-500/30"
                      : "bg-yellow-500/20 text-yellow-400 border border-yellow-500/30"
                  }`}
                >
                  {user.two_factor_enabled ? (
                    <>
                      <i className="fas fa-check mr-2"></i> Enabled
                    </>
                  ) : (
                    <>
                      <i className="fas fa-exclamation-triangle mr-2"></i>{" "}
                      Disabled
                    </>
                  )}
                </span>
              </div>
            </div>
          </div>
        </motion.div>
      </motion.div>
    </DashboardLayout>
  );
}
