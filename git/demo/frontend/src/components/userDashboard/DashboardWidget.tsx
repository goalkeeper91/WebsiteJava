import { motion, useDragControls, useMotionValue } from "motion/react";
import { useEffect } from "react";
import { GripVertical, X } from "lucide-react";
import type { RefObject } from "react";

interface DashboardWidgetProps {
  title: string;
  x: number;
  y: number;
  zIndex: number;
  width?: number;
  containerRef: RefObject<HTMLDivElement | null>;
  onDragEnd: (x: number, y: number) => void;
  onFocus: () => void;
  onHide: () => void;
  children: React.ReactNode;
}

// Draggable-by-titlebar window wrapper for the Live-Dashboard's freely
// arrangeable widgets. Drag is deliberately restricted to the titlebar
// (dragListener={false} + manual dragControls.start on pointerdown there)
// so clicking inputs/selects inside a widget's content never triggers a
// drag - the official Framer Motion pattern for "drag handle" UIs.
//
// Position is driven entirely through Motion's own x/y motion values (not
// plain CSS left/top) - mixing drag's transform with externally-controlled
// left/top causes the offset to double up on every subsequent drag, since
// Motion doesn't reset its transform when the underlying layout position
// changes. Keeping position and drag on the same mechanism avoids that.
export default function DashboardWidget({
  title,
  x,
  y,
  zIndex,
  width = 340,
  containerRef,
  onDragEnd,
  onFocus,
  onHide,
  children,
}: DashboardWidgetProps) {
  const dragControls = useDragControls();
  const motionX = useMotionValue(x);
  const motionY = useMotionValue(y);

  // Keeps the motion values in sync when the position changes from outside
  // (e.g. "reset layout", or switching to a different saved entry) - a
  // no-op when it's our own onDragEnd committing the same value back down.
  useEffect(() => {
    motionX.set(x);
    motionY.set(y);
  }, [x, y, motionX, motionY]);

  function handleDragEnd() {
    onDragEnd(motionX.get(), motionY.get());
  }

  return (
    <motion.div
      drag
      dragListener={false}
      dragControls={dragControls}
      dragMomentum={false}
      dragConstraints={containerRef}
      onDragEnd={handleDragEnd}
      onPointerDownCapture={onFocus}
      style={{ position: "absolute", x: motionX, y: motionY, zIndex, width }}
      className="max-w-full bg-gray-800 border border-gray-700 rounded-xl shadow-lg overflow-hidden"
    >
      <div
        onPointerDown={(e) => dragControls.start(e)}
        className="flex items-center justify-between gap-2 px-3 py-2 bg-gray-900 border-b border-gray-700 cursor-grab active:cursor-grabbing select-none"
      >
        <div className="flex items-center gap-2 min-w-0">
          <GripVertical className="w-4 h-4 text-gray-500 flex-shrink-0" />
          <span className="text-sm font-semibold text-white truncate">{title}</span>
        </div>
        <button
          onPointerDown={(e) => e.stopPropagation()}
          onClick={onHide}
          className="text-gray-500 hover:text-white transition-colors flex-shrink-0"
          title="Ausblenden"
        >
          <X className="w-4 h-4" />
        </button>
      </div>
      <div className="p-0">{children}</div>
    </motion.div>
  );
}
