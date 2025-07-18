'use client'

import { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import { api, User } from './api'

interface AuthContextType {
  user: User | null
  loading: boolean
  requires2FA: boolean
  login: (email: string, password: string, totpCode?: string) => Promise<any>
  logout: () => Promise<void>
  refreshUser: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)
  const [requires2FA, setRequires2FA] = useState(false)

  const refreshUser = async () => {
    try {
      const response = await api.getProfile()
      if (response.success && response.user) {
        setUser(response.user)
        // Check if user needs 2FA verification
        setRequires2FA(response.user.two_factor_enabled && !response.two_factor_verified)
      } else {
        setUser(null)
        setRequires2FA(false)
      }
    } catch (error) {
      setUser(null)
      setRequires2FA(false)
    } finally {
      setLoading(false)
    }
  }

  const login = async (email: string, password: string, totpCode?: string) => {
    const response = await api.login({ email, password, totp_code: totpCode })
    if (response.success) {
      await refreshUser()
      return response
    } else {
      throw new Error(response.error || response.message || 'Login failed')
    }
  }

  const logout = async () => {
    await api.logout()
    setUser(null)
    setRequires2FA(false)
  }

  useEffect(() => {
    refreshUser()
  }, [])

  return (
    <AuthContext.Provider value={{ user, loading, requires2FA, login, logout, refreshUser }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}