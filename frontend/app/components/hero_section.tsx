// メインビジュアル・キャッチコピー
export default function HeroSection() {
    return (
        <section id="hero" className="bg-gradient-to-r from-blue-400 to-blue-200 py-16 text-center">
            <h1 className="text-4xl font-bold text-white mb-4 drop-shadow">知られざる魚の魅力を発見しよう</h1>
            <p className="text-lg text-white mb-6">金沢の“有名じゃない”魚たちをもっと知って、もっと味わおう。</p>
            <img src="/fish-hero.png" alt="魚のイメージ" className="mx-auto w-64 h-40 object-contain" />
        </section>
    );
}
