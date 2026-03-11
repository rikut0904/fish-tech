"use client";

import React, { useState } from 'react';
import { Heart, ChevronLeft, Search, MapPin, ExternalLink } from 'lucide-react';
import { Shop } from "../types/shop";


// モックデータ (店舗)
const shops: Shop[] = [
    { id: 1, name: "近江町 魚中", address: "金沢市上近江町...", rating: "4.5", description: "地元の鮮魚を贅沢に使った寿司屋です。", tags: ["寿司", "刺身"], url: "https://sample.example.com/1", image: "/api/placeholder/100/100" },
    { id: 2, name: "海鮮居酒屋 まる", address: "金沢市片町...", rating: "4.2", description: "煮付けが絶品の人気店。", tags: ["煮付け", "地酒"], url: "https://sample.example.com/2", image: "/api/placeholder/100/100" },
];

// モックデータ (魚)
const allFishData = [
    {
        id: 1,
        name: "サンプル魚",
        scientificName: "Sampleus piscis",
        details: {
            alias: "テスト魚",
            habitat: "沿岸",
            length: "30cm",
            season: "通年",
            type: "白身",
            fishingMethod: "定置網",
            goodCompatibility: ["魚A", "魚B"],
            badCompatibility: ["魚C"]
        }
    }
];

export default function FishApp() {
    // list: 図鑑一覧, detail: 魚詳細, shopList: 店舗一覧, shopDetail: 店舗詳細
    const [view, setView] = useState<'list' | 'detail' | 'shopList' | 'shopDetail'>('list');
    const [selectedFish, setSelectedFish] = useState<any>(null);
    const [selectedShop, setSelectedShop] = useState<any>(null);

    // --- 魚詳細画面 ---
    if (view === 'detail') {
        return (
            <div className="min-h-screen bg-white">
                <header className="p-4 flex items-center border-b">
                    <button onClick={() => setView('list')}>
                        <ChevronLeft />
                    </button>
                    <div className="flex-1 text-center font-bold">
                        <span className="text-red-500 mr-2">
                            ●
                        </span>
                        {selectedFish.name} / {selectedFish.scientificName}
                    </div>
                </header>
                <main className="p-4 max-w-md mx-auto space-y-4">
                    <div className="bg-gray-200 aspect-video flex items-center justify-center text-2xl font-bold rounded-lg">
                        魚の画像を表示
                    </div>

                    <div className="grid grid-cols-2 gap-2 text-sm">
                        <p>別名：{selectedFish.details.alias}</p>
                        <p>生息地：{selectedFish.details.habitat}</p>
                        <p>体長：約{selectedFish.details.length}</p>
                        <p>旬：{selectedFish.details.season}</p>
                        <p>種類：{selectedFish.details.type}</p>
                        <p>漁法：{selectedFish.details.fishingMethod}</p>
                    </div>

                    <div className="flex gap-2">
                        <button className="flex-1 py-2 bg-blue-100 text-blue-700 rounded-lg text-sm font-bold">
                            おすすめレシピを探す
                        </button>
                        <button
                            onClick={() => setView('shopList')}
                            className="flex-1 py-2 bg-blue-600 text-white rounded-lg text-sm font-bold"
                        >
                            食べられる店舗を探す
                        </button>
                    </div>

                    <div className="space-y-4 pt-4 border-t">
                        <div>
                            <p className="font-bold border-l-4 border-blue-500 pl-2 mb-2">相性の良い魚</p>
                            <ul className="list-decimal list-inside text-sm text-gray-600">
                                {selectedFish.details.goodCompatibility.map((f: string) => <li key={f}>{f}</li>)}
                            </ul>
                        </div>
                        <div>
                            <p className="font-bold border-l-4 border-red-400 pl-2 mb-2 text-red-600">相性の悪い魚</p>
                            <ul className="list-decimal list-inside text-sm text-gray-600">
                                {selectedFish.details.badCompatibility.map((f: string) => <li key={f}>{f}</li>)}
                            </ul>
                        </div>
                    </div>
                </main>
            </div>
        );
    }

    // --- 店舗検索ページ (ヘッダーあり) ---
    if (view === 'shopList') {
        return (
            <div className="min-h-screen bg-gray-50">
                <header className="bg-blue-600 text-white p-6 text-center text-xl font-bold">
                    ヘッダー (別メンバー作成分)
                </header>
                <div className="p-4 max-w-md mx-auto space-y-4">
                    <button onClick={() => setView('detail')} className="text-sm text-gray-500 flex items-center gap-1">
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
                            onClick={() => { setSelectedShop(shop); setView('shopDetail'); }}
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

    // --- 店舗詳細ページ (ヘッダーなし) ---
    if (view === 'shopDetail' && selectedShop) {
        return (
            <div className="min-h-screen bg-white p-6 max-w-md mx-auto space-y-6">
                <button onClick={() => setView('shopList')} className="p-2 bg-gray-100 rounded-full"><ChevronLeft /></button>
                <div className="aspect-square bg-gray-100 rounded-2xl flex items-center justify-center text-gray-400 text-3xl font-bold">
                    画像
                </div>
                <div className="space-y-4">
                    <div className="flex justify-between items-start">
                        <h2 className="text-2xl font-bold">{selectedShop.name}</h2>
                        <span className="bg-yellow-100 text-yellow-700 px-3 py-1 rounded-full text-sm font-bold">★ {selectedShop.rating}</span>
                    </div>
                    <div className="space-y-2 text-gray-600">
                        <p className="flex items-center gap-2"><MapPin size={18} /> {selectedShop.address}</p>
                        <p className="pt-4 font-bold text-black">概要</p>
                        <p className="text-sm leading-relaxed">{selectedShop.description}</p>
                        <p className="pt-2 font-bold text-black">タグ</p>
                        <p className="text-sm">例：{selectedShop.tags.join('、')}</p>
                    </div>
                    <div className="pt-6">
                        <p className="text-xs text-gray-400 mb-1">店舗URL</p>
                        <a href={selectedShop.url} className="text-blue-600 underline flex items-center gap-1 break-all">
                            {selectedShop.url} <ExternalLink size={14} />
                        </a>
                    </div>
                </div>
            </div>
        );
    }

    // --- 初期表示：魚図鑑一覧 ---
    return (
        <div className="p-4 bg-white min-h-screen">
            <h1 className="text-2xl font-bold mb-6 text-center">魚図鑑</h1>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                {allFishData.map((fish) => (
                    <button
                        key={fish.id}
                        onClick={() => {
                            setSelectedFish(fish);
                            setView('detail');
                        }}
                        className="border rounded-xl overflow-hidden shadow-sm hover:shadow-md transition bg-gray-50 text-left"
                    >
                        <div className="aspect-video bg-gray-200 flex items-center justify-center text-gray-400">
                            {/* 本来はここに fish.image */}
                            画像
                        </div>
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