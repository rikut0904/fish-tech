// 装飾用の魚背景コンポーネント
export default function FishBackground() {
  return (
    <div className="pointer-events-none absolute inset-0 -z-10 overflow-hidden">
      <svg
        className="w-full h-full opacity-10 dark:opacity-20"
        viewBox="0 0 1200 800"
        preserveAspectRatio="xMidYMid slice"
        xmlns="http://www.w3.org/2000/svg"
        aria-hidden="true"
      >
        <defs>
          <linearGradient id="g1" x1="0" x2="1">
            <stop offset="0%" stopColor="#7DD3FC" />
            <stop offset="100%" stopColor="#60A5FA" />
          </linearGradient>
        </defs>
        <g transform="translate(0,50)">
          <path
            d="M50 300 C150 200, 250 200, 350 300 C450 400, 550 400, 650 300 C750 200, 850 200, 950 300 L950 450 L50 450 Z"
            fill="url(#g1)"
          />
          <g transform="translate(200,120) scale(0.9)">
            <ellipse cx="180" cy="220" rx="90" ry="45" fill="#fff" opacity="0.06" />
            <path d="M120 200 C140 160, 220 160, 260 200 C300 240, 260 280, 200 270 C160 260, 140 240, 120 200 Z" fill="#fff" opacity="0.08"/>
          </g>
          <g transform="translate(600,40) scale(1.1)">
            <path d="M80 120 C110 80, 200 80, 230 120 C260 160, 230 200, 180 190 C140 180, 120 160, 80 120 Z" fill="#fff" opacity="0.06"/>
          </g>
        </g>
      </svg>
    </div>
  );
}
