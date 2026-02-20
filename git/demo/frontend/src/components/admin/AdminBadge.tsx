import { Shield } from "lucide-react";

export default function AdminBadge() {
  return (
    <div className="inline-flex items-center gap-1 px-2 py-1 bg-gradient-to-r from-purple-600 to-pink-600 rounded text-xs font-bold">
      <Shield className="w-3 h-3" />
      ADMIN
    </div>
  );
}