// 魚のマッチングカードコンポーネント（A × B の形式で表示）
type Props = {
    left: string;
    right?: string;
    note?: string;
};

export default function MatchingCard({ left, right, note }: Props) {
    return (
        <div className="border rounded-md p-3 bg-white shadow-sm">
            <div className="font-semibold text-lg">{left} × {right ?? ""}</div>
            {note && <div className="text-sm text-gray-600 mt-1">{note}</div>}
        </div>
    );
}
