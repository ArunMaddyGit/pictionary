import { forwardRef, useEffect, useImperativeHandle } from "react";
import { DrawPayload } from "../../types/game";
import { useCanvas } from "./useCanvas";

export type DrawingCanvasHandle = {
  drawStroke: (points: [number, number][], color: string, brushSize: number) => void;
  clearCanvas: () => void;
  setColor: (value: string) => void;
  setBrushSize: (value: number) => void;
};

type DrawingCanvasProps = {
  isDrawer: boolean;
  onDraw: (payload: DrawPayload) => void;
  width?: number;
  height?: number;
  onCanvasStateChange?: (state: { color: string; brushSize: number }) => void;
};

const DrawingCanvas = forwardRef<DrawingCanvasHandle, DrawingCanvasProps>(function DrawingCanvas(
  { isDrawer, onDraw, width = 800, height = 500, onCanvasStateChange },
  ref
) {
  const { canvasRef, color, setColor, brushSize, setBrushSize, drawStroke, clearCanvas, handlers } = useCanvas({
    isDrawer,
    onDraw
  });

  useImperativeHandle(
    ref,
    () => ({
      drawStroke,
      clearCanvas,
      setColor,
      setBrushSize
    }),
    [drawStroke, clearCanvas]
  );

  useEffect(() => {
    onCanvasStateChange?.({ color, brushSize });
  }, [color, brushSize, onCanvasStateChange]);

  return (
    <canvas
      id="canvas-container"
      ref={canvasRef}
      width={width}
      height={height}
      style={{ cursor: isDrawer ? "crosshair" : "default", background: "#ffffff", width: "100%", height: "100%" }}
      onMouseDown={isDrawer ? handlers.onMouseDown : undefined}
      onMouseMove={isDrawer ? handlers.onMouseMove : undefined}
      onMouseUp={isDrawer ? handlers.onMouseUp : undefined}
      onMouseLeave={isDrawer ? handlers.onMouseLeave : undefined}
    />
  );
});

export default DrawingCanvas;
