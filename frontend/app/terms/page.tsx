import React from 'react';
import Link from 'next/link'

export default function TermsOfService() {
    const lastUpdated = "2026-03-11";

    return (
        // 1. 背景はライトブルー（#f0f7ff）
        <div className="min-h-screen bg-[#f0f7ff] py-12 px-4">

            {/* 2. 白い大きなカード（ボックス）で囲む */}
            <main className="max-w-4xl mx-auto bg-white rounded-3xl p-8 md:p-12 shadow-sm border border-blue-50">

                <div className="space-y-8 text-gray-700 leading-relaxed">
                    {/* 見出し部分 */}
                    <div>
                        <h1 className="text-3xl font-bold text-blue-700 mb-2">利用規約</h1>
                        <p className="text-sm text-gray-400">最終更新日: {lastUpdated}</p>
                    </div>

                    <section>
                        <h2 className="text-xl font-bold text-black mb-4 flex items-center">
                            1. はじめに
                        </h2>
                        <p>
                            Fish-Tech（以下「当アプリ」）は、ユーザーに魚の魅力や知識、関連する店舗情報を提供するためのサービスです。本規約は、利用者の皆様に本サービスを安心・安全にご利用いただくためのルールを定めたものです。
                        </p>
                    </section>

                    <section>
                        <h2 className="text-xl font-bold text-black mb-4">
                            2. サービスの提供内容
                        </h2>
                        <ul className="list-disc list-inside space-y-2 ml-2">
                            <li>魚の種類、生態、旬などの図鑑情報の提供</li>
                            <li>おすすめの魚料理レシピの掲載</li>
                            <li>魚料理を提供する店舗情報の提供および検索機能</li>
                        </ul>
                    </section>

                    <section>
                        <h2 className="text-xl font-bold text-black mb-4">
                            3. 免責事項
                        </h2>
                        <p className="mb-4 text-sm">
                            当アプリで提供される情報は細心の注意を払っておりますが、以下の点について保証するものではありません。
                        </p>
                        {/* 免責事項内のさらに強調するボックス */}
                        <div className="bg-blue-50/50 p-6 rounded-2xl border border-blue-100 text-sm space-y-3">
                            <p><strong>● 情報の正確性：</strong> 魚の生態や分布、旬の時期は環境や地域によって変動するため、常に最新かつ正確であることを保証しません。</p>
                            <p><strong>● 店舗情報：</strong> 掲載されている店舗の営業時間、メニュー、価格等は変更されている場合があります。訪問前に直接店舗へ確認することをお勧めします。</p>
                            <p><strong>● 取引の責任：</strong> 当アプリを通じて知った店舗での飲食や取引に関するトラブルについて、当アプリは一切の責任を負いません。</p>
                        </div>
                    </section>

                    <section>
                        <h2 className="text-xl font-bold text-black mb-4">
                            4. 禁止事項
                        </h2>
                        <p className="mb-2">利用者は、本サービスの利用にあたり以下の行為を行ってはなりません。</p>
                        <ul className="list-disc list-inside space-y-1 ml-2 text-sm">
                            <li>当アプリ内の画像、テキストの無断転載・再配布</li>
                            <li>他の利用者または第三者に不利益を与える行為</li>
                            <li>公序良俗に反する行為</li>
                        </ul>
                    </section>

                    <section className="pt-8 border-t border-gray-100">
                        <h2 className="text-xl font-bold text-black mb-4">
                            5. お問い合わせ
                        </h2>
                        <a>
                            本規約に関するお問い合わせは、
                        </a>
                        <Link href="/contact" className="underline">
                            お問い合わせフォーム
                        </Link>
                        <a>
                            よりご連絡ください
                        </a>
                    </section>
                </div>

                <footer className="mt-16 text-center text-gray-400 text-xs">
                    © 2026 FishTech / HOSA. All rights reserved.
                </footer>

            </main> {/* カード終了 */}
        </div> // 背景終了
    );
}