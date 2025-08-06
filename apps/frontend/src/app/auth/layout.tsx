"use client";

import { ReactNode } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";

interface AuthLayoutProps {
  children: ReactNode;
}

export default function AuthLayout({ children }: AuthLayoutProps) {
  const pathname = usePathname();

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
      <div className="flex-1 flex flex-col justify-center p-8 lg:p-12">
        <div className="w-full max-w-lg mx-auto">
          <div className="mb-8">
            <Link
              href="/"
              className="inline-flex items-center gap-2 text-white/60 hover:text-white transition-colors duration-200"
            >
              <i className="fas fa-arrow-left text-sm"></i>
              <span className="text-sm font-medium">Back to Home</span>
            </Link>
          </div>

          <div className="mb-10">
            <h1 className="text-4xl md:text-5xl font-extrabold mb-4 leading-tight">
              <span className="block text-white">Minecraft Hosting</span>
              <span className="block text-violet-400">Done Right.</span>
            </h1>
            <p className="text-lg md:text-xl text-gray-300 mb-8 max-w-xl leading-snug">
              <span className="text-white font-bold">$2.80/GB</span> Minecraft
              servers. Deploy in under{" "}
              <span className="font-semibold text-green-200">5 minutes</span>.
              Experience the difference.
            </p>
          </div>

          <div className="grid grid-cols-1 gap-4 mb-8">
            {features.map((feature) => (
              <div
                key={feature.title}
                className="flex items-center p-4 bg-white/5 rounded-xl border border-white/10 shadow-sm hover:bg-white/10 transition-colors duration-200"
              >
                <i
                  className={`${feature.icon} text-2xl ${feature.color} mr-4 flex-shrink-0`}
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

          <div className="text-center">
            <Link
              href="/"
              className="inline-flex items-center gap-2 px-6 py-3 bg-violet-500 hover:bg-violet-600 text-white font-semibold rounded-lg transition-all duration-200 shadow-lg hover:shadow-xl"
            >
              <span>Learn More</span>
              <i className="fas fa-arrow-right text-sm"></i>
            </Link>
          </div>
        </div>
      </div>

      <div className="flex-1 flex items-center justify-center p-8 lg:p-12">
        <div className="w-full max-w-md">
          <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-xl p-8 shadow-2xl">
            {children}
          </div>
        </div>
      </div>
    </div>
  );
}
