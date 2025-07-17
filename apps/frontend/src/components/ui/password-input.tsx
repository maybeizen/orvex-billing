'use client'

import { useState, useEffect } from 'react'
import { AnimatePresence, motion } from 'framer-motion'

interface PasswordRequirement {
  text: string
  met: boolean
}

interface PasswordInputProps {
  value: string
  onChange: (value: string) => void
  placeholder?: string
  showRequirements?: boolean
  className?: string
}

export default function PasswordInput({ 
  value, 
  onChange, 
  placeholder = "Enter password",
  showRequirements = false,
  className = ""
}: PasswordInputProps) {
  const [showPassword, setShowPassword] = useState(false)
  const [showPopup, setShowPopup] = useState(false)
  const [requirements, setRequirements] = useState<PasswordRequirement[]>([])
  const [strength, setStrength] = useState<'weak' | 'medium' | 'strong'>('weak')

  useEffect(() => {
    if (!showRequirements) return

    const reqs: PasswordRequirement[] = [
      { text: 'At least 8 characters', met: value.length >= 8 },
      { text: 'Contains uppercase letter', met: /[A-Z]/.test(value) },
      { text: 'Contains lowercase letter', met: /[a-z]/.test(value) },
      { text: 'Contains number', met: /\d/.test(value) },
      { text: 'Contains special character', met: /[!@#$%^&*(),.?":{}|<>]/.test(value) },
    ]

    setRequirements(reqs)

    // Calculate strength
    const metCount = reqs.filter(req => req.met).length
    if (metCount <= 2) setStrength('weak')
    else if (metCount <= 4) setStrength('medium')
    else setStrength('strong')
  }, [value, showRequirements])

  const getStrengthColor = () => {
    switch (strength) {
      case 'weak': return 'bg-red-500'
      case 'medium': return 'bg-yellow-500'
      case 'strong': return 'bg-green-500'
    }
  }

  const getStrengthWidth = () => {
    switch (strength) {
      case 'weak': return '33%'
      case 'medium': return '66%'
      case 'strong': return '100%'
    }
  }

  return (
    <div className="relative">
      <div className="relative">
        <input
          type={showPassword ? 'text' : 'password'}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          onFocus={() => showRequirements && setShowPopup(true)}
          onBlur={() => setTimeout(() => setShowPopup(false), 200)}
          className={`w-full px-4 py-3 bg-gray-800 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-violet-500 focus:border-transparent pr-12 ${className}`}
          placeholder={placeholder}
        />
        <button
          type="button"
          onClick={() => setShowPassword(!showPassword)}
          className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-white transition-colors"
        >
          <i className={`fas ${showPassword ? 'fa-eye-slash' : 'fa-eye'}`}></i>
        </button>
      </div>

      {/* Password strength bar */}
      {showRequirements && value && (
        <div className="mt-2">
          <div className="flex items-center justify-between text-xs text-gray-400 mb-1">
            <span>Password strength</span>
            <span className="capitalize text-white">{strength}</span>
          </div>
          <div className="w-full bg-gray-700 rounded-full h-1.5">
            <div 
              className={`h-1.5 rounded-full transition-all duration-300 ${getStrengthColor()}`}
              style={{ width: getStrengthWidth() }}
            ></div>
          </div>
        </div>
      )}

      {/* Requirements popup */}
      <AnimatePresence>
        {showRequirements && showPopup && (
          <motion.div
            initial={{ opacity: 0, y: -10, scale: 0.95 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            exit={{ opacity: 0, y: -10, scale: 0.95 }}
            transition={{ duration: 0.2 }}
            className="absolute top-full left-0 right-0 mt-2 p-4 bg-gray-800 border border-gray-600 rounded-lg shadow-xl z-50"
          >
            <h4 className="text-sm font-medium text-white mb-3">Password requirements:</h4>
            <div className="space-y-2">
              {requirements.map((req, index) => (
                <motion.div
                  key={index}
                  initial={{ opacity: 0, x: -10 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: index * 0.1 }}
                  className="flex items-center space-x-2 text-xs"
                >
                  <div className={`w-4 h-4 rounded-full flex items-center justify-center transition-colors ${
                    req.met ? 'bg-green-500' : 'bg-gray-600'
                  }`}>
                    {req.met && <i className="fas fa-check text-white text-[8px]"></i>}
                  </div>
                  <span className={req.met ? 'text-green-400' : 'text-gray-400'}>
                    {req.text}
                  </span>
                </motion.div>
              ))}
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  )
}