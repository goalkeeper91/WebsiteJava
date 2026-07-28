import { motion, useDragControls, useMotionValue } from "motion/react";
import { useEffect, useRef, useState } from "react";
import { GripVertical, X, Maximize2 } from "lucide-react";
import type { RefObject } from "react";

interface DashboardWidgetProps {
  title: string;
  x: number;
  y: number;
  zIndex: number;
  width?: number;
  height?: number;
  containerRef: RefObject<HTMLDivElement | null>;
  onDragEnd: (x: number, y: number) => void;
  onResizeEnd: (width: number, height: number) => void;
  onFocus: () => void;
  onHide: () => void;
  children: React.ReactNode;
}

const MIN_WIDTH = 260;
const MIN_HEIGHT = 120;

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
//
// Resize uses plain pointer events (not Motion's drag, which moves the
// whole element rather than growing/shrinking it) on a small corner handle.
// While dragging, size is tracked in local state for immediate visual
// feedback; the final size is only committed upward on pointer release.
export default function DashboardWidget({
  title,
  x,
  y,
  zIndex,
  width = 340,
  height,
  containerRef,
  onDragEnd,
  onResizeEnd,
  onFocus,
  onHide,
  children,
}: DashboardWidgetProps) {
  const dragControls = useDragControls();
  const motionX = useMotionValue(x);
  const motionY = useMotionValue(y);
  const bodyRef = useRef<HTMLDivElement>(null);
  const [liveSize, setLiveSize] = useState<{ width: number; height: number } | null>(null);

  useEffect(() => {
    motionX.set(x);
    motionY.set(y);
  }, [x, y, motionX, motionY]);

  function handleDragEnd() {
    onDragEnd(motionX.get(), motionY.get());
  }

  function handleResizeStart(e: React.PointerEvent) {
    e.stopPropagation();
    e.preventDefault();
    onFocus();

    const startX = e.clientX;
    const startY = e.clientY;
    const startWidth = width;
    const startHeight = height ?? bodyRef.current?.offsetHeight ?? MIN_HEIGHT;

    function handlePointerMove(moveEvent: PointerEvent) {
      const nextWidth = Math.max(MIN_WIDTH, startWidth + (moveEvent.clientX - startX));
      const nextHeight = Math.max(MIN_HEIGHT, startHeight + (moveEvent.clientY - startY));
      setLiveSize({ width: nextWidth, height: nextHeight });
    }

    function handlePointerUp() {
      window.removeEventListener("pointermove", handlePointerMove);
      window.removeEventListener("pointerup", handlePointerUp);
      setLiveSize((current) => {
        if (current) onResizeEnd(current.width, current.height);
        return null;
      });
    }

    window.addEventListener("pointermove", handlePointerMove);
    window.addEventListener("pointerup", handlePointerUp);
  }

  const displayWidth = liveSize?.width ?? width;
  const displayHeight = liveSize?.height ?? height;

  return (
    <motion.div
      drag
      dragListener={false}
      dragControls={dragControls}
      dragMomentum={false}
      dragConstraints={containerRef}
      onDragEnd={handleDragEnd}
      onPointerDownCapture={onFocus}
      style={{ position: "absolute", x: motionX, y: motionY, zIndex, width: displayWidth }}
      className="max-w-full bg-gray-800 border border-gray-700 rounded-xl shadow-lg overflow-hidden flex flex-col"
    >
      <div
        onPointerDown={(e) => dragControls.start(e)}
        className="flex items-center justify-between gap-2 px-3 py-2 bg-gray-900 border-b border-gray-700 cursor-grab active:cursor-grabbing select-none flex-shrink-0"
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
      <div
        ref={bodyRef}
        style={displayHeight ? { height: displayHeight, overflowY: "auto" } : undefined}
      >
        {children}
      </div>
      <div
        onPointerDown={handleResizeStart}
        className="absolute bottom-0 right-0 w-5 h-5 flex items-center justify-center text-gray-600 hover:text-gray-300 transition-colors"
        style={{ cursor: "nwse-resize" }}
        title="Größe ändern"
      >
        <Maximize2 className="w-3 h-3" />
      </div>
    </motion.div>
  );
}
