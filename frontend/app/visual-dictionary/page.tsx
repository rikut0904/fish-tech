"use client";

import React, { useState, useMemo } from 'react';
import { Heart, Search, ChevronLeft, MapPin, ExternalLink } from 'lucide-react';
import { Fish } from '../types/fish';
import { Shop } from '../types/shop'; // Shopはこちらから読み込む

// モック店舗データ
const shops: Shop[] = [
  { id: 1, name: "近江町 魚中", address: "金沢市上近江町...", rating: "4.5", description: "地元の鮮魚を贅沢に使った寿司屋です。", tags: ["寿司", "刺身"], url: "https://sample.example.com/1", image: "/api/placeholder/100/100" },
  { id: 2, name: "海鮮居酒屋 まる", address: "金沢市片町...", rating: "4.2", description: "煮付けが絶品の人気店。", tags: ["煮付け", "地酒"], url: "https://sample.example.com/2", image: "/api/placeholder/100/100" },
];

// モックデータ生成 (120種類)
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
    // 相性データを3種類ずつに修正
    goodCompatibility: ["アジ", "サヨリ", "キス"],
    badCompatibility: ["マグロ", "ブリ", "シャチ"]
  }
}));

export default function FishEncyclopedia() {
  // view管理に 'shopList' と 'shopDetail' を追加
  const [view, setView] = useState<'list' | 'detail' | 'shopList' | 'shopDetail'>('list');
  const [selectedFish, setSelectedFish] = useState<Fish | null>(null);
  const [selectedShop, setSelectedShop] = useState<Shop | null>(null);
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

  const totalPages = Math.ceil(filteredFish.length / itemsPerPage);
  const currentItems = filteredFish.slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage);

  const toggleFavorite = (id: number, e: React.MouseEvent) => {
    e.stopPropagation();
    setFavorites(prev => prev.includes(id) ? prev.filter(fid => fid !== id) : [...prev, id]);
  };

  // --- 魚詳細画面 (ワイヤーフレームに基づき修正) ---
  if (view === 'detail' && selectedFish) {
    return (
      <div className="min-h-screen bg-white pb-20">
        <header className="p-4 flex items-center border-b sticky top-0 bg-white z-10">
          <button onClick={() => setView('list')}><ChevronLeft size={24} /></button>
          <div className="flex-1 text-center font-bold">
            <span className="text-red-500 mr-2 text-xl">●</span>
            {selectedFish.name} / {selectedFish.scientificName}
          </div>
          <button onClick={(e) => toggleFavorite(selectedFish.id, e)}>
            <Heart className={favorites.includes(selectedFish.id) ? "fill-red-500 text-red-500" : "text-gray-400"} />
          </button>
        </header>

        <main className="max-w-md mx-auto p-4 space-y-6">
          <div className="aspect-video bg-gray-100 rounded-xl overflow-hidden flex items-center justify-center relative">
            <img src={selectedFish.image} alt={selectedFish.name} className="w-full h-full object-cover" />
            <div className="absolute inset-0 flex items-center justify-center bg-black/10 text-white text-2xl font-bold">
              魚の画像を表示
            </div>
          </div>
          
          <div className="grid grid-cols-2 gap-y-3 gap-x-4 text-sm border-b pb-6">
            <p><span className="text-gray-500">別名：</span>{selectedFish.details.alias}</p>
            <p><span className="text-gray-500">生息地：</span>{selectedFish.details.habitat}</p>
            <p><span className="text-gray-500">体長：</span>約{selectedFish.details.length}</p>
            <p><span className="text-gray-500">旬：</span>{selectedFish.details.season}</p>
            <p><span className="text-gray-500">種類：</span>{selectedFish.details.type}</p>
            <p><span className="text-gray-500">漁法：</span>{selectedFish.details.fishingMethod}</p>
          </div>

          <div className="flex gap-3">
            <button className="flex-1 py-3 bg-blue-50 text-blue-700 rounded-xl text-xs font-bold border border-blue-100 hover:bg-blue-100 transition">
              おすすめレシピを探す
            </button>
            <button 
              onClick={() => setView('shopList')}
              className="flex-1 py-3 bg-blue-600 text-white rounded-xl text-xs font-bold hover:bg-blue-700 transition shadow-md"
            >
              食べられる店舗を探す
            </button>
          </div>

          <div className="space-y-6 pt-2">
            <div>
              <p className="font-bold text-gray-800 mb-2">相性の良い魚</p>
              <ol className="space-y-1 text-sm text-gray-600">
                {selectedFish.details.goodCompatibility.map((name, i) => <li key={i}>{i+1}. {name}</li>)}
              </ol>
            </div>
            <div>
              <p className="font-bold text-gray-800 mb-2">相性の悪い魚</p>
              <ol className="space-y-1 text-sm text-gray-600">
                {selectedFish.details.badCompatibility.map((name, i) => <li key={i}>{i+1}. {name}</li>)}
              </ol>
            </div>
          </div>
        </main>
      </div>
    );
  }

  // --- 店舗検索画面 (ヘッダーあり) ---
  if (view === 'shopList') {
    return (
      <div className="min-h-screen bg-gray-50">
        <header className="bg-blue-600 text-white p-6 text-center text-xl font-bold">ヘッダー (別メンバー作成分)</header>
        <div className="p-4 max-w-md mx-auto space-y-4">
          <button onClick={() => setView('detail')} className="flex items-center text-sm text-gray-500 mb-2">
            <ChevronLeft size={18} /> 店舗検索
          </button>
          <div className="relative">
            <Search className="absolute left-3 top-2.5 text-gray-400" size={18} />
            <input type="text" placeholder="店舗名を検索" className="w-full pl-10 pr-4 py-2 border rounded-full bg-white outline-none focus:ring-2 ring-blue-400" />
          </div>
          <p className="text-sm font-bold">検索結果: {shops.length}件</p>
          {shops.map(shop => (
            <div 
              key={shop.id} 
              onClick={() => { setSelectedShop(shop); setView('shopDetail'); }}
              className="flex gap-4 bg-white p-3 rounded-xl shadow-sm border cursor-pointer hover:border-blue-400 transition"
            >
              <div className="w-20 h-20 bg-gray-200 rounded-md shrink-0 flex items-center justify-center text-[10px] text-gray-400">画像</div>
              <div className="text-[11px] space-y-1">
                <p className="font-bold text-sm text-blue-900">{shop.name}</p>
                <p className="text-gray-500 line-clamp-1">{shop.description}</p>
                <p className="text-yellow-600 font-bold text-xs">★ {shop.rating} <span className="text-gray-400 font-normal ml-2">{shop.address}</span></p>
                <div className="flex gap-1 pt-1">
                  
                {shop.tags.map((l: string) => (
                <span key={l} className="bg-blue-50 text-blue-600 ...">
                  {l}
                </span>
              ))}
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    );
  }

  // --- 店舗詳細画面 (ヘッダーなし) ---
  if (view === 'shopDetail' && selectedShop) {
    return (
      <div className="min-h-screen bg-white p-6 max-w-md mx-auto space-y-6">
        <button onClick={() => setView('shopList')} className="p-2 bg-gray-100 rounded-full hover:bg-gray-200 transition"><ChevronLeft /></button>
        <div className="aspect-square bg-gray-100 rounded-2xl flex items-center justify-center text-gray-300 text-3xl font-bold">画像</div>
        <div className="space-y-4">
          <div className="flex justify-between items-start">
            <h2 className="text-2xl font-bold text-blue-900">{selectedShop.name}</h2>
            <span className="bg-yellow-100 text-yellow-700 px-3 py-1 rounded-full text-sm font-bold">★ {selectedShop.rating}</span>
          </div>
          <div className="space-y-4 text-gray-600">
            <p className="flex items-center gap-2 text-sm"><MapPin size={18} className="text-blue-500" /> {selectedShop.address}</p>
            <div className="pt-2">
              <p className="font-bold text-black mb-1">概要</p>
              <p className="text-sm leading-relaxed">{selectedShop.description}</p>
            </div>
            <div>
              <p className="font-bold text-black mb-1">タグ</p>
              <p className="text-sm">例：{selectedShop.tags.join('、')}</p>
            </div>
            <div className="pt-4">
              <p className="text-xs text-gray-400 mb-1">店舗URL</p>
              <a href={selectedShop.url} target="_blank" className="text-blue-600 underline flex items-center gap-1 break-all text-sm">
                {selectedShop.url} <ExternalLink size={14} />
              </a>
            </div>
          </div>
        </div>
      </div>
    );
  }

  // --- 魚図鑑一覧 (初期表示) ---
  return (
    <div className="min-h-screen bg-white">
      <section className="bg-gradient-to-b from-blue-400 to-blue-300 py-12 px-4 text-center text-white">
        <h1 className="text-3xl font-bold mb-4">魚図鑑</h1>
        <p className="text-sm opacity-90">金沢の"有名じゃない"魚たちをもっと知って、もっと味わおう。</p>
      </section>

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
          <select className="p-2 border rounded-lg bg-white" onChange={(e) => setFilterSeason(e.target.value)}>
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

        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 pb-20">
          {currentItems.map((fish) => (
            <div 
              key={fish.id}
              onClick={() => { setSelectedFish(fish); setView('detail'); window.scrollTo(0,0); }}
              className="group cursor-pointer bg-white border rounded-xl overflow-hidden hover:shadow-xl transition duration-300"
            >
              <div className="relative aspect-square bg-gray-100">
                <img src={fish.image} alt={fish.name} className="object-cover w-full h-full group-hover:scale-110 transition duration-500" />
                <button onClick={(e) => toggleFavorite(fish.id, e)} className="absolute top-2 right-2 p-2 bg-white/80 backdrop-blur rounded-full">
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
        <div className="flex justify-center gap-2 mt-8 mb-24">
          {[...Array(totalPages)].map((_, i) => (
            <button key={i} onClick={() => setCurrentPage(i + 1)} className={`w-10 h-10 rounded-lg font-bold ${currentPage === i + 1 ? 'bg-blue-600 text-white' : 'bg-gray-100 text-gray-600'}`}>{i + 1}</button>
          ))}
        </div>
      </div>

      <nav className="fixed bottom-0 w-full bg-white border-t flex justify-around p-3 text-[10px] font-bold text-gray-600 z-20">
        <div className="text-center opacity-50"><div className="text-xl">🏠</div>ホーム</div>
        <div className="text-center opacity-50"><div className="text-xl">🍳</div>レシピ</div>
        <div className="text-center text-blue-600"><div className="text-xl">🐟</div>魚図鑑</div>
        <div className="text-center opacity-50"><div className="text-xl">👤</div>マイページ</div>
      </nav>
    </div>
  );
}