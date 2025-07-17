'use client'

import { useState, ReactNode } from 'react'
import Link from 'next/link'
import { useAuth } from '@/lib/auth-context'
import { useRouter, usePathname } from 'next/navigation'
import { motion } from 'framer-motion'

interface DashboardLayoutProps {
  children: ReactNode
}

export default function DashboardLayout({ children }: DashboardLayoutProps) {
  const [sidebarOpen, setSidebarOpen] = useState(false)
  const { user, logout } = useAuth()
  const router = useRouter()
  const pathname = usePathname()

  const handleLogout = async () => {
    await logout()
    router.push('/')
  }

  const navigation = [
    { name: 'Dashboard', href: '/dashboard', icon: 'fas fa-tachometer-alt' },
    { name: 'Services', href: '/dashboard/services', icon: 'fas fa-server' },
    { name: 'Invoices', href: '/dashboard/invoices', icon: 'fas fa-file-invoice' },
    { name: 'Profile', href: '/dashboard/profile', icon: 'fas fa-user' },
  ]

  const isActive = (href: string) => pathname === href

  return (
    <div className="min-h-screen bg-gray-900">
      <div className="flex">
        {/* Sidebar */}
        <div className={`fixed inset-y-0 left-0 z-50 w-64 bg-gray-800 border-r border-gray-700 transform ${sidebarOpen ? 'translate-x-0' : '-translate-x-full'} transition-transform duration-300 ease-in-out lg:translate-x-0 lg:static lg:inset-0`}>
          <div className="flex items-center justify-center h-16 bg-gradient-to-r from-violet-600 to-purple-600">
            <Link href="/dashboard" className="text-2xl font-bold bg-gradient-to-r from-violet-200 to-purple-200 bg-clip-text text-transparent">
              Orvex
            </Link>
          </div>
          <nav className="mt-8">
            <div className="px-4 space-y-2">
              {navigation.map((item) => (
                <motion.div
                  key={item.name}
                  whileHover={{ x: 4 }}
                  whileTap={{ scale: 0.98 }}
                >
                  <Link
                    href={item.href}
                    className={`group flex items-center px-4 py-3 text-sm font-medium rounded-lg transition-all duration-200 ${
                      isActive(item.href)
                        ? 'bg-violet-600 text-white shadow-lg shadow-violet-600/25'
                        : 'text-gray-300 hover:bg-gray-700 hover:text-white'
                    }`}
                  >
                    <i className={`${item.icon} mr-3 ${
                      isActive(item.href) ? 'text-white' : 'text-gray-400 group-hover:text-gray-300'
                    }`}></i>
                    {item.name}
                  </Link>
                </motion.div>
              ))}
            </div>
          </nav>
        </div>

        {/* Mobile sidebar overlay */}
        {sidebarOpen && (
          <motion.div 
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black bg-opacity-50 z-40 lg:hidden"
            onClick={() => setSidebarOpen(false)}
          ></motion.div>
        )}

        {/* Main content */}
        <div className="flex-1 lg:ml-0">
          {/* Top bar */}
          <header className="bg-gray-800 border-b border-gray-700">
            <div className="flex items-center justify-between px-4 py-4 sm:px-6 lg:px-8">
              <button
                type="button"
                className="lg:hidden -ml-0.5 -mt-0.5 h-12 w-12 inline-flex items-center justify-center rounded-md text-gray-400 hover:text-white hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-violet-500 transition-colors"
                onClick={() => setSidebarOpen(true)}
              >
                <i className="fas fa-bars text-xl"></i>
              </button>
              
              <div className="hidden lg:block">
                <h1 className="text-2xl font-semibold text-white">Dashboard</h1>
              </div>

              <div className="flex items-center space-x-4">
                <div className="flex items-center space-x-3">
                  <div className="text-sm text-gray-300">
                    Welcome, <span className="text-white font-medium">{user?.first_name}</span>!
                  </div>
                  {user?.avatar_url ? (
                    <img
                      className="h-10 w-10 rounded-full ring-2 ring-violet-500 ring-offset-2 ring-offset-gray-800"
                      src={user.avatar_url}
                      alt={user.username}
                    />
                  ) : (
                    <div className="h-10 w-10 rounded-full bg-gradient-to-r from-violet-500 to-purple-500 flex items-center justify-center ring-2 ring-violet-500 ring-offset-2 ring-offset-gray-800">
                      <span className="text-sm font-bold text-white">
                        {user?.first_name?.[0]?.toUpperCase()}
                      </span>
                    </div>
                  )}
                </div>
                <motion.button
                  whileHover={{ scale: 1.05 }}
                  whileTap={{ scale: 0.95 }}
                  onClick={handleLogout}
                  className="text-gray-400 hover:text-red-400 p-2 rounded-lg hover:bg-gray-700 transition-colors"
                  title="Logout"
                >
                  <i className="fas fa-sign-out-alt"></i>
                </motion.button>
              </div>
            </div>
          </header>

          {/* Page content */}
          <main className="p-4 sm:p-6 lg:p-8">
            {children}
          </main>
        </div>
      </div>
    </div>
  )
}