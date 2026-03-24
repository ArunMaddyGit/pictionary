type PaintingEssentialsProps = {
  isDrawer: boolean;
  color: string;
  onColorChange: (color: string) => void;
  brushSize: number;
  onBrushSizeChange: (size: number) => void;
  onClear: () => void;
};

const PRESET_COLORS = [
  "#000000",
  "#ffffff",
  "#ff0000",
  "#ff9800",
  "#ffeb3b",
  "#4caf50",
  "#00bcd4",
  "#2196f3",
  "#3f51b5",
  "#9c27b0",
  "#e91e63",
  "#795548"
];

const BRUSHES = [
  { label: "S", value: 2 },
  { label: "M", value: 6 },
  { label: "L", value: 14 }
];

export default function PaintingEssentials({
  isDrawer,
  color,
  onColorChange,
  brushSize,
  onBrushSizeChange,
  onClear
}: PaintingEssentialsProps) {
  if (!isDrawer) {
    return null;
  }

  return (
    <section id="painting-essentials-container" className="painting-essentials">
      <h4>PAINTING ESSENTIALS</h4>
      <div className="palette-row">
        {PRESET_COLORS.map((swatch) => (
          <button
            key={swatch}
            type="button"
            aria-label={`Color ${swatch}`}
            className={`color-swatch ${color === swatch ? "active" : ""}`}
            style={{ backgroundColor: swatch }}
            onClick={() => onColorChange(swatch)}
          />
        ))}
        <input type="color" value={color} onChange={(e) => onColorChange(e.target.value)} />
      </div>

      <div className="brush-row">
        {BRUSHES.map((b) => (
          <button
            key={b.label}
            type="button"
            className={`brush-btn ${brushSize === b.value ? "active" : ""}`}
            onClick={() => onBrushSizeChange(b.value)}
          >
            {b.label}
          </button>
        ))}
        <button type="button" className="clear-btn" onClick={onClear}>
          Clear
        </button>
      </div>
    </section>
  );
}
