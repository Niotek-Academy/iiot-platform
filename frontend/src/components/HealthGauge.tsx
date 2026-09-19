"use client";

function colorFor(score: number): { ring: string; text: string } {
  if (score >= 70) return { ring: "#16a34a", text: "text-green-600" };
  if (score >= 40) return { ring: "#ca8a04", text: "text-yellow-600" };
  return { ring: "#dc2626", text: "text-red-600" };
}

export default function HealthGauge({ score, size = 52 }: { score: number; size?: number }) {
  const { ring, text } = colorFor(score);
  const clamped = Math.max(0, Math.min(100, score));

  return (
    <div
      className="rounded-full flex items-center justify-center shrink-0"
      style={{
        width: size,
        height: size,
        background: `conic-gradient(${ring} ${clamped * 3.6}deg, #e2e8f0 0deg)`,
      }}
    >
      <div
        className="rounded-full bg-white flex items-center justify-center font-medium"
        style={{ width: size - 14, height: size - 14, fontSize: size < 60 ? 11 : 13 }}
      >
        <span className={text}>{Math.round(clamped)}%</span>
      </div>
    </div>
  );
}