import { RefObject, useState } from "react";
import { DrawPayload, GamePhase } from "../../types/game";
import AdBanner from "../game/AdBanner";
import TimerBadge from "../game/TimerBadge";
import WordDisplay from "../game/WordDisplay";
import DrawingCanvas, { DrawingCanvasHandle } from "../canvas/DrawingCanvas";
import PaintingEssentials from "../canvas/PaintingEssentials";

type CenterPanelProps = {
  phase: GamePhase;
  timer: number;
  maskedWord: string;
  isDrawer: boolean;
  canvasRef: RefObject<DrawingCanvasHandle>;
  onDraw: (payload: DrawPayload) => void;
  onClearCanvas: () => void;
};

export default function CenterPanel({ phase, timer, maskedWord, isDrawer, canvasRef, onDraw, onClearCanvas }: CenterPanelProps) {
  const [color, setColor] = useState("#000000");
  const [brushSize, setBrushSize] = useState(4);

  return (
    <main className="center-panel" data-phase={phase}>
      <div className="center-top">
        <TimerBadge seconds={timer} />
        <WordDisplay word={maskedWord} isDrawer={isDrawer} />
      </div>
      <DrawingCanvas
        ref={canvasRef}
        isDrawer={isDrawer}
        onDraw={onDraw}
        onCanvasStateChange={({ color: nextColor, brushSize: nextBrushSize }) => {
          setColor(nextColor);
          setBrushSize(nextBrushSize);
        }}
      />
      <PaintingEssentials
        isDrawer={isDrawer}
        color={color}
        onColorChange={(nextColor) => {
          setColor(nextColor);
          canvasRef.current?.setColor(nextColor);
        }}
        brushSize={brushSize}
        onBrushSizeChange={(nextBrush) => {
          setBrushSize(nextBrush);
          canvasRef.current?.setBrushSize(nextBrush);
        }}
        onClear={() => {
          canvasRef.current?.clearCanvas();
          onClearCanvas();
        }}
      />
      <AdBanner size="leaderboard" />
    </main>
  );
}
