import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
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
    faFilterCircleXmark
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import FormErrorBox from "../../Reusable/FormErrorBox";
import { getVideoCollectionListAPI } from "../../../API/VideoCollection";
import {
    topAlertMessageState,
    topAlertStatusState,
    currentUserState,
    videoCollectionsFilterShowState,
    videoCollectionsFilterTemporarySearchTextState,
    videoCollectionsFilterActualSearchTextState,
    videoCollectionsFilterStatusState,
    videoCollectionsFilterVideoTypeState,
    videoCollectionsFilterSortState
} from "../../../AppState";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import FormInputFieldWithButton from "../../Reusable/FormInputFieldWithButton";
import FormSelectField from "../../Reusable/FormSelectField";
import {
    PAGE_SIZE_OPTIONS,
    VIDEO_COLLECTION_STATUS_OPTIONS_WITH_EMPTY_OPTION,
    VIDEO_COLLECTION_TYPE_OPTIONS_WITH_EMPTY_OPTION
} from "../../../Constants/FieldOptions";
import MemberVideoCollectionListDesktop from "./ListDesktop";
import MemberVideoCollectionListMobile from "./ListMobile";
import Layout from '../../Menu/Layout';


function MemberVideoCollectionList() {
    ////
    //// Global state.
    ////

    const [topAlertMessage, setTopAlertMessage] = useRecoilState(topAlertMessageState);
    const [topAlertStatus, setTopAlertStatus] = useRecoilState(topAlertStatusState);
    const [currentUser] = useRecoilState(currentUserState);
    const [showFilter, setShowFilter] = useRecoilState(videoCollectionsFilterShowState); // Filtering + Searching
    const [temporarySearchText, setTemporarySearchText] = useRecoilState(videoCollectionsFilterTemporarySearchTextState); // Searching - The search field value as your writes their query.
    const [actualSearchText, setActualSearchText] = useRecoilState(videoCollectionsFilterActualSearchTextState); // Searching - The actual search query value to submit to the API.
    const [status, setStatus] = useRecoilState(videoCollectionsFilterStatusState);
    const [videoType, setVideoType] = useRecoilState(videoCollectionsFilterVideoTypeState);
    const [sort, setSort] = useRecoilState(videoCollectionsFilterSortState);

    ////
    //// Component states.
    ////

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

    function onVideoCollectionListSuccess(response) {
        console.log("onVideoCollectionListSuccess: Starting...");
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

    function onVideoCollectionListError(apiErr) {
        console.log("onVideoCollectionListError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onVideoCollectionListDone() {
        console.log("onVideoCollectionListDone: Starting...");
        setFetching(false);
    }

    function onVideoCollectionDeleteSuccess(response) {
        console.log("onVideoCollectionDeleteSuccess: Starting..."); // For debugging purposes only.

        // Update notification.
        setTopAlertStatus("success");
        setTopAlertMessage("Video collection deleted");
        setTimeout(() => {
            console.log(
                "onDeleteConfirmButtonClick: topAlertMessage, topAlertStatus:",
                topAlertMessage,
                topAlertStatus
            );
            setTopAlertMessage("");
        }, 2000);

        // Fetch again an updated list.
        fetchList(currentCursor, pageSize, actualSearchText, status, videoType, sort);
    }

    function onVideoCollectionDeleteError(apiErr) {
        console.log("onVideoCollectionDeleteError: Starting..."); // For debugging purposes only.
        setErrors(apiErr);

        // Update notification.
        setTopAlertStatus("danger");
        setTopAlertMessage("Failed deleting");
        setTimeout(() => {
            console.log(
                "onVideoCollectionDeleteError: topAlertMessage, topAlertStatus:",
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

    function onVideoCollectionDeleteDone() {
        console.log("onVideoCollectionDeleteDone: Starting...");
        setFetching(false);
    }

    ////
    //// BREADCRUMB
    ////
    const breadcrumbItems = {
        items: [
            { text: 'Dashboard', link: '/dashboard', isActive: false, icon: faGauge },
            { text: 'Video Collections', link: '#', icon: faVideoCamera, isActive: true }
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

    const fetchList = (cur, limit, keywords, st, vt, sbv) => {
        setFetching(true);
        setErrors({});

        let params = new Map();
        params.set("page_size", limit); // Pagination

        // DEVELOPERS NOTE: Our `sortByValue` is string with the sort field
        // and sort order combined with a comma seperation. Therefore we
        // need to split as follows.
        if (sbv !== undefined && sbv !== null && sbv !== "") {
            const sortArray = sbv.split(",");
            params.set("sort_field", sortArray[0]); // Sort (1 of 2)
            params.set("sort_order", sortArray[1]); // Sort (2 of 2)
        }

        if (cur !== "") {
            // Pagination
            params.set("cursor", cur);
        }

        params.set("type", vt);

        // Filtering
        if (keywords !== undefined && keywords !== null && keywords !== "") {
            // Searhcing
            params.set("search", keywords);
        }
        if (st !== undefined && st !== null && st !== "") {
            params.set("status", st);
        }
        if (vt !== undefined && vt !== null && vt !== "") {
            params.set("video_type", vt);
        }

        getVideoCollectionListAPI(
            params,
            onVideoCollectionListSuccess,
            onVideoCollectionListError,
            onVideoCollectionListDone
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
        setSort("created,-1");
    }

    ////
    //// Misc.
    ////

    useEffect(() => {
        let mounted = true;

        if (mounted) {
            window.scrollTo(0, 0); // Start the page at the top of the page.
            fetchList(currentCursor, pageSize, actualSearchText, status, videoType, sort);
        }

        return () => {
            mounted = false;
        };
    }, [currentCursor, pageSize, actualSearchText, status, videoType, sort]);

    ////
    //// Component rendering.
    ////

    return (
        <Layout breadcrumbItems={breadcrumbItems}>
            {/* Page */}
            <div className="box">
                <div className="columns">
                    <div className="column">
                        <h1 className="title is-4">
                            <FontAwesomeIcon className="fas" icon={faVideoCamera} />
                            &nbsp;Video Collections
                        </h1>
                    </div>
                    <div className="column has-text-right">
                        <button onClick={() => fetchList(currentCursor, pageSize, actualSearchText, status, videoType, sort)} class="is-fullwidth-mobile button is-link is-small" type="button">
                            <FontAwesomeIcon className="mdi" icon={faRefresh} />&nbsp;<span class="is-hidden-desktop is-hidden-tablet">Refresh</span>
                        </button>
                        &nbsp;
                        <button onClick={(e) => setShowFilter(!showFilter)} class="is-fullwidth-mobile button is-small is-primary" type="button">
                            <FontAwesomeIcon className="mdi" icon={faFilter} />&nbsp;Filter
                        </button>
                    </div>
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
                                <FormSelectField
                                    label="Video Type"
                                    name="videoType"
                                    placeholder="Pick"
                                    selectedValue={videoType}
                                    errorText={errors && errors.videoType}
                                    helpText=""
                                    onChange={(e) => setVideoType(e.target.value)}
                                    options={VIDEO_COLLECTION_TYPE_OPTIONS_WITH_EMPTY_OPTION}
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
                                    <MemberVideoCollectionListDesktop
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
                                    <MemberVideoCollectionListMobile
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
                                        &nbsp;No Video Collections
                                    </p>
                                    <p className="subtitle">
                                        No video collections.{" "}
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

export default MemberVideoCollectionList;
