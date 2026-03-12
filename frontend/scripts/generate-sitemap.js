const fs = require('fs');
const path = require('path');

const appDir = path.join(__dirname, '..', 'app');
const publicDir = path.join(__dirname, '..', 'public');
const OUT = path.join(publicDir, 'sitemap.xml');
// 優先順: ローカル変数 SITE_DOMAIN, Next.js の公開 env (もしあれば), 最後にデフォルト
const RAW_DOMAIN = process.env.SITE_DOMAIN || process.env.NEXT_PUBLIC_SITE_URL || process.env.NEXT_PUBLIC_SITE_DOMAIN || 'https://example.com';
// 末尾のスラッシュを除去して `/path` と連結したときに `//` にならないようにする
const DOMAIN = RAW_DOMAIN.replace(/\/$/, '');

function walk(dir) {
  const entries = [];
  const items = fs.readdirSync(dir, { withFileTypes: true });
  for (const it of items) {
    // ignore components and API routes and internal folders
    if (it.name === 'components' || it.name === 'api' || it.name === 'styles') continue;
    const full = path.join(dir, it.name);
    if (it.isDirectory()) {
      entries.push(...walk(full));
      // if directory contains a page file, include the directory path as a route
      const pageFiles = ['page.tsx', 'page.jsx', 'page.ts', 'page.js'];
      for (const pf of pageFiles) {
        if (fs.existsSync(path.join(full, pf))) {
          const rel = path.relative(appDir, full);
          const route = rel === '' ? '/' : `/${rel.replace(/\\\\/g, '/')}`;
          entries.push(route);
          break;
        }
      }
    } else {
      // root-level page file (app/page.tsx)
      if (dir === appDir) {
        const pageFiles = ['page.tsx', 'page.jsx', 'page.ts', 'page.js'];
        if (pageFiles.includes(it.name)) entries.push('/');
      }
    }
  }
  return entries;
}

function uniq(a) { return Array.from(new Set(a)); }

function buildSitemap(urls) {
  const now = new Date().toISOString().slice(0,10);
  const items = urls.map(u => `  <url>\n    <loc>${DOMAIN}${u}</loc>\n    <lastmod>${now}</lastmod>\n    <changefreq>weekly</changefreq>\n  </url>`).join('\n');
  return `<?xml version="1.0" encoding="UTF-8"?>\n<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">\n${items}\n</urlset>\n`;
}

function ensurePublic() {
  if (!fs.existsSync(publicDir)) fs.mkdirSync(publicDir, { recursive: true });
}

function main() {
  if (!fs.existsSync(appDir)) {
    console.error('app directory not found:', appDir);
    process.exit(1);
  }
  const routes = walk(appDir).filter(Boolean).map(r => r.replace(/\\/g, '/'));
  const unique = uniq(routes).sort();
  ensurePublic();
  const xml = buildSitemap(unique);
  fs.writeFileSync(OUT, xml, 'utf8');
  console.log('sitemap generated:', OUT);
}

if (require.main === module) main();
