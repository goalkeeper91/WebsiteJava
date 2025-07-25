import { useState } from 'react';
import YoutubeGrid from './YoutubeGrid';
import TiktokGrid from './TiktokGrid';

const VideoTabs = () => {
    const [activeTab, setActiveTab] = useState<'youtube' | 'tiktok'>('youtube');

    return (
        <section className='relative max-w-6xl mx-auto py-12 px-4 text-white overflow-hidden'>
            <div className="absolute inset-0 z-10">
                <div className="w-full h-full bg-gradient-to-br from-purple-500 via-pink-500 to-yellow-500 opacity-30 blur-3xl" />
            </div>
            <div className='relative z-15 flex justify-center space-x-4 mb-8'>
                <button
                    onClick={() => setActiveTab("youtube")}
                    className={`px-6 py-2 rounded-full font-semibold transition ${
                    activeTab === "youtube" ? "bg-goalyBlue text-white" : "bg-gray-700 hover:bg-gray-600"
                    }`}
                >
                    Youtube
                </button>
                <button
                    onClick={() => setActiveTab("tiktok")}
                    className={`px-6 py-2 rounded-full font-semibold transition ${
                    activeTab === "tiktok" ? "bg-goalyBlue text-white" : "bg-gray-700 hover:bg-gray-600"
                    }`}
                >
                    TikTok
                </button>
            </div>

            {activeTab === "youtube" ? <YoutubeGrid /> : <TiktokGrid />}
        </section>
    );
};

export default VideoTabs;