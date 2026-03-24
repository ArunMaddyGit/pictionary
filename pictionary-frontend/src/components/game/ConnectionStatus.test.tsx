import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import ConnectionStatus from "./ConnectionStatus";
import "../../styles/game.css";

describe("ConnectionStatus", () => {
  it("shows connected state with green class", () => {
    render(<ConnectionStatus status="connected" />);
    expect(screen.getByText("Connected")).toBeInTheDocument();
    expect(screen.getByText("Connected").closest(".connection-status")).toHaveClass("status-connected");
  });

  it("shows connecting state", () => {
    render(<ConnectionStatus status="connecting" />);
    expect(screen.getByText("Connecting")).toBeInTheDocument();
  });

  it("shows disconnected state", () => {
    render(<ConnectionStatus status="disconnected" />);
    expect(screen.getByText("Disconnected")).toBeInTheDocument();
  });

  it("shows error state", () => {
    render(<ConnectionStatus status="error" />);
    expect(screen.getByText("Error")).toBeInTheDocument();
  });
});
