import React from 'react';

interface TwitchChatEmbedProps {
    channel: string;
    darkMode?: boolean;
}

const TwitchChatEmbed: React.FC<TwitchChatEmbedProps> = ({ channel, darkMode = true }) => {
    const parent = window.location.hostname;
    const darkParam = darkMode ? '&darkpopout' : '';

    return (
        <div className="h-full w-full rounded-xl overflow-hidden shadow-lg">
            <iframe
                data-cookieconsent="marketing"
                src={`https://www.twitch.tv/embed/${channel}/chat?parent=${parent}${darkParam}`}
                className="size-full"
                title="Twitch Chat"
            ></iframe>
        </div>
    );
};

export default TwitchChatEmbed;
