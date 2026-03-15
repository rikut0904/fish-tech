// フッター
export default function Footer() {
    return (
        <footer className="bg-white border-t py-4 text-center text-xs md:text-sm text-gray-500" id="contact">
            &copy; 2026 FishTech / iwashikurukuru. All rights reserved.

            <br className="md:hidden" />
            <span className="hidden md:inline">&nbsp;|&nbsp;</span>
            <a href="/privacy">プライバシーポリシー</a>
            &nbsp;|&nbsp;
            <a href="/sitemap">サイトマップ</a>
            &nbsp;|&nbsp;
            <a href="mailto:sample@example.com">お問い合わせ</a>
        </footer>
    );
}
