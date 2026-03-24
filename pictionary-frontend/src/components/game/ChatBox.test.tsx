import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import ChatBox from "./ChatBox";
import "../../styles/game.css";

describe("ChatBox", () => {
  it("renders messages and styles correct guesses in green class", () => {
    render(
      <ChatBox
        messages={[
          { playerName: "Alice", text: "banana", isCorrect: false, timestamp: "1" },
          { playerName: "Bob", text: "guessed correctly", isCorrect: true, timestamp: "2" }
        ]}
      />
    );

    expect(screen.getByText(/Alice:/)).toBeInTheDocument();
    const correctMessage = screen.getByText(/Bob:/).closest("p");
    expect(correctMessage).toHaveClass("correct");
  });
});
