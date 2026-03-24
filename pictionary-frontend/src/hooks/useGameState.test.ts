import { act, renderHook } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { useGameState } from "./useGameState";

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

  send() {}
}

describe("useGameState", () => {
  beforeEach(() => {
    MockWebSocket.instances = [];
    vi.stubGlobal("WebSocket", MockWebSocket as unknown as typeof WebSocket);
  });

  afterEach(() => {
    vi.unstubAllGlobals();
    vi.useRealTimers();
  });

  it("ROOM_STATE updates players list", () => {
    const { result } = renderHook(() =>
      useGameState({ playerId: "p1", roomId: "r1", wsUrl: "ws://localhost:8080/ws" })
    );
    const socket = MockWebSocket.instances[0];
    act(() => {
      socket.receive({
        type: "ROOM_STATE",
        payload: {
          players: [{ id: "p1", name: "Alice", score: 10, isDrawer: false, hasGuessed: false }],
          round: 1,
          timer: 60
        }
      });
    });
    expect(result.current.players).toHaveLength(1);
    expect(result.current.players[0].name).toBe("Alice");
  });

  it("TURN_START sets isDrawer correctly", () => {
    const { result } = renderHook(() =>
      useGameState({ playerId: "p1", roomId: "r1", wsUrl: "ws://localhost:8080/ws" })
    );
    const socket = MockWebSocket.instances[0];
    act(() => {
      socket.receive({
        type: "TURN_START",
        payload: {
          drawerId: "p1",
          maskedWord: "a p p l e",
          isDrawer: true
        }
      });
    });
    expect(result.current.isDrawer).toBe(true);
    expect(result.current.maskedWord).toBe("a p p l e");
  });

  it("GUESS_MESSAGE appends chat messages", () => {
    const { result } = renderHook(() =>
      useGameState({ playerId: "p1", roomId: "r1", wsUrl: "ws://localhost:8080/ws" })
    );
    const socket = MockWebSocket.instances[0];
    act(() => {
      socket.receive({
        type: "GUESS_MESSAGE",
        payload: {
          playerId: "p2",
          playerName: "Bob",
          text: "cat"
        }
      });
    });
    expect(result.current.chatMessages).toHaveLength(1);
    expect(result.current.chatMessages[0].playerName).toBe("Bob");
  });

  it("WORD_OPTIONS sets chooser phase and drawer state", () => {
    const { result } = renderHook(() =>
      useGameState({ playerId: "p1", roomId: "r1", wsUrl: "ws://localhost:8080/ws" })
    );
    const socket = MockWebSocket.instances[0];
    act(() => {
      socket.receive({
        type: "WORD_OPTIONS",
        payload: {
          words: ["apple", "cat", "train"]
        }
      });
    });
    expect(result.current.phase).toBe("choosing_word");
    expect(result.current.isDrawer).toBe(true);
    expect(result.current.wordOptions).toHaveLength(3);
  });

  it("ROUND_END clears canvas via callback", () => {
    const onClearCanvas = vi.fn();
    renderHook(() =>
      useGameState({
        playerId: "p1",
        roomId: "r1",
        wsUrl: "ws://localhost:8080/ws",
        onClearCanvas
      })
    );
    const socket = MockWebSocket.instances[0];
    act(() => {
      socket.receive({
        type: "ROUND_END",
        payload: { word: "apple" }
      });
    });
    expect(onClearCanvas).toHaveBeenCalledTimes(1);
  });

  it("timer counts down in drawing phase", () => {
    vi.useFakeTimers();
    const { result } = renderHook(() =>
      useGameState({ playerId: "p1", roomId: "r1", wsUrl: "ws://localhost:8080/ws" })
    );
    const socket = MockWebSocket.instances[0];
    act(() => {
      socket.receive({
        type: "ROOM_STATE",
        payload: { players: [], round: 1, timer: 5 }
      });
      socket.receive({
        type: "TURN_START",
        payload: { drawerId: "p2", maskedWord: "_ _ _", isDrawer: false }
      });
    });
    expect(result.current.phase).toBe("drawing");

    act(() => {
      vi.advanceTimersByTime(1000);
    });
    expect(result.current.timer).toBe(4);
  });
});
