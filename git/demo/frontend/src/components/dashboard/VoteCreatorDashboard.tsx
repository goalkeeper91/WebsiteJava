import { useState } from "react";
import { Plus, Loader } from "lucide-react";

interface VoteOption {
  label: string;
}

export default function VoteCreatorDashboard() {
  const [category, setCategory] = useState("hero_class");
  const [title, setTitle] = useState("");
  const [options, setOptions] = useState<VoteOption[]>([
    { label: "" },
    { label: "" }
  ]);
  const [description, setDescription] = useState("");
  const [creating, setCreating] = useState(false);

  const categories = [
    { value: "hero_class", label: "Hero Class" },
    { value: "map_biome", label: "Map/Biome" },
    { value: "story_choice", label: "Story Decision" },
    { value: "game_feature", label: "Game Feature" },
    { value: "enemy_type", label: "Enemy Type" },
    { value: "weapon", label: "Weapon" },
    { value: "game_genre", label: "Game Genre" },
    { value: "custom", label: "Custom" },
  ];

  function addOption() {
    if (options.length < 9) {
      setOptions([...options, { label: "" }]);
    }
  }

  function removeOption(index: number) {
    if (options.length > 2) {
      setOptions(options.filter((_, i) => i !== index));
    }
  }

  function updateOption(index: number, label: string) {
    const newOptions = [...options];
    newOptions[index] = { label };
    setOptions(newOptions);
  }

  async function createVote() {
    // Validate
    const validOptions = options.filter(opt => opt.label.trim() !== "");
    if (validOptions.length < 2) {
      alert("Need at least 2 options!");
      return;
    }

    if (!title.trim()) {
      alert("Title is required!");
      return;
    }

    setCreating(true);

    try {
      // Get user from auth
      const authRes = await fetch("/auth/me", { credentials: "include" });
      const user = await authRes.json();

      // Format options for n8n
      const formattedOptions = validOptions.map((opt, i) => ({
        key: String(i + 1),
        label: opt.label.trim(),
        data: {}
      }));

      // Call n8n webhook
      const response = await fetch("http://localhost:5678/webhook/vote/create", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          user_id: user.twitch_id,
          category: category,
          title: title,
          description: description,
          options: formattedOptions
        })
      });

      const result = await response.json();

      if (result.success) {
        alert("✅ Vote created! Check Twitch chat!");

        // Reset form
        setTitle("");
        setDescription("");
        setOptions([{ label: "" }, { label: "" }]);
      } else {
        alert("❌ Failed to create vote");
      }
    } catch (err) {
      console.error("Vote creation error:", err);
      alert("❌ Error creating vote. Is n8n running?");
    } finally {
      setCreating(false);
    }
  }

  return (
    <div className="p-6 max-w-4xl mx-auto">
      {/* Header */}
      <div className="mb-6">
        <h1 className="text-2xl font-bold mb-2">Create Vote</h1>
        <p className="text-gray-400">
          Start a community vote for game development decisions
        </p>
      </div>

      {/* Form */}
      <div className="bg-gray-800 rounded-xl p-6 space-y-6">
        {/* Category */}
        <div>
          <label className="block text-sm font-medium mb-2">
            Vote Category
          </label>
          <select
            value={category}
            onChange={(e) => setCategory(e.target.value)}
            className="w-full bg-gray-900 rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-purple-500"
          >
            {categories.map((cat) => (
              <option key={cat.value} value={cat.value}>
                {cat.label}
              </option>
            ))}
          </select>
        </div>

        {/* Title */}
        <div>
          <label className="block text-sm font-medium mb-2">
            Vote Title <span className="text-red-400">*</span>
          </label>
          <input
            type="text"
            placeholder="e.g., Choose the first character class!"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            className="w-full bg-gray-900 rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-purple-500"
          />
        </div>

        {/* Description */}
        <div>
          <label className="block text-sm font-medium mb-2">
            Description (Optional)
          </label>
          <textarea
            placeholder="Additional context for voters..."
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            rows={3}
            className="w-full bg-gray-900 rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-purple-500"
          />
        </div>

        {/* Options */}
        <div>
          <label className="block text-sm font-medium mb-2">
            Vote Options <span className="text-red-400">*</span>
            <span className="text-gray-500 text-xs ml-2">
              (min 2, max 9)
            </span>
          </label>
          <div className="space-y-2">
            {options.map((option, index) => (
              <div key={index} className="flex items-center gap-2">
                <div className="w-8 h-8 rounded-full bg-purple-600 flex items-center justify-center text-sm font-bold flex-shrink-0">
                  {index + 1}
                </div>
                <input
                  type="text"
                  placeholder={`Option ${index + 1}`}
                  value={option.label}
                  onChange={(e) => updateOption(index, e.target.value)}
                  className="flex-1 bg-gray-900 rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-purple-500"
                />
                {options.length > 2 && (
                  <button
                    onClick={() => removeOption(index)}
                    className="px-3 py-2 bg-red-600 hover:bg-red-500 rounded-lg transition-colors"
                  >
                    ×
                  </button>
                )}
              </div>
            ))}
          </div>

          {/* Add Option Button */}
          {options.length < 9 && (
            <button
              onClick={addOption}
              className="mt-3 flex items-center gap-2 px-4 py-2 bg-gray-700 hover:bg-gray-600 rounded-lg transition-colors text-sm"
            >
              <Plus className="w-4 h-4" />
              Add Option
            </button>
          )}
        </div>

        {/* Preview */}
        <div className="bg-gray-900 rounded-lg p-4">
          <h3 className="text-sm font-semibold mb-3 text-gray-400">
            Chat Preview:
          </h3>
          <div className="font-mono text-sm space-y-1">
            <div className="text-green-400">🗳️ NEW VOTE: {title || "..."}</div>
            <div className="text-gray-400 text-xs">{description}</div>
            <div className="mt-2">
              {options.map((opt, i) => (
                opt.label && (
                  <div key={i} className="text-gray-300">
                    !vote {i + 1} - {opt.label}
                  </div>
                )
              ))}
            </div>
            <div className="text-purple-400 mt-2 text-xs">
              Type !vote &lt;number&gt; to cast your vote!
            </div>
          </div>
        </div>

        {/* Actions */}
        <div className="flex gap-3">
          <button
            onClick={createVote}
            disabled={creating}
            className="flex-1 py-3 bg-gradient-to-r from-purple-600 to-blue-600 hover:from-purple-500 hover:to-blue-500 disabled:from-gray-600 disabled:to-gray-600 rounded-lg font-semibold transition-all flex items-center justify-center gap-2"
          >
            {creating ? (
              <>
                <Loader className="w-5 h-5 animate-spin" />
                Creating...
              </>
            ) : (
              <>
                🗳️ Create Vote
              </>
            )}
          </button>
        </div>
      </div>

      {/* Quick Tips */}
      <div className="mt-6 bg-blue-900/20 border border-blue-500/30 rounded-xl p-4">
        <h3 className="font-semibold mb-2">💡 Quick Tips:</h3>
        <ul className="text-sm text-gray-400 space-y-1">
          <li>• Keep options short and clear</li>
          <li>• 2-4 options work best for quick votes</li>
          <li>• Use templates for common vote types (hero, map, genre)</li>
          <li>• Close votes with !vote close command in chat</li>
        </ul>
      </div>
    </div>
  );
}