type AdBannerProps = {
  size: "square" | "leaderboard";
};

export default function AdBanner({ size }: AdBannerProps) {
  return (
    <div className={`ad-banner ${size === "square" ? "ad-square" : "ad-leaderboard"}`}>
      ADS
    </div>
  );
}
