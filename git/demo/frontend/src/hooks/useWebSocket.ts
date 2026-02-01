import { useEffect, useRef, useState } from "react";

interface UseWebSocketOptions {
  url: string;
  onMessage: (data: any) => void;
  onOpen?: () => void;
  onClose?: () => void;
  onError?: (error: Event) => void;
  reconnectInterval?: number;
}

export function useWebSocket({
  url,
  onMessage,
  onOpen,
  onClose,
  onError,
  reconnectInterval = 3000,
}: UseWebSocketOptions) {
  const [isConnected, setIsConnected] = useState(false);
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<number | null>(null);

  useEffect(() => {
    let shouldReconnect = true;

    function connect() {
      try {
        const ws = new WebSocket(url);
        wsRef.current = ws;

        ws.onopen = () => {
          console.log("WebSocket verbunden:", url);
          setIsConnected(true);
          onOpen?.();
        };

        ws.onmessage = (event) => {
          try {
            const data = JSON.parse(event.data);
            onMessage(data);
          } catch (err) {
            console.error("Fehler beim Parsen der WebSocket-Nachricht:", err);
          }
        };

        ws.onclose = () => {
          console.log("WebSocket geschlossen");
          setIsConnected(false);
          onClose?.();

          // Automatisch reconnecten
          if (shouldReconnect) {
            reconnectTimeoutRef.current = setTimeout(() => {
              console.log("WebSocket reconnect...");
              connect();
            }, reconnectInterval);
          }
        };

        ws.onerror = (error) => {
          console.error("WebSocket Fehler:", error);
          onError?.(error);
        };
      } catch (err) {
        console.error("Fehler beim Verbinden:", err);

        // Retry
        if (shouldReconnect) {
          reconnectTimeoutRef.current = setTimeout(connect, reconnectInterval);
        }
      }
    }

    connect();

    // Cleanup
    return () => {
      shouldReconnect = false;

      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }

      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [url, reconnectInterval]);

  // Send message helper
  const send = (data: any) => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(typeof data === "string" ? data : JSON.stringify(data));
    } else {
      console.warn("WebSocket nicht verbunden");
    }
  };

  return { isConnected, send };
}