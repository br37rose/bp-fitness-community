import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
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
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import FormErrorBox from "../../../Reusable/FormErrorBox";
import { getDataPointListAPI } from "../../../../API/DataPoint";
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
} from "../../../../AppState";
import PageLoadingContent from "../../../Reusable/PageLoadingContent";
import FormInputFieldWithButton from "../../../Reusable/FormInputFieldWithButton";
import { PAGE_SIZE_OPTIONS } from "../../../../Constants/FieldOptions";
import MemberDataPointTabularHistoricalGraphicalListDesktop from "./TabularListDesktop";
import MemberDataPointTabularHistoricalGraphicalListMobile from "./TabularListMobile";
import Example from "../../../Reusable/Charts/MultiseriesCharts";
import FormCheckboxField from "../../../Reusable/FormCheckboxField";
import StepHeartChart from "../../../Reusable/Charts/MultiseriesCharts";
import Charts from "../../../Reusable/Charts/MultiseriesCharts";
import ChartBuilder from "../../../Reusable/Charts/MultiseriesCharts";


function MemberDataPointTabularHistoricalGraphicalList() {
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
  const [isHeartRate, setIsHeartRate] = useRecoilState(dataPointFilterIsHeartRateState);
  const [isStepsCounter, setIsStepsCounter] = useRecoilState(dataPointFilterIsStepsCounterState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [listData, setListData] = useState("");
  const [selectedFitnessPlanForDeletion, setSelectedFitnessPlanForDeletion] = useState("");
  const [isFetching, setFetching] = useState(false);
  const [pageSize, setPageSize] = useState(100); // Pagination
  const [previousCursors, setPreviousCursors] = useState([]); // Pagination
  const [nextCursor, setNextCursor] = useState(""); // Pagination
  const [currentCursor, setCurrentCursor] = useState(""); // Pagination

  ////
  //// API.
  ////

  function onFitnessPlanListSuccess(response) {
    console.log("onFitnessPlanListSuccess: Starting...");
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

  function onFitnessPlanListError(apiErr) {
    console.log("onFitnessPlanListError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onFitnessPlanListDone() {
    console.log("onFitnessPlanListDone: Starting...");
    setFetching(false);
  }

  ////
  //// Event handling.
  ////

  const fetchList = (user, cur, limit, keywords, stat, sbv, isHeartRate, isStepsCounter) => {
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

    // Filtering
    if (keywords !== undefined && keywords !== null && keywords !== "") {
      // Searhcing
      params.set("search", keywords);
    }

    params.set("status", stat);

    if (user !== undefined && user !== null && user !== "") {
      if (isHeartRate === true) {
        params.set("heart_rate_id", user.primaryHealthTrackingDeviceHeartRateMetricId);
      }
      if (isStepsCounter === true) {
        params.set("steps_counter_id", user.primaryHealthTrackingDeviceStepsCountMetricId);
      }
    }

    getDataPointListAPI(
      params,
      onFitnessPlanListSuccess,
      onFitnessPlanListError,
      onFitnessPlanListDone
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
  }

  ////
  //// Misc.
  ////

  console.log(listData);

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.
      fetchList(currentUser, currentCursor, pageSize, actualSearchText, status, sort, isHeartRate, isStepsCounter);
    }

    return () => {
      mounted = false;
    };
  }, [currentUser, currentCursor, pageSize, actualSearchText, status, sort, isHeartRate, isStepsCounter]);

  ////
  //// Component rendering.
  ////
  console.log(listData)
  return (
    <>
      <div className="container">
        <section className="section">

          {/* Desktop Breadcrumbs */}
          <nav className="breadcrumb is-hidden-touch" aria-label="breadcrumbs">
            <ul>
              <li className="">
                <Link to="/dashboard" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faGauge} />
                  &nbsp;Dashboard
                </Link>
              </li>
              <li className="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faHeartbeat} />
                  &nbsp;Biometrics
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to="/dashboard" aria-current="page"><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to Dashboard
                </Link>
              </li>
            </ul>
          </nav>

          {/* Page */}
          <nav className="box">
            <div className="columns">
              <div className="column">
                <h1 className="title is-4">
                  <FontAwesomeIcon className="fas" icon={faHeartbeat} />
                  &nbsp;Biometrics - My History
                </h1>
              </div>
              <div className="column has-text-right">
                <button onClick={() => fetchList(currentUser, currentCursor, pageSize, actualSearchText, status, sort, isHeartRate, isStepsCounter)} class="is-fullwidth-mobile button is-link is-small" type="button">
                  <FontAwesomeIcon className="mdi" icon={faRefresh} />&nbsp;<span class="is-hidden-desktop is-hidden-tablet">Refresh</span>
                </button>
                &nbsp;
                <button onClick={(e) => setShowFilter(!showFilter)} class="is-fullwidth-mobile button is-small is-primary" type="button">
                  <FontAwesomeIcon className="mdi" icon={faFilter} />&nbsp;Filter
                </button>
                &nbsp;
                {/*
                    DEVELOPERS NOTE:
                    - If the logged in user doesn't have a device registered then
                      show the following button to encourage them to register.
                */}
                {currentUser !== undefined && currentUser !== null && currentUser !== "" && <>
                    {currentUser.primaryHealthTrackingDeviceType === 0 && <>
                      &nbsp;
                      <Link to={`/account/wearable-tech`} className="is-fullwidth-mobile button is-small is-success" type="button">
                          <FontAwesomeIcon className="mdi" icon={faPlus} />&nbsp;Register Wearable
                      </Link>
                    </>}
                </>}
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
                </div>

                {/* Filter Options */}
                <div class="columns">
                  <div class="column">
                    <FormCheckboxField
                      label="Heart Rate"
                      name="isHeartRate"
                      checked={isHeartRate}
                      errorText={errors && errors.isHeartRate}
                      onChange={(e) => { setIsHeartRate(!isHeartRate) }}
                      maxWidth="180px"
                    />
                  </div>
                  <div class="column">
                    <FormCheckboxField
                      label="Steps Counter"
                      name="isStepsCounter"
                      checked={isStepsCounter}
                      errorText={errors && errors.isStepsCounter}
                      onChange={(e) => { setIsStepsCounter(!isStepsCounter) }}
                      maxWidth="180px"
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

                    <div class="column has-text-right">
                      <Link class="button is-small" type="button" to="/biometrics/history/tableview">
                        <FontAwesomeIcon className="mdi" icon={faTable} />
                      </Link>
                      <button class="button is-small is-info" type="button">
                        <FontAwesomeIcon className="mdi" icon={faChartLine} />
                      </button>
                      &nbsp;
                    </div>


                    <ChartBuilder
                      data={listData.results}
                    />


                  </div>
                ) : (
                  <section className="hero is-medium has-background-white-ter">
                    <div className="hero-body">
                      <p className="title">
                        <FontAwesomeIcon className="fas" icon={faTable} />
                        &nbsp;No Biometrics
                      </p>
                      <p className="subtitle">
                        You currently have no biometrics data.{" "}
                        <b>
                          <Link to="/account/wearable-tech">
                            Click here&nbsp;
                            <FontAwesomeIcon
                              className="mdi"
                              icon={faArrowRight}
                            />
                          </Link>
                        </b>{" "}
                        to get started registering your wearable tech!
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
                <Link to={`/account/wearable-tech`} class="button is-success is-fullwidth-mobile">
                  <FontAwesomeIcon className="fas" icon={faPlus} />&nbsp;Register Wearable
                </Link>
              </div>
            </div>

          </nav>
        </section>
      </div>
    </>
  );
}

export default MemberDataPointTabularHistoricalGraphicalList;
