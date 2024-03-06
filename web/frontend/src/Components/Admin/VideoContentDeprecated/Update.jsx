import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faImage, faVideo, faRepeat, faTasks, faTachometer, faPlus, faArrowLeft, faCheckCircle, faUserCircle, faGauge, faPencil, faVideoCamera, faEye, faIdCard, faAddressBook, faContactCard, faChartPie, faCogs } from '@fortawesome/free-solid-svg-icons'
import { useRecoilState } from 'recoil';
import { useParams } from 'react-router-dom';
import Vimeo from '@u-wave/react-vimeo';

import { postVideoCollectionCreateAPI, getVideoCollectionDetailAPI } from "../../../API/VideoCollection";
import FormErrorBox from "../../Reusable/FormErrorBox";
import FormAttachmentField from "../../Reusable/FormAttachmentField";
import FormInputField from "../../Reusable/FormInputField";
import FormTextareaField from "../../Reusable/FormTextareaField";
import FormRadioField from "../../Reusable/FormRadioField";
import FormMultiSelectField from "../../Reusable/FormMultiSelectField";
import FormSelectField from "../../Reusable/FormSelectField";
import FormCheckboxField from "../../Reusable/FormCheckboxField";
import FormCountryField from "../../Reusable/FormCountryField";
import FormRegionField from "../../Reusable/FormRegionField";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import { topAlertMessageState, topAlertStatusState } from "../../../AppState";
import {
    EXERCISE_MOMENT_TYPE_OPTIONS_WITH_EMPTY_OPTION,
    EXERCISE_CATEGORY_OPTIONS_WITH_EMPTY_OPTION,
    EXERCISE_TYPE_WITH_EMPTY_OPTIONS,
    EXERCISE_STATUS_OPTIONS_WITH_EMPTY_OPTION,
    EXERCISE_GENDER_OPTIONS_WITH_EMPTY_OPTION
} from "../../../Constants/FieldOptions";
import {
    EXERCISE_VIDEO_TYPE_SIMPLE_STORAGE_SERVER,
    EXERCISE_VIDEO_TYPE_YOUTUBE,
    EXERCISE_VIDEO_TYPE_VIMEO,
    EXERCISE_THUMBNAIL_TYPE_SIMPLE_STORAGE_SERVER,
    EXERCISE_THUMBNAIL_TYPE_EXTERNAL_URL
} from "../../../Constants/App";


function AdminVideoCollectionUpdate() {
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
    const [videoType, setVideoType] = useState(0);
    const [videoURL, setVideoURL] = useState("");
    const [videoAttachmentID, setVideoAttachmentID] = useState("");
    const [videoAttachmentFilename, setVideoAttachmentFilename] = useState("");
    const [thumbnailType, setThumbnailType] = useState(0);
    const [thumbnailURL, setThumbnailURL] = useState("");
    const [thumbnailAttachmentID, setThumbnailAttachmentID] = useState("");
    const [thumbnailAttachmentFilename, setThumbnailAttachmentFilename] = useState("");
    const [description, setDescription] = useState("");
    const [alternateName, setAlternateName] = useState("");
    const [name, setName] = useState("");
    const [gender, setGender] = useState("");
    const [movementType, setMovementType] = useState(0);
    const [category, setCategory] = useState("");
    const [type, setType] = useState(1); // 1=System

    ////
    //// Event handling.
    ////

    const onSubmitClick = (e) => {
        console.log("onSubmitClick: Starting...")
        setFetching(true);
        setErrors({});
        postVideoCollectionCreateAPI(
            {
                type: type,
                video_type: videoType,
                video_upload: videoAttachmentID,
                video_url: videoURL,
                thumbnail_type: thumbnailType,
                thumbnail_upload: thumbnailAttachmentID,
                thumbnail_url: thumbnailURL,
                name: name,
                alternate_name: alternateName,
                category: category,
                description: description,
                gender: gender,
                movement_type: movementType,
            },
            onUpdateSuccess,
            onUpdateError,
            onUpdateDone
        );
        console.log("onSubmitClick: Finished.")
    }

    ////
    //// API.
    ////

    // --- VideoCollection Update --- //

    function onUpdateSuccess(response){
        // For debugging purposes only.
        console.log("onUpdateSuccess: Starting...");
        console.log(response);

        // Update a temporary banner message in the app and then clear itself after 2 seconds.
        setTopAlertMessage("VideoCollection updated");
        setTopAlertStatus("success");
        setTimeout(() => {
            console.log("onUpdateSuccess: Delayed for 2 seconds.");
            console.log("onUpdateSuccess: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
            setTopAlertMessage("");
        }, 2000);

        // Redirect the organization to the organization attachments page.
        setForceURL("/admin/video-collection/"+response.id+"");
    }

    function onUpdateError(apiErr) {
        console.log("onUpdateError: Starting...");
        setErrors(apiErr);

        // Add a temporary banner message in the app and then clear itself after 2 seconds.
        setTopAlertMessage("Failed submitting");
        setTopAlertStatus("danger");
        setTimeout(() => {
            console.log("onUpdateError: Delayed for 2 seconds.");
            console.log("onUpdateError: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
            setTopAlertMessage("");
        }, 2000);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onUpdateDone() {
        console.log("onUpdateDone: Starting...");
        setFetching(false);
    }

    // --- VideoCollection Detail --- //

    function onDetailSuccess(response){
        // For debugging purposes only.
        console.log("onDetailSuccess: Starting...");
        console.log(response);

        // DEVELOPERS NOTE: If zero string then that's a MongoDB empty db record so handle the code as follows.
        if (response.videoAttachmentId === "000000000000000000000000") {
            response.videoAttachmentId = "";
        } else {
            setVideoAttachmentID(response.videoAttachmentId);
        }
        if (response.thumbnailAttachmentId === "000000000000000000000000") {
            response.thumbnailAttachmentId = "";
        } else {
            setThumbnailAttachmentID(response.thumbnailAttachmentId);
        }
        setVideoType(response.videoType);
        setVideoURL(response.videoURL);

        setVideoAttachmentFilename(response.videoAttachmentFilename);
        setThumbnailType(response.thumbnailType);
        setThumbnailURL(response.thumbnailURL);
        setThumbnailAttachmentFilename(response.thumbnailAttachmentFilename);
        setDescription(response.description);
        setAlternateName(response.alternateName);
        setName(response.name);
        setGender(response.gender);
        setMovementType(response.movementType);
        setCategory(response.category);
        setType(response.type);
    }

    function onDetailError(apiErr) {
        console.log("onDetailError: Starting...");
        setErrors(apiErr);

        // Detail a temporary banner message in the app and then clear itself after 2 seconds.
        setTopAlertMessage("Failed submitting");
        setTopAlertStatus("danger");
        setTimeout(() => {
            console.log("onDetailError: Delayed for 2 seconds.");
            console.log("onDetailError: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
            setTopAlertMessage("");
        }, 2000);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onDetailDone() {
        console.log("onDetailDone: Starting...");
        setFetching(false);
    }

    ////
    //// Misc.
    ////

    useEffect(() => {
        let mounted = true;

        if (mounted) {
            window.scrollTo(0, 0);  // Start the page at the top of the page.
            setFetching(true);
            getVideoCollectionDetailAPI(
                id,
                onDetailSuccess,
                onDetailError,
                onDetailDone
            );
        }

        return () => { mounted = false; }
    }, [id]);

    ////
    //// Component rendering.
    ////

    if (forceURL !== "") {
        return <Navigate to={forceURL}  />
    }

    const isVideoUploaded = (videoAttachmentID !== "" || videoURL !== "");
    const isThumbnailUploaded = (thumbnailAttachmentID !== "" || thumbnailURL !== "");

    return (
        <>
            <div class="container">
                <section class="section">

                    {/* Desktop Breadcrumbs */}
                    <nav class="breadcrumb is-hidden-touch" aria-label="breadcrumbs">
                        <ul>
                            <li class=""><Link to="/admin/dashboard" aria-current="page"><FontAwesomeIcon className="fas" icon={faGauge} />&nbsp;Dashboard</Link></li>
                            <li class=""><Link to="/admin/video-collections" aria-current="page"><FontAwesomeIcon className="fas" icon={faVideoCamera} />&nbsp;VideoCollections</Link></li>
                            <li class=""><Link to={`/admin/video-collection/${id}`} aria-current="page"><FontAwesomeIcon className="fas" icon={faEye} />&nbsp;Detail</Link></li>
                            <li class="is-active"><Link aria-current="page"><FontAwesomeIcon className="fas" icon={faPencil} />&nbsp;Edit</Link></li>
                        </ul>
                    </nav>

                    {/* Mobile Breadcrumbs */}
                    <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
                        <ul>
                            <li class=""><Link to={`/admin/video-collection/${id}`} aria-current="page"><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to Detail</Link></li>
                        </ul>
                    </nav>

                    {/* Modal */}
                    {/* Nothing ... */}

                    {/* Page */}
                    <nav class="box">
                        <div class="columns">
                            <div class="column">
                                <p class="title is-4"><FontAwesomeIcon className="fas" icon={faVideoCamera} />&nbsp;VideoCollection</p>
                            </div>
                            <div class="column has-text-right">
                            </div>
                        </div>
                        <FormErrorBox errors={errors} />

                        <p class="pb-4 mb-5 has-text-grey">Please fill out all the required fields before submitting this form.</p>

                        {isFetching
                            ?
                            <PageLoadingContent displayMessage={"Please wait..."} />
                            :
                            <>
                                <div class="container">

                                    {/*
                                        ------------------------
                                        VIDEO UPLOAD SECTION
                                        ------------------------
                                    */}
                                    <p class="subtitle is-6"><FontAwesomeIcon className="fas" icon={faVideo} />&nbsp;Video</p>
                                    <hr />

                                    <FormRadioField
                                        label="Video Type"
                                        name="videoType"
                                        placeholder="Pick"
                                        value={videoType}
                                        opt1Value={EXERCISE_VIDEO_TYPE_SIMPLE_STORAGE_SERVER}
                                        opt1Label="File Upload"
                                        opt2Value={EXERCISE_VIDEO_TYPE_YOUTUBE}
                                        opt2Label="YouTube"
                                        opt3Value={EXERCISE_VIDEO_TYPE_VIMEO}
                                        opt3Label="Vimeo"
                                        errorText={errors && errors.videoType}
                                        onChange={(e)=>setVideoType(parseInt(e.target.value))}
                                        maxWidth="180px"
                                        disabled={false}
                                    />

                                    {(() => {
                                        switch (videoType) {
                                            case EXERCISE_VIDEO_TYPE_SIMPLE_STORAGE_SERVER: return (
                                                <>
                                                    <FormAttachmentField
                                                        label="Video Upload"
                                                        name="videoUpload"
                                                        placeholder="Upload file"
                                                        errorText={errors && errors.videoUpload}
                                                        attachmentID={videoAttachmentID}
                                                        setAttachmentID={setVideoAttachmentID}
                                                        attachmentFilename={videoAttachmentFilename}
                                                        setAttachmentFilename={setVideoAttachmentFilename}
                                                    />

                                                </>
                                            );
                                            case EXERCISE_VIDEO_TYPE_YOUTUBE: return (
                                                <FormInputField
                                                    label="YouTube URL"
                                                    name="videoExternalURL"
                                                    placeholder="URL input"
                                                    value={videoURL}
                                                    errorText={errors && errors.videoUrl}
                                                    helpText=""
                                                    onChange={(e)=>setVideoURL(e.target.value)}
                                                    isRequired={true}
                                                    maxWidth="380px"
                                                />
                                            );
                                            case EXERCISE_VIDEO_TYPE_VIMEO: return (
                                                <FormInputField
                                                    label="Vimeo URL"
                                                    name="videoExternalURL"
                                                    placeholder="URL input"
                                                    value={videoURL}
                                                    errorText={errors && errors.videoUrl}
                                                    helpText=""
                                                    onChange={(e)=>setVideoURL(e.target.value)}
                                                    isRequired={true}
                                                    maxWidth="380px"
                                                />
                                            );
                                            default: return null;
                                        }
                                    })()}

                                    {/*
                                        ------------------------
                                        THUMBNAIL UPLOAD SECTION
                                        ------------------------
                                    */}
                                    <p class="subtitle is-6"><FontAwesomeIcon className="fas" icon={faImage} />&nbsp;Thumbnail</p>
                                    <hr />

                                    <FormRadioField
                                        label="Thumbnail Type"
                                        name="thumbnailType"
                                        placeholder="Pick"
                                        value={thumbnailType}
                                        opt1Value={EXERCISE_THUMBNAIL_TYPE_SIMPLE_STORAGE_SERVER}
                                        opt1Label="File Upload"
                                        opt2Value={EXERCISE_THUMBNAIL_TYPE_EXTERNAL_URL}
                                        opt2Label="External URL"
                                        errorText={errors && errors.thumbnailType}
                                        onChange={(e)=>setThumbnailType(parseInt(e.target.value))}
                                        maxWidth="180px"
                                        disabled={false}
                                    />

                                    {(() => {
                                        switch (thumbnailType) {
                                            case EXERCISE_THUMBNAIL_TYPE_SIMPLE_STORAGE_SERVER: return (
                                                <>
                                                    <FormAttachmentField
                                                        label="File Upload"
                                                        name="thumbnaiUpload"
                                                        placeholder="Upload file"
                                                        errorText={errors && errors.thumbnailUpload}
                                                        attachmentID={thumbnailAttachmentID}
                                                        setAttachmentID={setThumbnailAttachmentID}
                                                        attachmentFilename={thumbnailAttachmentFilename}
                                                        setAttachmentFilename={setThumbnailAttachmentFilename}
                                                    />

                                                </>
                                            );
                                            case EXERCISE_THUMBNAIL_TYPE_EXTERNAL_URL: return (
                                                <FormInputField
                                                    label="Thumbnail External URL"
                                                    name="thumbnailUrl"
                                                    placeholder="URL input"
                                                    value={thumbnailURL}
                                                    errorText={errors && errors.thumbnailUrl}
                                                    helpText=""
                                                    onChange={(e)=>setThumbnailURL(e.target.value)}
                                                    isRequired={true}
                                                    maxWidth="380px"
                                                />
                                            );
                                            default: return null;
                                        }
                                    })()}


                                    {/*
                                        ------------------------
                                       INFORMATION  SECTION
                                        ------------------------
                                    */}
                                    {isVideoUploaded && isThumbnailUploaded && <>
                                        <p class="subtitle is-6 mt-5"><FontAwesomeIcon className="fas" icon={faEye} />&nbsp;Information</p>
                                        <hr />

                                        <FormInputField
                                            label="Name"
                                            name="name"
                                            placeholder="Name input"
                                            value={name}
                                            errorText={errors && errors.name}
                                            helpText=""
                                            onChange={(e)=>setName(e.target.value)}
                                            isRequired={true}
                                            maxWidth="380px"
                                        />

                                        <FormInputField
                                            label="Alternate Name"
                                            name="alternateName"
                                            placeholder="Name input"
                                            value={alternateName}
                                            errorText={errors && errors.alternateName}
                                            helpText=""
                                            onChange={(e)=>setAlternateName(e.target.value)}
                                            isRequired={true}
                                            maxWidth="380px"
                                        />

                                        <FormTextareaField
                                            label="Description"
                                            name="description"
                                            placeholder="Description input"
                                            value={description}
                                            errorText={errors && errors.description}
                                            helpText=""
                                            onChange={(e)=>setDescription(e.target.value)}
                                            isRequired={true}
                                            maxWidth="380px"
                                        />

                                        <FormSelectField
                                            label="Gender"
                                            name="gender"
                                            placeholder="Pick"
                                            selectedValue={gender}
                                            errorText={errors && errors.gender}
                                            helpText=""
                                            onChange={(e) => setGender(e.target.value)}
                                            options={EXERCISE_GENDER_OPTIONS_WITH_EMPTY_OPTION}
                                        />

                                        <FormSelectField
                                            label="Movement Type"
                                            name="movementType"
                                            placeholder="Pick"
                                            selectedValue={movementType}
                                            errorText={errors && errors.movementType}
                                            helpText=""
                                            onChange={(e) => setMovementType(parseInt(e.target.value))}
                                            options={EXERCISE_MOMENT_TYPE_OPTIONS_WITH_EMPTY_OPTION}
                                        />

                                        <FormSelectField
                                            label="Category"
                                            name="category"
                                            placeholder="Pick"
                                            selectedValue={category}
                                            errorText={errors && errors.category}
                                            helpText=""
                                            onChange={(e) => setCategory(parseInt(e.target.value))}
                                            options={EXERCISE_CATEGORY_OPTIONS_WITH_EMPTY_OPTION}
                                        />

                                        <FormSelectField
                                            label="Type"
                                            name="type"
                                            placeholder="Pick"
                                            selectedValue={type}
                                            errorText={errors && errors.type}
                                            helpText=""
                                            onChange={(e) => setType(e.target.value)}
                                            options={EXERCISE_TYPE_WITH_EMPTY_OPTIONS}
                                            disabled={true}
                                        />
                                    </>}



                                    <div class="columns pt-5">
                                        <div class="column is-half">
                                            <Link class="button is-fullwidth-mobile" to={`/admin/video-collection/${id}`}><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to Detail</Link>
                                        </div>
                                        <div class="column is-half has-text-right">
                                            <button onClick={onSubmitClick} class="button is-success is-fullwidth-mobile" type="button"><FontAwesomeIcon className="fas" icon={faPlus}/>&nbsp;Submit</button>
                                        </div>
                                    </div>

                                </div>
                            </>
                        }
                    </nav>
                </section>
            </div>
        </>
    );
}

export default AdminVideoCollectionUpdate;
