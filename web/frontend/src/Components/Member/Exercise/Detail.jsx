import React, { useState, useEffect } from 'react';
import Scroll from 'react-scroll';
import { useRecoilState } from 'recoil';
import { useParams, Link, Navigate } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faArrowLeft, faDumbbell, faUser, faHeart, faGauge, faCog, faPlusCircle, faSearch } from '@fortawesome/free-solid-svg-icons';
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
} from "../../../Constants/App";

////
//// Custom Component
////

const VideoPlayer = ({ name, videoType, videoUrl, videoObjectUrl, thumbnailType, thumbnailUrl, thumbnailObjectUrl }) => {
    const [showVideo, setShowVideo] = useState(false);

    return (
        <>
            {showVideo ? (
                <div className="video-container is-16by9">
                    {(() => {
                        switch (videoType) {
                            case EXERCISE_VIDEO_TYPE_SIMPLE_STORAGE_SERVER:
                                return <iframe src={videoObjectUrl} allowFullScreen></iframe>;
                            case EXERCISE_VIDEO_TYPE_YOUTUBE:
                                return <iframe src={`https://www.youtube.com/embed/${videoUrl}`} allowFullScreen></iframe>;
                            case EXERCISE_VIDEO_TYPE_VIMEO:
                                return <Vimeo video={videoUrl} autoplay={true} />;
                            default:
                                return null;
                        }
                    })()}
                </div>
            ) : (
                <div className="thumbnail-container" onClick={() => setShowVideo(true)}>
                    <img className="border-radius image is-256x256"
                        src={thumbnailType === EXERCISE_THUMBNAIL_TYPE_SIMPLE_STORAGE_SERVER ? thumbnailObjectUrl : thumbnailUrl}
                        alt="Thumbnail"
                        style={{ cursor: 'pointer' }}
                    />
                    <div className="play-button"></div>
                </div>
            )}
            <h4 class="is-size-5 has-text-weight-bold mb-2 has-text-black youtube-title my-4">{name}</h4>
            <div class="tags is-size-5">
                <span class="tag">
                    <FontAwesomeIcon className="fas" icon={faUser} />&nbsp;User
                </span>
                <span class="tag">
                    <FontAwesomeIcon className="fas" icon={faHeart} />&nbsp;Favorites
                </span>
                <span class="tag">
                    <FontAwesomeIcon className="fas" icon={faCog} />&nbsp;Settings
                </span>
            </div>
        </>
    );
};

const VideoDescription = ({ description, gender, category, type, movementType }) => (
    <div>
        <h4 class="is-size-6 has-text-weight-bold mb-1">DESCRIPTION</h4>
        <p>{description}</p>
        <h4 class="is-size-6 has-text-weight-bold mb-1 mt-4">CATEGORY</h4>
        <p>{EXERCISE_CATEGORY_OPTIONS.find(item => item.value === category)?.label}</p>
        <h4 class="is-size-6 has-text-weight-bold mb-1 mt-4">MOVEMENT TYPE</h4>
        <p>{EXERCISE_MOMENT_TYPE_OPTIONS.find(item => item.value === movementType)?.label}</p>
        <h4 class="is-size-6 has-text-weight-bold mb-1 mt-4">GENDER</h4>
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

    ////
    //// BREADCRUMB
    ////
    const breadcrumbItems = {
        items: [
            { text: 'Dashboard', link: '/dashboard', isActive: false, icon: faGauge },
            { text: 'Exercises', link: '#', icon: faDumbbell, isActive: true }
        ],
        mobileBackLinkItems: {
            link: "/exercises",
            text: "Back to Exercises",
            icon: faArrowLeft
        }
    }

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

    console.log(datum)

    return (
        <Layout breadcrumbItems={breadcrumbItems}>
            <div className="box">
                {/* <div class="columns is-multiline border-bottom">
                    <div class="column">
                        <h4 class="is-size-4 has-text-weight-bold mb-2">{datum.name}</h4>
                    </div>
                </div> */}
                <div class="columns pt-5 pb-4 border-bottom">
                    <div className="column">
                        <VideoPlayer
                            name={datum.name}
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
