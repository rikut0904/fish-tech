// マッチング一覧コンポーネント（良い組み合わせ／悪い組み合わせ）
import MatchingCard from "./matching_card";

type Pair = { left: string; right: string; note?: string };

export default function MatchingList() {
    // サンプルデータ（必要に応じてAPIから取得するよう拡張可能）
    const good: Pair[] = [
        { left: "サバ", right: "大根", note: "脂とさっぱりの相性が良い" },
        { left: "ブリ", right: "ねぎ", note: "甘辛だれと相性が良い" },
    ];
    const bad: Pair[] = [
        { left: "タイ", right: "こってりソース", note: "淡白なので重い味付けは不向き" },
        { left: "タコ", right: "クリーム", note: "食感とソース感がぶつかる場合がある" },
    ];

    return (
        <section className="container mx-auto px-4 py-8">
            <h2 className="text-xl font-bold mb-4">魚のマッチング</h2>
            <div className="grid md:grid-cols-2 gap-6">
                <div>
                    <h3 className="text-lg font-semibold mb-3 text-green-600">良い組み合わせ</h3>
                    <div className="grid gap-3">
                        {good.map((g) => (
                            <MatchingCard key={`${g.left}-${g.right}`} left={g.left} right={g.right} note={g.note} />
                        ))}
                    </div>
                </div>
                <div>
                    <h3 className="text-lg font-semibold mb-3 text-red-600">悪い組み合わせ</h3>
                    <div className="grid gap-3">
                        {bad.map((b) => (
                            <MatchingCard key={`${b.left}-${b.right}`} left={b.left} right={b.right} note={b.note} />
                        ))}
                    </div>
                </div>
            </div>
        </section>
    );
}
