import React, { useState, useEffect } from "react";
import { Link, Navigate, useParams } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faImage, faVideo, faRepeat, faTasks, faTachometer, faPlus, faArrowLeft, faCheckCircle, faUserCircle, faGauge, faPencil, faVideoCamera, faEye, faIdCard, faAddressBook, faContactCard, faChartPie, faCogs } from '@fortawesome/free-solid-svg-icons'
import { useRecoilState } from 'recoil';
import Vimeo from '@u-wave/react-vimeo';

import { putVideoCollectionUpdateAPI, getVideoCollectionDetailAPI } from "../../../API/VideoCollection";
import FormErrorBox from "../../Reusable/FormErrorBox";
import FormAttachmentField from "../../Reusable/FormAttachmentField";
import FormInputField from "../../Reusable/FormInputField";
import FormTextareaField from "../../Reusable/FormTextareaField";
import FormRadioField from "../../Reusable/FormRadioField";
import FormMultiSelectField from "../../Reusable/FormMultiSelectField";
import FormSelectField from "../../Reusable/FormSelectField";
import FormCheckboxField from "../../Reusable/FormCheckboxField";
import FormSelectFieldForVideoCategory from "../../Reusable/FormSelectFieldForVideoCategory";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import { topAlertMessageState, topAlertStatusState } from "../../../AppState";
import {
    VIDEO_COLLECTION_TYPE_OPTIONS_WITH_EMPTY_OPTION
} from "../../../Constants/FieldOptions";
import {
    EXERCISE_THUMBNAIL_TYPE_SIMPLE_STORAGE_SERVER,
    EXERCISE_THUMBNAIL_TYPE_EXTERNAL_URL,
    VIDEO_COLLECTION_TYPE_MANY_VIDEOS
} from "../../../Constants/App";


function AdminVideoCollectionUpdate() {
    ////
    //// URL Parameters.
    ////

    const { vcid } = useParams()

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
    const [thumbnailType, setThumbnailType] = useState(0);
    const [thumbnailURL, setThumbnailURL] = useState("");
    const [thumbnailAttachmentID, setThumbnailAttachmentID] = useState("");
    const [thumbnailAttachmentName, setThumbnailAttachmentName] = useState("");
    const [alternateName, setAlternateName] = useState("");
    const [name, setName] = useState("");
    const [summary, setSummary] = useState("");
    const [description, setDescription] = useState("");
    const [type, setType] = useState(VIDEO_COLLECTION_TYPE_MANY_VIDEOS);
    const [videoCategoryID, setVideoCategoryID] = useState("");
    const [isVideoCategoryOther, setIsVideoCategoryOther] = useState("");

    ////
    //// Event handling.
    ////

    const onSubmitClick = (e) => {
        console.log("onSubmitClick: Starting...")
        setFetching(true);
        setErrors({});
        putVideoCollectionUpdateAPI(
            {
                id: vcid,
                type: type,
                thumbnail_type: thumbnailType,
                thumbnail_upload: thumbnailAttachmentID,
                thumbnail_url: thumbnailURL,
                name: name,
                summary: summary,
                description: description,
                video_category_id: videoCategoryID,
                is_video_category_other: isVideoCategoryOther,
                category_id: videoCategoryID,
                status: 1, //1=Active
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

    // --- Detail --- //

    function onVideoCollectionDetailSuccess(response){
        console.log("onVideoCollectionDetailSuccess: Starting...");

        setThumbnailType(response.thumbnailType);
        setThumbnailURL(response.thumbnailUrl);
        setThumbnailAttachmentID(response.thumbnailAttachmentId);
        setThumbnailAttachmentName(response.thumbnailAttachmentName);
        setName(response.name);
        setSummary(response.summary);
        setDescription(response.description);
        setType(response.type);
        setVideoCategoryID(response.categoryId);
        setIsVideoCategoryOther(response.isVideoCategoryOther);
    }

    function onVideoCollectionDetailError(apiErr) {
        console.log("onVideoCollectionDetailError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onVideoCollectionDetailDone() {
        console.log("onVideoCollectionDetailDone: Starting...");
        setFetching(false);
    }

    // --- Update --- //

    function onUpdateSuccess(response){
        // For debugging purposes only.
        console.log("onUpdateSuccess: Starting...");
        console.log(response);

        // Update a temporary banner message in the app and then clear itself after 2 seconds.
        setTopAlertMessage("Video collection update");
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

        // Update a temporary banner message in the app and then clear itself after 2 seconds.
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

    ////
    //// Misc.
    ////

    useEffect(() => {
        let mounted = true;

        if (mounted) {
            window.scrollTo(0, 0);  // Start the page at the top of the page.

            setFetching(true);
            getVideoCollectionDetailAPI(
                vcid,
                onVideoCollectionDetailSuccess,
                onVideoCollectionDetailError,
                onVideoCollectionDetailDone
            );
        }

        return () => { mounted = false; }
    }, []);

    ////
    //// Component rendering.
    ////

    if (forceURL !== "") {
        return <Navigate to={forceURL}  />
    }

    const isThumbnailUploaded = (thumbnailAttachmentID !== "" || thumbnailURL !== "");

    return (
        <>
            <div class="container">
                <section class="section">

                    {/* Desktop Breadcrumbs */}
                    <nav class="breadcrumb is-hidden-touch" aria-label="breadcrumbs">
                        <ul>
                            <li class=""><Link to="/admin/dashboard" aria-current="page"><FontAwesomeIcon className="fas" icon={faGauge} />&nbsp;Dashboard</Link></li>
                            <li class=""><Link to="/admin/video-collections" aria-current="page"><FontAwesomeIcon className="fas" icon={faVideoCamera} />&nbsp;Video Collections</Link></li>
                            <li class=""><Link to={`/admin/video-collection/${vcid}`} aria-current="page"><FontAwesomeIcon className="fas" icon={faEye} />&nbsp;Detail</Link></li>
                            <li class="is-active"><Link aria-current="page"><FontAwesomeIcon className="fas" icon={faPencil} />&nbsp;Update</Link></li>
                        </ul>
                    </nav>

                    {/* Mobile Breadcrumbs */}
                    <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
                        <ul>
                            <li class=""><Link to="/admin/video-collections" aria-current="page"><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to Video Collections</Link></li>
                        </ul>
                    </nav>

                    {/* Modal */}
                    {/* Nothing ... */}

                    {/* Page */}
                    <nav class="box">
                        <div class="columns">
                            <div class="column">
                                <p class="title is-4"><FontAwesomeIcon className="fas" icon={faVideoCamera} />&nbsp;Video Collection</p>
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
                                                        attachmentFilename={thumbnailAttachmentName}
                                                        setAttachmentFilename={setThumbnailAttachmentName}
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
                                                    maxWidth="100%"
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
                                    {isThumbnailUploaded && <>
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
                                            label="Summary"
                                            name="summary"
                                            placeholder="Text input"
                                            value={summary}
                                            errorText={errors && errors.summary}
                                            helpText=""
                                            onChange={(e)=>setSummary(e.target.value)}
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
                                            label="Type"
                                            name="type"
                                            placeholder="Pick"
                                            selectedValue={type}
                                            errorText={errors && errors.type}
                                            helpText=""
                                            onChange={(e) => setType(parseInt(e.target.value))}
                                            options={VIDEO_COLLECTION_TYPE_OPTIONS_WITH_EMPTY_OPTION}
                                        />

                                        <FormSelectFieldForVideoCategory
                                            videoCategoryID={videoCategoryID}
                                            setVideoCategoryID={setVideoCategoryID}
                                            isVideoCategoryOther={isVideoCategoryOther}
                                            setIsVideoCategoryOther={setIsVideoCategoryOther}
                                            errorText={errors && errors.videoCategoryID}
                                            helpText=""
                                            isRequired={true}
                                            maxWidth="520px"
                                        />
                                    </>}

                                    <div class="columns pt-5">
                                        <div class="column is-half">
                                            <Link class="button is-fullwidth-mobile" to={`/admin/video-collection/${vcid}`}><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to Detail</Link>
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
