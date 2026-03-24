import { act, renderHook } from "@testing-library/react";
import { describe, expect, it, vi, beforeEach, afterEach } from "vitest";
import { useWebSocket } from "./useWebSocket";

class MockWebSocket {
  static instances: MockWebSocket[] = [];
  static OPEN = 1;
  static CONNECTING = 0;
  static CLOSED = 3;

  url: string;
  readyState = MockWebSocket.CONNECTING;
  onopen: (() => void) | null = null;
  onmessage: ((event: { data: string }) => void) | null = null;
  onclose: (() => void) | null = null;
  onerror: (() => void) | null = null;
  sent: string[] = [];

  constructor(url: string) {
    this.url = url;
    MockWebSocket.instances.push(this);
  }

  open() {
    this.readyState = MockWebSocket.OPEN;
    this.onopen?.();
  }

  receive(data: unknown) {
    this.onmessage?.({ data: JSON.stringify(data) });
  }

  close() {
    this.readyState = MockWebSocket.CLOSED;
    this.onclose?.();
  }

  send(data: string) {
    this.sent.push(data);
  }
}

describe("useWebSocket", () => {
  beforeEach(() => {
    MockWebSocket.instances = [];
    vi.stubGlobal("WebSocket", MockWebSocket as unknown as typeof WebSocket);
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it("connects when url provided", () => {
    const onMessage = vi.fn();
    renderHook(() =>
      useWebSocket({
        url: "ws://localhost:8080/ws",
        onMessage
      })
    );
    expect(MockWebSocket.instances.length).toBe(1);
    expect(MockWebSocket.instances[0].url).toBe("ws://localhost:8080/ws");
  });

  it("calls onMessage when message received", () => {
    const onMessage = vi.fn();
    renderHook(() =>
      useWebSocket({
        url: "ws://localhost:8080/ws",
        onMessage
      })
    );
    const socket = MockWebSocket.instances[0];
    act(() => {
      socket.receive({ type: "ROOM_STATE", payload: { players: [] } });
    });
    expect(onMessage).toHaveBeenCalledWith({ type: "ROOM_STATE", payload: { players: [] } });
  });

  it("sendMessage serializes to JSON", () => {
    const onMessage = vi.fn();
    const { result } = renderHook(() =>
      useWebSocket({
        url: "ws://localhost:8080/ws",
        onMessage
      })
    );
    const socket = MockWebSocket.instances[MockWebSocket.instances.length - 1];
    act(() => {
      socket.open();
    });
    act(() => {
      result.current.sendMessage("GUESS", { text: "cat" });
    });
    expect(socket.sent[0]).toBe(JSON.stringify({ type: "GUESS", payload: { text: "cat" } }));
  });

  it("status becomes disconnected on close", () => {
    const onMessage = vi.fn();
    const { result } = renderHook(() =>
      useWebSocket({
        url: "ws://localhost:8080/ws",
        onMessage
      })
    );
    const socket = MockWebSocket.instances[MockWebSocket.instances.length - 1];
    act(() => {
      socket.open();
    });
    expect(result.current.status).toBe("connected");
    act(() => {
      socket.close();
    });
    expect(result.current.status).toBe("disconnected");
  });
});
