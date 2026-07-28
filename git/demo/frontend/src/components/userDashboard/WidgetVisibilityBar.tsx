import { Eye, RotateCcw } from "lucide-react";

interface HiddenWidget {
  id: string;
  title: string;
}

interface WidgetVisibilityBarProps {
  hiddenWidgets: HiddenWidget[];
  onShow: (id: string) => void;
  onReset: () => void;
}

export default function WidgetVisibilityBar({ hiddenWidgets, onShow, onReset }: WidgetVisibilityBarProps) {
  if (hiddenWidgets.length === 0) {
    return (
      <div className="flex justify-end">
        <button
          onClick={onReset}
          className="flex items-center gap-1.5 text-xs text-gray-400 hover:text-white transition-colors"
        >
          <RotateCcw className="w-3.5 h-3.5" />
          Layout zurücksetzen
        </button>
      </div>
    );
  }

  return (
    <div className="flex flex-wrap items-center gap-2">
      {hiddenWidgets.map((widget) => (
        <button
          key={widget.id}
          onClick={() => onShow(widget.id)}
          className="flex items-center gap-1.5 px-3 py-1.5 bg-gray-800 border border-gray-700 rounded-full text-xs text-gray-300 hover:bg-gray-700 transition-colors"
        >
          <Eye className="w-3.5 h-3.5" />
          {widget.title} einblenden
        </button>
      ))}
      <button
        onClick={onReset}
        className="flex items-center gap-1.5 px-3 py-1.5 text-xs text-gray-400 hover:text-white transition-colors"
      >
        <RotateCcw className="w-3.5 h-3.5" />
        Layout zurücksetzen
      </button>
    </div>
  );
}
