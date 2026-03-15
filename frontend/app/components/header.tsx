"use client";

// ヘッダーコンポーネント（モバイルでハンバーガーメニュー）
import Link from "next/link";
import { useState } from "react";

export default function Header() {
    const [menuOpen, setMenuOpen] = useState(false);

    return (
        <header className="bg-white shadow-md py-4 px-6 flex items-center justify-between relative">
            <div className="flex items-center gap-2">
                <Link href="/">
                    <span className="text-2xl md:text-3xl font-bold text-blue-700">FishTech</span>
                </Link>
            </div>

            {/* デスクトップ・タブレット: 横並びナビ */}
            <nav className="hidden md:flex space-x-4">
                <Link href="/">
                    <span className="text-base md:text-lg text-blue-700 hover:underline">ホーム</span>
                </Link>
                <Link href="/recipe">
                    <span className="text-base md:text-lg text-blue-700 hover:underline">レシピ</span>
                </Link>
                <Link href="/fish-book">
                    <span className="text-blue-700 hover:underline">魚図鑑</span>

                </Link>
                <Link href="/about">
                    <span className="text-base md:text-lg text-blue-700 hover:underline">このサイトについて</span>
                </Link>
                <Link href="/contact">
                    <span className="text-base md:text-lg text-blue-700 hover:underline">お問い合わせ</span>
                </Link>
                {/* <Link href="/login">
                    <span className="text-blue-700 hover:underline">ログイン</span>
                </Link> */}
            </nav>

            {/* モバイル: ハンバーガーボタン */}
            <button
                type="button"
                aria-expanded={menuOpen}
                aria-label={menuOpen ? "メニューを閉じる" : "メニューを開く"}
                className="md:hidden p-2 rounded text-blue-700 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
                onClick={() => setMenuOpen((prev) => !prev)}
            >
                <span className="material-icons text-4xl" aria-hidden>
                    {menuOpen ? "close" : "menu"}
                </span>
            </button>

            {/* モバイル: 開いたメニュー */}
            {menuOpen && (
                <>
                    <div
                        className="md:hidden fixed inset-0 bg-black/20 z-10 top-[57px]"
                        aria-hidden
                        onClick={() => setMenuOpen(false)}
                    />
                    <nav
                        className="md:hidden absolute top-full left-0 right-0 bg-white shadow-lg border-t z-20 py-2"
                        role="navigation"
                    >
                        <Link href="/" className="block px-6 py-3 text-base text-blue-700 hover:bg-gray-50" onClick={() => setMenuOpen(false)}>
                            ホーム
                        </Link>
                        <Link href="/recipe" className="block px-6 py-3 text-base text-blue-700 hover:bg-gray-50" onClick={() => setMenuOpen(false)}>
                            レシピ
                        </Link>
                        <Link href="/#encyclopedia" className="block px-6 py-3 text-base text-blue-700 hover:bg-gray-50" onClick={() => setMenuOpen(false)}>
                            魚図鑑
                        </Link>
                        <Link href="/about" className="block px-6 py-3 text-base text-blue-700 hover:bg-gray-50" onClick={() => setMenuOpen(false)}>
                            このサイトについて
                        </Link>
                        <Link href="/contact" className="block px-6 py-3 text-base text-blue-700 hover:bg-gray-50" onClick={() => setMenuOpen(false)}>
                            お問い合わせ
                        </Link>
                        {/* <Link href="/login">
                            <span className="text-blue-700 hover:underline">ログイン</span>
                        </Link> */}
                    </nav>
                </>
            )}
        </header>
    );
}