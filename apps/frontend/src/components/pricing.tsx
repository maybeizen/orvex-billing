"use client";

import { euPlans, usPlans } from "../data/pricing";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "./ui/button";

const Pricing: React.FC = () => {
  const [selectedRegion, setSelectedRegion] = useState<"eu" | "us">("eu");
  const router = useRouter();

  const currentPlans = selectedRegion === "eu" ? euPlans : usPlans;
  const topPlans = currentPlans.slice(0, 3);
  const bottomPlans = currentPlans.slice(3, 5);

  return (
    <section id="pricing" className="py-16">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center mb-12">
          <h2 className="text-3xl md:text-4xl font-bold text-white mb-4">
            Simple, Transparent Pricing
          </h2>
          <p className="text-lg text-gray-300 max-w-2xl mx-auto">
            Choose the perfect plan for your Minecraft server. All plans include
            premium features.
          </p>
        </div>

        <div className="flex justify-center mb-8">
          <div className="bg-zinc-900 border border-white/10 rounded-lg p-1">
            <button
              onClick={() => setSelectedRegion("eu")}
              className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
                selectedRegion === "eu"
                  ? "bg-white text-black"
                  : "text-gray-300 hover:text-white"
              }`}
            >
              <i
                className={`fas fa-globe inline mr-2 ${
                  selectedRegion === "eu" ? "text-gray-900" : "text-white"
                }`}
              />
              EU Region
            </button>
            <button
              onClick={() => setSelectedRegion("us")}
              className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
                selectedRegion === "us"
                  ? "bg-white text-black"
                  : "text-gray-300 hover:text-white"
              }`}
            >
              <i
                className={`fas fa-globe inline mr-2 ${
                  selectedRegion === "us" ? "text-gray-900" : "text-white"
                }`}
              />
              US Region
            </button>
          </div>
        </div>

        {selectedRegion === "us" && (
          <div className="text-center mb-6">
            <div className="inline-flex items-center px-3 py-1 bg-yellow-500/10 rounded-full border border-yellow-500/20">
              <span className="text-xs text-yellow-400">
                US plans have a +$0.17/GB upcharge due to higher operating costs
              </span>
            </div>
          </div>
        )}

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
          {topPlans.map((plan: any, index: number) => (
            <div className="relative group" key={index}>
              {plan.isPopular && (
                <div className="absolute -inset-1 bg-gradient-to-r from-blue-500 via-purple-500 to-cyan-500 rounded-xl blur opacity-30 group-hover:opacity-50 transition duration-300" />
              )}

              <div
                className={`relative bg-zinc-950 border ${
                  plan.isPopular ? "border-purple-500/50" : "border-white/10"
                } rounded-xl p-6 hover:border-white/20 transition-all duration-300 h-full`}
              >
                {plan.isPopular && (
                  <div className="absolute -top-3 left-1/2 transform -translate-x-1/2">
                    <div className="bg-gradient-to-r from-blue-500 to-purple-500 text-white px-3 py-1 rounded-full text-xs font-semibold flex items-center">
                      <i className="fas fa-crown text-white mr-1" />
                      Most Popular
                    </div>
                  </div>
                )}

                <div className="mb-6">
                  <div className="flex items-center justify-between mb-3">
                    <h3 className="text-lg font-bold text-white">
                      {plan.name}
                    </h3>
                    {plan.isPopular && (
                      <i className="fas fa-star text-purple-400" />
                    )}
                  </div>
                  <p className="text-gray-400 text-sm mb-4">
                    {plan.description}
                  </p>
                  <div className="flex items-baseline mb-1">
                    <span className="text-3xl font-bold text-white">
                      ${plan.price}
                    </span>
                    <span className="text-gray-400 ml-1">/month</span>
                  </div>
                  <div className="text-xs text-gray-500">
                    ${(plan.price / plan.ram).toFixed(2)}/GB
                  </div>
                </div>

                <ul className="space-y-3 mb-6">
                  {plan.features.map((feature: string, i: number) => (
                    <li className="flex items-start" key={i}>
                      <i className="fas fa-check text-green-400 mr-2 mt-1 flex-shrink-0" />
                      <span className="text-gray-300 text-sm">{feature}</span>
                    </li>
                  ))}
                </ul>

                <Button
                  variant={plan.isPopular ? "primary" : "glass"}
                  rounded="md"
                  className="w-full block text-center text-white"
                  onClick={() => router.push("/auth/login")}
                >
                  Get Started
                </Button>
              </div>
            </div>
          ))}
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 max-w-4xl mx-auto">
          {bottomPlans.map((plan: any, index: number) => (
            <div
              className={`relative bg-zinc-950 border border-white/10 rounded-xl p-6 hover:border-white/20 transition-all duration-300`}
              key={index}
            >
              <div className="mb-6">
                <div className="flex items-center justify-between mb-3">
                  <h3 className="text-lg font-bold text-white">{plan.name}</h3>
                  <i className="fas fa-bolt text-yellow-400" />
                </div>
                <p className="text-gray-400 text-sm mb-4">{plan.description}</p>
                <div className="flex items-baseline mb-1">
                  <span className="text-3xl font-bold text-white">
                    ${plan.price}
                  </span>
                  <span className="text-gray-400 ml-1">/month</span>
                </div>
                <div className="text-xs text-gray-500">
                  ${(plan.price / plan.ram).toFixed(2)}/GB
                </div>
              </div>

              <ul className="space-y-3 mb-6">
                {plan.features.map((feature: string, i: number) => (
                  <li className="flex items-start" key={i}>
                    <i className="fas fa-check text-green-400 mr-2 mt-1 flex-shrink-0" />
                    <span className="text-gray-300 text-sm">{feature}</span>
                  </li>
                ))}
              </ul>

              <Button
                variant="glass"
                rounded="md"
                className="w-full block text-center text-white"
                onClick={() => router.push("/auth/login")}
              >
                Get Started
              </Button>
            </div>
          ))}
        </div>

        <div className="text-center mt-12">
          <p className="text-gray-400 mb-4">
            Need a custom solution?{" "}
            <a
              href="#contact"
              className="text-white hover:text-gray-300 underline"
            >
              Contact us
            </a>{" "}
            for enterprise pricing.
          </p>
          <div className="inline-flex items-center px-3 py-1 bg-white/5 rounded-full border border-white/10">
            <span className="text-xs text-gray-400">
              Payments processed securely via Stripe
            </span>
          </div>
        </div>
      </div>
    </section>
  );
};

export default Pricing;
