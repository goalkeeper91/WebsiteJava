import { useState, useEffect } from "react";
import { Plus, Edit, Trash2, Eye, Download } from "lucide-react";

interface Workflow {
  id: number;
  user_id: string;
  name: string;
  category: string;
  n8n_workflow_json: any;
  enabled: boolean;
  created_at: string;
}

export default function WorkflowManager() {
  const [workflows, setWorkflows] = useState<Workflow[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedWorkflow, setSelectedWorkflow] = useState<Workflow | null>(null);
  const [searchTerm, setSearchTerm] = useState("");

  useEffect(() => {
    loadWorkflows();
  }, []);

  async function loadWorkflows() {
    setLoading(true);
    try {
      // TODO: Implement API endpoint
      // const res = await fetch('/api/admin/workflows', { credentials: 'include' });
      // const data = await res.json();
      // setWorkflows(data);

      // Mock data for now
      setWorkflows([]);
    } catch (err) {
      console.error("Failed to load workflows:", err);
    } finally {
      setLoading(false);
    }
  }

  const filteredWorkflows = workflows.filter((w) =>
    w.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    w.category.toLowerCase().includes(searchTerm.toLowerCase())
  );

  function handleViewJson(workflow: Workflow) {
    setSelectedWorkflow(workflow);
  }

  function handleDownloadJson(workflow: Workflow) {
    const blob = new Blob([JSON.stringify(workflow.n8n_workflow_json, null, 2)], {
      type: "application/json",
    });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `${workflow.name.replace(/\s+/g, "_")}.json`;
    a.click();
    URL.revokeObjectURL(url);
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
          <h2 className="text-2xl font-bold">Workflow Management</h2>
          <p className="text-gray-400 mt-1">
            View and manage all n8n workflows across all users
          </p>
        </div>
        <button className="flex items-center gap-2 px-4 py-2 bg-purple-600 hover:bg-purple-500 rounded-lg transition-colors">
          <Plus className="w-5 h-5" />
          New Workflow
        </button>
      </div>

      {/* Search */}
      <input
        type="text"
        placeholder="Search workflows..."
        value={searchTerm}
        onChange={(e) => setSearchTerm(e.target.value)}
        className="w-full bg-gray-800 rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-purple-500"
      />

      {/* Workflows List */}
      {filteredWorkflows.length === 0 ? (
        <div className="bg-gray-800 rounded-xl p-12 text-center">
          <p className="text-gray-500">No workflows found</p>
          <p className="text-sm text-gray-600 mt-2">
            Workflows will appear here once users create them
          </p>
        </div>
      ) : (
        <div className="grid md:grid-cols-2 gap-4">
          {filteredWorkflows.map((workflow) => (
            <div
              key={workflow.id}
              className="bg-gray-800 rounded-xl p-4 hover:bg-gray-750 transition-colors"
            >
              <div className="flex items-start justify-between mb-3">
                <div>
                  <h3 className="font-bold">{workflow.name}</h3>
                  <p className="text-sm text-gray-400">{workflow.category}</p>
                </div>
                <span
                  className={`text-xs px-2 py-1 rounded ${
                    workflow.enabled
                      ? "bg-green-600/20 text-green-400"
                      : "bg-gray-700 text-gray-400"
                  }`}
                >
                  {workflow.enabled ? "Active" : "Disabled"}
                </span>
              </div>

              <div className="text-xs text-gray-500 mb-3">
                User: {workflow.user_id}
              </div>

              <div className="flex items-center gap-2">
                <button
                  onClick={() => handleViewJson(workflow)}
                  className="flex-1 flex items-center justify-center gap-1 px-3 py-1.5 bg-blue-600 hover:bg-blue-500 rounded text-sm transition-colors"
                >
                  <Eye className="w-4 h-4" />
                  View
                </button>
                <button
                  onClick={() => handleDownloadJson(workflow)}
                  className="flex-1 flex items-center justify-center gap-1 px-3 py-1.5 bg-gray-700 hover:bg-gray-600 rounded text-sm transition-colors"
                >
                  <Download className="w-4 h-4" />
                  Export
                </button>
                <button className="px-3 py-1.5 bg-gray-700 hover:bg-gray-600 rounded text-sm transition-colors">
                  <Edit className="w-4 h-4" />
                </button>
                <button className="px-3 py-1.5 bg-gray-700 hover:bg-red-600 rounded text-sm transition-colors">
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* JSON Viewer Modal */}
      {selectedWorkflow && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-xl p-6 max-w-4xl w-full max-h-[80vh] overflow-y-auto">
            <div className="flex items-start justify-between mb-4">
              <div>
                <h3 className="text-xl font-bold">{selectedWorkflow.name}</h3>
                <p className="text-sm text-gray-400">{selectedWorkflow.category}</p>
              </div>
              <button
                onClick={() => setSelectedWorkflow(null)}
                className="text-gray-400 hover:text-white text-2xl"
              >
                ×
              </button>
            </div>

            <div className="bg-gray-900 rounded-lg p-4">
              <pre className="text-xs text-gray-400 overflow-x-auto">
                {JSON.stringify(selectedWorkflow.n8n_workflow_json, null, 2)}
              </pre>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}