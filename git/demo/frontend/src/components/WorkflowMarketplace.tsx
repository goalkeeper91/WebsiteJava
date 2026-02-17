import { useEffect, useState } from "react";
import { Package, Download, CheckCircle, Lock, Search, Filter } from "lucide-react";

interface WorkflowTemplate {
  id: number;
  name: string;
  description: string;
  category: string;
  tierRequired: string;
  usageCount: number;
  isPublic: boolean;
  n8nWorkflowJSON?: any;
}

interface WorkflowTemplateSummary {
  id: number;
  name: string;
  description: string;
  category: string;
  tierRequired: string;
  usageCount: number;
}

export default function WorkflowMarketplace() {
  const [templates, setTemplates] = useState<WorkflowTemplateSummary[]>([]);
  const [availableTemplates, setAvailableTemplates] = useState<WorkflowTemplate[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedTemplate, setSelectedTemplate] = useState<WorkflowTemplate | null>(null);
  const [searchTerm, setSearchTerm] = useState("");
  const [categoryFilter, setCategoryFilter] = useState<string>("all");
  const [currentTier, setCurrentTier] = useState<string>("free");

  const categories = ["all", "stream_interaction", "moderation", "fun", "utility"];

  useEffect(() => {
    loadTemplates();
    checkSubscription();
  }, [categoryFilter]);

  async function checkSubscription() {
    try {
      const res = await fetch("/api/subscription", { credentials: "include" });
      if (res.ok) {
        const data = await res.json();
        setCurrentTier(data.tierID || "free");
      }
    } catch (err) {
      console.error("Failed to check subscription:", err);
    }
  }

  async function loadTemplates() {
    setLoading(true);
    try {
      // Load public templates (summaries)
      const publicRes = await fetch("/api/workflows/templates");
      if (publicRes.ok) {
        const publicData = await publicRes.json();
        setTemplates(publicData);
      }

      // Load available templates (full) for authenticated user
      const params = categoryFilter !== "all" ? `?category=${categoryFilter}` : "";
      const availableRes = await fetch(`/api/workflows/templates/available${params}`, {
        credentials: "include",
      });
      if (availableRes.ok) {
        const availableData = await availableRes.json();
        setAvailableTemplates(availableData);
      }
    } catch (err) {
      console.error("Failed to load templates:", err);
    } finally {
      setLoading(false);
    }
  }

  async function handleUseTemplate(templateId: number) {
    try {
      const res = await fetch(`/api/workflows/templates/${templateId}`, {
        credentials: "include",
      });

      if (res.ok) {
        const template = await res.json();
        setSelectedTemplate(template);
      } else if (res.status === 403) {
        alert("Dieses Template erfordert ein höheres Abo");
      } else {
        alert("Fehler beim Laden des Templates");
      }
    } catch (err) {
      console.error("Use template error:", err);
      alert("Fehler beim Laden des Templates");
    }
  }

  function canAccessTemplate(tierRequired: string): boolean {
    const tierOrder = { free: 0, pro: 1, premium: 2 };
    return tierOrder[currentTier as keyof typeof tierOrder] >= tierOrder[tierRequired as keyof typeof tierOrder];
  }

  function getTierBadge(tierRequired: string) {
    const colors = {
      free: "bg-gray-600",
      pro: "bg-blue-600",
      premium: "bg-purple-600",
    };

    return (
      <span className={`text-xs px-2 py-1 rounded ${colors[tierRequired as keyof typeof colors] || "bg-gray-600"}`}>
        {tierRequired.toUpperCase()}
      </span>
    );
  }

  function getCategoryIcon(category: string) {
    return "📦"; // Placeholder - könnte durch spezifische Icons ersetzt werden
  }

  const filteredTemplates = templates.filter((t) =>
    t.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    t.description.toLowerCase().includes(searchTerm.toLowerCase())
  );

  if (loading && templates.length === 0) {
    return (
      <div className="p-6 max-w-7xl mx-auto">
        <div className="animate-pulse space-y-4">
          <div className="h-8 bg-gray-800 rounded w-1/3"></div>
          <div className="grid md:grid-cols-3 gap-4">
            {[1, 2, 3].map((i) => (
              <div key={i} className="h-64 bg-gray-800 rounded"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="p-6 max-w-7xl mx-auto">
      {/* Header */}
      <div className="mb-6">
        <h1 className="text-3xl font-bold mb-2">Workflow Templates</h1>
        <p className="text-gray-400">
          Vorgefertigte n8n Workflows zum Importieren und Anpassen
        </p>
      </div>

      {/* Filters */}
      <div className="flex flex-col sm:flex-row gap-4 mb-6">
        {/* Search */}
        <div className="flex-1 relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
          <input
            type="text"
            placeholder="Templates durchsuchen..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="w-full bg-gray-800 rounded-lg pl-10 pr-4 py-2 focus:outline-none focus:ring-2 focus:ring-green-500"
          />
        </div>

        {/* Category Filter */}
        <div className="flex items-center gap-2">
          <Filter className="w-5 h-5 text-gray-400" />
          <select
            value={categoryFilter}
            onChange={(e) => setCategoryFilter(e.target.value)}
            className="bg-gray-800 rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-green-500"
          >
            <option value="all">Alle Kategorien</option>
            {categories.slice(1).map((cat) => (
              <option key={cat} value={cat}>
                {cat.replace("_", " ").replace(/\b\w/g, (l) => l.toUpperCase())}
              </option>
            ))}
          </select>
        </div>
      </div>

      {/* Templates Grid */}
      {filteredTemplates.length === 0 ? (
        <div className="text-center py-12 bg-gray-800 rounded-xl">
          <Package className="w-12 h-12 mx-auto mb-3 text-gray-600" />
          <p className="text-gray-500">Keine Templates gefunden</p>
        </div>
      ) : (
        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredTemplates.map((template) => {
            const hasAccess = canAccessTemplate(template.tierRequired);

            return (
              <div
                key={template.id}
                className="bg-gray-800 rounded-xl p-6 hover:bg-gray-750 transition-colors relative"
              >
                {/* Header */}
                <div className="flex items-start justify-between mb-3">
                  <div className="flex items-center gap-2">
                    <span className="text-2xl">{getCategoryIcon(template.category)}</span>
                    <h3 className="font-bold text-lg">{template.name}</h3>
                  </div>
                  {getTierBadge(template.tierRequired)}
                </div>

                {/* Description */}
                <p className="text-sm text-gray-400 mb-4 line-clamp-3">
                  {template.description}
                </p>

                {/* Stats */}
                <div className="flex items-center gap-4 text-xs text-gray-500 mb-4">
                  <span className="flex items-center gap-1">
                    <Download className="w-3 h-3" />
                    {template.usageCount} Uses
                  </span>
                  <span className="capitalize">{template.category.replace("_", " ")}</span>
                </div>

                {/* Action Button */}
                {hasAccess ? (
                  <button
                    onClick={() => handleUseTemplate(template.id)}
                    className="w-full py-2 bg-green-600 hover:bg-green-500 rounded-lg transition-colors flex items-center justify-center gap-2"
                  >
                    <Package className="w-4 h-4" />
                    Template verwenden
                  </button>
                ) : (
                  <button
                    disabled
                    className="w-full py-2 bg-gray-700 rounded-lg cursor-not-allowed flex items-center justify-center gap-2 opacity-60"
                  >
                    <Lock className="w-4 h-4" />
                    Upgrade erforderlich
                  </button>
                )}
              </div>
            );
          })}
        </div>
      )}

      {/* Template Detail Modal */}
      {selectedTemplate && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-xl p-6 max-w-4xl w-full max-h-[80vh] overflow-y-auto">
            <div className="flex items-start justify-between mb-6">
              <div>
                <h2 className="text-2xl font-bold mb-2">{selectedTemplate.name}</h2>
                <div className="flex items-center gap-2 mb-3">
                  {getTierBadge(selectedTemplate.tierRequired)}
                  <span className="text-sm text-gray-400 capitalize">
                    {selectedTemplate.category.replace("_", " ")}
                  </span>
                </div>
                <p className="text-gray-400">{selectedTemplate.description}</p>
              </div>
              <button
                onClick={() => setSelectedTemplate(null)}
                className="text-gray-400 hover:text-white text-2xl"
              >
                ×
              </button>
            </div>

            {/* Workflow JSON Preview */}
            <div className="bg-gray-900 rounded-lg p-4 mb-6">
              <h3 className="text-sm font-bold mb-2 text-gray-400">n8n Workflow JSON:</h3>
              <pre className="text-xs overflow-x-auto text-gray-500 max-h-64">
                {JSON.stringify(selectedTemplate.n8nWorkflowJSON, null, 2)}
              </pre>
            </div>

            {/* Instructions */}
            <div className="bg-blue-900/20 border border-blue-500/30 rounded-lg p-4 mb-4">
              <h3 className="font-bold mb-2 flex items-center gap-2">
                <CheckCircle className="w-5 h-5 text-blue-400" />
                So verwendest du dieses Template:
              </h3>
              <ol className="text-sm text-gray-400 space-y-2 ml-7">
                <li>1. Kopiere das JSON oben</li>
                <li>2. Öffne deine n8n Instanz</li>
                <li>3. Klicke auf "Import from JSON"</li>
                <li>4. Füge das JSON ein und passe die Parameter an</li>
                <li>5. Aktiviere den Workflow und teste ihn</li>
              </ol>
            </div>

            {/* Actions */}
            <div className="flex gap-3">
              <button
                onClick={() => {
                  navigator.clipboard.writeText(
                    JSON.stringify(selectedTemplate.n8nWorkflowJSON, null, 2)
                  );
                  alert("JSON in Zwischenablage kopiert!");
                }}
                className="flex-1 py-2 bg-green-600 hover:bg-green-500 rounded-lg transition-colors"
              >
                JSON kopieren
              </button>
              <button
                onClick={() => setSelectedTemplate(null)}
                className="px-6 py-2 bg-gray-700 hover:bg-gray-600 rounded-lg transition-colors"
              >
                Schließen
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}