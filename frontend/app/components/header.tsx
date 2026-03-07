// ヘッダーコンポーネント
import Link from 'next/link';
export default function Header() {
    return (
        <header className="bg-white shadow-md py-4 px-6 flex items-center justify-between">
            <div className="flex items-center gap-2">
                <Link href="/">
                    <span className="text-2xl font-bold text-blue-700">FishTech</span>
                </Link>
            </div>
            <nav className="space-x-4">
                <Link href="/">
                    <span className="text-blue-700 hover:underline">ホーム</span>
                </Link>
                <Link href="/recipe">
                    <span className="text-blue-700 hover:underline">レシピ</span>
                </Link>
                <Link href="/contact">
                    <span className="text-blue-700 hover:underline">お問い合わせ</span>
                </Link>
                <Link href="/login">
                    <span className="text-blue-700 hover:underline">ログイン</span>
                </Link>
            </nav>
        </header>
    );
}