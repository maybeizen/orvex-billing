"use client";

import { useState } from "react";
import { euPlans, usPlans } from "../data/pricing";
import { Button } from "./ui/button";

const Pricing: React.FC = () => {
  const [selectedRegion, setSelectedRegion] = useState<"eu" | "us" | "apac">(
    "eu"
  );
  const plans = selectedRegion === "eu" ? euPlans : usPlans;

  return (
    <section
      id="pricing"
      className="py-20 px-8 bg-neutral-950 text-white relative"
    >
      <div className="max-w-6xl mx-auto">
        <div className="text-center mb-16">
          <h2 className="text-4xl font-semibold mb-4 text-transparent bg-clip-text bg-gradient-to-r from-neutral-400 via-white to-neutral-300">
            Simple, Transparent Pricing
          </h2>
          <p className="text-lg font-normal text-neutral-400 max-w-2xl mx-auto mb-8">
            Choose the perfect plan for your Minecraft server. All plans include
            DDoS protection, 24/7 support, and instant setup.
          </p>

          <div className="inline-flex items-center bg-neutral-900/50 rounded-lg p-1 border border-white/10">
            <button
              onClick={() => setSelectedRegion("eu")}
              className={`px-6 py-2 rounded-md text-sm font-medium transition-all ${
                selectedRegion === "eu"
                  ? "bg-violet-600 text-white"
                  : "text-neutral-400 hover:text-white"
              }`}
            >
              <i className="fas fa-flag mr-2"></i>
              Europe
            </button>
            <button
              onClick={() => setSelectedRegion("us")}
              className={`px-6 py-2 rounded-md text-sm font-medium transition-all ${
                selectedRegion === "us"
                  ? "bg-violet-600 text-white"
                  : "text-neutral-400 hover:text-white"
              }`}
            >
              <i className="fas fa-flag-usa mr-2"></i>
              United States
            </button>
            <button
              className={`px-6 py-2 rounded-md text-sm font-medium transition-all text-neutral-500 cursor-not-allowed`}
            >
              <i className="fas fa-question mr-2"></i>
              APAC
            </button>
          </div>

          {selectedRegion === "us" && (
            <div className="mt-4 p-3 bg-amber-500/10 border border-amber-500/20 rounded-lg max-w-lg mx-auto">
              <div className="flex items-center justify-center text-amber-400">
                <i className="fas fa-exclamation-triangle mr-2"></i>
                <span className="text-sm font-medium">
                  US servers include a +$0.17/GB upcharge due to higher
                  operating costs
                </span>
              </div>
            </div>
          )}
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-12">
          {plans.map((plan, index) => (
            <div
              key={index}
              className={`p-6 bg-neutral-900/50 backdrop-blur-sm border relative rounded-3xl ${
                plan.isPopular
                  ? "border-violet-500/50 ring-1 ring-violet-500/20"
                  : "border-white/10"
              }`}
            >
              {plan.isPopular && (
                <div className="absolute -top-3 left-1/2 transform -translate-x-1/2">
                  <div className="bg-gradient-to-r from-violet-600 to-purple-600 text-white text-xs font-semibold px-3 py-1 rounded-full">
                    Most Popular
                  </div>
                </div>
              )}

              <div className="flex flex-col h-full">
                <div className="mb-6">
                  <h3 className="text-xl font-semibold text-white mb-2">
                    {plan.name}
                  </h3>
                  <p className="text-neutral-400 text-sm mb-4">
                    {plan.description}
                  </p>
                  <div className="flex items-baseline">
                    <span className="text-3xl font-bold text-white">
                      ${plan.price}
                    </span>
                    <span className="text-neutral-400 text-sm ml-1">
                      /month
                    </span>
                  </div>
                </div>

                <div className="mb-6">
                  <div className="grid grid-cols-2 gap-4 text-sm">
                    <div className="text-center p-2 bg-white/5 rounded-lg">
                      <div className="text-white font-semibold">
                        {plan.ram}GB
                      </div>
                      <div className="text-neutral-400 text-xs">RAM</div>
                    </div>
                    <div className="text-center p-2 bg-white/5 rounded-lg">
                      <div className="text-white font-semibold">
                        {plan.cores}
                      </div>
                      <div className="text-neutral-400 text-xs">CPU Cores</div>
                    </div>
                  </div>
                </div>

                <div className="flex-grow mb-6">
                  <ul className="space-y-2">
                    {plan.features.map((feature, featureIndex) => (
                      <li
                        key={featureIndex}
                        className="flex items-center text-sm"
                      >
                        <i className="fas fa-check text-green-400 mr-2"></i>
                        <span className="text-neutral-300">{feature}</span>
                      </li>
                    ))}
                  </ul>
                </div>

                <Button
                  variant={plan.isPopular ? "primary" : "glass"}
                  size="md"
                  rounded="lg"
                  fullWidth
                >
                  Get Started
                </Button>
              </div>
            </div>
          ))}
        </div>

        <div className="text-center">
          <div className="bg-gradient-to-r from-violet-600/20 to-purple-600/20 rounded-2xl p-8 border border-white/10">
            <h3 className="text-2xl font-semibold mb-4 text-white">
              Need Something Custom?
            </h3>
            <p className="text-neutral-300 text-lg mb-6">
              Looking for enterprise solutions or have specific requirements?
              Contact us for custom hosting solutions tailored to your needs.
            </p>
            <Button
              variant="white"
              size="md"
              rounded="lg"
              icon="fas fa-envelope"
              iconPosition="left"
            >
              Contact Sales
            </Button>
          </div>
        </div>
      </div>
    </section>
  );
};

export default Pricing;
