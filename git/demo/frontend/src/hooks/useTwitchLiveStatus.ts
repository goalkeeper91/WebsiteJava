import { useEffect, useState } from 'react';

const useTwitchLiveStatus = () => {
    const [isLive, setIsLive] = useState<boolean>(false);

    useEffect(() => {
        const fetchLiveStatus = async () => {
            try {
                const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/twitch/status`);
                if (!response.ok) throw new Error('Network response not ok');
                const data = await response.json();
                setIsLive(data.live);
            } catch (error) {
                console.log(error)
                console.error('Fehler beim Abrufen des Twitch-Status', error);
            }
        };
        fetchLiveStatus();
            const interval = setInterval(fetchLiveStatus, 1000);

            return () => clearInterval(interval);
    }, []);

    return isLive;
};

export default useTwitchLiveStatus;