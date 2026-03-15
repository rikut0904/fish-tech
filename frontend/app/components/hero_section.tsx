import Image from "next/image";
// メインビジュアル・キャッチコピー
export default function HeroSection() {
    return (
        <section id="hero" className="bg-gradient-to-r from-blue-400 to-blue-200 py-16 text-center">
            <h1 className="text-4xl font-bold text-white mb-4 drop-shadow">知られざる魚の魅力を発見しよう</h1>
            <p className="text-lg text-white mb-6">金沢の“有名じゃない”魚たちをもっと知って、もっと味わおう。</p>
                <Image src="/hero/fish_01.jpg" alt="魚のイメージ" width={256} height={160} className="mx-auto" />
        </section>
    );
}
