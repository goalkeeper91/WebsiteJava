import { useState, useEffect } from 'react'
import { tiktokVideos } from '../../hooks/tiktokVideos';

const TiktokGrid = () => {
    const [featuredVideo, setFeaturedVideo] = useState<typeof tiktokVideos[0] | null>(null);

    useEffect(() => {
        if (tiktokVideos.length > 0) {
            const random = tiktokVideos[Math.floor(Math.random() * tiktokVideos.length)];
            setFeaturedVideo(random);
        }
    }, [tiktokVideos]);

    return (
        <section className='relative max-w-6xl mx-auto py-2 px-4 text-white z-20'>
            <h1 className="text-4xl font-bold mb-8 text-center">Video Highlights</h1>

            {featuredVideo && (
                <div className="mb-12 mx-auto">

                    <h2 className="text-2xl font-semibold mb-4 text-center">Empfohlenes Video</h2>
                    <div className="aspect-video place-self-center">
                        <iframe
                            src={featuredVideo.embedUrl}
                            className="w-82 h-190 rounded-lg bg-gray-600"
                            allowFullScreen
                            loading="lazy"
                        ></iframe>
                    </div>
                </div>
            )}
            <h2 className='text-3xl font-bold mb-8 text-center'>TikTok Highlights</h2>

            <div className='grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-8'>
                {tiktokVideos.map((video) => (
                    <div key={video.id} className="aspect-[9/16]">
                        <iframe
                          src={video.embedUrl}
                          title={video.title}
                          className="w-full h-full rounded-lg"
                          allowFullScreen
                          loading="lazy"
                        />
                    </div>
                ))}
            </div>
        </section>
    );
};

export default TiktokGrid;