"use client";

import { useState } from "react";
import { useAuth } from "@/lib/auth-context";
import DashboardLayout from "@/components/dashboard-layout";
import { LoadingScreen } from "@/components/ui/loading-spinner";
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
    return <LoadingScreen message="Loading profile..." />;
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

        <div className="grid grid-cols-1 lg:grid-cols-4 gap-8">
          {/* Sidebar Navigation */}
          <div className="lg:col-span-1">
            <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-xl p-1">
              <nav className="space-y-1">
                {tabs.map((tab) => (
                  <button
                    key={tab.id}
                    onClick={() => setActiveTab(tab.id)}
                    className={`w-full text-left p-4 rounded-lg transition-all duration-200 group ${
                      activeTab === tab.id
                        ? "bg-white/10 text-white border border-white/20"
                        : "text-gray-300 hover:text-white hover:bg-white/5"
                    }`}
                  >
                    <div className="flex items-center space-x-3">
                      <div
                        className={`w-8 h-8 rounded-lg flex items-center justify-center transition-colors duration-200 ${
                          activeTab === tab.id
                            ? "bg-gradient-to-r from-blue-500 to-purple-500"
                            : "bg-gray-700 group-hover:bg-gray-600"
                        }`}
                      >
                        <i
                          className={`${tab.icon} text-sm ${
                            activeTab === tab.id ? "text-white" : "text-gray-400"
                          }`}
                        ></i>
                      </div>
                      <div className="flex-1 min-w-0">
                        <h3 className="font-semibold text-sm">{tab.name}</h3>
                        <p className="text-xs text-gray-400 mt-1 truncate">
                          {tab.description}
                        </p>
                      </div>
                    </div>
                  </button>
                ))}
              </nav>
            </div>
          </div>

          {/* Content Area */}
          <div className="lg:col-span-3">
            <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-xl p-8">
              {renderTabContent()}
            </div>
          </div>
        </div>
      </div>
    </DashboardLayout>
  );
}