import { act, renderHook } from "@testing-library/react";
import { MouseEvent } from "react";
import { describe, expect, it, vi } from "vitest";
import { useCanvas } from "./useCanvas";

function mouseEvent(x: number, y: number, target: HTMLCanvasElement) {
  return {
    clientX: x,
    clientY: y,
    currentTarget: target
  } as unknown as MouseEvent<HTMLCanvasElement>;
}

describe("useCanvas", () => {
  it("non-drawer does not emit onDraw", () => {
    const onDraw = vi.fn();
    const { result } = renderHook(() => useCanvas({ isDrawer: false, onDraw }));
    const canvas = document.createElement("canvas");
    canvas.getBoundingClientRect = () => ({ left: 0, top: 0, width: 100, height: 100 } as DOMRect);

    act(() => {
      result.current.handlers.onMouseDown(mouseEvent(10, 10, canvas));
    });
    act(() => {
      result.current.handlers.onMouseMove(mouseEvent(20, 20, canvas));
    });

    expect(onDraw).not.toHaveBeenCalled();
  });

  it("drawer mouse down + move emits onDraw payload", () => {
    const onDraw = vi.fn();
    const { result, rerender } = renderHook(() => useCanvas({ isDrawer: true, onDraw }));
    const canvas = document.createElement("canvas");
    canvas.getBoundingClientRect = () => ({ left: 0, top: 0, width: 100, height: 100 } as DOMRect);
    const ctx = {
      beginPath: vi.fn(),
      moveTo: vi.fn(),
      lineTo: vi.fn(),
      stroke: vi.fn(),
      strokeStyle: "",
      lineWidth: 1,
      lineCap: "round",
      lineJoin: "round"
    };
    canvas.getContext = vi.fn().mockReturnValue(ctx);
    (result.current.canvasRef as unknown as { current: HTMLCanvasElement | null }).current = canvas;

    act(() => {
      result.current.handlers.onMouseDown(mouseEvent(10, 10, canvas));
    });
    rerender();
    act(() => {
      result.current.handlers.onMouseMove(mouseEvent(20, 20, canvas));
    });

    expect(onDraw).toHaveBeenCalledTimes(1);
    expect(onDraw).toHaveBeenCalledWith({
      points: [
        [10, 10],
        [20, 20]
      ],
      color: "#000000",
      brushSize: 4
    });
  });

  it("drawStroke draws on canvas context", () => {
    const onDraw = vi.fn();
    const { result } = renderHook(() => useCanvas({ isDrawer: true, onDraw }));
    const canvas = document.createElement("canvas");
    const ctx = {
      beginPath: vi.fn(),
      moveTo: vi.fn(),
      lineTo: vi.fn(),
      stroke: vi.fn(),
      strokeStyle: "",
      lineWidth: 1,
      lineCap: "round",
      lineJoin: "round"
    };
    canvas.getContext = vi.fn().mockReturnValue(ctx);
    (result.current.canvasRef as unknown as { current: HTMLCanvasElement | null }).current = canvas;

    act(() => {
      result.current.drawStroke(
        [
          [0, 0],
          [10, 10]
        ],
        "#ff0000",
        6
      );
    });

    expect(ctx.beginPath).toHaveBeenCalled();
    expect(ctx.moveTo).toHaveBeenCalledWith(0, 0);
    expect(ctx.lineTo).toHaveBeenCalledWith(10, 10);
    expect(ctx.stroke).toHaveBeenCalled();
  });
});
