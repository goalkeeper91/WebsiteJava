import { useState, useEffect } from "react";
import { fetchStreamInfo, updateStreamInfo, searchCategories } from "../../features/stream/api";
import type { StreamInfo, Category } from "../../features/stream/types";

export default function StreamInfoPanel() {
  const [streamInfo, setStreamInfo] = useState<StreamInfo | null>(null);
  const [editing, setEditing] = useState(false);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  // Edit state
  const [title, setTitle] = useState("");
  const [gameId, setGameId] = useState("");
  const [gameName, setGameName] = useState("");

  // Category search
  const [categorySearch, setCategorySearch] = useState("");
  const [categories, setCategories] = useState<Category[]>([]);
  const [searchingCategories, setSearchingCategories] = useState(false);

  useEffect(() => {
    loadStreamInfo();
  }, []);

  // Debounced category search
  useEffect(() => {
    if (categorySearch.length < 2) {
      setCategories([]);
      return;
    }

    const timer = setTimeout(async () => {
      setSearchingCategories(true);
      try {
        const results = await searchCategories(categorySearch);
        setCategories(results);
      } catch (err) {
        console.error("Fehler beim Suchen:", err);
      } finally {
        setSearchingCategories(false);
      }
    }, 300);

    return () => clearTimeout(timer);
  }, [categorySearch]);

  async function loadStreamInfo() {
    setLoading(true);
    try {
      const data = await fetchStreamInfo();
      setStreamInfo(data);
      setTitle(data.title);
      setGameId(data.gameId);
      setGameName(data.gameName);
    } catch (err) {
      console.error("Fehler beim Laden:", err);
      alert("Fehler beim Laden der Stream-Info");
    } finally {
      setLoading(false);
    }
  }

  async function handleSave() {
    setSaving(true);
    try {
      const updated = await updateStreamInfo({
        title,
        gameId,
      });

      setStreamInfo(updated);
      setEditing(false);
    } catch (err) {
      console.error("Fehler beim Speichern:", err);
      alert("Fehler beim Speichern der Stream-Info");
    } finally {
      setSaving(false);
    }
  }

  function handleCategorySelect(category: Category) {
    setGameId(category.id);
    setGameName(category.name);
    setCategorySearch("");
    setCategories([]);
  }

  if (loading) {
    return (
      <div className="bg-gray-800 rounded-lg p-4">
        <div className="animate-pulse space-y-3">
          <div className="h-4 bg-gray-700 rounded w-1/2"></div>
          <div className="h-8 bg-gray-700 rounded"></div>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-gray-800 rounded-lg p-4">
      <h2 className="text-lg font-bold mb-3 flex items-center gap-2">
        <span className="text-red-400">🔴</span> Stream Info
      </h2>

      {!editing ? (
        <div className="space-y-3">
          <div>
            <p className="text-xs text-gray-400 mb-1">Titel</p>
            <p className="text-sm font-medium break-words">{streamInfo?.title}</p>
          </div>
          <div>
            <p className="text-xs text-gray-400 mb-1">Kategorie</p>
            <p className="text-sm font-medium">{streamInfo?.gameName || "Keine Kategorie"}</p>
          </div>
          <button
            onClick={() => setEditing(true)}
            className="w-full mt-2 px-4 py-2 bg-blue-600 rounded hover:bg-blue-500 transition-colors text-sm"
          >
            ✏️ Bearbeiten
          </button>
        </div>
      ) : (
        <div className="space-y-3">
          {/* Titel */}
          <div>
            <label className="block text-xs text-gray-400 mb-1">Titel</label>
            <textarea
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              className="w-full bg-gray-900 rounded px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 resize-none"
              maxLength={140}
              rows={3}
            />
            <p className="text-xs text-gray-500 mt-1">{title.length}/140</p>
          </div>

          {/* Kategorie */}
          <div>
            <label className="block text-xs text-gray-400 mb-1">Kategorie</label>
            <div className="relative">
              <input
                type="text"
                value={categorySearch || gameName}
                onChange={(e) => {
                  setCategorySearch(e.target.value);
                  setGameName(e.target.value);
                }}
                className="w-full bg-gray-900 rounded px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="Suche Spiel/Kategorie..."
              />

              {/* Category Dropdown */}
              {categories.length > 0 && (
                <div className="absolute z-10 w-full mt-1 bg-gray-900 rounded shadow-lg max-h-48 overflow-y-auto">
                  {categories.map((cat) => (
                    <button
                      key={cat.id}
                      onClick={() => handleCategorySelect(cat)}
                      className="w-full text-left px-3 py-2 hover:bg-gray-800 transition-colors text-sm"
                    >
                      {cat.name}
                    </button>
                  ))}
                </div>
              )}

              {searchingCategories && (
                <div className="absolute right-3 top-2">
                  <div className="animate-spin rounded-full h-4 w-4 border-2 border-gray-600 border-t-blue-500"></div>
                </div>
              )}
            </div>
          </div>

          {/* Buttons */}
          <div className="flex gap-2 mt-2">
            <button
              onClick={() => {
                setEditing(false);
                setTitle(streamInfo?.title || "");
                setGameId(streamInfo?.gameId || "");
                setGameName(streamInfo?.gameName || "");
                setCategorySearch("");
              }}
              disabled={saving}
              className="flex-1 px-4 py-2 bg-gray-700 rounded hover:bg-gray-600 transition-colors text-sm disabled:opacity-50"
            >
              Abbrechen
            </button>
            <button
              onClick={handleSave}
              disabled={saving}
              className="flex-1 px-4 py-2 bg-green-600 rounded hover:bg-green-500 transition-colors text-sm disabled:opacity-50"
            >
              {saving ? "Speichern..." : "💾 Speichern"}
            </button>
          </div>
        </div>
      )}
    </div>
  );
}