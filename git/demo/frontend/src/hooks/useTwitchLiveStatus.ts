import { useEffect, useState } from 'react';

const CHANNEL_DUMMY = 'goalkeeper91';

const useTwitchLiveStatus = () => {
    const[isLive, setIsLive] = useState<boolean>(false);

    useEffect(() => {
        const fetchLiveStatus = async () => {
            try {
                const response = await fetch(`https://decapi.me/twitch/status/${CHANNEL_DUMMY}`);
                const text = await response.text();
                setIsLive(!text.includes('offline'));
            } catch (error) {
                console.error('Fehler beim Abrufen des Twitch-Status', error);
            }
        };
        fetchLiveStatus();
            const interval = setInterval(fetchLiveStatus, 30000);

            return () => clearInterval(interval);
    }, []);

    return isLive;
};

export default useTwitchLiveStatus;