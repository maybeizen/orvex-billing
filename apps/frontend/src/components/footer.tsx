import { FC } from "react";

const Footer: FC = () => {
  return (
    <footer className="bg-black border-t border-white/10">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
          <div className="col-span-1 md:col-span-2">
            <h3 className="text-2xl font-bold text-white mb-4">Orvex</h3>
            <p className="text-gray-400 mb-6 max-w-md leading-relaxed text-sm">
              Premium Minecraft hosting in the EU/US, starting at just $0.70/GB.
              Experience unmatched performance and reliability for your gaming
              community.
            </p>
            <div className="flex space-x-3">
              <a
                href="mailto:support@orvex.cc"
                className="bg-white/10 w-10 h-10 rounded-lg flex items-center justify-center hover:bg-white/20 transition-colors"
              >
                <i className="fas fa-envelope text-white" />
              </a>
              <a
                href="https://discord.gg/orvex"
                className="bg-white/10 w-10 h-10 rounded-lg flex items-center justify-center hover:bg-white/20 transition-colors"
              >
                <i className="fab fa-discord text-white" />
              </a>
              <a
                href="https://orvex.cc"
                className="bg-white/10 w-10 h-10 rounded-lg flex items-center justify-center hover:bg-white/20 transition-colors"
              >
                <i className="fas fa-globe text-white" />
              </a>
            </div>
          </div>

          <div>
            <h4 className="text-white font-semibold mb-4">Services</h4>
            <ul className="space-y-2">
              <li>
                <a
                  href="https://my.orvex.cc"
                  className="text-gray-400 hover:text-white transition-colors text-sm"
                >
                  Minecraft Hosting
                </a>
              </li>
              <li>
                <a
                  href="https://my.orvex.cc/products/discord-bots"
                  className="text-gray-400 hover:text-white transition-colors text-sm"
                >
                  Discord Bot Hosting
                </a>
              </li>
              <li>
                <a
                  href="https://my.orvex.cc"
                  className="text-gray-400 hover:text-white transition-colors text-sm"
                >
                  Custom Solutions
                </a>
              </li>
            </ul>
          </div>

          <div>
            <h4 className="text-white font-semibold mb-4">Support</h4>
            <ul className="space-y-2">
              <li>
                <a
                  href="https://orvex.cc"
                  className="text-neutral-600 hover:cursor-not-allowed transition-colors text-sm"
                >
                  Knowledge Base
                </a>
              </li>
              <li>
                <a
                  href="#contact"
                  className="text-gray-400 hover:text-white transition-colors text-sm"
                >
                  Contact Support
                </a>
              </li>
              <li>
                <a
                  href="https://status.orvex.cc"
                  className="text-gray-400 hover:text-white transition-colors text-sm"
                >
                  Server Status
                </a>
              </li>
            </ul>
          </div>
        </div>

        <div className="border-t border-white/10 mt-8 pt-6 flex flex-col sm:flex-row justify-between items-center">
          <p className="text-gray-400 text-sm">
            © 2025 Orvex. All rights reserved.
          </p>
          <div className="flex space-x-6 mt-3 sm:mt-0">
            <a
              href="/legal/privacy"
              className="text-gray-400 hover:text-white transition-colors text-sm"
            >
              Privacy Policy
            </a>
            <a
              href="/legal/terms"
              className="text-gray-400 hover:text-white transition-colors text-sm"
            >
              Terms of Service
            </a>
          </div>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
