"use client";

import React, { useState } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { ChefHat, Clock, Star, ExternalLink, Search } from 'lucide-react';
import { useSearchParams } from 'next/navigation';
import { allFishData } from '../../fish-book/data';

const sampleRecipes = [
    {
        id: 's1',
        title: 'サンプル魚のソテー 〜バター醤油〜',
        image: '/recipes-sample/image.png',
        time: '15分',
        cost: '約500円',
        score: 5,
        desc: 'シンプルに塩胡椒してバター醤油で香ばしく仕上げます。',
        url: 'https://cookpad.com/jp/recipes/17858696',
    },
    {
        id: 's2',
        title: 'サンプル魚の南蛮漬け',
        image: '/recipes-sample/image.png',
        time: '30分',
        cost: '約400円',
        score: 4,
        desc: 'さっぱりとした味わいでご飯にもおつまみにも合います。',
        url: 'https://cookpad.com/jp/recipes/20238212',
    }
];

export default function RecipeSampleClient() {
    const [q, setQ] = useState('');
    const searchParams = useSearchParams();
    const fishIdParam = searchParams?.get('fish');
    const fishId = fishIdParam ? Number(fishIdParam) : null;
    const fish = fishId ? allFishData.find(f => f.id === fishId) : null;
    const fishName = fish ? fish.name : 'サンプル魚';
    const filtered = sampleRecipes.filter(r => r.title.includes(q) || r.desc.includes(q));

    return (
        <div className="min-h-screen bg-gradient-to-b from-blue-50 to-white p-6">
            <header className="mb-8 flex items-center gap-3">
                <ChefHat className="text-blue-600" />
                <div>
                    <h1 className="text-2xl font-bold">{fishName}に関する料理（{filtered.length}件）</h1>
                    <p className="text-sm text-gray-500">{fishName}向けのおすすめレシピを掲載しています</p>
                </div>
                <div className="ml-auto">
                    <Link href="/fish-book/detail/1" className="text-sm text-blue-600 underline">魚ページに戻る</Link>
                </div>
            </header>

            <div className="mb-6">
                <div className="relative max-w-md">
                    <Search className="absolute left-3 top-3 text-gray-400" />
                    <input value={q} onChange={(e) => setQ(e.target.value)} placeholder="レシピを検索" className="w-full pl-10 pr-4 py-2 border rounded-lg" />
                </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {filtered.map(r => (
                    <article key={r.id} className="bg-white rounded-lg shadow-md overflow-hidden">
                        <div className="relative h-44 bg-gray-100">
                            <Image src={r.image} alt={r.title} fill className="object-cover" />
                        </div>
                        <div className="p-4">
                            <h3 className="text-lg font-semibold mb-2 flex items-center justify-between">
                                {r.title}
                                <span className="text-sm text-gray-500">{r.time}</span>
                            </h3>
                            <p className="text-sm text-gray-700 mb-3">{r.desc}</p>
                            <div className="flex items-center justify-between">
                                <div className="flex items-center gap-2 text-sm text-gray-600">
                                    <Star className="text-yellow-400" />
                                    <span>{r.score}点</span>
                                </div>
                                <div className="flex items-center gap-2">
                                    <span className="text-sm text-gray-500">{r.cost}</span>
                                    <Link href={r.url} target="_blank" className="text-blue-600 underline flex items-center gap-1">作り方 <ExternalLink size={14} /></Link>
                                </div>
                            </div>
                        </div>
                    </article>
                ))}
            </div>
        </div>
    );
}
