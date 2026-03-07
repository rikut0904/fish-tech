"use client";
import { useEffect, useRef, useState } from "react";
import { usePathname } from "next/navigation";

export default function RouteLoader() {
  const pathname = usePathname();
  const [loading, setLoading] = useState(false);
  const prevPathRef = useRef(pathname);
  const timeoutRef = useRef<number | null>(null);
  const startTimeRef = useRef<number | null>(null);

  useEffect(() => {
    // When pathname changes, navigation finished — hide loader after minimum duration
    if (prevPathRef.current !== pathname) {
      prevPathRef.current = pathname;

      const minMs = 3000;
      const started = startTimeRef.current ?? 0;
      const elapsed = started ? Date.now() - started : minMs;

      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
        timeoutRef.current = null;
      }

      if (elapsed >= minMs) {
        setLoading(false);
        startTimeRef.current = null;
      } else {
        timeoutRef.current = window.setTimeout(() => {
          setLoading(false);
          timeoutRef.current = null;
          startTimeRef.current = null;
        }, minMs - elapsed);
      }
    }
  }, [pathname]);

  useEffect(() => {
    const onClick = (e: MouseEvent) => {
      const target = e.target as HTMLElement | null;
      if (!target) return;

      // find nearest anchor
      const anchor = target.closest("a") as HTMLAnchorElement | null;
      if (!anchor) return;
      const href = anchor.getAttribute("href");
      if (!href) return;

      // ignore external links and hash/mailto/tel
      if (href.startsWith("http") || href.startsWith("mailto:") || href.startsWith("tel:") || href.startsWith("#")) return;
      // ignore modified clicks (open in new tab/window)
      if (e.metaKey || e.ctrlKey || e.shiftKey || e.altKey) return;

      // internal navigation — show loader and record start time
      setLoading(true);
      startTimeRef.current = Date.now();

      // safety timeout to avoid stuck loader (longer than min display)
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
      timeoutRef.current = window.setTimeout(() => {
        setLoading(false);
        timeoutRef.current = null;
        startTimeRef.current = null;
      }, 10000);
    };

    document.addEventListener("click", onClick);
    return () => {
      document.removeEventListener("click", onClick);
      if (timeoutRef.current) clearTimeout(timeoutRef.current);
    };
  }, []);

  if (!loading) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-white/50 backdrop-blur-sm">
      <img src="/iwashi.png" alt="ローディング中のイワシ" className="w-24 h-24 object-contain animate-spin" />
    </div>
  );
}
