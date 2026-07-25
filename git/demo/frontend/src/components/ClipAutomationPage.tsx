import AutomationSettingsPanel from "./clips/AutomationSettingsPanel";
import ClipHistoryPanel from "./clips/ClipHistoryPanel";

export default function ClipAutomationPage() {
  return (
    <div className="min-h-screen bg-gray-900 text-white">
      <div className="p-4 sm:p-6 max-w-4xl mx-auto space-y-8">
        <div>
          <h1 className="text-2xl sm:text-3xl font-bold mb-2 break-words">Clip-Automatisierung</h1>
          <p className="text-gray-400">
            Erkenne automatisch Highlights aus deinem Stream und teile sie auf weiteren Plattformen.
          </p>
        </div>

        <AutomationSettingsPanel />
        <ClipHistoryPanel />
      </div>
    </div>
  );
}
