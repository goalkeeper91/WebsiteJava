import { useState } from 'react';
import YoutubeGrid from './YoutubeGrid';
import TiktokGrid from './TiktokGrid';

const VideoTabs = () => {
    const [activeTab, setActiveTab] = useState<'youtube' | 'tiktok'>('youtube');

    return (
        <section className='relative w-full mx-auto py-12 text-white overflow-hidden bg-slate-950'>
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