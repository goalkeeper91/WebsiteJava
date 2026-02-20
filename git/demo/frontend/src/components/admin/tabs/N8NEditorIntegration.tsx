import { useState, useEffect } from "react";
import { ExternalLink, RefreshCw, AlertCircle, CheckCircle } from "lucide-react";

export default function N8NEditorIntegration() {
  const [n8nStatus, setN8nStatus] = useState<"checking" | "online" | "offline">("checking");
  const n8nUrl = "http://localhost:5678";

  useEffect(() => {
    checkN8NStatus();
  }, []);

  async function checkN8NStatus() {
    setN8nStatus("checking");
    try {
      // Try to reach n8n health endpoint
      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), 3000);

      const res = await fetch(`${n8nUrl}/healthz`, {
        signal: controller.signal,
        mode: 'no-cors' // n8n might not have CORS enabled
      });

      clearTimeout(timeoutId);
      setN8nStatus("online");
    } catch (err) {
      console.error("n8n not reachable:", err);
      setN8nStatus("offline");
    }
  }

  return (
    <div className="space-y-6">
      {/* Status Header */}
      <div className="bg-gray-800 rounded-xl p-6">
        <div className="flex items-center justify-between mb-4">
          <div>
            <h2 className="text-2xl font-bold mb-1">n8n Workflow Editor</h2>
            <p className="text-gray-400">
              Create and edit workflows directly in n8n
            </p>
          </div>
          <button
            onClick={checkN8NStatus}
            className="p-2 hover:bg-gray-700 rounded-lg transition-colors"
          >
            <RefreshCw className={`w-5 h-5 ${n8nStatus === "checking" ? "animate-spin" : ""}`} />
          </button>
        </div>

        {/* Status Badge */}
        <div className="flex items-center gap-2">
          {n8nStatus === "checking" && (
            <>
              <div className="w-3 h-3 bg-yellow-500 rounded-full animate-pulse" />
              <span className="text-yellow-400">Checking connection...</span>
            </>
          )}
          {n8nStatus === "online" && (
            <>
              <CheckCircle className="w-5 h-5 text-green-500" />
              <span className="text-green-400 font-semibold">n8n Editor Online</span>
              <span className="text-gray-500 text-sm ml-2">{n8nUrl}</span>
            </>
          )}
          {n8nStatus === "offline" && (
            <>
              <AlertCircle className="w-5 h-5 text-red-500" />
              <span className="text-red-400 font-semibold">n8n Editor Offline</span>
            </>
          )}
        </div>
      </div>

      {/* Editor Actions */}
      {n8nStatus === "online" ? (
        <div className="grid md:grid-cols-2 gap-4">
          {/* Open Editor */}
          <a
            href={n8nUrl}
            target="_blank"
            rel="noopener noreferrer"
            className="bg-gradient-to-r from-purple-600 to-blue-600 hover:from-purple-500 hover:to-blue-500 rounded-xl p-6 transition-all flex items-center justify-between group"
          >
            <div>
              <h3 className="text-xl font-bold mb-2">Open n8n Editor</h3>
              <p className="text-sm text-gray-200">
                Create new workflows in n8n
              </p>
            </div>
            <ExternalLink className="w-6 h-6 group-hover:translate-x-1 group-hover:-translate-y-1 transition-transform" />
          </a>

          {/* Quick Links */}
          <div className="bg-gray-800 rounded-xl p-6">
            <h3 className="text-lg font-bold mb-4">Quick Links</h3>
            <div className="space-y-2">
              <a
                href={`${n8nUrl}/workflows`}
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center justify-between p-3 bg-gray-900 hover:bg-gray-750 rounded-lg transition-colors group"
              >
                <span className="text-sm">View All Workflows</span>
                <ExternalLink className="w-4 h-4 opacity-50 group-hover:opacity-100" />
              </a>
              <a
                href={`${n8nUrl}/workflows/new`}
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center justify-between p-3 bg-gray-900 hover:bg-gray-750 rounded-lg transition-colors group"
              >
                <span className="text-sm">Create New Workflow</span>
                <ExternalLink className="w-4 h-4 opacity-50 group-hover:opacity-100" />
              </a>
              <a
                href={`${n8nUrl}/executions`}
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center justify-between p-3 bg-gray-900 hover:bg-gray-750 rounded-lg transition-colors group"
              >
                <span className="text-sm">View Executions</span>
                <ExternalLink className="w-4 h-4 opacity-50 group-hover:opacity-100" />
              </a>
            </div>
          </div>
        </div>
      ) : (
        <div className="bg-red-900/20 border border-red-500/30 rounded-xl p-8 text-center">
          <AlertCircle className="w-12 h-12 text-red-400 mx-auto mb-4" />
          <h3 className="text-xl font-bold mb-2">n8n Editor Not Running</h3>
          <p className="text-gray-400 mb-4">
            Make sure n8n is running on {n8nUrl}
          </p>
          <div className="bg-gray-900 rounded-lg p-4 text-left text-sm font-mono mb-4">
            <p className="text-gray-500 mb-2"># Start n8n with Docker:</p>
            <p className="text-green-400">docker-compose up -d n8n</p>
            <p className="text-gray-500 mt-4 mb-2"># Or with npx:</p>
            <p className="text-green-400">npx n8n</p>
          </div>
          <button
            onClick={checkN8NStatus}
            className="px-4 py-2 bg-gray-700 hover:bg-gray-600 rounded-lg transition-colors"
          >
            Check Again
          </button>
        </div>
      )}

      {/* Workflow Creation Guide */}
      <div className="bg-gray-800 rounded-xl p-6">
        <h3 className="text-lg font-bold mb-4">Creating Workflows for Users</h3>
        <div className="space-y-3 text-sm text-gray-400">
          <div className="flex gap-3">
            <span className="flex-shrink-0 w-6 h-6 rounded-full bg-purple-600 text-white flex items-center justify-center text-xs font-bold">
              1
            </span>
            <div>
              <strong className="text-white">Create Workflow in n8n</strong>
              <p>Design your workflow with nodes like Webhook, HTTP Request, etc.</p>
            </div>
          </div>
          <div className="flex gap-3">
            <span className="flex-shrink-0 w-6 h-6 rounded-full bg-purple-600 text-white flex items-center justify-center text-xs font-bold">
              2
            </span>
            <div>
              <strong className="text-white">Export Workflow JSON</strong>
              <p>Click the 3-dot menu → Download → Copy the JSON</p>
            </div>
          </div>
          <div className="flex gap-3">
            <span className="flex-shrink-0 w-6 h-6 rounded-full bg-purple-600 text-white flex items-center justify-center text-xs font-bold">
              3
            </span>
            <div>
              <strong className="text-white">Create Template</strong>
              <p>Go to Templates tab and paste the JSON to create a reusable template</p>
            </div>
          </div>
          <div className="flex gap-3">
            <span className="flex-shrink-0 w-6 h-6 rounded-full bg-purple-600 text-white flex items-center justify-center text-xs font-bold">
              4
            </span>
            <div>
              <strong className="text-white">Users Import</strong>
              <p>Users can now import your template from the Workflow Marketplace</p>
            </div>
          </div>
        </div>
      </div>

      {/* Note about iframe */}
      <div className="bg-blue-900/20 border border-blue-500/30 rounded-xl p-4">
        <p className="text-sm text-gray-400">
          💡 <strong className="text-white">Tip:</strong> n8n Editor must be opened in a new tab due to security restrictions.
          Use the "Open n8n Editor" button above to work with workflows.
        </p>
      </div>
    </div>
  );
}