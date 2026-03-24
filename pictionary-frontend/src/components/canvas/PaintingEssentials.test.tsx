import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import PaintingEssentials from "./PaintingEssentials";
import "../../styles/game.css";

describe("PaintingEssentials", () => {
  it("does not render when isDrawer is false", () => {
    const { container } = render(
      <PaintingEssentials
        isDrawer={false}
        color="#000000"
        onColorChange={vi.fn()}
        brushSize={4}
        onBrushSizeChange={vi.fn()}
        onClear={vi.fn()}
      />
    );
    expect(container).toBeEmptyDOMElement();
  });

  it("calls onColorChange when swatch is clicked", () => {
    const onColorChange = vi.fn();
    render(
      <PaintingEssentials
        isDrawer={true}
        color="#000000"
        onColorChange={onColorChange}
        brushSize={4}
        onBrushSizeChange={vi.fn()}
        onClear={vi.fn()}
      />
    );
    fireEvent.click(screen.getByLabelText("Color #ff0000"));
    expect(onColorChange).toHaveBeenCalledWith("#ff0000");
  });

  it("calls onBrushSizeChange when brush button clicked", () => {
    const onBrush = vi.fn();
    render(
      <PaintingEssentials
        isDrawer={true}
        color="#000000"
        onColorChange={vi.fn()}
        brushSize={4}
        onBrushSizeChange={onBrush}
        onClear={vi.fn()}
      />
    );
    fireEvent.click(screen.getByRole("button", { name: "L" }));
    expect(onBrush).toHaveBeenCalledWith(14);
  });

  it("calls onClear when clear button clicked", () => {
    const onClear = vi.fn();
    render(
      <PaintingEssentials
        isDrawer={true}
        color="#000000"
        onColorChange={vi.fn()}
        brushSize={4}
        onBrushSizeChange={vi.fn()}
        onClear={onClear}
      />
    );
    fireEvent.click(screen.getByRole("button", { name: "Clear" }));
    expect(onClear).toHaveBeenCalled();
  });
});
