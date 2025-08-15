import React from "react";

interface UptimeDisplayProps {
  uptimeSeconds: number;
  maxSeconds?: number; // Optional: z.B. 86400 für 24h
}

const UptimeDisplay: React.FC<UptimeDisplayProps> = ({ uptimeSeconds, maxSeconds = 86400 }) => {
  const hours = Math.floor(uptimeSeconds / 3600);
  const minutes = Math.floor((uptimeSeconds % 3600) / 60);
  const seconds = uptimeSeconds % 60;

  const percent = Math.min(100, (uptimeSeconds / maxSeconds) * 100);
  const radius = 40;
  const circumference = 2 * Math.PI * radius;
  const offset = circumference - (percent / 100) * circumference;

  return (
    <div className="flex flex-col items-center mt-4">
      <svg width="100" height="100" className="transform -rotate-90">
        <circle
          cx="50"
          cy="50"
          r={radius}
          stroke="#e5e7eb"
          strokeWidth="8"
          fill="transparent"
        />
        <circle
          cx="50"
          cy="50"
          r={radius}
          stroke="#10b981"
          strokeWidth="8"
          fill="transparent"
          strokeDasharray={circumference}
          strokeDashoffset={offset}
          strokeLinecap="round"
        />
      </svg>
      <div className="mt-2 text-center">
        <p className="text-sm text-gray-500">Uptime</p>
        <p className="font-mono text-lg">
          {String(hours).padStart(2, "0")}:
          {String(minutes).padStart(2, "0")}:
          {String(seconds).padStart(2, "0")}
        </p>
      </div>
    </div>
  );
};

export default UptimeDisplay;
