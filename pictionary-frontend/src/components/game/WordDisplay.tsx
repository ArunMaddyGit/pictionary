type WordDisplayProps = {
  word: string;
  isDrawer: boolean;
};

function maskWord(word: string): string {
  return word
    .split("")
    .map((ch) => (ch === " " ? " " : "_"))
    .join(" ");
}

export default function WordDisplay({ word, isDrawer }: WordDisplayProps) {
  const visible = isDrawer ? word : maskWord(word);
  return <div className="word-display">{visible}</div>;
}
