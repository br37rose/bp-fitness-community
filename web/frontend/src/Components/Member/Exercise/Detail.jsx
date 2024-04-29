import React, { useState, useEffect } from 'react';
import Scroll from 'react-scroll';
import { useRecoilState } from 'recoil';
import { useParams, Link, Navigate } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faArrowLeft, faDumbbell, faEye, faFilter, faGauge, faPlus, faPlusCircle, faSearch } from '@fortawesome/free-solid-svg-icons';
import Vimeo from '@u-wave/react-vimeo';

import Layout from '../../Menu/Layout';
import { getExerciseDetailAPI } from "../../../API/Exercise";
import FormErrorBox from "../../Reusable/FormErrorBox";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import { topAlertMessageState, topAlertStatusState } from "../../../AppState";
import {
    EXERCISE_CATEGORY,
    EXERCISE_CATEGORY_OPTIONS,
    EXERCISE_MOMENT_TYPE_OPTIONS,
    EXERCISE_TYPE_SYSTEM
} from '../../../Constants/FieldOptions';
import {
    EXERCISE_THUMBNAIL_TYPE_SIMPLE_STORAGE_SERVER,
    EXERCISE_THUMBNAIL_TYPE_EXTERNAL_URL,
    EXERCISE_VIDEO_TYPE_SIMPLE_STORAGE_SERVER,
    EXERCISE_VIDEO_TYPE_YOUTUBE,
    EXERCISE_VIDEO_TYPE_VIMEO,
}  from "../../../Constants/App";

////
//// Custom Component
////

const VideoPlayer = ({ videoType, videoUrl, videoObjectUrl, thumbnailType, thumbnailUrl, thumbnailObjectUrl }) => {
    const [showVideo, setShowVideo] = useState(false);
    return (
        <>
            {showVideo ? (
                <>
                {(() => {
                    switch (videoType) {
                        case EXERCISE_VIDEO_TYPE_SIMPLE_STORAGE_SERVER: return (
                            <>
                                <video style={{width:'100%', height:'100%'}} controls 
                                >
                                    <source src={videoObjectUrl}
                                        type="video/mp4"
                                    />
                                </video>
                            </>
                        );
                        case EXERCISE_VIDEO_TYPE_YOUTUBE: return (
                            <>yt</>
                        );
                        case EXERCISE_VIDEO_TYPE_VIMEO: return (
                            <Vimeo
                              video={`${videoUrl}`}
                              autoplay
                            />
                        );
                        default: return null;
                    }
                })()}
                </>
            ) : (
                <div className="video-thumbnail-container" onClick={() => setShowVideo(true)}>
                    {thumbnailType === EXERCISE_THUMBNAIL_TYPE_SIMPLE_STORAGE_SERVER && <img className="border-radius image is-256x256" src={thumbnailObjectUrl} alt={`Thumbnail`} style={{ cursor: 'pointer' }} />}
                    {thumbnailType === EXERCISE_THUMBNAIL_TYPE_EXTERNAL_URL && <img className="border-radius image is-256x256" src={thumbnailUrl} alt={`Thumbnail`} style={{ cursor: 'pointer' }} />}
                    <div className="play-button"></div>
                </div>
            )}
        </>
    );
};

const VideoDescription = ({ description, gender, category, type, movementType }) => (
    <div>
        <h4 class="is-size-6 has-text-grey-light has-text-weight-bold mb-1">DESCRIPTION</h4>
        <p>{description}</p>
        <h4 class="is-size-6 has-text-grey-light has-text-weight-bold mb-1 mt-4">CATEGORY</h4>
        <p>{EXERCISE_CATEGORY_OPTIONS.find(item => item.value === category)?.label}</p>
        <h4 class="is-size-6 has-text-grey-light has-text-weight-bold mb-1 mt-4">MOVEMENT TYPE</h4>
        <p>{EXERCISE_MOMENT_TYPE_OPTIONS.find(item => item.value === movementType)?.label}</p>
        <h4 class="is-size-6 has-text-grey-light has-text-weight-bold mb-1 mt-4">GENDER</h4>
        <p>{gender}</p>
    </div>

);

const MemberExerciseDetail = () => {

    ////
    //// URL Parameters.
    ////

    const { id } = useParams()

    ////
    //// Global state.
    ////

    const [topAlertMessage, setTopAlertMessage] = useRecoilState(topAlertMessageState);
    const [topAlertStatus, setTopAlertStatus] = useRecoilState(topAlertStatusState);

    ////
    //// Component states.
    ////

    const [errors, setErrors] = useState({});
    const [isFetching, setFetching] = useState(false);
    const [forceURL, setForceURL] = useState("");
    const [datum, setDatum] = useState({});
    const [tabIndex, setTabIndex] = useState(1);
    const [selectedExerciseForDeletion, setSelectedExerciseForDeletion] = useState(null);

    const breadcrumbItems = [
        { text: 'Dashboard', link: '/dashboard', icon: faGauge },
        { text: 'Exercises', link: '/exercises', icon: faDumbbell },
        { text: 'Detail', link: '#', icon: faDumbbell, isActive: true }
    ];

    ////
    //// API.
    ////

    // --- Detail --- //

    function onExerciseDetailSuccess(response) {
        console.log("onExerciseDetailSuccess: Starting...");
        setDatum(response);
    }

    function onExerciseDetailError(apiErr) {
        console.log("onExerciseDetailError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onExerciseDetailDone() {
        console.log("onExerciseDetailDone: Starting...");
        setFetching(false);
    }

    ////
    //// Event handling.
    ////


    ////
    //// Misc.
    ////

    useEffect(() => {
        let mounted = true;

        if (mounted) {
            window.scrollTo(0, 0);  // Start the page at the top of the page.

            setFetching(true);
            getExerciseDetailAPI(
                id,
                onExerciseDetailSuccess,
                onExerciseDetailError,
                onExerciseDetailDone
            );
        }

        return () => { mounted = false; }
    }, []);

    ////
    //// Component rendering.
    ////

    if (forceURL !== "") {
        return <Navigate to={forceURL} />
    }

    return (
        <Layout breadcrumbItems={breadcrumbItems}>
            <div className="box">
                <div class="columns is-multiline border-bottom">
                    <div class="column">
                        <h4 class="is-size-4 has-text-weight-bold mb-2">{datum.name}</h4>
                    </div>
                </div>
                <div class="columns pt-5 pb-4 border-bottom">
                    <div className="column">
                        <VideoPlayer
                            videoType={datum.videoType}
                            videoUrl={datum.videoUrl}
                            videoObjectUrl={datum.videoObjectUrl}
                            thumbnailType={datum.thumbnailType}
                            thumbnailUrl={datum.thumbnailUrl}
                            thumbnailObjectUrl={datum.thumbnailObjectUrl}
                        />
                    </div>
                    <div className="column is-7 is-flex is-flex-wrap-wrap is-align-content-space-between">
                        <VideoDescription
                            name={datum.name}
                            description={datum.description}
                            gender={datum.gender}
                            category={datum.category}
                            type={datum.type}
                            movementType={datum.movementType}
                        />
                    </div>
                </div>
                <div className="columns pt-5 pb-4">
                    <div className="column">
                        <div className="column is-2 has-text-right p-0">
                            <Link to="/exercises" aria-current="page">
                                <button className="button is-fullwidth">
                                    <FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back
                                </button>
                            </Link>
                        </div>
                    </div>
                </div>
            </div>
        </Layout>
    );
}

export default MemberExerciseDetail;
