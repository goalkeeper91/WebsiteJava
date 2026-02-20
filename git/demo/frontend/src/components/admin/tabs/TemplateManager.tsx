import { useState, useEffect } from "react";
import { Plus, Edit, Trash2, Eye, Upload, Package } from "lucide-react";

interface Template {
  id: number;
  name: string;
  description: string;
  category: string;
  tier_required: string;
  usage_count: number;
  is_public: boolean;
  n8n_workflow_json: any;
  created_at: string;
}

export default function TemplateManager() {
  const [templates, setTemplates] = useState<Template[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedTemplate, setSelectedTemplate] = useState<Template | null>(null);
  const [showCreateModal, setShowCreateModal] = useState(false);

  useEffect(() => {
    loadTemplates();
  }, []);

  async function loadTemplates() {
    setLoading(true);
    try {
      // TODO: Implement API endpoint
      // const res = await fetch('/api/admin/templates', { credentials: 'include' });
      // const data = await res.json();
      // setTemplates(data);

      // Mock data for now
      setTemplates([]);
    } catch (err) {
      console.error("Failed to load templates:", err);
    } finally {
      setLoading(false);
    }
  }

  async function handleDelete(templateId: number) {
    if (!confirm("Template wirklich löschen?")) {
      return;
    }

    try {
      // await fetch(`/api/admin/templates/${templateId}`, { method: 'DELETE' });
      setTemplates((prev) => prev.filter((t) => t.id !== templateId));
    } catch (err) {
      console.error("Delete failed:", err);
    }
  }

  function getTierBadge(tier: string) {
    const colors = {
      free: "bg-gray-600",
      pro: "bg-blue-600",
      premium: "bg-purple-600",
    };
    return (
      <span className={`text-xs px-2 py-1 rounded ${colors[tier as keyof typeof colors]}`}>
        {tier.toUpperCase()}
      </span>
    );
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="animate-spin rounded-full h-8 w-8 border-4 border-gray-600 border-t-purple-500"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">Template Management</h2>
          <p className="text-gray-400 mt-1">
            Create and manage workflow templates for users
          </p>
        </div>
        <button
          onClick={() => setShowCreateModal(true)}
          className="flex items-center gap-2 px-4 py-2 bg-purple-600 hover:bg-purple-500 rounded-lg transition-colors"
        >
          <Plus className="w-5 h-5" />
          New Template
        </button>
      </div>

      {/* Templates Grid */}
      {templates.length === 0 ? (
        <div className="bg-gray-800 rounded-xl p-12 text-center">
          <Package className="w-16 h-16 mx-auto mb-4 text-gray-600" />
          <p className="text-gray-500 mb-2">No templates yet</p>
          <p className="text-sm text-gray-600 mb-4">
            Create workflow templates that users can import
          </p>
          <button
            onClick={() => setShowCreateModal(true)}
            className="px-4 py-2 bg-purple-600 hover:bg-purple-500 rounded-lg transition-colors"
          >
            Create First Template
          </button>
        </div>
      ) : (
        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
          {templates.map((template) => (
            <div
              key={template.id}
              className="bg-gray-800 rounded-xl p-4 hover:bg-gray-750 transition-colors"
            >
              <div className="flex items-start justify-between mb-2">
                <h3 className="font-bold">{template.name}</h3>
                {getTierBadge(template.tier_required)}
              </div>

              <p className="text-sm text-gray-400 mb-3 line-clamp-2">
                {template.description}
              </p>

              <div className="flex items-center gap-4 text-xs text-gray-500 mb-3">
                <span>{template.usage_count} uses</span>
                <span className="capitalize">{template.category}</span>
                {template.is_public && (
                  <span className="text-green-400">Public</span>
                )}
              </div>

              <div className="flex items-center gap-2">
                <button
                  onClick={() => setSelectedTemplate(template)}
                  className="flex-1 flex items-center justify-center gap-1 px-3 py-1.5 bg-blue-600 hover:bg-blue-500 rounded text-sm transition-colors"
                >
                  <Eye className="w-4 h-4" />
                  View
                </button>
                <button className="px-3 py-1.5 bg-gray-700 hover:bg-gray-600 rounded text-sm transition-colors">
                  <Edit className="w-4 h-4" />
                </button>
                <button
                  onClick={() => handleDelete(template.id)}
                  className="px-3 py-1.5 bg-gray-700 hover:bg-red-600 rounded text-sm transition-colors"
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Create Modal */}
      {showCreateModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-xl p-6 max-w-2xl w-full">
            <h3 className="text-xl font-bold mb-4">Create New Template</h3>

            <div className="space-y-4 mb-6">
              <div>
                <label className="block text-sm font-medium mb-2">Name</label>
                <input
                  type="text"
                  placeholder="Hero Voting System"
                  className="w-full bg-gray-900 rounded px-4 py-2 focus:outline-none focus:ring-2 focus:ring-purple-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium mb-2">Description</label>
                <textarea
                  placeholder="Allows chat to vote on hero picks..."
                  rows={3}
                  className="w-full bg-gray-900 rounded px-4 py-2 focus:outline-none focus:ring-2 focus:ring-purple-500"
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium mb-2">Category</label>
                  <select className="w-full bg-gray-900 rounded px-4 py-2 focus:outline-none focus:ring-2 focus:ring-purple-500">
                    <option value="stream_interaction">Stream Interaction</option>
                    <option value="moderation">Moderation</option>
                    <option value="fun">Fun</option>
                    <option value="utility">Utility</option>
                  </select>
                </div>

                <div>
                  <label className="block text-sm font-medium mb-2">Tier Required</label>
                  <select className="w-full bg-gray-900 rounded px-4 py-2 focus:outline-none focus:ring-2 focus:ring-purple-500">
                    <option value="free">Free</option>
                    <option value="pro">Pro</option>
                    <option value="premium">Premium</option>
                  </select>
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium mb-2">n8n Workflow JSON</label>
                <textarea
                  placeholder='{"nodes": [...], "connections": {...}}'
                  rows={8}
                  className="w-full bg-gray-900 rounded px-4 py-2 font-mono text-sm focus:outline-none focus:ring-2 focus:ring-purple-500"
                />
              </div>

              <label className="flex items-center gap-2">
                <input type="checkbox" className="rounded" />
                <span className="text-sm">Make public</span>
              </label>
            </div>

            <div className="flex gap-3">
              <button
                onClick={() => setShowCreateModal(false)}
                className="flex-1 py-2 bg-gray-700 hover:bg-gray-600 rounded-lg transition-colors"
              >
                Cancel
              </button>
              <button className="flex-1 py-2 bg-purple-600 hover:bg-purple-500 rounded-lg transition-colors">
                Create Template
              </button>
            </div>
          </div>
        </div>
      )}

      {/* View Modal */}
      {selectedTemplate && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-xl p-6 max-w-4xl w-full max-h-[80vh] overflow-y-auto">
            <div className="flex items-start justify-between mb-4">
              <div>
                <h3 className="text-xl font-bold mb-1">{selectedTemplate.name}</h3>
                <p className="text-sm text-gray-400">{selectedTemplate.description}</p>
              </div>
              <button
                onClick={() => setSelectedTemplate(null)}
                className="text-gray-400 hover:text-white text-2xl"
              >
                ×
              </button>
            </div>

            <div className="bg-gray-900 rounded-lg p-4">
              <pre className="text-xs text-gray-400 overflow-x-auto">
                {JSON.stringify(selectedTemplate.n8n_workflow_json, null, 2)}
              </pre>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}