'use client'

import { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import { api, User } from './api'

interface AuthContextType {
  user: User | null
  loading: boolean
  login: (email: string, password: string, totpCode?: string) => Promise<void>
  logout: () => Promise<void>
  refreshUser: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  const refreshUser = async () => {
    try {
      const response = await api.getProfile()
      if (response.success && response.user) {
        setUser(response.user)
      } else {
        setUser(null)
      }
    } catch (error) {
      setUser(null)
    } finally {
      setLoading(false)
    }
  }

  const login = async (email: string, password: string, totpCode?: string) => {
    const response = await api.login({ email, password, totp_code: totpCode })
    if (response.success) {
      await refreshUser()
    } else {
      throw new Error(response.error || response.message || 'Login failed')
    }
  }

  const logout = async () => {
    await api.logout()
    setUser(null)
  }

  useEffect(() => {
    refreshUser()
  }, [])

  return (
    <AuthContext.Provider value={{ user, loading, login, logout, refreshUser }}>
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