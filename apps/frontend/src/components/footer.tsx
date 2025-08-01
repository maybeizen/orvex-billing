import { FC } from "react";

const Footer: FC = () => {
  return (
    <footer className="bg-neutral-950 border-t border-white/10 py-12 px-8">
      <div className="max-w-6xl mx-auto">
        <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-8">
          <div>
            <h3 className="text-2xl font-semibold text-white mb-2">Orvex</h3>
            <p className="text-neutral-400 text-sm max-w-md">
              Premium Minecraft server hosting built by gamers, for gamers.
            </p>
          </div>
          
          <div className="flex space-x-4">
            <a
              href="mailto:support@orvex.cc"
              className="text-neutral-400 hover:text-white transition-colors text-sm"
            >
              Contact
            </a>
            <a
              href="https://discord.gg/orvex"
              className="text-neutral-400 hover:text-white transition-colors text-sm"
            >
              Discord
            </a>
            <a
              href="https://status.orvex.cc"
              className="text-neutral-400 hover:text-white transition-colors text-sm"
            >
              Status
            </a>
          </div>
        </div>

        <div className="border-t border-white/10 mt-8 pt-6 flex flex-col sm:flex-row justify-between items-center gap-4">
          <p className="text-neutral-400 text-sm">
            © 2025 Orvex. All rights reserved.
          </p>
          <div className="flex space-x-6">
            <a
              href="/legal/privacy"
              className="text-neutral-400 hover:text-white transition-colors text-sm"
            >
              Privacy
            </a>
            <a
              href="/legal/terms"
              className="text-neutral-400 hover:text-white transition-colors text-sm"
            >
              Terms
            </a>
          </div>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
