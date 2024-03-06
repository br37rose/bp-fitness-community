import React, { useState, useEffect } from 'react';
import Scroll from "react-scroll";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faDroplet, faDumbbell, faEye, faFilter, faFire, faGauge, faGenderless, faMap, faRefresh, faSearch, faTable } from '@fortawesome/free-solid-svg-icons';
import { useParams } from 'react-router-dom';
import { useRecoilState } from "recoil";

import Layout from '../../../../Menu/Layout';
import FormErrorBox from "../../../../Reusable/FormErrorBox";
import {
    topAlertMessageState,
    topAlertStatusState,
    currentUserState,
} from "../../../../../AppState";
import PageLoadingContent from "../../../../Reusable/PageLoadingContent";
import { getVideoCollectionListAPI } from '../../../../../API/VideoCollection';


const MemberVideoCollectionList = () => {
    ////
    //// URL Parameters.
    ////

    const { vcatid } = useParams()

    ////
    //// Global state.
    ////

    const [topAlertMessage, setTopAlertMessage] =
        useRecoilState(topAlertMessageState);
    const [topAlertStatus, setTopAlertStatus] =
        useRecoilState(topAlertStatusState);
    const [currentUser] = useRecoilState(currentUserState);

    ////
    //// Component states.
    ////

    const [errors, setErrors] = useState({});
    const [listData, setListData] = useState("");
    const [category, setCategory] = useState({
        name: "",
    });
    const [isFetching, setFetching] = useState(false);
    const [pageSize, setPageSize] = useState(10); // Pagination
    const [previousCursors, setPreviousCursors] = useState([]); // Pagination
    const [nextCursor, setNextCursor] = useState(""); // Pagination
    const [currentCursor, setCurrentCursor] = useState(""); // Pagination
    const [showFilter, setShowFilter] = useState(false); // Filtering + Searching
    const [sortField, setSortField] = useState("created"); // Sorting
    const [temporarySearchText, setTemporarySearchText] = useState(""); // Searching - The search field value as your writes their query.
    const [actualSearchText, setActualSearchText] = useState(""); // Searching - The actual search query value to submit to the API.


    const breadcrumbItems = [
        { text: 'Dashboard', link: '/dashboard', icon: faGauge },
        { text: 'Categories', link: '/video-categories', icon: faMap },
        { text: `${category.name}`, link: '#', icon: faDumbbell, isActive: true }
    ];

    ////
    //// API.
    ////

    function onVideoCollectionListSuccess(response) {
        console.log("onVideoCollectionListSuccess: Starting...");
        if (response.results !== null) {
            let categoryName = "";
            let categorizedResponse = response.results.filter((collection) => {
                if (collection.categoryId === vcatid) {
                    categoryName = collection.categoryName;
                    return collection;
                }
            });
            setCategory({
                name: categoryName
            })
            setListData(categorizedResponse);

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

    ////
    //// Event handling.
    ////

    const fetchList = (cur, limit, keywords) => {
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

            <div class="box">
                <div class="container is-fluid p-0">
                    <div class="columns is-multiline">
                        <div class="hero-body p-4">
                            <div class="columns is-multiline mb-0">
                                <div class="column">
                                    <h1 className="title is-4">
                                        <FontAwesomeIcon className="fas" icon={faDumbbell} />
                                        &nbsp;{category.name}
                                    </h1>
                                </div>
                            </div>
                        </div>
                    </div>

                </div>

                {isFetching ? (
                    <PageLoadingContent displayMessage={"Please wait..."} />
                ) : (
                    <>
                        <FormErrorBox errors={errors} />
                        {listData &&
                            listData &&
                            (listData.length > 0 || previousCursors.length > 0) ? (


                            <div class="columns is-multiline mb-0">
                                {listData.map((collection) => (
                                    <div class="column mb-0 is-3">
                                        <a href={`/video-collection/${collection.categoryId}/video-content/${collection.id}`} class="border-radius has-text-black">
                                            <img class="border-radius m-w-100" src={collection.thumbnailUrl} alt={collection.name} />
                                            <h5 class="is-size-5 has-text-weight-semibold">{collection.name}</h5>
                                        </a>
                                    </div>
                                ))}
                            </div>


                        ) : (
                            <section className="hero is-medium has-background-white-ter">
                                <div className="hero-body">
                                    <p className="title">
                                        <FontAwesomeIcon className="fas" icon={faTable} />
                                        &nbsp;No Exercises
                                    </p>
                                    <p className="subtitle">
                                        No exercises found at the moment. Please check back later!
                                    </p>
                                </div>
                            </section>
                        )}
                    </>
                )}
            </div>
        </Layout>
    );
}

export default MemberVideoCollectionList;