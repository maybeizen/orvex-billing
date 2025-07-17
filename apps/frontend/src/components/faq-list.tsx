"use client";

import { useState } from "react";

interface FAQItem {
  question: string;
  answer: string;
}

interface FAQListProps {
  faqs: FAQItem[];
}

const FAQList = ({ faqs }: FAQListProps) => {
  const [openIndex, setOpenIndex] = useState<number | null>(null);

  const toggleFAQ = (index: number) => {
    setOpenIndex(openIndex === index ? null : index);
  };

  return (
    <div className="space-y-3">
      {faqs.map((faq, index) => (
        <div
          key={index}
          className="bg-black/30 border border-white/10 rounded-xl overflow-hidden hover:border-white/20 transition-all duration-300"
        >
          <button
            onClick={() => toggleFAQ(index)}
            className="w-full px-6 py-4 text-left flex items-center justify-between hover:bg-white/5 transition-colors"
          >
            <span className="text-white font-medium">{faq.question}</span>
            <i className="fas fa-chevron-down text-gray-400 transition-transform" />
          </button>

          {openIndex === index && (
            <div className="px-6 pb-4 border-t border-white/10">
              <p className="text-gray-300 pt-4 leading-relaxed text-sm">
                {faq.answer}
              </p>
            </div>
          )}
        </div>
      ))}
    </div>
  );
};

export default FAQList;
