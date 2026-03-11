import FishBackground from '@/app/components/fish_background';

export const metadata = {
  title: 'プライバシーポリシー - Fish-Tech',
  description: 'Fish-Tech のプライバシーポリシー',
};

export default function PrivacyPage() {
  return (
    <main className="relative min-h-screen bg-gradient-to-b from-sky-50 to-white dark:from-slate-900 dark:to-slate-800">
      <FishBackground />
      <section className="container mx-auto px-4 py-16">
        <div className="max-w-3xl mx-auto bg-white/95 dark:bg-slate-900/75 backdrop-blur rounded-xl shadow-lg p-8">
          <h1 className="text-3xl font-bold text-blue-700 mb-4">プライバシーポリシー</h1>
          <p className="text-sm text-slate-600 dark:text-slate-300 mb-6">最終更新日: 2026-03-09</p>

          <section className="mb-4">
            <h2 className="text-xl font-semibold mb-2">1. 基本方針</h2>
            <p className="text-slate-700 dark:text-slate-200">
              Fish-Tech（以下「当アプリ」）は、利用者の個人情報を適切に取り扱い、その保護に努めます。
            </p>
          </section>

          <section className="mb-4">
            <h2 className="text-xl font-semibold mb-2">2. 収集する情報</h2>
            <ul className="list-disc pl-5 text-slate-700 dark:text-slate-200">
              <li>ユーザーが自発的に提供する情報（例：登録情報、問い合わせ内容）</li>
              <li>端末・利用履歴に関する情報（クッキー、アクセスログ、IPアドレス等）</li>
            </ul>
          </section>

          <section className="mb-4">
            <h2 className="text-xl font-semibold mb-2">3. 利用目的</h2>
            <p className="text-slate-700 dark:text-slate-200">
              収集した情報は、サービス提供・改善、ユーザーサポート、セキュリティ向上、法令遵守のために利用します。
            </p>
          </section>

          <section className="mb-4">
            <h2 className="text-xl font-semibold mb-2">4. 第三者提供</h2>
            <p className="text-slate-700 dark:text-slate-200">
              法令に基づく場合や、利用目的の達成に必要な範囲で外部サービス（解析ツール等）に提供することがあります。
            </p>
          </section>

          <section className="mb-4">
            <h2 className="text-xl font-semibold mb-2">5. クッキー等</h2>
            <p className="text-slate-700 dark:text-slate-200">
              当アプリでは利便性向上や利用解析のためにクッキーを使用することがあります。ブラウザ設定により無効化できますが、機能の一部が制限される場合があります。
            </p>
          </section>

          <section className="mb-4">
            <h2 className="text-xl font-semibold mb-2">6. セキュリティ</h2>
            <p className="text-slate-700 dark:text-slate-200">
              情報の漏えい・改ざん防止のために合理的な安全対策を実施しますが、完全な安全性を保証するものではありません。
            </p>
          </section>

          <section className="mb-4">
            <h2 className="text-xl font-semibold mb-2">7. 個人情報の開示・訂正</h2>
            <p className="text-slate-700 dark:text-slate-200">
              利用者本人からの請求があった場合、法令に従い対応いたします。手続きについてはお問い合わせください。
            </p>
          </section>

          <section className="mb-4">
            <h2 className="text-xl font-semibold mb-2">8. お問い合わせ</h2>
            <p className="text-slate-700 dark:text-slate-200">
              プライバシーに関するお問い合わせは、当サイトのお問い合わせページからご連絡ください。
            </p>
          </section>

          <p className="text-xs text-slate-500 mt-6">※ 本ポリシーは予告なく改定されることがあります。最新の内容は本ページでご確認ください。</p>
        </div>
      </section>
    </main>
  );
}
