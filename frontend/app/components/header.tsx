// ヘッダーコンポーネント
export default function Header() {
    return (
        <header className="bg-white shadow-md py-4 px-6 flex items-center justify-between">
            <div className="flex items-center gap-2">
                <span className="text-2xl font-bold text-blue-700">FishTech</span>
            </div>
            <nav className="space-x-4">
                <a href="#features" className="text-blue-700 hover:underline">特徴</a>
                <a href="#fish" className="text-blue-700 hover:underline">魚図鑑</a>
                <a href="#contact" className="text-blue-700 hover:underline">お問い合わせ</a>
            </nav>
        </header>
    );
}
