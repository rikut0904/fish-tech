import React from 'react';
import Link from 'next/link';
import { ChevronLeft } from 'lucide-react';
import { allFishData } from '../../data';

export default function FishDetail1() {
    const fish = allFishData.find(f => f.id === 1);

    if (!fish) {
        return (
            <div className="p-6">
                <p>魚が見つかりません。</p>
                <Link href="/fish-book" className="text-blue-600 underline">一覧に戻る</Link>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-white">
            <header className="p-4 flex items-center border-b">
                <Link href="/fish-book" className="p-2 bg-gray-100 rounded-full"><ChevronLeft /></Link>
                <div className="flex-1 text-center font-bold">{fish.name} / {fish.scientificName}</div>
            </header>
            <main className="p-4 max-w-md mx-auto space-y-4">
                <div className="bg-gray-200 aspect-video flex items-center justify-center text-2xl font-bold rounded-lg">
                    魚の画像を表示
                </div>

                <div className="grid grid-cols-2 gap-2 text-sm">
                    <p>別名：{fish.details.alias}</p>
                    <p>生息地：{fish.details.habitat}</p>
                    <p>体長：約{fish.details.length}</p>
                    <p>旬：{fish.details.season}</p>
                    <p>種類：{fish.details.type}</p>
                    <p>漁法：{fish.details.fishingMethod}</p>
                </div>

                <div className="flex gap-2">
                    <Link
                        href="/recipe/sample"
                        className="flex-1 py-2 bg-blue-100 text-blue-700 rounded-lg text-sm font-bold flex items-center justify-center filter"
                    >
                        おすすめレシピを探す
                    </Link>
                    <Link
                        href="/shopList"
                        className="flex-1 py-2 bg-blue-600 text-white rounded-lg text-sm font-bold flex items-center justify-center grayscale opacity-50 cursor-not-allowed"
                    >
                        食べられる店舗を探す
                    </Link>
                </div>

                <div className="space-y-4 pt-4 border-t">
                    <div>
                        <p className="font-bold border-l-4 border-blue-500 pl-2 mb-2">相性の良い魚</p>
                        <ul className="list-decimal list-inside text-sm text-gray-600">
                            {fish.details.goodCompatibility.map((f: string) => <li key={f}>{f}</li>)}
                        </ul>
                    </div>
                    <div>
                        <p className="font-bold border-l-4 border-red-400 pl-2 mb-2 text-red-600">相性の悪い魚</p>
                        <ul className="list-decimal list-inside text-sm text-gray-600">
                            {fish.details.badCompatibility.map((f: string) => <li key={f}>{f}</li>)}
                        </ul>
                    </div>
                </div>
            </main>
        </div>
    );
}
