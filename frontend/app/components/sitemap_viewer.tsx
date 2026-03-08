"use client"
import React, { useEffect, useState } from "react";

type Entry = {
  loc: string;
  lastmod?: string;
  changefreq?: string;
  priority?: string;
};

export default function SitemapViewer() {
  const [entries, setEntries] = useState<Entry[] | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    fetch('/sitemap.xml')
      .then((res) => {
        if (!res.ok) throw new Error('sitemap fetch failed');
        return res.text();
      })
      .then((text) => {
        if (cancelled) return;
        const parser = new DOMParser();
        const xml = parser.parseFromString(text, 'application/xml');
        const urlNodes = Array.from(xml.getElementsByTagName('url'));
        const parsed = urlNodes.map((u) => ({
          loc: u.getElementsByTagName('loc')[0]?.textContent || '',
          lastmod: u.getElementsByTagName('lastmod')[0]?.textContent || undefined,
          changefreq: u.getElementsByTagName('changefreq')[0]?.textContent || undefined,
          priority: u.getElementsByTagName('priority')[0]?.textContent || undefined,
        }));
        setEntries(parsed);
      })
      .catch((e) => {
        if (cancelled) return;
        setError(String(e));
      });
    return () => {
      cancelled = true;
    };
  }, []);

  if (error) return <div className="text-red-600">読み込みエラー: {error}</div>;
  if (!entries) return <div>読み込み中…</div>;

  return (
    <div className="prose max-w-none">
      <ul>
        {entries.map((e, i) => (
          <li key={i} className="mb-2">
            <a href={e.loc} className="text-blue-600 hover:underline">
              {e.loc}
            </a>
            {e.lastmod && <div className="text-sm text-gray-500">更新日: {e.lastmod}</div>}
          </li>
        ))}
      </ul>
    </div>
  );
}
