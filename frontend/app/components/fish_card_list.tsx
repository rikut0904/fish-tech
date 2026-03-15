// 魚カードリスト
import FishCard from "@/app/components/fish_card";

const fishList = [
    { name: "カマス", desc: "淡白で上品な味わい。塩焼きや干物が人気。", image: "/recomend/kamasu.jpg" },
    { name: "アカガレイ", desc: "煮付けや唐揚げで美味しい白身魚。", image: "/recomend/aka-karei.jpg" },
    { name: "メギス", desc: "すり身や天ぷらで親しまれる地魚。", image: "/recomend/meigsu.jpg" },
];

export default function FishCardList() {
    return (
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
            {fishList.map((fish) => (
                <FishCard key={fish.name} {...fish} />
            ))}
        </div>
    );
}
