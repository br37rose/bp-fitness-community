import React, { useState, useEffect } from 'react';
import Scroll from "react-scroll";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faArrowLeft, faDroplet, faDumbbell, faEye, faFilter, faFire, faGauge, faGenderless, faMap, faRefresh, faSearch, faTable } from '@fortawesome/free-solid-svg-icons';
import { Link, useNavigate } from 'react-router-dom';
import { useRecoilState } from "recoil";

import Layout from '../../../Menu/Layout';
import FormErrorBox from "../../../Reusable/FormErrorBox";
import { getVideoCategoryListAPI } from '../../../../API/VideoCategory';
import {
    topAlertMessageState,
    topAlertStatusState,
    currentUserState,
} from "../../../../AppState";
import PageLoadingContent from "../../../Reusable/PageLoadingContent";
import FormInputFieldWithButton from "../../../Reusable/FormInputFieldWithButton";
import VideoSection from '../../../Reusable/VideoSection'
import { PAGE_SIZE_OPTIONS } from '../../../../Constants/FieldOptions';
import { getVideoCollectionListAPI } from '../../../../API/VideoCollection';


////
//// Custom Component
////

const MemberCategoriesList = () => {

    ////
    //// Global state.
    ////

    const [topAlertMessage, setTopAlertMessage] = useRecoilState(topAlertMessageState);
    const [topAlertStatus, setTopAlertStatus] = useRecoilState(topAlertStatusState);
    const [currentUser] = useRecoilState(currentUserState);

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
    const [showFilter, setShowFilter] = useState(false); // Filtering + Searching
    const [sortField, setSortField] = useState("created"); // Sorting
    const [temporarySearchText, setTemporarySearchText] = useState(""); // Searching - The search field value as your writes their query.
    const [actualSearchText, setActualSearchText] = useState(""); // Searching - The actual search query value to submit to the API.

    const [modal, setModal] = useState(false);

    const navigate = useNavigate();

    ////
    //// API.
    ////

    const onVideoCollectionListSuccess = (response) => {
        console.log("onVideoCollectionListSuccess: Starting...");
        if (response.results !== null) {
            const groupedCollections = response.results.reduce((acc, collection) => {
                const { categoryId, id, categoryName, thumbnailUrl, name, type } = collection;
                if (!acc[categoryId]) {
                    acc[categoryId] = {
                        title: categoryName,
                        categoryId: categoryId,
                        videoCollectionId: id,
                        items: []
                    };
                }
                acc[categoryId].items.push({
                    src: thumbnailUrl,
                    alt: name,
                    text: name,
                    typeId: type
                });
                return acc;
            }, {});

            const listObjArr = Object.values(groupedCollections);
            setListData(listObjArr);

            console.log(listObjArr)
        } else {
            setListData([]);
        }
    };

    const onVideoCollectionListError = (apiErr) => {
        console.log("onVideoCollectionListError: Starting...");
        // ... Handle the error response here
    }

    const onVideoCollectionListDone = () => {
        console.log("onVideoCollectionListDone: Starting...");
        setFetching(false);
    }

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

    const fetchList = async (cur, limit, keywords) => {
        try {
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

            // Filtering
            if (keywords !== undefined && keywords !== null && keywords !== "") {
                // Searhcing
                params.set("search", keywords);
            }

            await getVideoCollectionListAPI(params, onVideoCollectionListSuccess, onVideoCollectionListError, onVideoCollectionListDone);
        } catch (error) {
            console.error('Error fetching data:', error);
            setErrors(error);
        } finally {
            setFetching(false);
        }
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

    ////
    //// Misc.
    ////

    useEffect(() => {
        let mounted = true;

        if (mounted) {
            window.scrollTo(0, 0); // Start the page at the top of the page.
            fetchList(currentCursor, pageSize, actualSearchText);
        }

        return () => {
            mounted = false;
        };
    }, [currentCursor, pageSize, actualSearchText]);

    ////
    //// Component rendering.
    ////

    return (
        <Layout breadcrumbItems={breadcrumbItems}>
            <div className="box">
                <div className="columns">
                    <div className="column">
                        <h1 className="title is-4">
                            <FontAwesomeIcon className="fas" icon={faDumbbell} />
                            &nbsp;Videos
                        </h1>
                    </div>

                    <hr class="mt-0 mb-2" />
                    <div className="column has-text-right">
                        <button onClick={() => fetchList(currentCursor, pageSize, actualSearchText)} class="is-fullwidth-mobile button is-link is-small" type="button">
                            <FontAwesomeIcon className="mdi" icon={faRefresh} />&nbsp;Refresh
                        </button>
                        &nbsp;
                        <button onClick={(e) => setShowFilter(!showFilter)} class="is-fullwidth-mobile button is-small is-primary" type="button">
                            <FontAwesomeIcon className="mdi" icon={faFilter} />&nbsp;Filter
                        </button>
                    </div>
                </div>

                {showFilter && (
                    <div
                        class="columns has-background-white-bis"
                        style={{ borderRadius: "15px", padding: "20px" }}
                    >
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

                        </div>
                    </div>
                )}

                {isFetching ? (
                    <PageLoadingContent displayMessage={"Please wait..."} />
                ) : (
                    <>
                        <FormErrorBox errors={errors} />
                        {listData &&
                            (listData.length > 0 || previousCursors.length > 0) ? (
                            <div className="container">

                                {/* Map over sections data and render VideoSection for each */}
                                {listData.map((section, idx) => (
                                    <VideoSection
                                        key={idx}
                                        {...section}
                                    />
                                ))}

                                <div class="columns">
                                    <div class="column is-half">
                                        <span class="select">
                                            <select
                                                class={`input has-text-grey-light`}
                                                name="pageSize"
                                                onChange={(e) =>
                                                    setPageSize(parseInt(e.target.value))
                                                }
                                            >
                                                {PAGE_SIZE_OPTIONS.map(function (option, i) {
                                                    return (
                                                        <option
                                                            selected={pageSize === option.value}
                                                            value={option.value}
                                                        >
                                                            {option.label}
                                                        </option>
                                                    );
                                                })}
                                            </select>
                                        </span>
                                    </div>
                                    <div class="column is-half has-text-right">
                                        {previousCursors.length > 0 && (
                                            <button
                                                class="button"
                                                onClick={onPreviousClicked}
                                            >
                                                Previous
                                            </button>
                                        )}
                                        {listData.hasNextPage && (
                                            <>
                                                <button class="button" onClick={onNextClicked}>
                                                    Next
                                                </button>
                                            </>
                                        )}
                                    </div>
                                </div>
                            </div>


                        ) : (
                            <section className="hero is-medium has-background-white-ter">
                                <div className="hero-body">
                                    <p className="title">
                                        <FontAwesomeIcon className="fas" icon={faMap} />
                                        &nbsp;No Videos
                                    </p>
                                    <p className="subtitle">
                                        No videos were found at the moment. Please check back later!
                                    </p>
                                </div>
                            </section>
                        )}
                    </>)}


            </div>
        </Layout>
    );
}

export default MemberCategoriesList;