import React, { useState, useEffect } from "react";
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
    faPercent,
    faAdd,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import FormErrorBox from "../../../Reusable/FormErrorBox";
import { getGoogleFitDataPointListAPI } from "../../../../API/GoogleFitDataPoint";
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
    dataPointFilterIsStepCountDeltaState
} from "../../../../AppState";
import FormMultiSelectField from "../../../Reusable/FormMultiSelectField";
import FormCheckboxField from "../../../Reusable/FormCheckboxField";
import PageLoadingContent from "../../../Reusable/PageLoadingContent";
import FormInputFieldWithButton from "../../../Reusable/FormInputFieldWithButton";
import { PAGE_SIZE_OPTIONS } from "../../../../Constants/FieldOptions";
import { RANK_POINT_PERIOD_DAY, RANK_POINT_PERIOD_WEEK, RANK_POINT_PERIOD_MONTH, RANK_POINT_PERIOD_YEAR, RANK_POINT_FUNCTION_AVERAGE,
RANK_POINT_FUNCTION_SUM } from "../../../../Constants/App";
import MemberHistoricalDataTabularListDesktop from "./TabularListDesktop";
import MemberHistoricalDataTabularListMobile from "./TabularListMobile";
import {
  RANK_POINT_METRIC_TYPE_HEART_RATE,
  RANK_POINT_METRIC_TYPE_STEP_COUNTER,
} from "../../../../Constants/App";


function MemberHistoricalDataTabularList() {
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
    const [period, setPeriod] = useState(RANK_POINT_PERIOD_DAY);
    const [selectedDataTypes, setSelectedDataTypes] = useState([]);

    ////
    //// API.
    ////

    function onDataPointistSuccess(response) {
        console.log("onDataPointistSuccess: Starting...");
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

    function onDataPointistError(apiErr) {
        console.log("onDataPointistError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onDataPointistDone() {
        console.log("onDataPointistDone: Starting...");
        setFetching(false);
    }

    ////
    //// Event handling.
    ////

    const fetchList = (user, cur, limit, keywords, stat, sbv, selectedDataTypes=[], p) => {
        setFetching(true);
        setErrors({});

        let params = new Map();
        if (cur !== "") {
            params.set("cursor", cur); // Pagination
        }
        params.set("page_size", limit);
        params.set("sort_field", "start_at");
        params.set("sort_order","DESC");
        params.set("period", p);
        params.set("user_id", user.id);
        var values = selectedDataTypes.map(item => item).join(',');
        params.set("metric_ids",values);

        console.log("params:", params);

        // Make the submission to the API backend.
        getGoogleFitDataPointListAPI(
          params,
          onDataPointistSuccess,
          onDataPointistError,
          onDataPointistDone
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
        setSort("created_at,DESC");
        setStatus(0);
    }

    ////
    //// Misc.
    ////

    useEffect(() => {
        let mounted = true;

        if (mounted) {
          window.scrollTo(0, 0); // Start the page at the top of the page.
          fetchList(currentUser, currentCursor, pageSize, actualSearchText, status, sort, selectedDataTypes, period);
        }

        return () => {
          mounted = false;
        };
    }, [currentUser, currentCursor, pageSize, actualSearchText, status, sort, selectedDataTypes, period]);

    ////
    //// Component rendering.
    ////

    // The following block of code will generate the dataTypes we can filter by.
    let dataTypes = [];
    if (currentUser) {
        dataTypes.push({
            label: "Activity",
            value: currentUser.primaryHealthTrackingDevice.activitySegmentMetricId
        });
        dataTypes.push({
            label: "Basal Metabolic Rate",
            value: currentUser.primaryHealthTrackingDevice.basalMetabolicRateMetricId
        });
        dataTypes.push({
            label: "Blood Glucose",
            value: currentUser.primaryHealthTrackingDevice.bloodGlucoseMetricId
        });
        dataTypes.push({
            label: "Blood Pressure",
            value: currentUser.primaryHealthTrackingDevice.bloodPressureMetricId
        });
        dataTypes.push({
            label: "Steps Delta",
            value: currentUser.primaryHealthTrackingDevice.stepCountDeltaMetricId
        });
        dataTypes.push({
            label: "Body Temperature",
            value: currentUser.primaryHealthTrackingDevice.bodyTemperaturePercentageMetricId
        });
        dataTypes.push({
            label: "Calories Burned",
            value: currentUser.primaryHealthTrackingDevice.caloriesBurnedMetricId
        });
        dataTypes.push({
            label: "Cycling Pedaling Cadence",
            value: currentUser.primaryHealthTrackingDevice.cyclingPedalingCadenceMetricId
        });
        dataTypes.push({
            label: "Cycling Pedaling Cumulative",
            value: currentUser.primaryHealthTrackingDevice.cyclingPedalingCumulativeMetricId
        });
        dataTypes.push({
            label: "cycling Wheel Revolution Cumulative",
            value: currentUser.primaryHealthTrackingDevice.cyclingWheelRevolutionCumulativeMetricId
        });
        dataTypes.push({
            label: "Cycling Wheel Revolution RPM",
            value: currentUser.primaryHealthTrackingDevice.cyclingWheelRevolutionRpmMetricId
        });
        dataTypes.push({
            label: "Distance Delta",
            value: currentUser.primaryHealthTrackingDevice.distanceDeltaMetricId
        });
        dataTypes.push({
            label: "Heart Points",
            value: currentUser.primaryHealthTrackingDevice.heartPointsId
        });
        dataTypes.push({
            label: "Heart Rate (BPM)",
            value: currentUser.primaryHealthTrackingDevice.heartRateBpmMetricId
        });
        dataTypes.push({
            label: "Height",
            value: currentUser.primaryHealthTrackingDevice.heightMetricId
        });
        dataTypes.push({
            label: "Hydration",
            value: currentUser.primaryHealthTrackingDevice.hydrationMetricId
        });
        dataTypes.push({
            label: "Location Sample",
            value: currentUser.primaryHealthTrackingDevice.locationSampleMetricId
        });
        dataTypes.push({
            label: "Move Minutes",
            value: currentUser.primaryHealthTrackingDevice.moveMinutesMetricId
        });
        dataTypes.push({
            label: "Nutrition",
            value: currentUser.primaryHealthTrackingDevice.nutritionMetricId
        });
        dataTypes.push({
            label: "Oxygen Saturation",
            value: currentUser.primaryHealthTrackingDevice.oxygenSaturationMetricId
        });
        dataTypes.push({
            label: "Power",
            value: currentUser.primaryHealthTrackingDevice.powerMetricId
        });
        dataTypes.push({
            label: "Sleep",
            value: currentUser.primaryHealthTrackingDevice.sleepMetricId
        });
        dataTypes.push({
            label: "Speed",
            value: currentUser.primaryHealthTrackingDevice.speedMetricId
        });
        dataTypes.push({
            label: "Steps Counter (Cadence)",
            value: currentUser.primaryHealthTrackingDevice.stepCountCadenceMetricId
        });
        dataTypes.push({
            label: "Steps Counter (Delta)",
            value: currentUser.primaryHealthTrackingDevice.stepCountDeltaMetricId
        });
        dataTypes.push({
            label: "Weight",
            value: currentUser.primaryHealthTrackingDevice.weightMetricId
        });
        dataTypes.push({
            label: "Workout",
            value: currentUser.primaryHealthTrackingDevice.workoutMetricId
        });
    }
    console.log("dataTypes:", dataTypes);

    return (
    <>
        <div className="container">
            <section className="section">

              {/* Desktop Breadcrumbs */}
              <nav className="breadcrumb has-background-light is-hidden-touch p-4" aria-label="breadcrumbs">
                <ul>
                  <li className=""><Link to="/dashboard" aria-current="page"><FontAwesomeIcon className="fas" icon={faGauge} />&nbsp;Dashboard</Link></li>
                  <li className=""><Link to="/biometrics"><FontAwesomeIcon className="fas" icon={faHeartbeat} />&nbsp;Biometrics</Link></li>
                  <li className="is-active"><Link aria-current="page"><FontAwesomeIcon className="fas" icon={faRankingStar} />&nbsp;Leadboard</Link></li>
                </ul>
              </nav>

              {/* Mobile Breadcrumbs */}
              <nav class="breadcrumb has-background-light is-hidden-desktop p-4">
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
                      &nbsp;My History
                    </h1>
                  </div>
                  <div className="column has-text-right">
                      <button onClick={
                          ()=>fetchList(currentUser, currentCursor, pageSize, actualSearchText, status, sort, selectedDataTypes, period)
                      } class="is-fullwidth-mobile button is-link is-small" type="button">
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
                  <>
                    <FormErrorBox errors={errors} />

                    {/* Section for selecting `metric type` */}
                    <div class="column ">

                        <FormMultiSelectField
                          label="Data Types"
                          name="selectedDataTypes"
                          placeholder="Text input"
                          options={dataTypes}
                          selectedValues={selectedDataTypes}
                          onChange={(e) => {
                            let values = [];
                            for (let option of e) {
                              values.push(option.value);
                            }
                            setSelectedDataTypes(values);
                          }}
                          errorText={errors && errors.paymentMethods}
                          helpText=""
                          isRequired={true}
                          maxWidth="320px"
                        />

                    {/*
                                DEVELOPERS NOTE:
                                - As we add more sensors, add your new sensors here...
                            */}
                    </div>{/* Section for selecting `function` */}



                    {listRank &&
                    listRank.results &&
                    (listRank.results.length > 0 || previousCursors.length > 0) ? (
                      <div className="container">


                        {/*
                            ##################################################################
                            EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A DESKTOP SCREEN.
                            ##################################################################
                        */}
                        <div class="is-hidden-touch" >
                            <MemberHistoricalDataTabularListDesktop
                                listRank={listRank}
                                setPageSize={setPageSize}
                                pageSize={pageSize}
                                previousCursors={previousCursors}
                                onPreviousClicked={onPreviousClicked}
                                onNextClicked={onNextClicked}
                                currentUser={currentUser}
                            />
                        </div>

                        {/*
                            ###########################################################################
                            EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A TABLET OR MOBILE SCREEN.
                            ###########################################################################
                        */}
                        <div class="is-fullwidth is-hidden-desktop">
                            <MemberHistoricalDataTabularListMobile
                                listRank={listRank}
                                setPageSize={setPageSize}
                                pageSize={pageSize}
                                previousCursors={previousCursors}
                                onPreviousClicked={onPreviousClicked}
                                onNextClicked={onNextClicked}
                                currentUser={currentUser}
                            />
                        </div>

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

export default MemberHistoricalDataTabularList;
