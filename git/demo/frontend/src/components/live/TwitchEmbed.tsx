import React from 'react';

interface TwitchEmbedProp {
    channel: string;
}

const TwitchEmbed: React.FC<TwitchEmbedProp> = ({ channel }) => {
    // parent must match the domain actually serving this page (Twitch's
    // embed gate checks it) - computed at render time so this works for
    // localhost dev and the real production domain without an env var.
    const parent = window.location.hostname;

    return (
        <div className="aspect-video w-full rounded-xl overflow-hidden shadow-lg">
            <iframe
                data-cookieconsent="marketing"
                src={`https://player.twitch.tv/?channel=${channel}&parent=${parent}&muted=true`}
                allowFullScreen
                allow="autoplay; fullscreen"
                className="size-full"
                title="Twitch Stream"
            ></iframe>
        </div>
    );
};

export default TwitchEmbed;