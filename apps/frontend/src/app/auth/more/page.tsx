"use client";

import Link from "next/link";
import AuthLayout from "@/components/auth-layout";

export default function AuthMorePage() {
  const socialOptions = [
    {
      title: "Continue with Discord",
      description: "Sign in using your Discord account",
      icon: "fab fa-discord",
      href: "/auth/discord",
      color: "text-indigo-400",
      bgColor: "bg-indigo-500/10",
      borderColor: "border-indigo-500/20",
    },
    {
      title: "Continue with Google",
      description: "Sign in using your Google account",
      icon: "fab fa-google",
      href: "/auth/google",
      color: "text-red-400",
      bgColor: "bg-red-500/10",
      borderColor: "border-red-500/20",
    },
    {
      title: "Continue with GitHub",
      description: "Sign in using your GitHub account",
      icon: "fab fa-github",
      href: "/auth/github",
      color: "text-gray-400",
      bgColor: "bg-gray-500/10",
      borderColor: "border-gray-500/20",
    },
  ];

  return (
    <AuthLayout>
      <div className="space-y-6">
        <div className="text-center mb-8">
          <h2 className="text-2xl md:text-3xl font-black mb-4 leading-tight">
            Social Login
          </h2>
          <p className="text-gray-300">
            Sign in with your preferred social account
          </p>
        </div>

        <div className="grid grid-cols-1 gap-4">
          {socialOptions.map((option) => (
            <Link
              key={option.title}
              href={option.href}
              className={`block p-6 ${option.bgColor} rounded-xl border ${option.borderColor} hover:bg-opacity-20 transition-all duration-200 group`}
            >
              <div className="flex items-center gap-4">
                <div className="w-12 h-12 rounded-lg bg-white/10 flex items-center justify-center">
                  <i className={`${option.icon} ${option.color} text-xl`}></i>
                </div>
                <div className="flex-1">
                  <h3 className="text-lg font-bold text-white mb-1">
                    {option.title}
                  </h3>
                  <p className="text-gray-300 text-sm">{option.description}</p>
                </div>
                <i className="fas fa-arrow-right text-gray-400 group-hover:text-gray-300 transition-colors"></i>
              </div>
            </Link>
          ))}
        </div>

        <div className="text-center pt-6 border-t border-white/10">
          <p className="text-gray-400 text-sm">
            Prefer email and password?{" "}
            <Link
              href="/auth/login"
              className="text-blue-400 hover:text-blue-300 transition-colors font-medium"
            >
              Use traditional login
            </Link>
          </p>
        </div>
      </div>
    </AuthLayout>
  );
}
