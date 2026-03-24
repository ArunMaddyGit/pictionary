import { useCallback, useEffect, useRef, useState } from "react";
import { Message } from "../types/game";

type UseWebSocketParams = {
  url: string | null;
  onMessage: (msg: Message<unknown>) => void;
};

type SocketStatus = "connecting" | "connected" | "disconnected" | "error";

export function useWebSocket({ url, onMessage }: UseWebSocketParams) {
  const socketRef = useRef<WebSocket | null>(null);
  const onMessageRef = useRef(onMessage);
  const [status, setStatus] = useState<SocketStatus>("disconnected");

  useEffect(() => {
    onMessageRef.current = onMessage;
  }, [onMessage]);

  const disconnect = useCallback(() => {
    const socket = socketRef.current;
    if (socket) {
      socket.close();
      socketRef.current = null;
    }
    setStatus("disconnected");
  }, []);

  useEffect(() => {
    if (!url) {
      disconnect();
      return;
    }

    if (socketRef.current) {
      socketRef.current.close();
      socketRef.current = null;
    }

    setStatus("connecting");
    const socket = new WebSocket(url);
    socketRef.current = socket;

    socket.onopen = () => {
      setStatus("connected");
    };

    socket.onmessage = (event) => {
      try {
        const parsed = JSON.parse(String(event.data)) as Message<unknown>;
        onMessageRef.current(parsed);
      } catch {
        // Ignore invalid frames for MVP.
      }
    };

    socket.onerror = () => {
      setStatus("error");
    };

    socket.onclose = () => {
      if (socketRef.current === socket) {
        socketRef.current = null;
      }
      setStatus("disconnected");
    };

    return () => {
      socket.close();
      if (socketRef.current === socket) {
        socketRef.current = null;
      }
    };
  }, [url, disconnect]);

  const sendMessage = useCallback((type: string, payload: unknown) => {
    const socket = socketRef.current;
    if (!socket || socket.readyState !== WebSocket.OPEN) {
      return;
    }
    socket.send(
      JSON.stringify({
        type,
        payload
      })
    );
  }, []);

  return {
    status,
    sendMessage,
    disconnect
  };
}
