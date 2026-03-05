"use client";

import React, { useState, useMemo } from 'react';
import { Heart, Search, Filter, ChevronLeft } from 'lucide-react';
// app/fish-book/page.tsx の場合
import { Fish } from '../types/fish';

const allFishData: Fish[] = Array.from({ length: 120 }, (_, i) => ({
  id: i + 1,
  name: i === 0 ? "カマス" : i === 1 ? "アカガレイ" : i === 2 ? "メギス" : `魚の名前 ${i + 1}`,
  scientificName: "Sphyraena pinguis",
  image: "/api/placeholder/400/320",
  isFavorite: false,
  details: {
    alias: "サヨリボウ", habitat: "石川県近海", length: "30cm",
    season: i % 4 === 0 ? "春" : i % 4 === 1 ? "夏" : i % 4 === 2 ? "秋" : "冬",
    type: i % 2 === 0 ? "白身" : "青物",
    fishingMethod: "定置網", recipe: "塩焼き、干物",
    compatibilityGood: ["スダチ", "大根おろし"],
    compatibilityBad: ["生クリーム"]
  }
}));

export default function FishEncyclopedia() {
  const [view, setView] = useState<'list' | 'detail'>('list');
  const [selectedFish, setSelectedFish] = useState<Fish | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [searchQuery, setSearchQuery] = useState("");
  const [filterSeason, setFilterSeason] = useState("");
  const [showOnlyFavorites, setShowOnlyFavorites] = useState(false);
  const [favorites, setFavorites] = useState<number[]>([]);

  const itemsPerPage = 40;

  // フィルタリング処理
  const filteredFish = useMemo(() => {
    return allFishData.filter(fish => {
      const matchesSearch = fish.name.includes(searchQuery);
      const matchesSeason = filterSeason ? fish.details.season === filterSeason : true;
      const matchesFavorite = showOnlyFavorites ? favorites.includes(fish.id) : true;
      return matchesSearch && matchesSeason && matchesFavorite;
    });
  }, [searchQuery, filterSeason, showOnlyFavorites, favorites]);

  // ページネーション処理
  const totalPages = Math.ceil(filteredFish.length / itemsPerPage);
  const currentItems = filteredFish.slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage);

  const toggleFavorite = (id: number, e: React.MouseEvent) => {
    e.stopPropagation();
    setFavorites(prev => prev.includes(id) ? prev.filter(fid => fid !== id) : [...prev, id]);
  };

  if (view === 'detail' && selectedFish) {
    return (
      <div className="min-h-screen bg-gray-50 pb-20">
        <header className="bg-white border-b p-4 flex items-center justify-between sticky top-0 z-10">
          <button onClick={() => setView('list')} className="flex items-center text-blue-600">
            <ChevronLeft size={24} /> 戻る
          </button>
          <div className="text-center">
            <h1 className="text-xl font-bold">{selectedFish.name}</h1>
            <p className="text-xs text-gray-500">{selectedFish.scientificName}</p>
          </div>
          <button onClick={(e) => toggleFavorite(selectedFish.id, e)}>
            <Heart className={favorites.includes(selectedFish.id) ? "fill-red-500 text-red-500" : "text-gray-400"} />
          </button>
        </header>

        <main className="max-w-md mx-auto p-4 space-y-6">
          <div className="bg-white rounded-2xl shadow-sm overflow-hidden border border-gray-100">
             <img src={selectedFish.image} alt={selectedFish.name} className="w-full h-64 object-cover" />
             <div className="p-6 space-y-4 text-sm">
                {[
                  { label: "別名", value: selectedFish.details.alias },
                  { label: "生息地", value: selectedFish.details.habitat },
                  { label: "体長", value: selectedFish.details.length },
                  { label: "旬", value: selectedFish.details.season },
                  { label: "種類", value: selectedFish.details.type },
                  { label: "漁法", value: selectedFish.details.fishingMethod },
                  { label: "おすすめレシピ", value: selectedFish.details.recipe },
                ].map((item, idx) => (
                  <div key={idx} className="flex border-b border-gray-50 pb-2">
                    <span className="w-32 font-bold text-gray-600">{item.label}</span>
                    <span className="text-gray-800">～ {item.value}</span>
                  </div>
                ))}
                <div className="pt-2">
                  <p className="font-bold text-blue-600 mb-1">相性の良い食材</p>
                  <p className="text-gray-600">{selectedFish.details.compatibilityGood.join('、')}</p>
                </div>
                <button className="w-full py-4 bg-blue-600 text-white rounded-xl font-bold hover:bg-blue-700 transition shadow-lg shadow-blue-200">
                  食べるまでのストーリーを見る
                </button>
             </div>
          </div>
        </main>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-white">
      {/* ヒーローセクション（FishTech風） */}
      <section className="bg-gradient-to-b from-blue-400 to-blue-300 py-12 px-4 text-center text-white">
        <h1 className="text-3xl font-bold mb-4">魚図鑑</h1>
        <p className="text-sm opacity-90">金沢の"有名じゃない"魚たちをもっと知って、もっと味わおう。</p>
      </section>

      {/* 検索・絞り込みエリア */}
      <div className="max-w-6xl mx-auto p-4 space-y-4">
        <div className="flex flex-wrap gap-3 items-center bg-gray-50 p-4 rounded-xl">
          <div className="relative flex-1 min-w-[200px]">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" size={18} />
            <input 
              type="text" placeholder="魚の名前で検索..." 
              className="w-full pl-10 pr-4 py-2 rounded-lg border focus:ring-2 ring-blue-500 outline-none"
              onChange={(e) => setSearchQuery(e.target.value)}
            />
          </div>
          <select 
            className="p-2 border rounded-lg bg-white"
            onChange={(e) => setFilterSeason(e.target.value)}
          >
            <option value="">旬を選択</option>
            <option value="春">春</option><option value="夏">夏</option>
            <option value="秋">秋</option><option value="冬">冬</option>
          </select>
          <button 
            onClick={() => setShowOnlyFavorites(!showOnlyFavorites)}
            className={`flex items-center gap-2 px-4 py-2 rounded-lg border transition ${showOnlyFavorites ? 'bg-red-50 border-red-200 text-red-600' : 'bg-white'}`}
          >
            <Heart size={18} className={showOnlyFavorites ? "fill-red-600" : ""} /> お気に入り
          </button>
        </div>

        {/* 魚カードグリッド */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          {currentItems.map((fish) => (
            <div 
              key={fish.id}
              onClick={() => { setSelectedFish(fish); setView('detail'); window.scrollTo(0,0); }}
              className="group cursor-pointer bg-white border rounded-xl overflow-hidden hover:shadow-xl transition duration-300"
            >
              <div className="relative aspect-square bg-gray-100">
                <img src={fish.image} alt={fish.name} className="object-cover w-full h-full group-hover:scale-110 transition duration-500" />
                <button 
                  onClick={(e) => toggleFavorite(fish.id, e)}
                  className="absolute top-2 right-2 p-2 bg-white/80 backdrop-blur rounded-full hover:bg-white"
                >
                  <Heart size={20} className={favorites.includes(fish.id) ? "fill-red-500 text-red-500" : "text-gray-400"} />
                </button>
              </div>
              <div className="p-3">
                <h3 className="font-bold text-blue-900">{fish.name}</h3>
                <p className="text-[10px] text-gray-500 truncate">{fish.scientificName}</p>
              </div>
            </div>
          ))}
        </div>

        {/* ページネーション */}
        <div className="flex justify-center gap-2 mt-12 mb-20">
          {[...Array(totalPages)].map((_, i) => (
            <button
              key={i}
              onClick={() => setCurrentPage(i + 1)}
              className={`w-10 h-10 rounded-lg font-bold transition ${currentPage === i + 1 ? 'bg-blue-600 text-white' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'}`}
            >
              {i + 1}
            </button>
          ))}
        </div>
      </div>

      {/* フッターナビゲーション */}
      <nav className="fixed bottom-0 w-full bg-white border-t flex justify-around p-3 text-[10px] font-bold text-gray-600">
        <div className="text-center opacity-50"><div className="text-xl">🏠</div>ホーム</div>
        <div className="text-center opacity-50"><div className="text-xl">🍳</div>レシピ</div>
        <div className="text-center text-blue-600"><div className="text-xl">🐟</div>魚図鑑</div>
        <div className="text-center opacity-50"><div className="text-xl">👤</div>マイページ</div>
      </nav>
    </div>
  );
}