"use client";

import { Button } from "./ui/button";
import { useRouter } from "next/navigation";

const Hero: React.FC = () => {
  const router = useRouter();
  return (
    <section
      id="home"
      className="min-h-screen flex items-center justify-center text-white pt-14"
    >
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
        <div className="mb-8">
          <div className="inline-flex items-center px-4 py-2 bg-white/5 rounded-full border border-white/10 mb-6">
            <i className="fas fa-bolt text-white mr-2" />
            <span className="text-sm text-white/80">
              Premium EU/US Minecraft Hosting
            </span>
          </div>

          <h1 className="text-4xl md:text-6xl lg:text-7xl font-black mb-6 leading-tight">
            Lightning Fast
            <span className="block text-transparent bg-clip-text bg-gradient-to-r from-blue-400 via-purple-400 to-cyan-400">
              Minecraft Servers
            </span>
          </h1>

          <p className="text-lg md:text-xl text-gray-300 mb-8 max-w-2xl mx-auto leading-none">
            Premium Minecraft hosting, starting at{" "}
            <span className="text-white font-bold">$2.80/GB</span>. Deploy and
            start playing in under 5 minutes.
          </p>
        </div>

        <div className="flex flex-col sm:flex-row gap-4 justify-center items-center mb-12">
          <Button
            variant="primary"
            icon="fas fa-arrow-right"
            iconPosition="right"
            rounded="md"
            onClick={() => router.push("/auth/login")}
          >
            Get Started
          </Button>
          <Button
            variant="glass"
            rounded="md"
            onClick={() => router.push("/pricing")}
          >
            View Pricing
          </Button>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 max-w-3xl mx-auto">
          <div className="flex items-center justify-center p-4 bg-white/5 rounded-xl border border-white/10">
            <i className="fas fa-shield-alt text-2xl text-green-400 mr-3" />
            <div className="text-left">
              <div className="text-lg font-bold">99.9%</div>
              <div className="text-gray-400 text-sm">Uptime</div>
            </div>
          </div>

          <div className="flex items-center justify-center p-4 bg-white/5 rounded-xl border border-white/10">
            <i className="fas fa-clock text-2xl text-blue-400 mr-3" />
            <div className="text-left">
              <div className="text-lg font-bold">5 Min</div>
              <div className="text-gray-400 text-sm">Setup</div>
            </div>
          </div>

          <div className="flex items-center justify-center p-4 bg-white/5 rounded-xl border border-white/10">
            <i className="fas fa-bolt text-2xl text-purple-400 mr-3" />
            <div className="text-left">
              <div className="text-lg font-bold">24/7</div>
              <div className="text-gray-400 text-sm">Support</div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
};

export default Hero;
