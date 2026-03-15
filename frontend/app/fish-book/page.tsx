"use client";

import React from 'react';
import { useRouter } from 'next/navigation';
import { allFishData } from './data';
import { Shop } from "../types/shop";
import Image from 'next/image';


// モックデータ (店舗)
const shops: Shop[] = [
    { id: 1, name: "近江町 魚中", address: "金沢市上近江町...", rating: "4.5", description: "地元の鮮魚を贅沢に使った寿司屋です。", tags: ["寿司", "刺身"], url: "https://sample.example.com/1", image: "/api/placeholder/100/100" },
    { id: 2, name: "海鮮居酒屋 まる", address: "金沢市片町...", rating: "4.2", description: "煮付けが絶品の人気店。", tags: ["煮付け", "地酒"], url: "https://sample.example.com/2", image: "/api/placeholder/100/100" },
];

// 魚データは `data.ts` から読み込み

// 一覧ページのみを表示（詳細は別ルートへ）
export default function FishApp() {
    const router = useRouter();

    return (
        <div className="p-4 bg-white min-h-screen">
            <h1 className="text-2xl font-bold mb-6 text-center">魚図鑑</h1>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                {allFishData.map((fish) => (
                    <button
                        key={fish.id}
                        onClick={() => router.push(`/fish-book/detail/${fish.id}`)}
                        className="border rounded-xl overflow-hidden shadow-sm hover:shadow-md transition bg-gray-50 text-left"
                    >
                        <Image src="/hero/fish_01.jpg" alt={fish.name} width={400} height={300} className="w-full h-32 object-cover" />
                        <div className="p-3">
                            <p className="font-bold text-blue-800">{fish.name}</p>
                            <p className="text-xs text-gray-500 italic">{fish.scientificName}</p>
                        </div>
                    </button>
                ))}
            </div>
        </div>
    );
}