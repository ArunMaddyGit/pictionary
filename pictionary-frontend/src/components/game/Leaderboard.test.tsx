import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import Leaderboard from "./Leaderboard";
import "../../styles/game.css";

describe("Leaderboard", () => {
  it("renders players sorted by score desc", () => {
    render(
      <Leaderboard
        players={[
          { id: "1", name: "Alice", score: 30, isDrawer: false, hasGuessed: false },
          { id: "2", name: "Bob", score: 90, isDrawer: false, hasGuessed: false },
          { id: "3", name: "Cara", score: 50, isDrawer: false, hasGuessed: false }
        ]}
        currentUserId="1"
      />
    );

    const rows = screen.getAllByRole("listitem");
    expect(rows[0]).toHaveTextContent("Bob");
    expect(rows[1]).toHaveTextContent("Cara");
    expect(rows[2]).toHaveTextContent("Alice");
  });

  it("highlights current user", () => {
    render(
      <Leaderboard
        players={[
          { id: "1", name: "Alice", score: 30, isDrawer: false, hasGuessed: false }
        ]}
        currentUserId="1"
      />
    );

    expect(screen.getByText("Alice").closest("li")).toHaveClass("current-user");
  });
});
