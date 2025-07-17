"use client";

import { ReactNode } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";

interface AuthLayoutProps {
  children: ReactNode;
}

export default function AuthLayout({ children }: AuthLayoutProps) {
  const pathname = usePathname();

  const navItems = [
    { name: "Login", href: "/auth/login", icon: "fas fa-sign-in-alt" },
    { name: "Register", href: "/auth/register", icon: "fas fa-user-plus" },
    { name: "More", href: "/auth/more", icon: "fas fa-ellipsis-h" },
  ];

  const features = [
    {
      icon: "fas fa-shield-alt",
      title: "Enterprise Security",
      description:
        "99.9% uptime with DDoS protection and advanced security measures",
      color: "text-green-400",
    },
    {
      icon: "fas fa-bolt",
      title: "Lightning Fast",
      description:
        "Ultra-low latency servers in EU and US for optimal performance",
      color: "text-blue-400",
    },
    {
      icon: "fas fa-headset",
      title: "24/7 Support",
      description: "Expert support team available around the clock to help you",
      color: "text-purple-400",
    },
  ];

  return (
    <div className="min-h-screen flex bg-black text-white">
      <div className="flex-1 flex items-center justify-center p-8">
        <div className="w-full max-w-md">
          <div className="mb-8">
            <div className="flex bg-white/5 rounded-full p-2 border border-white/10">
              {navItems.map((item) => (
                <Link
                  key={item.name}
                  href={item.href}
                  className={`flex-1 px-4 py-2 rounded-full text-center transition-all duration-200 ${
                    pathname === item.href
                      ? "bg-white/10 text-white border border-white/20"
                      : "text-white/60 hover:text-white/80 hover:bg-white/5"
                  }`}
                >
                  <div className="flex items-center justify-center gap-2">
                    <i className={`${item.icon} text-sm`}></i>
                    <span className="text-sm font-medium">{item.name}</span>
                  </div>
                </Link>
              ))}
            </div>
          </div>

          <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-xl p-8">
            {children}
          </div>
        </div>
      </div>

      <div className="flex-1 flex items-center justify-center p-8">
        <div className="w-full max-w-lg">
          <div className="text-center mb-10">
            <h2 className="text-4xl md:text-5xl font-extrabold mb-4 leading-tight">
              <span className="block text-white">Minecraft Hosting</span>
              <span className="block text-transparent bg-clip-text bg-gradient-to-r from-blue-400 via-purple-400 to-cyan-400">
                Done Right.
              </span>
            </h2>
            <p className="text-lg md:text-xl text-gray-300 mb-8 max-w-xl mx-auto leading-snug">
              <span className="text-white font-bold">$2.80/GB</span> Minecraft
              servers. Deploy in under{" "}
              <span className="font-semibold text-blue-300">5 minutes</span>.
              Experience the difference.
            </p>
          </div>

          <div className="grid grid-cols-1 gap-4 max-w-3xl mx-auto">
            {features.map((feature) => (
              <div
                key={feature.title}
                className="flex items-center p-4 bg-white/5 rounded-xl border border-white/10 shadow-sm"
              >
                <i
                  className={`${feature.icon} text-2xl ${feature.color} mr-4`}
                />
                <div>
                  <div className="text-base font-bold text-white">
                    {feature.title}
                  </div>
                  <div className="text-gray-400 text-sm">
                    {feature.description}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
