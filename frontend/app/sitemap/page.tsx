import SitemapViewer from '../components/sitemap_viewer'

export const metadata = {
    title: 'サイトマップ',
    description: 'サイトマップ（sitemap.xml）を読み込んで表示します。',
}

export default function SitemapPage() {
    return (
        <main className="container mx-auto p-6">
            <h1 className="text-2xl font-bold mb-4">サイトマップ</h1>
            <p className="mb-4 text-sm text-gray-600">
                サイトの全URLを一覧で確認できます。検索エンジン向けのファイルは
                {' '}
                <a href="/sitemap.xml" className="text-blue-600 hover:underline">sitemap.xml</a>
                にあります。
            </p>
            <SitemapViewer />
        </main>
    );
}
