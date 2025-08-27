import { useEffect, useState, useRef } from 'react';

const videoIds = [
    'JDSXB-n1QLQ',
    'uPzpZbB4mr8&pp=0gcJCccJAYcqIYzv',
    'AA_m8rLcIx0',
    'MWrnLqYwfjE&pp=0gcJCccJAYcqIYzv',
    'PzggVsIcSuY&pp=0gcJCccJAYcqIYzv',
    'sokBPz93o5g',
    'hNwE7r25Amv8',
    'iN0P_5q8Iso',
    '6jvlln4p3o0',
    'gaZXybW9_Ew'
];

const getRandomVideoId = (excludeId?: string) => {
  const filtered = excludeId ? videoIds.filter(id => id !== excludeId) : videoIds;
  return filtered[Math.floor(Math.random() * filtered.length)];
};

const RandomYoutubePlayer = () => {
  const [currentVideoId, setCurrentVideoId] = useState(getRandomVideoId());
  const playerRef = useRef<any>(null);
  const [consentGiven, setConsentGiven] = useState(false);

  // Cookiebot Consent prüfen
  useEffect(() => {
    const checkConsent = () => {
      if (window.Cookiebot && window.Cookiebot.consent) {
        setConsentGiven(!!window.Cookiebot.consent.marketing);
      }
    };

    window.addEventListener("CookieConsentDeclaration", checkConsent);

    // Polling für initialen Consent, falls Cookiebot noch nicht geladen ist
    const interval = setInterval(() => {
      if (window.Cookiebot && window.Cookiebot.consent) {
        checkConsent();
        clearInterval(interval);
      }
    }, 100);

    return () => {
      window.removeEventListener("CookieConsentDeclaration", checkConsent);
      clearInterval(interval);
    };
  }, []);

  // YouTube-Player nur laden, wenn Consent gegeben
  useEffect(() => {
    if (!consentGiven) return;

    const tag = document.createElement("script");
    tag.src = "https://www.youtube.com/iframe_api";
    document.body.appendChild(tag);

    (window as any).onYouTubeIframeAPIReady = () => {
      playerRef.current = new (window as any).YT.Player("youtube-player", {
        height: "390",
        width: "640",
        videoId: currentVideoId,
        events: {
          onStateChange: (event: any) => {
            if (event.data === 0) {
              setCurrentVideoId(getRandomVideoId(currentVideoId));
            }
          },
        },
      });
    };

    return () => {
      if (playerRef.current?.destroy) playerRef.current.destroy();
    };
  }, [consentGiven]);

  // Video wechseln, wenn currentVideoId sich ändert
  useEffect(() => {
    if (playerRef.current?.loadVideoById) {
      playerRef.current.loadVideoById(currentVideoId);
    }
  }, [currentVideoId]);

  if (!consentGiven) {
    return (
      <div className="w-full flex justify-center items-center py-6">
        <p className="text-center text-gray-400">
          YouTube-Videos werden nach Zustimmung zu Marketing-Cookies angezeigt.
        </p>
      </div>
    );
  }

  return (
    <div className="w-full flex justify-center items-center py-6">
      <div id="youtube-player" className="rounded-lg shadow-lg" />
    </div>
  );
};

export default RandomYoutubePlayer;
