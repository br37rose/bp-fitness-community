import React, { useState, useEffect } from "react";
import { Link, useParams } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
    faVideo,
    faArrowLeft,
    faVideoCamera,
    faEye,
    faPencil,
    faTrashCan,
    faPlus,
    faGauge,
    faArrowRight,
    faTable,
    faArrowUpRightFromSquare,
    faRefresh,
    faFilter,
    faSearch,
    faFilterCircleXmark,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import FormErrorBox from "../../../Reusable/FormErrorBox";
import { getVideoContentListAPI } from "../../../../API/VideoContent";
import {
    topAlertMessageState,
    topAlertStatusState,
    currentUserState,
    videoContentsFilterShowState,
    videoContentsFilterSortState,
    videoContentsFilterTemporarySearchTextState,
    videoContentsFilterActualSearchTextState,
    videoContentsFilterVideoTypeState,
    videoContentsFilterGenderState,
    videoContentsFilterStatusState,
    videoContentsFilterOfferIDState,
    videoContentsFilterCategoryIDState
} from "../../../../AppState";
import PageLoadingContent from "../../../Reusable/PageLoadingContent";
import FormInputFieldWithButton from "../../../Reusable/FormInputFieldWithButton";
import FormSelectFieldForOffer from "../../../Reusable/FormSelectFieldForOffer";
import FormSelectFieldForVideoCategory from "../../../Reusable/FormSelectFieldForVideoCategory";
import FormSelectField from "../../../Reusable/FormSelectField";
import {
    PAGE_SIZE_OPTIONS,
    VIDEO_COLLECTION_STATUS_OPTIONS_WITH_EMPTY_OPTION,
    VIDEO_CONTENT_VIDEO_TYPE_WITH_EMPTY_OPTIONS
} from "../../../../Constants/FieldOptions";
import MemberVideoContentListDesktop from "./ListDesktop";
import MemberVideoContentListMobile from "./ListMobile";
import Layout from "../../../Menu/Layout";


function MemberVideoContentList() {
    ////
    //// URL Parameters.
    ////

    const { vcid } = useParams()

    ////
    //// Global state.
    ////

    const [topAlertMessage, setTopAlertMessage] = useRecoilState(topAlertMessageState);
    const [topAlertStatus, setTopAlertStatus] = useRecoilState(topAlertStatusState);
    const [currentUser] = useRecoilState(currentUserState);
    const [showFilter, setShowFilter] = useRecoilState(videoContentsFilterShowState); // Filtering + Searching
    const [sort, setSort] = useRecoilState(videoContentsFilterSortState); // Sorting
    const [temporarySearchText, setTemporarySearchText] = useRecoilState(videoContentsFilterTemporarySearchTextState); // Searching - The search field value as your writes their query.
    const [actualSearchText, setActualSearchText] = useRecoilState(videoContentsFilterActualSearchTextState); // Searching - The actual search query value to submit to the API.
    const [status, setStatus] = useRecoilState(videoContentsFilterStatusState);
    const [videoType, setVideoType] = useRecoilState(videoContentsFilterVideoTypeState);
    const [offerID, setOfferID] = useRecoilState(videoContentsFilterOfferIDState);
    const [videoCategoryID, setVideoCategoryID] = useRecoilState(videoContentsFilterCategoryIDState);

    ////
    //// Component states.
    ////

    const [isVideoCategoryOther, setIsVideoCategoryOther] = useState("");
    const [isOfferOther, setIsOfferOther] = useState("");
    const [errors, setErrors] = useState({});
    const [listData, setListData] = useState("");
    const [isFetching, setFetching] = useState(false);
    const [pageSize, setPageSize] = useState(10); // Pagination
    const [previousCursors, setPreviousCursors] = useState([]); // Pagination
    const [nextCursor, setNextCursor] = useState(""); // Pagination
    const [currentCursor, setCurrentCursor] = useState(""); // Pagination

    ////
    //// API.
    ////

    function onVideoContentListSuccess(response) {
        console.log("onVideoContentListSuccess: Starting...");
        if (response.results !== null) {
            setListData(response);
            if (response.hasNextPage) {
                setNextCursor(response.nextCursor); // For pagination purposes.
            }
        } else {
            setListData([]);
            setNextCursor("");
        }
    }

    function onVideoContentListError(apiErr) {
        console.log("onVideoContentListError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onVideoContentListDone() {
        console.log("onVideoContentListDone: Starting...");
        setFetching(false);
    }

    function onVideoContentDeleteSuccess(response) {
        console.log("onVideoContentDeleteSuccess: Starting..."); // For debugging purposes only.

        // Update notification.
        setTopAlertStatus("success");
        setTopAlertMessage("VideoContent deleted");
        setTimeout(() => {
            console.log(
                "onDeleteConfirmButtonClick: topAlertMessage, topAlertStatus:",
                topAlertMessage,
                topAlertStatus
            );
            setTopAlertMessage("");
        }, 2000);

        // Fetch again an updated list.
        fetchList(currentCursor, pageSize, actualSearchText, status, videoType, offerID, videoCategoryID, vcid);
    }

    function onVideoContentDeleteError(apiErr) {
        console.log("onVideoContentDeleteError: Starting..."); // For debugging purposes only.
        setErrors(apiErr);

        // Update notification.
        setTopAlertStatus("danger");
        setTopAlertMessage("Failed deleting");
        setTimeout(() => {
            console.log(
                "onVideoContentDeleteError: topAlertMessage, topAlertStatus:",
                topAlertMessage,
                topAlertStatus
            );
            setTopAlertMessage("");
        }, 2000);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onVideoContentDeleteDone() {
        console.log("onVideoContentDeleteDone: Starting...");
        setFetching(false);
    }

    ////
    //// BREADCRUMB
    ////
    const breadcrumbItems = {
        items: [
            { text: 'Dashboard', link: '/dashboard', isActive: false, icon: faGauge },
            { text: 'Video Collections', link: '/video-collections', icon: faVideo, isActive: false },
            { text: 'Detail (Video Content)', link: '#', icon: faEye, isActive: true }
        ],
        mobileBackLinkItems: {
            link: '/video-collections',
            text: 'Back to Video Collections',
            icon: faArrowLeft
        }
    }

    ////
    //// Event handling.
    ////

    const fetchList = (cur, limit, keywords, st, vct, oid, catid, vcid) => {
        setFetching(true);
        setErrors({});

        let params = new Map();
        params.set("page_size", limit); // Pagination
        params.set("sort_field", "created"); // Sorting
        params.set("sort_order", -1)         // Sorting - descending, meaning most recent start date to oldest start date.

        if (cur !== "") {
            // Pagination
            params.set("cursor", cur);
        }

        params.set("video_collection_id", vcid);

        // Filtering
        if (keywords !== undefined && keywords !== null && keywords !== "") {
            // Searhcing
            params.set("search", keywords);
        }
        if (st !== undefined && st !== null && st !== "") {
            params.set("status", st);
        }
        if (vct !== undefined && vct !== null && vct !== "") {
            params.set("video_type", vct);
        }
        if (oid !== undefined && oid !== null && oid !== "") {
            params.set("offer_id", oid);
        }
        if (catid !== undefined && catid !== null && catid !== "") {
            params.set("category_id", catid);
        }

        console.log("params:", params);

        getVideoContentListAPI(
            params,
            onVideoContentListSuccess,
            onVideoContentListError,
            onVideoContentListDone
        );
    };

    const onNextClicked = (e) => {
        console.log("onNextClicked");
        let arr = [...previousCursors];
        arr.push(currentCursor);
        setPreviousCursors(arr);
        setCurrentCursor(nextCursor);
    };

    const onPreviousClicked = (e) => {
        console.log("onPreviousClicked");
        let arr = [...previousCursors];
        const previousCursor = arr.pop();
        setPreviousCursors(arr);
        setCurrentCursor(previousCursor);
    };

    const onSearchButtonClick = (e) => {
        // Searching
        console.log("Search button clicked...");
        setActualSearchText(temporarySearchText);
    };

    // Function resets the filter state to its default state.
    const onClearFilterClick = (e) => {
        setShowFilter(false);
        setActualSearchText("");
        setTemporarySearchText("");
        setVideoType(0);
        setStatus(0);
        setOfferID(null);
        setVideoCategoryID("");
        setSort("created,-1");
    }

    ////
    //// Misc.
    ////

    useEffect(() => {
        let mounted = true;

        if (mounted) {
            window.scrollTo(0, 0); // Start the page at the top of the page.
            fetchList(currentCursor, pageSize, actualSearchText, status, videoType, offerID, videoCategoryID, vcid);
        }

        return () => {
            mounted = false;
        };
    }, [currentCursor, pageSize, actualSearchText, status, videoType, offerID, videoCategoryID, vcid]);

    ////
    //// Component rendering.
    ////

    return (
        <Layout breadcrumbItem={breadcrumbItems}>
            {/* Page */}
            <div className="box">

                <div className="columns">
                    <div className="column">
                        <h1 className="title is-4">
                            <FontAwesomeIcon className="fas" icon={faVideoCamera} />
                            &nbsp;Video Collection - Contents
                        </h1>
                    </div>
                    <div className="column has-text-right">
                        <button onClick={() => fetchList(currentCursor, pageSize, actualSearchText, status, videoType, offerID, videoCategoryID, vcid)} class="is-fullwidth-mobile button is-link is-small" type="button">
                            <FontAwesomeIcon className="mdi" icon={faRefresh} />&nbsp;<span class="is-hidden-desktop is-hidden-tablet">Refresh</span>
                        </button>
                        &nbsp;
                        <button onClick={(e) => setShowFilter(!showFilter)} class="is-fullwidth-mobile button is-small is-primary" type="button">
                            <FontAwesomeIcon className="mdi" icon={faFilter} />&nbsp;Filter
                        </button>
                    </div>
                </div>

                {/* Tab Navigation */}
                <div class="tabs is-medium is-size-7-mobile">
                    <ul>
                        <li>
                            <Link to={`/video-collection/${vcid}`}>Detail</Link>
                        </li>
                        <li class="is-active">
                            <Link to={`/video-collection/${vcid}/video-contents`}><strong>Contents</strong></Link>
                        </li>
                    </ul>
                </div>

                {/* FILTER */}
                {showFilter && (
                    <div class="has-background-white-bis" style={{ borderRadius: "15px", padding: "20px" }}>

                        {/* Filter Title + Clear Button */}
                        <div class="columns is-mobile">
                            <div class="column is-half">
                                <strong><u><FontAwesomeIcon className="mdi" icon={faFilter} />&nbsp;Filter</u></strong>
                            </div>
                            <div class="column is-half has-text-right">
                                <Link onClick={onClearFilterClick}><FontAwesomeIcon className="mdi" icon={faFilterCircleXmark} />&nbsp;Clear Filter</Link>
                            </div>
                        </div>

                        {/* Filter Options */}
                        <div class="columns">

                            <div class="column">
                                <FormInputFieldWithButton
                                    label={"Search"}
                                    name="temporarySearchText"
                                    type="text"
                                    placeholder="Search by name"
                                    value={temporarySearchText}
                                    helpText=""
                                    onChange={(e) => setTemporarySearchText(e.target.value)}
                                    isRequired={true}
                                    maxWidth="100%"
                                    buttonLabel={
                                        <>
                                            <FontAwesomeIcon className="fas" icon={faSearch} />
                                        </>
                                    }
                                    onButtonClick={onSearchButtonClick}
                                />
                            </div>
                            <div class="column">
                                <FormSelectField
                                    label="Status"
                                    name="status"
                                    placeholder="Pick"
                                    selectedValue={status}
                                    errorText={errors && errors.status}
                                    helpText=""
                                    onChange={(e) => setStatus(parseInt(e.target.value))}
                                    options={VIDEO_COLLECTION_STATUS_OPTIONS_WITH_EMPTY_OPTION}
                                />
                            </div>
                            <div class="column">
                                <FormSelectFieldForOffer
                                    label={`Enrollment`}
                                    isSubscription={true}
                                    offerID={offerID}
                                    setOfferID={setOfferID}
                                    isOfferOther={isOfferOther}
                                    setIsOfferOther={setIsOfferOther}
                                    errorText={errors && errors.offerId}
                                />
                            </div>
                            <div class="column">
                                <FormSelectField
                                    label="Video Type"
                                    name="videoType"
                                    placeholder="Pick"
                                    selectedValue={videoType}
                                    errorText={errors && errors.videoType}
                                    helpText=""
                                    onChange={(e) => setVideoType(e.target.value)}
                                    options={VIDEO_CONTENT_VIDEO_TYPE_WITH_EMPTY_OPTIONS}
                                />
                            </div>
                            <div class="column">
                                <FormSelectFieldForVideoCategory
                                    label="Video Category"
                                    videoCategoryID={videoCategoryID}
                                    setVideoCategoryID={setVideoCategoryID}
                                    isVideoCategoryOther={isVideoCategoryOther}
                                    setIsVideoCategoryOther={setIsVideoCategoryOther}
                                    errorText={errors && errors.videoCategoryID}
                                    helpText=""
                                    isRequired={true}
                                    maxWidth="520px"
                                />
                            </div>
                        </div>
                    </div>
                )}

                {isFetching ? (
                    <PageLoadingContent displayMessage={"Please wait..."} />
                ) : (
                    <>
                        <FormErrorBox errors={errors} />
                        {listData &&
                            listData.results &&
                            (listData.results.length > 0 || previousCursors.length > 0) ? (
                            <div className="container">

                                {/*
                                    ##################################################################
                                    EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A DESKTOP SCREEN.
                                    ##################################################################
                                */}
                                <div class="is-hidden-touch" >
                                    <MemberVideoContentListDesktop
                                        listData={listData}
                                        setPageSize={setPageSize}
                                        pageSize={pageSize}
                                        previousCursors={previousCursors}
                                        onPreviousClicked={onPreviousClicked}
                                        onNextClicked={onNextClicked}
                                    />
                                </div>

                                {/*
                                    ###########################################################################
                                    EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A TABLET OR MOBILE SCREEN.
                                    ###########################################################################
                                */}
                                <div class="is-fullwidth is-hidden-desktop">
                                    <MemberVideoContentListMobile
                                        listData={listData}
                                        setPageSize={setPageSize}
                                        pageSize={pageSize}
                                        previousCursors={previousCursors}
                                        onPreviousClicked={onPreviousClicked}
                                        onNextClicked={onNextClicked}
                                    />
                                </div>

                            </div>
                        ) : (
                            <section className="hero is-medium has-background-white-ter">
                                <div className="hero-body">
                                    <p className="title">
                                        <FontAwesomeIcon className="fas" icon={faTable} />
                                        &nbsp;No Video Contents
                                    </p>
                                    <p className="subtitle">
                                        No class types.{" "}
                                        <b>
                                            <Link to="/video-collections/add">
                                                Click here&nbsp;
                                                <FontAwesomeIcon
                                                    className="mdi"
                                                    icon={faArrowRight}
                                                />
                                            </Link>
                                        </b>{" "}
                                        to get started creating your first video collections.
                                    </p>
                                </div>
                            </section>
                        )}
                    </>
                )}

                <div class="columns pt-5">
                    <div class="column is-half">
                        <Link class="button is-fullwidth-mobile" to={`/dashboard`}><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to Dashboard</Link>
                    </div>
                    <div class="column is-half has-text-right">

                    </div>
                </div>
            </div>
        </Layout>
    );
}

export default MemberVideoContentList;