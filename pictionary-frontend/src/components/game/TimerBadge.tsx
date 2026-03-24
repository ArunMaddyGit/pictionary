type TimerBadgeProps = {
  seconds: number;
};

export default function TimerBadge({ seconds }: TimerBadgeProps) {
  let colorClass = "timer-green";
  if (seconds < 15) {
    colorClass = "timer-red";
  } else if (seconds <= 30) {
    colorClass = "timer-orange";
  }
  const pulseClass = seconds < 10 ? "timer-pulse" : "";

  return <div className={`timer-badge ${colorClass} ${pulseClass}`}>{seconds}</div>;
}
