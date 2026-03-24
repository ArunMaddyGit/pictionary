import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import WordDisplay from "./WordDisplay";
import "../../styles/game.css";

describe("WordDisplay", () => {
  it("drawer sees actual word", () => {
    render(<WordDisplay word="hot dog" isDrawer={true} />);
    expect(screen.getByText("hot dog")).toBeInTheDocument();
  });

  it("non-drawer sees masked word", () => {
    render(<WordDisplay word="hot dog" isDrawer={false} />);
    expect(screen.getByText((text) => text.replace(/\s+/g, " ").trim() === "_ _ _ _ _ _")).toBeInTheDocument();
  });
});
