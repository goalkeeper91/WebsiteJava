import { useState } from 'react';
import YoutubeGrid from './YoutubeGrid';
import TiktokGrid from './TiktokGrid';

const VideoTabs = () => {
    const [activeTab, setActiveTab] = useState<'youtube' | 'tiktok'>('youtube');

    return (
        <section className='relative max-w-6xl mx-auto py-12 px-4 text-white overflow-hidden'>
            <div className='flex justify-center space-x-4 mb-8'>
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