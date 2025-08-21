import React from 'react';

interface TwitchEmbedProp {
    channel: string;
}

const TwitchEmbed: React.FC<TwitchEmbedProp> = ({ channel }) => {
    return (
        <div className="aspect-video w-full rounded-xl overflow-hidden shadow-lg">
            <iframe
                src={`https://player.twitch.tv/?channel=${channel}&parent=localhost&muted=true`}
                allowFullScreen
                allow="autoplay; fullscreen"
                className="size-full"
                title="Twitch Stream"
            ></iframe>
        </div>
    );
};

export default TwitchEmbed;