import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import TimerBadge from "./TimerBadge";
import "../../styles/game.css";

describe("TimerBadge", () => {
  it("shows the seconds number", () => {
    render(<TimerBadge seconds={42} />);
    expect(screen.getByText("42")).toBeInTheDocument();
  });

  it("uses green class when > 30", () => {
    render(<TimerBadge seconds={31} />);
    expect(screen.getByText("31")).toHaveClass("timer-green");
  });

  it("uses orange class when 15-30", () => {
    render(<TimerBadge seconds={20} />);
    expect(screen.getByText("20")).toHaveClass("timer-orange");
  });

  it("uses red class when < 15", () => {
    render(<TimerBadge seconds={14} />);
    expect(screen.getByText("14")).toHaveClass("timer-red");
  });
});
