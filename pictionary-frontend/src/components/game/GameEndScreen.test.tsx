import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import GameEndScreen from "./GameEndScreen";

describe("GameEndScreen", () => {
  it("renders sorted leaderboard", () => {
    render(
      <GameEndScreen
        leaderboard={[
          { id: "1", name: "Alice", score: 20 },
          { id: "2", name: "Bob", score: 80 },
          { id: "3", name: "Cara", score: 50 }
        ]}
        onPlayAgain={vi.fn()}
      />
    );

    const rows = screen.getAllByRole("listitem");
    expect(rows[0]).toHaveTextContent("Bob");
    expect(rows[1]).toHaveTextContent("Cara");
    expect(rows[2]).toHaveTextContent("Alice");
  });

  it("Play Again calls callback", () => {
    const onPlayAgain = vi.fn();
    render(<GameEndScreen leaderboard={[]} onPlayAgain={onPlayAgain} />);
    fireEvent.click(screen.getByRole("button", { name: "Play Again" }));
    expect(onPlayAgain).toHaveBeenCalled();
  });
});
