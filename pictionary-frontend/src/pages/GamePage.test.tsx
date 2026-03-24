import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi, beforeEach } from "vitest";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import GamePage from "./GamePage";
import { useGameState } from "../hooks/useGameState";

vi.mock("../hooks/useGameState", () => ({
  useGameState: vi.fn()
}));

const useGameStateMock = vi.mocked(useGameState);

const baseState: ReturnType<typeof useGameState> = {
  playerId: "p1",
  roomId: "r1",
  status: "connected",
  players: [{ id: "p1", name: "Alice", score: 10, isDrawer: false, hasGuessed: false }],
  phase: "drawing",
  timer: 50,
  maskedWord: "_ _ _ _",
  isDrawer: false,
  wordOptions: [],
  chatMessages: [],
  leaderboard: [],
  gameEnded: false,
  round: 1,
  disconnect: vi.fn(),
  sendDraw: vi.fn(),
  sendGuess: vi.fn(),
  sendSelectWord: vi.fn(),
  sendClearCanvas: vi.fn()
};

function renderWithRoutes() {
  return render(
    <MemoryRouter initialEntries={["/game"]}>
      <Routes>
        <Route path="/" element={<div>Landing</div>} />
        <Route path="/game" element={<GamePage />} />
      </Routes>
    </MemoryRouter>
  );
}

describe("GamePage", () => {
  beforeEach(() => {
    localStorage.clear();
    useGameStateMock.mockReturnValue(baseState);
  });

  it('redirects to "/" if wsUrl missing', () => {
    localStorage.setItem("playerId", "p1");
    localStorage.setItem("roomId", "r1");
    renderWithRoutes();
    expect(screen.getByText("Landing")).toBeInTheDocument();
  });

  it("renders GameLayout when wsUrl present", () => {
    localStorage.setItem("playerId", "p1");
    localStorage.setItem("roomId", "r1");
    localStorage.setItem("wsUrl", "ws://localhost:8080/ws?playerId=p1&roomId=r1");

    renderWithRoutes();
    expect(screen.getByText("Leaderboard")).toBeInTheDocument();
  });

  it('shows WordOptions for drawer in "choosing_word"', () => {
    useGameStateMock.mockReturnValue({
      ...baseState,
      phase: "choosing_word",
      isDrawer: true,
      wordOptions: ["apple", "dog", "chair"]
    } as ReturnType<typeof useGameState>);

    localStorage.setItem("playerId", "p1");
    localStorage.setItem("roomId", "r1");
    localStorage.setItem("wsUrl", "ws://localhost:8080/ws?playerId=p1&roomId=r1");

    renderWithRoutes();
    expect(screen.getByRole("dialog", { name: "Word options" })).toBeInTheDocument();
  });

  it("shows GameEndScreen when gameEnded is true", () => {
    useGameStateMock.mockReturnValue({
      ...baseState,
      gameEnded: true,
      leaderboard: [{ id: "p1", name: "Alice", score: 100 }]
    } as ReturnType<typeof useGameState>);

    localStorage.setItem("playerId", "p1");
    localStorage.setItem("roomId", "r1");
    localStorage.setItem("wsUrl", "ws://localhost:8080/ws?playerId=p1&roomId=r1");

    renderWithRoutes();
    expect(screen.getByText("Game Over!")).toBeInTheDocument();
  });

  it("shows connection lost overlay on disconnect", () => {
    useGameStateMock.mockReturnValueOnce({
      ...baseState,
      status: "disconnected"
    } as ReturnType<typeof useGameState>);

    localStorage.setItem("playerId", "p1");
    localStorage.setItem("roomId", "r1");
    localStorage.setItem("wsUrl", "ws://localhost:8080/ws?playerId=p1&roomId=r1");

    renderWithRoutes();
    expect(screen.queryByText("Connection lost. Please refresh.")).not.toBeInTheDocument();
  });

  it("shows connection lost overlay on error", () => {
    useGameStateMock.mockReturnValueOnce({
      ...baseState,
      status: "error"
    } as ReturnType<typeof useGameState>);

    localStorage.setItem("playerId", "p1");
    localStorage.setItem("roomId", "r1");
    localStorage.setItem("wsUrl", "ws://localhost:8080/ws?playerId=p1&roomId=r1");

    renderWithRoutes();
    expect(screen.getByText("Connection lost. Please refresh.")).toBeInTheDocument();
  });
});
