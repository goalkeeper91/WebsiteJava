import React, { useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { CheckCircle, AlertTriangle, Loader } from "lucide-react";

const DiscordCallback: React.FC = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const [status, setStatus] = useState<"loading" | "success" | "error">("loading");
  const [message, setMessage] = useState("");

  useEffect(() => {
    const handleCallback = async () => {
      const code = searchParams.get("code");
      const state = searchParams.get("state");
      const error = searchParams.get("error");

      // Check for OAuth error
      if (error) {
        setStatus("error");
        setMessage(`Discord OAuth error: ${error}`);
        setTimeout(() => navigate("/dashboard/discord"), 3000);
        return;
      }

      // Verify code exists
      if (!code) {
        setStatus("error");
        setMessage("Missing authorization code");
        setTimeout(() => navigate("/dashboard/discord"), 3000);
        return;
      }

      // Verify state matches (CSRF protection)
      const savedState = localStorage.getItem("discord_oauth_state");
      if (state !== savedState) {
        setStatus("error");
        setMessage("Invalid state parameter. Possible CSRF attack.");
        setTimeout(() => navigate("/dashboard/discord"), 3000);
        return;
      }

      // Clear saved state
      localStorage.removeItem("discord_oauth_state");

      try {
        // Send callback to backend
        const res = await fetch(
          `/api/discord/auth/callback?code=${encodeURIComponent(code || '')}&state=${encodeURIComponent(state)}`,
          {
            method: "GET",
            credentials: "include",
          }
        );

        if (res.ok) {
          setStatus("success");
          setMessage("Discord connected successfully!");
          setTimeout(() => navigate("/dashboard/discord"), 2000);
        } else {
          const errorText = await res.text();
          setStatus("error");
          setMessage(errorText || "Failed to connect Discord");
          setTimeout(() => navigate("/dashboard/discord"), 3000);
        }
      } catch (err) {
        console.error("Callback error:", err);
        setStatus("error");
        setMessage("Failed to connect Discord. Please try again.");
        setTimeout(() => navigate("/dashboard/discord"), 3000);
      }
    };

    handleCallback();
  }, [searchParams, navigate]);

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center p-6">
      <div className="max-w-md w-full p-8 bg-white rounded-xl shadow-lg text-center space-y-6">
        {status === "loading" && (
          <>
            <Loader className="w-16 h-16 mx-auto text-indigo-600 animate-spin" />
            <h2 className="text-2xl font-bold text-gray-900">Connecting Discord...</h2>
            <p className="text-gray-600">Please wait while we connect your account</p>
          </>
        )}

        {status === "success" && (
          <>
            <CheckCircle className="w-16 h-16 mx-auto text-green-600" />
            <h2 className="text-2xl font-bold text-gray-900">Success!</h2>
            <p className="text-gray-600">{message}</p>
            <p className="text-sm text-gray-500">Redirecting to dashboard...</p>
          </>
        )}

        {status === "error" && (
          <>
            <AlertTriangle className="w-16 h-16 mx-auto text-red-600" />
            <h2 className="text-2xl font-bold text-gray-900">Connection Failed</h2>
            <p className="text-gray-600">{message}</p>
            <p className="text-sm text-gray-500">Redirecting back...</p>
          </>
        )}
      </div>
    </div>
  );
};

export default DiscordCallback;