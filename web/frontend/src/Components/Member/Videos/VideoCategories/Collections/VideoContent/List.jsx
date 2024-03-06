import React, { useState, useEffect, useMemo } from 'react';
import { useRecoilState } from 'recoil';
import { useParams } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faArrowLeft, faDumbbell, faGauge, faMap, faPlay, faTable } from '@fortawesome/free-solid-svg-icons';
import ReactPlayer from 'react-player';

import Layout from '../../../../../Menu/Layout';
import { getVideoContentListAPI } from "../../../../../../API/VideoContent";
import FormErrorBox from "../../../../../Reusable/FormErrorBox";
import PageLoadingContent from "../../../../../Reusable/PageLoadingContent";
import { topAlertMessageState, topAlertStatusState } from "../../../../../../AppState";
import { EXERCISE_CATEGORY_OPTIONS, EXERCISE_MOMENT_TYPE_OPTIONS } from '../../../../../../Constants/FieldOptions';


const MemberVideoCollectionContentList = () => {

    ////
    //// URL Parameter.
    ////

    const { vcid, vconid } = useParams();

    console.log("Video collection id: " + vconid + ", Video Content Id: " + vcid)

    ////
    //// Global state.
    ////

    const [topAlertMessage, setTopAlertMessage] = useRecoilState(topAlertMessageState);
    const [topAlertStatus, setTopAlertStatus] = useRecoilState(topAlertStatusState);

    ////
    //// Component states.
    ////

    const [errors, setErrors] = useState({});
    const [videoList, setVideoList] = useState([]);
    const [isFetching, setFetching] = useState(false);
    // State for the selected video and handleThumbnailClick function
    const [selectedVideo, setSelectedVideo] = useState({
        videoUrl: '',
        videoName: ''
    });

    ////
    //// API.
    ////

    const onVideoCollectionContentListSuccess = (response) => {

        if (response && response.results) {
            // console.log(response.results)
            setVideoList(response.results)
        }

    };

    const onVideoCollectionContentListError = (apiErr) => {
        console.log(apiErr);
    };

    const onVideoCollectionContentListDone = () => {
        setFetching(false);
    };

    ////
    //// BREADCRUMB
    ////
    const breadcrumbItems = {
        items: [
            { text: 'Dashboard', link: '/dashboard', isActive: false, icon: faGauge },
            { text: 'Categories', link: '#', icon: faMap, isActive: true }
        ],
        mobileBackLinkItems: {
            link: '/dashboard',
            text: 'Back to Dashboard',
            icon: faArrowLeft
        }
    }

    ////
    //// Event handling.
    ////

    const handleThumbnailClick = (collection) => {
        setSelectedVideo({
            videoUrl: collection.videoUrl,
            videoName: collection.videoName
        });
    };

    const fetchList = () => {
        setFetching(true);
        setErrors({});

        let params = new Map();
        params.set("sort_field", "created"); // Sorting
        params.set("sort_order", -1)         // Sorting - descending, meaning most recent start date to oldest start date.

        params.set("video_collection_id", vconid);

        console.log("params:", params);

        getVideoContentListAPI(
            params,
            onVideoCollectionContentListSuccess,
            onVideoCollectionContentListError,
            onVideoCollectionContentListDone
        );
    };

    ////
    //// Misc.
    ////

    console.log(videoList)

    useEffect(() => {
        let mounted = true;

        if (mounted) {
            window.scrollTo(0, 0);
            fetchList();
        }

        return () => {
            mounted = false;
        };
    }, [vcid]);

    ////
    //// Custom Component
    ////

    const VideoComponent = ({ videoList, handleThumbnailClick }) => {
        const [currentVideo, setCurrentVideo] = useState(videoList.find(video => video.no === 1));

        const handleVideoClick = (video) => {
            setCurrentVideo(video);
        };

        return (
            <>

                <div className="columns is-multiline">
                    <div className="hero-body p-4">
                        <div className="columns is-multiline mb-0">
                            <div className="column">
                                <div className="is-flex-desktop is-justify-content-space-between">
                                    <div className="">
                                        <h5 className="is-size-7 has-text-weight-semibold">COLLECTION</h5>
                                        <h5 className="is-size-4 has-text-weight-bold">{currentVideo.categoryName}</h5>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div className="columns is-multiline mb-0">
                    <div className="column mb-0 is-12">
                        <ReactPlayer
                            url={currentVideo.videoUrl}
                            controls={true}
                            playing={true} // Auto play
                            width="100%"
                            height="400px"
                        />
                    </div>
                </div>
                <div className="columns is-multiline mb-0">
                    <div className="column mb-0 is-full">
                        <div className="has-background-black border-radius"></div>
                        <h5 className="is-size-5 has-text-weight-semibold mt-5">{currentVideo.name || "No title found"}</h5>
                    </div>
                </div>

                {/* Display thumbnails for other videos */}
                <div className="columns is-multiline mb-0 mt-5">
                    <div className="column mb-0">
                        <div className="is-flex is-justify-content-space-between">
                            <h3 className="is-size-6 has-text-weight-semibold has-text-grey-light">videos</h3>
                        </div>
                    </div>
                </div>
                <hr className="mt-0 mb-1" />
                <div className="columns is-multiline mb-0">
                    {videoList.map((video) => (
                        <div
                            className="column mb-0 is-3"
                            key={video.id}
                            onClick={() => handleVideoClick(video)}
                            style={{ cursor: 'pointer' }}
                        >
                            <div className={`border-radius ${video === currentVideo ? 'playing-video' : ''}`}>
                                <div>
                                    {video === currentVideo && <FontAwesomeIcon className='icon' icon={faPlay} />}
                                    <img
                                        className="border-radius"
                                        src={video.thumbnailUrl}
                                        alt={video.name}
                                        style={{ width: '100%', height: '200px', objectFit: 'cover' }}
                                    />
                                </div>
                                {/* <h5 className="is-size-5 has-text-weight-semibold">{video.name}</h5> */}
                            </div>
                        </div>
                    ))}
                </div>
            </>
        );
    };


    ////
    //// Component rendering.
    ////

    return (
        <Layout breadcrumbItems={breadcrumbItems}>
            <div className="box">
                <div className="container is-fluid p-0">
                    {isFetching ? (
                        <PageLoadingContent displayMessage={"Please wait..."} />
                    ) : (
                        <>
                            <FormErrorBox errors={errors} />
                            {videoList.length > 0 ? (
                                <>
                                    <VideoComponent videoList={videoList} />
                                </>
                            ) : (
                                <section className="hero is-medium has-background-white-ter">
                                    <div className="hero-body">
                                        <p className="title">
                                            <FontAwesomeIcon className="fas" icon={faTable} />
                                            &nbsp;No Videos
                                        </p>
                                        <p className="subtitle">
                                            No videos found at the moment. Please check back later!
                                        </p>
                                    </div>
                                </section>
                            )}
                        </>

                    )}
                </div>
            </div>
        </Layout>
    );
};

export default MemberVideoCollectionContentList;
