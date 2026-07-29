import { useState } from "react";
import { Link } from "react-router-dom";
import { ArrowLeft, Users, Workflow, Package, BarChart, Wrench } from "lucide-react";
import AdminBadge from "../../components/admin/AdminBadge";
import WorkflowManager from "../../components/admin/tabs/WorkflowManager";
import TemplateManager from "../../components/admin/tabs/TemplateManager";
import N8NEditorIntegration from "../../components/admin/tabs/N8NEditorIntegration";
import BotManager from "../../components/admin/tabs/bots/BotManager";
import VoteCreatorDashboard from '../../components/dashboard/VoteCreatorDashboard';
import CustomerSubscriptionManager from "../../components/admin/tabs/customers/CustomerSubscriptionManager";

type AdminTab = "n8n-editor" | "workflows" | "templates" | "votes" | "bots" | "customers" | "stats";

export default function AdminDashboard() {
  const [activeTab, setActiveTab] = useState<AdminTab>("n8n-editor");

  const tabs = [
    { id: "n8n-editor" as AdminTab, label: "n8n Editor", icon: Wrench },
    { id: "workflows" as AdminTab, label: "Workflows", icon: Workflow },
    { id: "templates" as AdminTab, label: "Templates", icon: Package },
    { id: "votes" as AdminTab, label: "Votes", icon: Package },
    { id: "bots" as AdminTab, label: "Bot Management", icon: Users },
    { id: "customers" as AdminTab, label: "Kunden/Abos", icon: Users },
    { id: "stats" as AdminTab, label: "Statistics", icon: BarChart },
  ];

  return (
    <div className="min-h-screen bg-gray-900 text-white">
      {/* Header */}
      <div className="bg-gray-800 border-b border-gray-700 sticky top-0 z-10">
        <div className="max-w-7xl mx-auto px-6 py-4">
          <div className="flex items-center justify-between mb-4">
            <div className="flex items-center gap-4">
              <Link
                to="/dashboard"
                className="flex items-center gap-2 text-gray-400 hover:text-white transition-colors"
              >
                <ArrowLeft className="w-5 h-5" />
                <span>Back to Dashboard</span>
              </Link>
            </div>
            <AdminBadge />
          </div>

          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-bold">Admin Panel</h1>
              <p className="text-sm text-gray-400 mt-1">
                Manage workflows, templates, and system configuration
              </p>
            </div>
          </div>

          {/* Tabs */}
          <div className="flex gap-2 mt-4">
            {tabs.map((tab) => {
              const Icon = tab.icon;
              const isActive = activeTab === tab.id;

              return (
                <button
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id)}
                  className={`flex items-center gap-2 px-4 py-2 rounded-t-lg font-medium transition-all ${
                    isActive
                      ? "bg-gray-900 text-white border-t-2 border-purple-500"
                      : "bg-gray-700 text-gray-400 hover:bg-gray-600 hover:text-gray-200"
                  }`}
                >
                  <Icon className="w-4 h-4" />
                  <span>{tab.label}</span>
                </button>
              );
            })}
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="max-w-7xl mx-auto p-6">
        {activeTab === "n8n-editor" && <N8NEditorIntegration />}
        {activeTab === "workflows" && <WorkflowManager />}
        {activeTab === "templates" && <TemplateManager />}
        {activeTab === "votes" && <VoteCreatorDashboard />}
        {activeTab === "bots" && <BotManager />}
        {activeTab === "customers" && <CustomerSubscriptionManager />}
        {activeTab === "stats" && (
          <div className="bg-gray-800 rounded-xl p-8 text-center">
            <BarChart className="w-16 h-16 mx-auto mb-4 text-gray-600" />
            <h2 className="text-xl font-bold mb-2">Statistics Dashboard</h2>
            <p className="text-gray-400">Coming soon...</p>
          </div>
        )}
      </div>
    </div>
  );
}