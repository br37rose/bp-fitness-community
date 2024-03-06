import React, { useState, useEffect, useRef } from "react";
import { Link } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
    faHeart,
    faPersonWalking,
    faRankingStar,
    faChartLine,
    faHeartbeat,
    faFilterCircleXmark,
    faArrowLeft,
    faUsers,
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
    faUser,
    faShoePrints,
    faMapMarkerAlt,
    faLock,
    faShoppingBag,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import FormErrorBox from "../../../../Reusable/FormErrorBox";
import { getRankPointListAPI } from "../../../../../API/RankPoint";
import {
    topAlertMessageState,
    topAlertStatusState,
    currentUserState,
    dataPointFilterShowState,
    dataPointFilterTemporarySearchTextState,
    dataPointFilterActualSearchTextState,
    dataPointFilterSortState,
    dataPointFilterStatusState,
    dataPointFilterIsHeartRateState,
    dataPointFilterIsStepsCounterState
} from "../../../../../AppState";
import FormCheckboxField from "../../../../Reusable/FormCheckboxField";
import PageLoadingContent from "../../../../Reusable/PageLoadingContent";
import FormInputFieldWithButton from "../../../../Reusable/FormInputFieldWithButton";
import { PAGE_SIZE_OPTIONS } from "../../../../../Constants/FieldOptions";
import { RANK_POINT_METRIC_TYPE_HEART_RATE, RANK_POINT_METRIC_TYPE_STEP_COUNTER, RANK_POINT_PERIOD_DAY, RANK_POINT_PERIOD_WEEK, RANK_POINT_PERIOD_MONTH, RANK_POINT_PERIOD_YEAR } from "../../../../../Constants/App";
import MemberLeaderboardPersonalDesktop from "./PersonalDesktop";
import MemberLeaderboardPersonalMobile from "./PersonalMobile";
import Card from '../../../../Reusable/Dashboard/Card'
// import { Bar, Doughnut } from 'react-chartjs-2';
import BarChart from "../../../../Reusable/Charts/Bar";
import DoughnutChart from "../../../../Reusable/Charts/Doughnut";

function MemberLeaderboardPersonal() {
    ////
    //// Global state.
    ////

    const [topAlertMessage, setTopAlertMessage] = useRecoilState(topAlertMessageState);
    const [topAlertStatus, setTopAlertStatus] = useRecoilState(topAlertStatusState);
    const [currentUser] = useRecoilState(currentUserState);
    const [showFilter, setShowFilter] = useRecoilState(dataPointFilterShowState); // Filtering + Searching
    const [sort, setSort] = useRecoilState(dataPointFilterSortState); // Sorting
    const [temporarySearchText, setTemporarySearchText] = useRecoilState(dataPointFilterTemporarySearchTextState); // Searching - The search field value as your writes their query.
    const [actualSearchText, setActualSearchText] = useRecoilState(dataPointFilterActualSearchTextState); // Searching - The actual search query value to submit to the API.
    const [status, setStatus] = useRecoilState(dataPointFilterStatusState);
    // const [isHeartRate, setIsHeartRate] = useRecoilState(dataPointFilterIsHeartRateState);
    // const [isStepsCounter, setIsStepsCounter] = useRecoilState(dataPointFilterIsStepsCounterState);

    ////
    //// Component states.
    ////

    const [errors, setErrors] = useState({});
    const [listRank, setListRank] = useState("");
    const [selectedFitnessPlanForDeletion, setSelectedFitnessPlanForDeletion] = useState("");
    const [isFetching, setFetching] = useState(false);
    const [pageSize, setPageSize] = useState(100); // Pagination
    const [previousCursors, setPreviousCursors] = useState([]); // Pagination
    const [nextCursor, setNextCursor] = useState(""); // Pagination
    const [currentCursor, setCurrentCursor] = useState(""); // Pagination
    const [isHeartRate, setIsHeartRate] = useState(true);
    const [isStepsCounter, setIsStepsCounter] = useState(false);
    const [period, setPeriod] = useState(RANK_POINT_PERIOD_DAY);

    ////
    //// API.
    ////

    function onRankPointistSuccess(response) {
        console.log("onRankPointistSuccess: Starting...");
        if (response.results !== null) {
            setListRank(response);
            if (response.hasNextPage) {
                setNextCursor(response.nextCursor); // For pagination purposes.
            }
        } else {
            setListRank([]);
            setNextCursor("");
        }
    }

    function onRankPointistError(apiErr) {
        console.log("onRankPointistError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onRankPointistDone() {
        console.log("onRankPointistDone: Starting...");
        setFetching(false);
    }

    ////
    //// Event handling.
    ////



    const fetchList = (user, cur, limit, keywords, stat, sbv, isHeartRate, isStepsCounter, p) => {
        setFetching(true);
        setErrors({});

        let params = new Map();
        params.set("page_size", limit); // Pagination

        // Always sort by place in ascending order so the list will come like:
        // #1, #2, #3, etc.
        params.set("sort_field", "place");
        params.set("sort_order", "ASC");

        if (cur !== "") {
            // Pagination
            params.set("cursor", cur);
        }

        // Filtering
        if (keywords !== undefined && keywords !== null && keywords !== "") {
            // Searhcing
            params.set("search", keywords);
        }

        params.set("status", stat);

        // Set the filtering by devices.
        if (user !== undefined && user !== null && user !== "") {
            if (isHeartRate === true) {
                params.set("metric_types", RANK_POINT_METRIC_TYPE_HEART_RATE);
            }
            if (isStepsCounter === true) {
                params.set("metric_types", RANK_POINT_METRIC_TYPE_STEP_COUNTER);
            }
        }

        // Set the filtering by time period.
        params.set("period", parseInt(period));

        // Add extra filter if its for today.
        if (period === RANK_POINT_PERIOD_DAY) {
            params.set("is_today_only", true);
        }

        // Make the submission to the API backend.
        getRankPointListAPI(
            params,
            onRankPointistSuccess,
            onRankPointistError,
            onRankPointistDone
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

    const onSelectFitnessPlanForDeletion = (e, datum) => {
        console.log("onSelectFitnessPlanForDeletion", datum);
        setSelectedFitnessPlanForDeletion(datum);
    };

    const onDeselectFitnessPlanForDeletion = (e) => {
        console.log("onDeselectFitnessPlanForDeletion");
        setSelectedFitnessPlanForDeletion("");
    };

    // Function resets the filter state to its default state.
    const onClearFilterClick = (e) => {
        setShowFilter(false);
        setActualSearchText("");
        setTemporarySearchText("");
        setSort("timestamp,DESC");
        setStatus(0);
        setIsHeartRate(true);
        setIsStepsCounter(true);
    }

    ////
    //// Misc.
    ////

    useEffect(() => {
        let mounted = true;

        if (mounted) {
            window.scrollTo(0, 0); // Start the page at the top of the page.
            fetchList(currentUser, currentCursor, pageSize, actualSearchText, status, sort, isHeartRate, isStepsCounter, period);
        }

        return () => {
            mounted = false;
        };
    }, [currentUser, currentCursor, pageSize, actualSearchText, status, sort, isHeartRate, isStepsCounter, period]);

    ////
    //// Component rendering.
    ////
    const barChartData = {
        labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
        datasets: [
            {
                label: 'Total Spent',
                data: [12, 19, 3, 5, 2, 3, 6, 4, 7, 8, 9, 10],
                backgroundColor: 'rgba(54, 162, 235, 0.2)',
                borderColor: 'rgba(54, 162, 235, 1)',
                borderWidth: 1,
            },
        ],
    };

    const doughnutChartData = {

        labels: ['Sales', 'Orders', 'Returns'],
        datasets: [
            {
                label: 'Analytics',
                data: [300, 50, 100],
                backgroundColor: [
                    'rgba(255, 99, 132, 0.2)',
                    'rgba(54, 162, 235, 0.2)',
                    'rgba(255, 206, 86, 0.2)',
                ],
                borderColor: [
                    'rgba(255,99,132,1)',
                    'rgba(54, 162, 235, 1)',
                    'rgba(255, 206, 86, 1)',
                ],
                borderWidth: 1,
            },
        ],
    };

    const sampleTopSellingData = [
        { name: 'Nike Airmax 170', price: 567, rating: 5 },
        { name: 'Nike Airmax 170', price: 200, rating: 4 },
        { name: 'Nike Airmax 170', price: 400, rating: 5 }
    ]

    const formatCurrency = (amount) => {
        return `$${amount.toFixed(2).replace(/\d(?=(\d{3})+\.)/g, '$&,')}`;
    };

    return (
        <>
            <div className="container">
                <section className="section">

                    {/* Desktop Breadcrumbs */}
                    <nav className="breadcrumb is-hidden-touch" aria-label="breadcrumbs">
                        <ul>
                            <li className=""><Link to="/dashboard" aria-current="page"><FontAwesomeIcon className="fas" icon={faGauge} />&nbsp;Dashboard</Link></li>
                            <li className=""><Link to="/biometrics"><FontAwesomeIcon className="fas" icon={faHeartbeat} />&nbsp;Biometrics</Link></li>
                            <li className="is-active"><Link aria-current="page"><FontAwesomeIcon className="fas" icon={faRankingStar} />&nbsp;Leadboard</Link></li>
                        </ul>
                    </nav>

                    {/* Mobile Breadcrumbs */}
                    <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
                        <ul>
                            <li class="">
                                <Link to="/dashboard" aria-current="page"><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to Dashboard</Link>
                            </li>
                        </ul>
                    </nav>

                    {/* Page */}
                    <nav className="box">
                        <div className="columns">
                            <div className="column">
                                <h1 className="title is-4">
                                    <FontAwesomeIcon className="fas" icon={faRankingStar} />
                                    &nbsp;Leaderboard
                                </h1>
                            </div>
                            <div className="column has-text-right">
                                <button onClick={() => fetchList(currentUser, currentCursor, pageSize, actualSearchText, status, sort, isHeartRate, isStepsCounter)} class="is-fullwidth-mobile button is-link is-small" type="button">
                                    <FontAwesomeIcon className="mdi" icon={faRefresh} />&nbsp;<span class="is-hidden-desktop is-hidden-tablet">Refresh</span>
                                </button>
                                &nbsp;
                                {/*
                      <button onClick={(e)=>setShowFilter(!showFilter)} class="is-fullwidth-mobile button is-small is-primary" type="button">
                          <FontAwesomeIcon className="mdi" icon={faFilter} />&nbsp;Filter
                      </button>

                      &nbsp;
                      <Link to={`/account/wearable-tech`} className="is-fullwidth-mobile button is-small is-success" type="button">
                          <FontAwesomeIcon className="mdi" icon={faPlus} />&nbsp;Register Wearable
                      </Link>
                      */}
                            </div>
                        </div>

                        {isFetching ? (
                            <PageLoadingContent displayMessage={"Please wait..."} />
                        ) : (
                            <div className="section">
                                <div className="columns">
                                    <div className="column is-one-fourth">
                                        <div className="box">
                                            <div className="media">
                                                <div className="media-left">
                                                    <span className="icon">
                                                        <FontAwesomeIcon className="fas" icon={faShoppingBag} />
                                                    </span>
                                                </div>
                                                <div className="media-content ">
                                                    <p className="title has-text-weight-semibold is-7">All Spending</p>
                                                    <p className="has-text-weight-semibold subtitle is-6">$574</p>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                    <div className="column is-one-fourth">
                                        <div className="box">
                                            <div className="media">
                                                <div className="media-left">
                                                    <span className="icon">
                                                        <FontAwesomeIcon className="fas" icon={faShoppingBag} />
                                                    </span>
                                                </div>
                                                <div className="media-content ">
                                                    <p className="title has-text-weight-semibold is-7">All Spending</p>
                                                    <p className="has-text-weight-semibold subtitle is-6">$574</p>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                    <div className="column is-one-fourth">
                                        <div className="box">
                                            <div className="media">
                                                <div className="media-left">
                                                    <span className="icon">
                                                        <FontAwesomeIcon className="fas" icon={faShoppingBag} />
                                                    </span>
                                                </div>
                                                <div className="media-content ">
                                                    <p className="title has-text-weight-semibold is-7">All Spending</p>
                                                    <p className="has-text-weight-semibold subtitle is-6">$574</p>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                    <div className="column is-one-fourth">
                                        <div className="box">
                                            <div className="media">
                                                <div className="media-left">
                                                    <span className="icon">
                                                        <FontAwesomeIcon className="fas" icon={faShoppingBag} />
                                                    </span>
                                                </div>
                                                <div className="media-content ">
                                                    <p className="title has-text-weight-semibold is-7">All Spending</p>
                                                    <p className="has-text-weight-semibold subtitle is-6">$574</p>
                                                </div>
                                            </div>
                                        </div>
                                    </div>

                                </div> {/* Graphs Row */}


                                <div className="columns">
                                    {/* Doughnut Chart */}
                                    <div className="column is-one-third">
                                        <BarChart data={doughnutChartData} options={{ maintainAspectRatio: false }} />
                                    </div>
                                    {/* Another Chart or Content */}
                                    <div className="column is-one-third">
                                        <BarChart data={doughnutChartData} options={{ maintainAspectRatio: false }} />
                                    </div>
                                    <div className="column is-one-third">
                                        <BarChart data={doughnutChartData} options={{ maintainAspectRatio: false }} />
                                    </div>
                                </div>

                                {/* Table */}
                                <div className="columns">
                                    <div className="column is-half">
                                        <DoughnutChart data={doughnutChartData} options={{ maintainAspectRatio: false }} />
                                    </div>
                                    <div className="column">
                                        <div className="box">
                                            <h3 className="title is-4">Top Selling Products</h3>
                                            <table className="table is-fullwidth is-striped">
                                                <thead>
                                                    <tr>
                                                        <th>Product Name</th>
                                                        <th>Price</th>
                                                        <th>Rating</th>
                                                    </tr>
                                                </thead>
                                                <tbody>
                                                    {sampleTopSellingData.map((product, index) => (
                                                        <tr key={index}>
                                                            <td>{product.name}</td>
                                                            <td>{formatCurrency(product.price)}</td>
                                                            <td>{'★'.repeat(product.rating)}</td>
                                                        </tr>
                                                    ))}
                                                </tbody>
                                            </table>
                                        </div>
                                    </div>

                                </div>
                            </div>

                        )}

                        <div class="columns pt-5">
                            <div class="column is-half">
                                <Link class="button is-fullwidth-mobile" to={`/biometrics`}><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to Biometrics</Link>
                            </div>
                            <div class="column is-half has-text-right">
                                {/*
                        <Link to={`/account/wearable-tech`} class="button is-success is-fullwidth-mobile">
                            <FontAwesomeIcon className="fas" icon={faPlus} />&nbsp;Register Wearable
                        </Link>
                        */}
                            </div>
                        </div>

                    </nav>
                </section>
            </div>
        </>
    );
}

export default MemberLeaderboardPersonal;
