'use client'

import { useEffect, useState } from 'react'
import { useAuth } from '@/lib/auth-context'
import { useRouter } from 'next/navigation'
import { motion } from 'framer-motion'
import DashboardLayout from '@/components/dashboard-layout'
import Card, { CardContent, CardHeader } from '@/components/ui/card'
import { LoadingScreen } from '@/components/ui/loading-spinner'
import Button from '@/components/ui/button'
import { api } from '@/lib/api'

interface DashboardStats {
  services: number
  invoices: number
  revenue: number
}

export default function DashboardPage() {
  const { user, loading } = useAuth()
  const router = useRouter()
  const [stats, setStats] = useState<DashboardStats>({ services: 0, invoices: 0, revenue: 0 })
  const [loadingStats, setLoadingStats] = useState(true)

  useEffect(() => {
    if (!loading && !user) {
      router.push('/auth/login')
    }
  }, [user, loading, router])

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const [servicesRes, invoicesRes] = await Promise.all([
          api.getServices().catch(() => ({ success: false, data: [] })),
          api.getInvoices().catch(() => ({ success: false, data: [] }))
        ])

        const servicesData = servicesRes.success ? (servicesRes.data || []) : []
        const invoicesData = invoicesRes.success ? (invoicesRes.data || []) : []

        setStats({
          services: servicesData.length,
          invoices: invoicesData.length,
          revenue: invoicesData.reduce((sum: number, inv: any) => sum + (inv.amount || 0), 0)
        })
      } catch (error) {
        console.error('Failed to fetch dashboard stats:', error)
      } finally {
        setLoadingStats(false)
      }
    }

    if (user) {
      fetchStats()
    }
  }, [user])

  if (loading) {
    return <LoadingScreen message="Loading dashboard..." />
  }

  if (!user) {
    return null
  }

  return (
    <DashboardLayout>
      <motion.div 
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="space-y-8"
      >
        {/* Welcome section */}
        <Card variant="highlight">
          <CardContent className="py-8">
            <h1 className="text-4xl font-bold text-white mb-2">
              Welcome back, {user.first_name}!
            </h1>
            <p className="text-gray-300 text-lg">
              Here's an overview of your account and services.
            </p>
          </CardContent>
        </Card>

        {/* Stats cards */}
        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {[
            {
              title: 'Active Services',
              value: stats.services,
              icon: 'fas fa-server',
              gradient: 'from-violet-500 to-purple-500'
            },
            {
              title: 'Total Invoices',
              value: stats.invoices,
              icon: 'fas fa-file-invoice',
              gradient: 'from-blue-500 to-cyan-500'
            },
            {
              title: 'Total Revenue',
              value: `$${stats.revenue.toFixed(2)}`,
              icon: 'fas fa-dollar-sign',
              gradient: 'from-green-500 to-emerald-500'
            }
          ].map((stat, index) => (
            <motion.div
              key={stat.title}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: index * 0.1 }}
            >
              <Card variant="glass" className="hover:scale-105 transition-transform duration-200">
                <CardContent className="p-6">
                  <div className="flex items-center">
                    <div className={`w-12 h-12 rounded-lg bg-gradient-to-r ${stat.gradient} flex items-center justify-center`}>
                      <i className={`${stat.icon} text-white text-xl`}></i>
                    </div>
                    <div className="ml-4 flex-1">
                      <p className="text-sm font-medium text-gray-400">
                        {stat.title}
                      </p>
                      <p className="text-2xl font-bold text-white">
                        {loadingStats ? (
                          <div className="animate-pulse bg-gray-600 h-8 w-16 rounded"></div>
                        ) : (
                          stat.value
                        )}
                      </p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </motion.div>
          ))}
        </div>

        {/* Quick actions */}
        <Card variant="glass">
          <CardHeader>
            <h3 className="text-xl font-semibold text-white">Quick Actions</h3>
            <p className="text-gray-400">Shortcuts to common tasks</p>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
              {[
                {
                  title: 'New Service',
                  description: 'Create a new service for your account',
                  icon: 'fas fa-plus',
                  gradient: 'from-violet-500 to-purple-500',
                  onClick: () => router.push('/dashboard/services')
                },
                {
                  title: 'View Invoices',
                  description: 'Check your billing history and invoices',
                  icon: 'fas fa-file-invoice',
                  gradient: 'from-blue-500 to-cyan-500',
                  onClick: () => router.push('/dashboard/invoices')
                },
                {
                  title: 'Manage Profile',
                  description: 'Update your account settings and preferences',
                  icon: 'fas fa-user-cog',
                  gradient: 'from-green-500 to-emerald-500',
                  onClick: () => router.push('/dashboard/profile')
                },
                {
                  title: 'Support',
                  description: 'Get help from our support team',
                  icon: 'fas fa-headset',
                  gradient: 'from-orange-500 to-red-500',
                  onClick: () => {}
                }
              ].map((action, index) => (
                <motion.div
                  key={action.title}
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: 0.3 + index * 0.1 }}
                  whileHover={{ y: -2 }}
                  className="cursor-pointer"
                  onClick={action.onClick}
                >
                  <div className="bg-gray-800/50 p-6 rounded-lg border border-gray-700 hover:border-gray-600 transition-all duration-200 group">
                    <div className={`w-12 h-12 rounded-lg bg-gradient-to-r ${action.gradient} flex items-center justify-center mb-4 group-hover:scale-110 transition-transform`}>
                      <i className={`${action.icon} text-white text-xl`}></i>
                    </div>
                    <h4 className="text-lg font-semibold text-white mb-2">
                      {action.title}
                    </h4>
                    <p className="text-sm text-gray-400">
                      {action.description}
                    </p>
                  </div>
                </motion.div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Account status */}
        <Card variant="glass">
          <CardHeader>
            <h3 className="text-xl font-semibold text-white">Account Status</h3>
            <p className="text-gray-400">Security and verification status</p>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <div className="flex items-center">
                  <i className="fas fa-envelope text-gray-400 mr-3"></i>
                  <span className="text-gray-300">Email Verification</span>
                </div>
                <span className={`inline-flex items-center px-3 py-1 rounded-full text-xs font-medium ${
                  user.email_verified 
                    ? 'bg-green-500/20 text-green-400 border border-green-500/30' 
                    : 'bg-red-500/20 text-red-400 border border-red-500/30'
                }`}>
                  {user.email_verified ? (
                    <><i className="fas fa-check mr-1"></i> Verified</>
                  ) : (
                    <><i className="fas fa-exclamation-triangle mr-1"></i> Not Verified</>
                  )}
                </span>
              </div>
              <div className="flex items-center justify-between">
                <div className="flex items-center">
                  <i className="fas fa-shield-alt text-gray-400 mr-3"></i>
                  <span className="text-gray-300">Two-Factor Authentication</span>
                </div>
                <span className={`inline-flex items-center px-3 py-1 rounded-full text-xs font-medium ${
                  user.two_factor_enabled 
                    ? 'bg-green-500/20 text-green-400 border border-green-500/30' 
                    : 'bg-yellow-500/20 text-yellow-400 border border-yellow-500/30'
                }`}>
                  {user.two_factor_enabled ? (
                    <><i className="fas fa-check mr-1"></i> Enabled</>
                  ) : (
                    <><i className="fas fa-exclamation-triangle mr-1"></i> Disabled</>
                  )}
                </span>
              </div>
            </div>
          </CardContent>
        </Card>
      </motion.div>
    </DashboardLayout>
  )
}