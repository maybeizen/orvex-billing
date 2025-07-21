"use client";

import { useState } from "react";
import { useAuth } from "@/lib/auth-context";
import DashboardLayout from "@/components/dashboard-layout";
import { LoadingScreen } from "@/components/ui/loading-spinner";
import { Skeleton } from "@/components/ui/skeleton";
import { ProfileGeneral } from "@/components/profile/ProfileGeneral";
import { ProfileSecurity } from "@/components/profile/ProfileSecurity";
import { ProfileAvatar } from "@/components/profile/ProfileAvatar";
import { ProfileTwoFactor } from "@/components/profile/ProfileTwoFactor";

type ProfileTab = "general" | "security" | "avatar" | "2fa";

const tabs = [
  {
    id: "general" as ProfileTab,
    name: "General",
    icon: "fas fa-user",
    description: "Name, email, and basic information",
  },
  {
    id: "security" as ProfileTab,
    name: "Security",
    icon: "fas fa-shield-alt",
    description: "Password and account security",
  },
  {
    id: "avatar" as ProfileTab,
    name: "Avatar",
    icon: "fas fa-image",
    description: "Profile picture and display settings",
  },
  {
    id: "2fa" as ProfileTab,
    name: "Two-Factor Auth",
    icon: "fas fa-mobile-alt",
    description: "Enhanced account security",
  },
];

export default function ProfilePage() {
  const { user, loading } = useAuth();
  const [activeTab, setActiveTab] = useState<ProfileTab>("general");

  if (loading) {
    return (
      <DashboardLayout>
        <div className="space-y-8">
          <div>
            <Skeleton className="h-10 w-64 mb-2" />
            <Skeleton className="h-6 w-96" />
          </div>
          <Skeleton className="h-16 w-full rounded-xl" />
          <Skeleton className="h-96 w-full rounded-xl" />
        </div>
      </DashboardLayout>
    );
  }

  if (!user) {
    return null;
  }

  const renderTabContent = () => {
    switch (activeTab) {
      case "general":
        return <ProfileGeneral />;
      case "security":
        return <ProfileSecurity />;
      case "avatar":
        return <ProfileAvatar />;
      case "2fa":
        return <ProfileTwoFactor />;
      default:
        return <ProfileGeneral />;
    }
  };

  return (
    <DashboardLayout>
      <div className="space-y-8">
        {/* Header */}
        <div>
          <h1 className="text-4xl font-black text-white mb-2">
            Profile Settings
          </h1>
          <p className="text-gray-300 text-lg">
            Manage your account settings and preferences
          </p>
        </div>

        {/* Navbar-style Navigation */}
        <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-xl p-1">
          <nav className="flex flex-wrap gap-1">
            {tabs.map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`flex items-center space-x-2 px-4 py-3 rounded-lg transition-all duration-200 ${
                  activeTab === tab.id
                    ? "bg-white/10 text-white border border-white/20"
                    : "text-gray-400 hover:text-white hover:bg-white/5"
                }`}
              >
                <i className={`${tab.icon} text-sm`}></i>
                <span className="font-medium text-sm">{tab.name}</span>
              </button>
            ))}
          </nav>
        </div>

        {/* Content Area */}
        <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-xl p-8">
          {renderTabContent()}
        </div>
      </div>
    </DashboardLayout>
  );
}
