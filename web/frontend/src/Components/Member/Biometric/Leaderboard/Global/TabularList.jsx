import React, { useState, useEffect, useMemo } from "react";
import { Link } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faAdd,
  faPercent,
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
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import FormErrorBox from "../../../../Reusable/FormErrorBox";
import { getLeaderboardListAPI } from "../../../../../API/RankPoint";
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
import {
  RANK_POINT_METRIC_TYPE_HEART_RATE,
  RANK_POINT_METRIC_TYPE_STEP_COUNTER,
  RANK_POINT_PERIOD_DAY,
  RANK_POINT_PERIOD_WEEK,
  RANK_POINT_PERIOD_MONTH,
  RANK_POINT_PERIOD_YEAR,
  RANK_POINT_FUNCTION_AVERAGE,
  RANK_POINT_FUNCTION_SUM
} from "../../../../../Constants/App";
import MemberLeaderboardGlobalTabularListDesktop from "./TabularListDesktop";
import MemberLeaderboardGlobalTabularListMobile from "./TabularListMobile";
import Layout from "../../../../Menu/Layout";


function MemberLeaderboardGlobalTabularList() {
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
  const [calcFunction, setCalcFunction] = useState(RANK_POINT_FUNCTION_AVERAGE);
  const [firstApiResponseData, setFirstApiResponseData] = useState(null);
  const [secondApiResponseData, setSecondApiResponseData] = useState(null);

  ////
  //// API.
  ////

  function onRankPointistSuccess(firstApiResponse, secondApiResponse) {
    console.log("onRankPointistSuccess: Starting...");
    if (firstApiResponse !== null && secondApiResponse !== null) {
      setFirstApiResponseData(firstApiResponse);
      setSecondApiResponseData(secondApiResponse);
      //   if (response.hasNextPage) {
      //     setNextCursor(response.nextCursor); // For pagination purposes.
      //   }
      // } else {
      //   setListRank([]);
      //   setNextCursor("");
      // }
      setListRank([]);
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
  //// BREADCRUMB
  ////
  const breadcrumbItems = {
    items: [
      { text: 'Dashboard', link: '/dashboard', isActive: false, icon: faGauge },
      { text: 'Biometrics', link: '/biometrics', icon: faHeartbeat, isActive: false },
      { text: 'Leaderboard', link: '#', icon: faRankingStar, isActive: true }
    ],
    mobileBackLinkItems: {
      link: "/biometrics",
      text: "Back to Biometrics",
      icon: faArrowLeft
    }
  }

  const fetchList = (user, cur, limit, keywords, stat, sbv, isHeartRate, isStepsCounter, p, cf) => {
    setFetching(true);
    setErrors({});

    let metricType;
    if (isHeartRate) {
      metricType = RANK_POINT_METRIC_TYPE_HEART_RATE;
    } else if (isStepsCounter) {
      metricType = RANK_POINT_METRIC_TYPE_STEP_COUNTER;
    }

    let firstParams = new Map([
      ["page_size", limit],
      ["sort_field", "place"],
      ["sort_order", "ASC"],
      // ... other parameters ...
      ["metric_type", metricType],
      ["period", p],
      ["function", cf],
    ]);

    let secondParams = new Map([
      ["page_size", limit],
      ["sort_field", "place"],
      ["sort_order", "ASC"],
      // ... other parameters ...
      ["metric_type", metricType],
      ["period", p],
      ["function", cf],
    ]);

    // Return the API call promise
    return Promise.all([
      getLeaderboardListAPI(firstParams),
      getLeaderboardListAPI(secondParams),
    ]);
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

  const onHeartRateButtonClick = (e) => {
    e.preventDefault(); // Do not remove this line!
    setIsStepsCounter(!isStepsCounter);
    setIsHeartRate(!isHeartRate)
    setCalcFunction(RANK_POINT_FUNCTION_AVERAGE);
  }

  const onStepCounterButtonClick = (e) => {
    e.preventDefault(); // Do not remove this line!
    setIsStepsCounter(!isStepsCounter);
    setIsHeartRate(!isHeartRate);
    setCalcFunction(RANK_POINT_FUNCTION_AVERAGE);
  }

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0);
      fetchList(currentUser, currentCursor, pageSize, actualSearchText, status, sort, isHeartRate, isStepsCounter, period, calcFunction)
        .then(([firstApiResponse, secondApiResponse]) => {
          if (mounted) {
            onRankPointistSuccess(firstApiResponse, secondApiResponse)
          }
        })
        .catch(error => {
          if (mounted) {
            onRankPointistError(error);
          }
        })
        .finally(() => {
          if (mounted) {
            onRankPointistDone();
          }
        });
    }
  }, [currentUser, currentCursor, pageSize, actualSearchText, status, sort, isHeartRate, isStepsCounter, period, calcFunction]);

  ////
  //// Component rendering.
  ////

  const mergedData = useMemo(() => {
    if (firstApiResponseData && secondApiResponseData) {
      return firstApiResponseData.results.map(item => {
        const weeklyAvgItem = secondApiResponseData.results.find(secondItem => secondItem.userId === item.userId);

        return {
          ...item, // all data from the first API response
          weeklyAvg: weeklyAvgItem ? weeklyAvgItem.value : null, // value from the second API response
          weeklyAvgStart: weeklyAvgItem ? weeklyAvgItem.start : null,
          weeklyAvgEnd: weeklyAvgItem ? weeklyAvgItem.end : null
        };
      });
    }
    return [];
  }, [firstApiResponseData, secondApiResponseData]);

  return (
    <Layout breadcrumbItems={breadcrumbItems}>
      <div className="box">
        <div className="columns">
          <div className="column">
            <h1 className="title is-4">
              <FontAwesomeIcon className="fas" icon={faRankingStar} />
              &nbsp;Leaderboard
            </h1>
          </div>
          <div className="column has-text-right">
            <button onClick={() => {
              setFetching(true);
              setErrors({});
              fetchList(currentUser, currentCursor, pageSize, actualSearchText, status, sort, isHeartRate, isStepsCounter, period, calcFunction)
                .then(([firstApiResponse, secondApiResponse]) => {
                  onRankPointistSuccess(firstApiResponse, secondApiResponse)
                })
                .catch(error => {
                  onRankPointistError(error);
                })
                .finally(() => {
                  onRankPointistDone();
                });
            }} className="is-fullwidth-mobile button is-link is-small" type="button">
              <FontAwesomeIcon className="mdi" icon={faRefresh} />&nbsp;<span className="is-hidden-desktop is-hidden-tablet">Refresh</span>
            </button>
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

        {isFetching ? (
          <PageLoadingContent displayMessage={"Please wait..."} />
        ) : (
          <>
            <FormErrorBox errors={errors} />

            {/* Section for selecting `metric type` */}
            <div class="column has-text-right">
              <button class={`button is-small ${isHeartRate && `is-info`}`} type="button" onClick={(e) => { onHeartRateButtonClick(e) }}>
                <FontAwesomeIcon className="mdi" icon={faHeart} />&nbsp;Heart Rate
              </button>
              <Link class={`button is-small ${isStepsCounter && `is-info`}`} type="button" onClick={(e) => { onStepCounterButtonClick(e) }}>
                <FontAwesomeIcon className="mdi" icon={faPersonWalking} />&nbsp;Steps Count
              </Link>&nbsp;
              {/*
                        DEVELOPERS NOTE:
                        - As we add more sensors, add your new sensors here...
                    */}
            </div>

            {/* Section for selecting `function` */}
            {/* <div class="column has-text-right">

                        DEVELOPERS NOTE:
                        - Some functions are not available for some metric types because it makes sense. Why would you keep a summation of heart rate? This
                        is the reason for code restrictions below.
                        - The only functions you should use are `RANK_POINT_FUNCTION_AVERAGE` and `RANK_POINT_FUNCTION_SUM`.

                  <button class={`button is-small ${calcFunction === RANK_POINT_FUNCTION_AVERAGE && `is-info`}`} type="button" onClick={(e) => { setCalcFunction(RANK_POINT_FUNCTION_AVERAGE) }}>
                    <FontAwesomeIcon className="mdi" icon={faPercent} />&nbsp;Average
                  </button>
                  {!isHeartRate && <Link class={`button is-small ${calcFunction === RANK_POINT_FUNCTION_SUM && `is-info`}`} type="button" onClick={(e) => { setCalcFunction(RANK_POINT_FUNCTION_SUM) }}>
                    <FontAwesomeIcon className="mdi" icon={faAdd} />&nbsp;Sum
                  </Link>}&nbsp;
                </div> */}

            {/* Section for selecting `period` */}
            <div class="column has-text-right">
              {/*
                        DEVELOPERS NOTE:
                        - Period refers to the period of time the ranking is between. For example `day` would mean ranking for today.
                        - Week is ISO week, meaning it the week starts on Sunday and ends on Saturday.
                        - Month or year ranking are only ment for THIS month or year.
                    */}
              <button class={`button is-small ${period === RANK_POINT_PERIOD_DAY && `is-info`}`} type="button" onClick={(e) => { e.preventDefault(); setPeriod(RANK_POINT_PERIOD_DAY); }}>
                Today
              </button>
              <Link class={`button is-small ${period === RANK_POINT_PERIOD_WEEK && `is-info`}`} type="button" onClick={(e) => { e.preventDefault(); setPeriod(RANK_POINT_PERIOD_WEEK); }}>
                Week
              </Link>
              <Link class={`button is-small ${period === RANK_POINT_PERIOD_MONTH && `is-info`}`} type="button" onClick={(e) => { e.preventDefault(); setPeriod(RANK_POINT_PERIOD_MONTH); }}>
                Month
              </Link>
              <Link class={`button is-small ${period === RANK_POINT_PERIOD_YEAR && `is-info`}`} type="button" onClick={(e) => { e.preventDefault(); setPeriod(RANK_POINT_PERIOD_YEAR); }}>
                Year
              </Link>
              &nbsp;
            </div>
            {firstApiResponseData &&
              firstApiResponseData.results &&
              secondApiResponseData &&
              secondApiResponseData.results &&
              (secondApiResponseData.results.length > 0 && firstApiResponseData.results.length > 0 || previousCursors.length > 0) ? (
              <div className="container">
                {/*
                            ##################################################################
                            EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A DESKTOP SCREEN.
                            ##################################################################
                        */}
                <div class="is-hidden-mobile" >
                  <MemberLeaderboardGlobalTabularListDesktop
                    data={mergedData}
                    setPageSize={setPageSize}
                    pageSize={pageSize}
                    previousCursors={previousCursors}
                    onPreviousClicked={onPreviousClicked}
                    onNextClicked={onNextClicked}
                    currentUser={currentUser}
                    period={period}
                    calcFunction={calcFunction}
                  />
                </div>

                {/*
                            ###########################################################################
                            EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A TABLET OR MOBILE SCREEN.
                            ###########################################################################
                        */}
                <div class="is-fullwidth is-hidden-tablet">
                  <MemberLeaderboardGlobalTabularListMobile
                    data={mergedData}
                    setPageSize={setPageSize}
                    pageSize={pageSize}
                    previousCursors={previousCursors}
                    onPreviousClicked={onPreviousClicked}
                    onNextClicked={onNextClicked}
                    currentUser={currentUser}
                    period={period}
                    calcFunction={calcFunction}
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

      </div>
    </Layout>
  );
}

export default MemberLeaderboardGlobalTabularList;
