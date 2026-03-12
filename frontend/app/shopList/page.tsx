"use client";

import React from 'react';
import { ChevronLeft, Search } from 'lucide-react';
import { useRouter } from 'next/navigation';
import { Shop } from "../types/shop";

// モックデータ (店舗)
const shops: Shop[] = [
    { id: 1, name: "近江町 魚中", address: "金沢市上近江町...", rating: "4.5", description: "地元の鮮魚を贅沢に使った寿司屋です。", tags: ["寿司", "刺身"], url: "https://sample.example.com/1", image: "/api/placeholder/100/100" },
    { id: 2, name: "海鮮居酒屋 まる", address: "金沢市片町...", rating: "4.2", description: "煮付けが絶品の人気店。", tags: ["煮付け", "地酒"], url: "https://sample.example.com/2", image: "/api/placeholder/100/100" },
];

export default function ShopListPage() {
    const router = useRouter();

    return (
        <div className="min-h-screen bg-gray-50">
            <header className="bg-blue-600 text-white p-6 text-center text-xl font-bold">
                ヘッダー (別メンバー作成分)
            </header>
            <div className="p-4 max-w-md mx-auto space-y-4">
                <button onClick={() => router.back()} className="text-sm text-gray-500 flex items-center gap-1">
                    <ChevronLeft size={16} /> 店舗検索
                </button>
                <div className="relative">
                    <Search className="absolute left-3 top-2.5 text-gray-400" size={18} />
                    <input type="text" placeholder="店舗名を検索" className="w-full pl-10 pr-4 py-2 border rounded-full bg-white" />
                </div>
                <p className="text-sm font-bold">検索結果: {shops.length}件</p>
                {shops.map(shop => (
                    <div
                        key={shop.id}
                        onClick={() => router.push('/shopList')}
                        className="flex gap-4 bg-white p-3 rounded-xl shadow-sm border cursor-pointer hover:border-blue-400"
                    >
                        <div className="w-20 h-20 bg-gray-200 rounded-md shrink-0">画像</div>
                        <div className="text-xs space-y-1">
                            <p className="font-bold text-sm">{shop.name}</p>
                            <p className="text-gray-500 line-clamp-1">{shop.description}</p>
                            <p className="text-yellow-600">評価: ★{shop.rating}</p>
                            <p className="text-gray-400">{shop.address}</p>
                            <div className="flex gap-1">
                                {shop.tags.map(t => <span key={t} className="bg-gray-100 px-1 rounded">#{t}</span>)}
                            </div>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
}
