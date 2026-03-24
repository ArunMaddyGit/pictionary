import { MouseEvent, RefObject, useRef, useState } from "react";
import { DrawPayload } from "../../types/game";

type UseCanvasParams = {
  isDrawer: boolean;
  onDraw: (payload: DrawPayload) => void;
};

type CanvasHandlers = {
  onMouseDown: (event: MouseEvent<HTMLCanvasElement>) => void;
  onMouseMove: (event: MouseEvent<HTMLCanvasElement>) => void;
  onMouseUp: () => void;
  onMouseLeave: () => void;
};

type UseCanvasResult = {
  canvasRef: RefObject<HTMLCanvasElement>;
  color: string;
  setColor: (value: string) => void;
  brushSize: number;
  setBrushSize: (value: number) => void;
  drawStroke: (points: [number, number][], strokeColor: string, strokeSize: number) => void;
  clearCanvas: () => void;
  handlers: CanvasHandlers;
  isDrawing: boolean;
};

export function useCanvas({ isDrawer, onDraw }: UseCanvasParams): UseCanvasResult {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const [color, setColor] = useState("#000000");
  const [brushSize, setBrushSize] = useState(4);
  const [isDrawing, setIsDrawing] = useState(false);
  const lastPointRef = useRef<[number, number] | null>(null);
  const lastEmittedAtRef = useRef(0);
  const pendingPointsRef = useRef<[number, number][]>([]);

  const drawStroke = (points: [number, number][], strokeColor: string, strokeSize: number) => {
    const canvas = canvasRef.current;
    if (!canvas || points.length < 2) {
      return;
    }

    const ctx = canvas.getContext("2d");
    if (!ctx) {
      return;
    }

    ctx.beginPath();
    ctx.moveTo(points[0][0], points[0][1]);
    for (let i = 1; i < points.length; i += 1) {
      ctx.lineTo(points[i][0], points[i][1]);
    }
    ctx.strokeStyle = strokeColor;
    ctx.lineWidth = strokeSize;
    ctx.lineCap = "round";
    ctx.lineJoin = "round";
    ctx.stroke();
  };

  const clearCanvas = () => {
    const canvas = canvasRef.current;
    if (!canvas) {
      return;
    }
    const ctx = canvas.getContext("2d");
    if (!ctx) {
      return;
    }
    ctx.clearRect(0, 0, canvas.width, canvas.height);
  };

  const toCanvasPoint = (event: MouseEvent<HTMLCanvasElement>): [number, number] => {
    const rect = event.currentTarget.getBoundingClientRect();
    return [event.clientX - rect.left, event.clientY - rect.top];
  };

  const onMouseDown = (event: MouseEvent<HTMLCanvasElement>) => {
    if (!isDrawer) {
      return;
    }
    setIsDrawing(true);
    const start = toCanvasPoint(event);
    lastPointRef.current = start;
    pendingPointsRef.current = [start];
  };

  const onMouseMove = (event: MouseEvent<HTMLCanvasElement>) => {
    if (!isDrawer || !isDrawing || !lastPointRef.current) {
      return;
    }
    const point = toCanvasPoint(event);
    const points: [number, number][] = [lastPointRef.current, point];
    drawStroke(points, color, brushSize);
    lastPointRef.current = point;
    pendingPointsRef.current.push(point);

    const now = Date.now();
    if (now - lastEmittedAtRef.current >= 40) {
      lastEmittedAtRef.current = now;
      if (pendingPointsRef.current.length >= 2) {
        onDraw({ points: [...pendingPointsRef.current], color, brushSize });
        pendingPointsRef.current = [point];
      }
    }
  };

  const stopDrawing = () => {
    if (isDrawer && pendingPointsRef.current.length >= 2) {
      onDraw({ points: [...pendingPointsRef.current], color, brushSize });
    }
    pendingPointsRef.current = [];
    setIsDrawing(false);
    lastPointRef.current = null;
  };

  return {
    canvasRef,
    color,
    setColor,
    brushSize,
    setBrushSize,
    drawStroke,
    clearCanvas,
    handlers: {
      onMouseDown,
      onMouseMove,
      onMouseUp: stopDrawing,
      onMouseLeave: stopDrawing
    },
    isDrawing
  };
}
