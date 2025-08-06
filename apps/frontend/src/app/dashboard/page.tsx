"use client";

import React from "react";
import { useAuth } from "@/contexts/auth-context";

export default function Dashboard() {
  const { user, isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-white"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-lg p-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-white">
              Welcome back, {user?.firstName}! =K
            </h1>
            <p className="text-white/60 mt-2">
              Here's what's happening with your servers today.
            </p>
          </div>
          <div className="text-right">
            <div className="text-2xl font-bold text-green-400">
              <i className="fas fa-server mr-2" />3 Active
            </div>
            <p className="text-white/40 text-sm">Servers running</p>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-lg p-6">
          <div className="flex items-center">
            <div className="p-3 bg-blue-500/20 rounded-lg">
              <i className="fas fa-server text-blue-400 text-xl" />
            </div>
            <div className="ml-4">
              <p className="text-white/60 text-sm">Total Servers</p>
              <p className="text-2xl font-semibold text-white">12</p>
            </div>
          </div>
        </div>

        <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-lg p-6">
          <div className="flex items-center">
            <div className="p-3 bg-green-500/20 rounded-lg">
              <i className="fas fa-chart-line text-green-400 text-xl" />
            </div>
            <div className="ml-4">
              <p className="text-white/60 text-sm">Monthly Revenue</p>
              <p className="text-2xl font-semibold text-white">$2,847</p>
            </div>
          </div>
        </div>

        <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-lg p-6">
          <div className="flex items-center">
            <div className="p-3 bg-yellow-500/20 rounded-lg">
              <i className="fas fa-users text-yellow-400 text-xl" />
            </div>
            <div className="ml-4">
              <p className="text-white/60 text-sm">Active Users</p>
              <p className="text-2xl font-semibold text-white">1,247</p>
            </div>
          </div>
        </div>

        <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-lg p-6">
          <div className="flex items-center">
            <div className="p-3 bg-purple-500/20 rounded-lg">
              <i className="fas fa-clock text-purple-400 text-xl" />
            </div>
            <div className="ml-4">
              <p className="text-white/60 text-sm">Uptime</p>
              <p className="text-2xl font-semibold text-white">99.9%</p>
            </div>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-lg p-6">
          <h2 className="text-xl font-semibold text-white mb-4">
            Recent Activity
          </h2>
          <div className="space-y-4">
            <div className="flex items-center p-3 bg-white/5 rounded-lg">
              <div className="w-2 h-2 bg-green-400 rounded-full"></div>
              <div className="ml-3 flex-1">
                <p className="text-white text-sm">
                  Server "Web-01" started successfully
                </p>
                <p className="text-white/40 text-xs">2 minutes ago</p>
              </div>
            </div>
            <div className="flex items-center p-3 bg-white/5 rounded-lg">
              <div className="w-2 h-2 bg-blue-400 rounded-full"></div>
              <div className="ml-3 flex-1">
                <p className="text-white text-sm">
                  Backup completed for "Database-01"
                </p>
                <p className="text-white/40 text-xs">15 minutes ago</p>
              </div>
            </div>
            <div className="flex items-center p-3 bg-white/5 rounded-lg">
              <div className="w-2 h-2 bg-yellow-400 rounded-full"></div>
              <div className="ml-3 flex-1">
                <p className="text-white text-sm">
                  High CPU usage detected on "API-02"
                </p>
                <p className="text-white/40 text-xs">1 hour ago</p>
              </div>
            </div>
          </div>
        </div>

        <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-lg p-6">
          <h2 className="text-xl font-semibold text-white mb-4">
            Quick Actions
          </h2>
          <div className="grid grid-cols-2 gap-3">
            <button className="flex flex-col items-center p-4 bg-white/5 hover:bg-white/10 rounded-lg transition-colors">
              <i className="fas fa-plus text-blue-400 text-2xl mb-2" />
              <span className="text-white text-sm">Create Server</span>
            </button>
            <button className="flex flex-col items-center p-4 bg-white/5 hover:bg-white/10 rounded-lg transition-colors">
              <i className="fas fa-download text-green-400 text-2xl mb-2" />
              <span className="text-white text-sm">Backup All</span>
            </button>
            <button className="flex flex-col items-center p-4 bg-white/5 hover:bg-white/10 rounded-lg transition-colors">
              <i className="fas fa-chart-bar text-purple-400 text-2xl mb-2" />
              <span className="text-white text-sm">View Analytics</span>
            </button>
            <button className="flex flex-col items-center p-4 bg-white/5 hover:bg-white/10 rounded-lg transition-colors">
              <i className="fas fa-cog text-gray-400 text-2xl mb-2" />
              <span className="text-white text-sm">Settings</span>
            </button>
          </div>
        </div>
      </div>

      <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-lg p-6">
        <h2 className="text-xl font-semibold text-white mb-4">
          Server Overview
        </h2>
        <div className="overflow-x-auto">
          <table className="w-full text-left">
            <thead>
              <tr className="border-b border-white/10">
                <th className="text-white/60 text-sm font-medium pb-3">
                  Server Name
                </th>
                <th className="text-white/60 text-sm font-medium pb-3">
                  Status
                </th>
                <th className="text-white/60 text-sm font-medium pb-3">CPU</th>
                <th className="text-white/60 text-sm font-medium pb-3">
                  Memory
                </th>
                <th className="text-white/60 text-sm font-medium pb-3">
                  Uptime
                </th>
                <th className="text-white/60 text-sm font-medium pb-3">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="space-y-2">
              <tr className="border-b border-white/5">
                <td className="py-3">
                  <div className="flex items-center">
                    <i className="fas fa-server text-blue-400 mr-3" />
                    <span className="text-white">Web-01</span>
                  </div>
                </td>
                <td className="py-3">
                  <span className="inline-flex items-center px-2 py-1 rounded-full text-xs bg-green-500/20 text-green-400">
                    <div className="w-1.5 h-1.5 bg-green-400 rounded-full mr-1"></div>
                    Online
                  </span>
                </td>
                <td className="py-3 text-white">23%</td>
                <td className="py-3 text-white">1.2GB</td>
                <td className="py-3 text-white">47d 12h</td>
                <td className="py-3">
                  <button className="text-white/40 hover:text-white mr-2">
                    <i className="fas fa-cog" />
                  </button>
                  <button className="text-white/40 hover:text-red-400">
                    <i className="fas fa-power-off" />
                  </button>
                </td>
              </tr>
              <tr className="border-b border-white/5">
                <td className="py-3">
                  <div className="flex items-center">
                    <i className="fas fa-database text-purple-400 mr-3" />
                    <span className="text-white">Database-01</span>
                  </div>
                </td>
                <td className="py-3">
                  <span className="inline-flex items-center px-2 py-1 rounded-full text-xs bg-green-500/20 text-green-400">
                    <div className="w-1.5 h-1.5 bg-green-400 rounded-full mr-1"></div>
                    Online
                  </span>
                </td>
                <td className="py-3 text-white">67%</td>
                <td className="py-3 text-white">3.8GB</td>
                <td className="py-3 text-white">23d 5h</td>
                <td className="py-3">
                  <button className="text-white/40 hover:text-white mr-2">
                    <i className="fas fa-cog" />
                  </button>
                  <button className="text-white/40 hover:text-red-400">
                    <i className="fas fa-power-off" />
                  </button>
                </td>
              </tr>
              <tr className="border-b border-white/5">
                <td className="py-3">
                  <div className="flex items-center">
                    <i className="fas fa-cloud text-green-400 mr-3" />
                    <span className="text-white">API-02</span>
                  </div>
                </td>
                <td className="py-3">
                  <span className="inline-flex items-center px-2 py-1 rounded-full text-xs bg-yellow-500/20 text-yellow-400">
                    <div className="w-1.5 h-1.5 bg-yellow-400 rounded-full mr-1"></div>
                    Warning
                  </span>
                </td>
                <td className="py-3 text-white">89%</td>
                <td className="py-3 text-white">2.1GB</td>
                <td className="py-3 text-white">12d 18h</td>
                <td className="py-3">
                  <button className="text-white/40 hover:text-white mr-2">
                    <i className="fas fa-cog" />
                  </button>
                  <button className="text-white/40 hover:text-red-400">
                    <i className="fas fa-power-off" />
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
