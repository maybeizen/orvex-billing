import type { Metadata } from "next";
import { Geist, Geist_Mono, Figtree, Sora } from "next/font/google";
import "./globals.css";
import "../../public/fa/css/all.min.css";
import { AuthProvider } from "@/lib/auth-context";
import Navbar from "@/components/navbar";
import Footer from "@/components/footer";

const figtree = Figtree({
  subsets: ["latin"],
  variable: "--font-figtree",
});

const sora = Sora({
  subsets: ["latin"],
  variable: "--font-sora",
});

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title:
    "Orvex | Premium Minecraft Server Hosting - Fast, Reliable, Affordable",
  description:
    "Orvex offers premium Minecraft server hosting in Europe and the US with 99.9% uptime, ultra-low latency, DDoS protection, instant setup, and 24/7 support. Start your high-performance Minecraft server today!",
  keywords: [
    "minecraft server hosting",
    "premium minecraft hosting",
    "minecraft hosting provider",
    "minecraft server rental",
    "EU minecraft hosting",
    "US minecraft hosting",
    "free minecraft hosting",
    "cheap minecraft hosting",
    "best minecraft server hosting",
    "minecraft server hosting EU",
    "minecraft server hosting US",
    "low latency minecraft hosting",
    "high performance minecraft servers",
    "minecraft hosting 99.9% uptime",
    "game server hosting",
    "minecraft hosting services",
    "minecraft server setup",
    "minecraft hosting plans",
    "dedicated minecraft servers",
    "DDoS protected minecraft hosting",
    "instant setup minecraft server",
    "24/7 minecraft server support",
    "affordable minecraft hosting",
    "minecraft server control panel",
    "minecraft java hosting",
    "minecraft bedrock hosting",
    "modded minecraft server hosting",
    "minecraft server with backups",
    "orvex minecraft hosting",
    "orvex server hosting",
  ],
  authors: [{ name: "Orvex", url: "https://orvex.cc" }],
  creator: "Orvex",
  publisher: "Orvex",
  metadataBase: new URL("https://orvex.cc"),
  alternates: {
    canonical: "https://orvex.cc/",
  },
  openGraph: {
    title:
      "Orvex | Premium Minecraft Server Hosting - Fast, Reliable, Affordable",
    description:
      "Host your Minecraft server with Orvex for unbeatable performance, 99.9% uptime, DDoS protection, and 24/7 support. EU & US locations. Start your server instantly!",
    url: "https://orvex.cc",
    siteName: "Orvex",
    locale: "en_US",
    type: "website",
    images: [],
  },
  twitter: {
    card: "summary_large_image",
    title:
      "Orvex | Premium Minecraft Server Hosting - Fast, Reliable, Affordable",
    description:
      "Orvex provides high-performance Minecraft server hosting in EU & US with 99.9% uptime, DDoS protection, and instant setup. Start your server today!",
    images: [],
    site: "@orvexcc",
    creator: "@orvexcc",
  },
  icons: {
    icon: [
      { url: "/favicon.png", type: "image/png" },
      { url: "/favicon.ico", type: "image/x-icon" },
    ],
    apple: "/favicon.png",
    shortcut: "/favicon.png",
  },
  robots: {
    index: true,
    follow: true,
    nocache: false,
    googleBot: {
      index: true,
      follow: true,
      "max-image-preview": "large",
      "max-snippet": -1,
      "max-video-preview": -1,
    },
  },
  category: "Minecraft Server Hosting",
  applicationName: "Orvex",
  generator: "Next.js",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={`${figtree.variable} ${sora.variable} antialiased`}>
        <AuthProvider>{children}</AuthProvider>
      </body>
    </html>
  );
}
