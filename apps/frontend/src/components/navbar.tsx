"use client";

import { useState } from "react";
import Link from "next/link";
import { Button } from "./ui/button";
import { useRouter } from "next/navigation";
import { useAuth } from "@/lib/auth-context";

const Navbar: React.FC = () => {
  const router = useRouter();
  const { user, logout } = useAuth();
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  const toggleMenu = () => setIsMenuOpen((prev) => !prev);
  const closeMenu = () => setIsMenuOpen(false);

  const handleLogout = async () => {
    await logout();
    router.push("/");
    closeMenu();
  };

  return (
    <nav className="fixed top-0 left-0 right-0 z-50 bg-black/70 backdrop-blur-2xl border-b border-white/5">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-14">
          <div className="flex items-center space-x-2">
            <Link href="/" className="text-xl font-bold text-white">
              Orvex
            </Link>
          </div>

          <div className="hidden md:flex items-center space-x-1">
            {[
              ["Home", "/#home"],
              ["Services", "/#services"],
              ["Pricing", "/#pricing"],
              ["About", "/#about"],
              ["Contact", "/#contact"],
            ].map(([label, href]) => (
              <Link
                key={label}
                href={href}
                className="text-white/70 hover:text-white px-3 py-2 text-sm font-medium transition-colors rounded-md hover:bg-white/5"
              >
                {label}
              </Link>
            ))}
          </div>

          <div className="hidden md:block">
            {user ? (
              <div className="flex items-center space-x-3">
                <span className="text-white/70 text-sm">
                  Welcome, {user.first_name}
                </span>
                <Button
                  variant="glass"
                  size="sm"
                  rounded="md"
                  onClick={() => router.push("/dashboard")}
                >
                  Dashboard
                </Button>
                <Button
                  variant="ghost"
                  size="sm"
                  rounded="md"
                  onClick={handleLogout}
                >
                  Logout
                </Button>
              </div>
            ) : (
              <div className="flex items-center space-x-2">
                <Button
                  variant="ghost"
                  size="sm"
                  rounded="md"
                  onClick={() => router.push("/auth/login")}
                >
                  Login
                </Button>
                <Button
                  variant="glass"
                  size="sm"
                  rounded="md"
                  onClick={() => router.push("/auth/register")}
                >
                  Get Started
                </Button>
              </div>
            )}
          </div>

          <div className="md:hidden">
            <button
              onClick={toggleMenu}
              className="text-white hover:text-white/80 p-2"
            >
              {isMenuOpen ? (
                <i className="fas fa-times text-white" />
              ) : (
                <i className="fas fa-bars text-white" />
              )}
            </button>
          </div>
        </div>
      </div>

      {isMenuOpen && (
        <div className="md:hidden border-t border-white/5">
          <div className="px-4 pt-2 pb-3 space-y-1 bg-black/90">
            {[
              ["Home", "/#home"],
              ["Services", "/#services"],
              ["Pricing", "/#pricing"],
              ["About", "/#about"],
              ["Contact", "/#contact"],
            ].map(([label, href]) => (
              <a
                key={label}
                href={href}
                onClick={closeMenu}
                className="text-white/70 hover:text-white block px-3 py-2 text-sm font-medium"
              >
                {label}
              </a>
            ))}
            <div className="pt-2">
              {user ? (
                <div className="space-y-2">
                  <div className="text-white/70 text-sm px-3">
                    Welcome, {user.first_name}
                  </div>
                  <Button
                    variant="glass"
                    rounded="md"
                    className="w-full block text-center text-white"
                    onClick={() => {
                      router.push("/dashboard");
                      closeMenu();
                    }}
                  >
                    Dashboard
                  </Button>
                  <Button
                    variant="ghost"
                    rounded="md"
                    className="w-full block text-center text-white"
                    onClick={handleLogout}
                  >
                    Logout
                  </Button>
                </div>
              ) : (
                <div className="space-y-2">
                  <Button
                    variant="ghost"
                    rounded="md"
                    className="w-full block text-center text-white"
                    onClick={() => {
                      router.push("/auth/login");
                      closeMenu();
                    }}
                  >
                    Login
                  </Button>
                  <Button
                    variant="glass"
                    rounded="md"
                    className="w-full block text-center text-white"
                    onClick={() => {
                      router.push("/auth/register");
                      closeMenu();
                    }}
                  >
                    Get Started
                  </Button>
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </nav>
  );
};

export default Navbar;
