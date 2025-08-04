"use client";

import { useState } from "react";
import { faqs } from "../data/faq";
import { Button } from "./ui/button";

const FAQ: React.FC = () => {
  const [openIndex, setOpenIndex] = useState<number | null>(null);

  const toggleFaq = (index: number) => {
    setOpenIndex(openIndex === index ? null : index);
  };

  return (
    <section id="faq" className="py-20 px-8 bg-neutral-950 text-white relative">
      <div className="max-w-4xl mx-auto">
        <div className="text-center mb-16">
          <h2 className="text-4xl font-semibold mb-4 text-transparent bg-clip-text bg-gradient-to-r from-neutral-400 via-white to-neutral-300">
            Frequently Asked Questions
          </h2>
          <p className="text-lg font-normal text-neutral-400 max-w-2xl mx-auto">
            Got questions? We've got answers. Find everything you need to know about our Minecraft hosting services.
          </p>
        </div>

        <div className="space-y-4">
          {faqs.map((faq, index) => (
            <div
              key={index}
              className="bg-neutral-900/50 backdrop-blur-sm border border-white/10 rounded-lg overflow-hidden transition-all duration-200 hover:border-white/20"
            >
              <button
                onClick={() => toggleFaq(index)}
                className="w-full px-6 py-4 text-left flex items-center justify-between hover:bg-white/5 transition-colors"
              >
                <h3 className="text-lg font-medium text-white pr-4">
                  {faq.question}
                </h3>
                <div className="flex-shrink-0">
                  <i
                    className={`fas transition-transform duration-200 text-violet-400 ${
                      openIndex === index
                        ? "fa-minus"
                        : "fa-plus"
                    }`}
                  ></i>
                </div>
              </button>
              
              <div
                className={`overflow-hidden transition-all duration-200 ${
                  openIndex === index ? "max-h-96 opacity-100" : "max-h-0 opacity-0"
                }`}
              >
                <div className="px-6 pb-4">
                  <div className="h-px bg-white/10 mb-4"></div>
                  <p className="text-neutral-300 leading-relaxed">
                    {faq.answer}
                  </p>
                </div>
              </div>
            </div>
          ))}
        </div>

        <div className="text-center mt-12">
          <div className="bg-gradient-to-r from-violet-600/20 to-purple-600/20 rounded-2xl p-8 border border-white/10">
            <h3 className="text-2xl font-semibold mb-4 text-white">
              Still Have Questions?
            </h3>
            <p className="text-neutral-300 text-lg mb-6">
              Our support team is here to help 24/7. Reach out to us through our ticketing system 
              or join our Discord community for quick assistance.
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Button 
                variant="white" 
                size="md" 
                rounded="lg" 
                icon="fas fa-ticket-alt"
                iconPosition="left"
              >
                Open Support Ticket
              </Button>
              <Button 
                variant="primary" 
                size="md" 
                rounded="lg" 
                icon="fab fa-discord"
                iconPosition="left"
                className="bg-indigo-600 hover:bg-indigo-700"
              >
                Join Discord
              </Button>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
};

export default FAQ;