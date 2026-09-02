import Seo from '../components/Seo';
import Hero from '../components/hero/Hero';
import CommunitySection from '../components/community/CommunitySection';

const Home = () => {
    return (
        <>
            <Seo
                title="Twitch Bot SaaS, Clip-Automatisierung & Live Entertainment"
                description="Goalkeeper91: Twitch-Chatbot mit Automod, Loyalty-Punkten, Giveaways und automatischer Clip-Erstellung als SaaS - plus Livestream und Community."
                path="/"
                structuredData={{
                    "@context": "https://schema.org",
                    "@type": "SoftwareApplication",
                    name: "Goalkeeper91 Twitch Bot",
                    applicationCategory: "BusinessApplication",
                    operatingSystem: "Web",
                    description:
                        "Twitch-Chatbot als SaaS mit Automod, Loyalty-Punkten, Giveaways, Chat-Commands und automatischer Clip-Erstellung.",
                    url: "https://goalkeeper91.de/pricing",
                    featureList: [
                        "Automod",
                        "Loyalty-Punkte",
                        "Giveaways",
                        "Chat-Commands",
                        "Discord-Integration",
                        "Clip-Automatisierung",
                    ],
                }}
            />
            <Hero />
            <CommunitySection />
        </>
        );
    };

export default Home;
