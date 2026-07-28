// hooks/usePricePreview.ts
import { useEffect, useState } from "react";
import { getPaddle } from "../lib/paddleSdk";

// Maps a Paddle price_id to the exact formatted total string Paddle itself
// returns - no frontend price math, no re-formatting, no rounding, per the
// requirement to only ever display Paddle's own formattedTotals values.
export function usePricePreview(priceIds: string[]) {
  const [totals, setTotals] = useState<Record<string, string>>({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const key = priceIds.join(",");

  useEffect(() => {
    let cancelled = false;
    if (priceIds.length === 0) {
      setLoading(false);
      return;
    }

    setLoading(true);
    setError("");

    (async () => {
      try {
        const paddle = await getPaddle();
        if (!paddle) {
          throw new Error("Paddle.js konnte nicht initialisiert werden");
        }

        // No address.countryCode passed - Paddle auto-detects the
        // visitor's location from their browser IP when it's omitted.
        const response = await paddle.PricePreview({
          items: priceIds.map((priceId) => ({ priceId, quantity: 1 })),
        });

        if (cancelled) return;

        const next: Record<string, string> = {};
        for (const lineItem of response.data.details.lineItems) {
          next[lineItem.price.id] = lineItem.formattedTotals.total;
        }
        setTotals(next);
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : "Preise konnten nicht geladen werden");
        }
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();

    return () => {
      cancelled = true;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [key]);

  return { totals, loading, error };
}
