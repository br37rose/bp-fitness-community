import React, { useState, useEffect } from "react";
import { Link, useParams } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faRankingStar,
  faArrowLeft,
  faGauge,
  faTable,
  faRefresh,
  faEye,
  faBolt,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import {
  topAlertMessageState,
  topAlertStatusState,
  currentUserState,
  dataPointFilterShowState,
  dataPointFilterTemporarySearchTextState,
  dataPointFilterActualSearchTextState,
  dataPointFilterSortState,
  dataPointFilterStatusState,
} from "../../../AppState";
import {
  RANK_POINT_PERIOD_DAY,
  RANK_POINT_FUNCTION_AVERAGE,
} from "../../../Constants/App";
import MemberLeaderboardGlobalTabularListDesktop from "./TabularListDesktop";
import MemberLeaderboardGlobalTabularListMobile from "./TabularListMobile";
import FormErrorBox from "../../Reusable/FormErrorBox";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import { getFitnessChalengeLeaderboard } from "../../../API/FitnessChallenge";

function MemberLeaderboardGlobalTabularListForChallenge() {
  ////
  //// Global state.
  ////

  const [topAlertMessage, setTopAlertMessage] =
    useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] =
    useRecoilState(topAlertStatusState);
  const [currentUser] = useRecoilState(currentUserState);
  const [showFilter, setShowFilter] = useRecoilState(dataPointFilterShowState); // Filtering + Searching
  const [sort, setSort] = useRecoilState(dataPointFilterSortState); // Sorting
  const [temporarySearchText, setTemporarySearchText] = useRecoilState(
    dataPointFilterTemporarySearchTextState
  ); // Searching - The search field value as your writes their query.
  const [actualSearchText, setActualSearchText] = useRecoilState(
    dataPointFilterActualSearchTextState
  ); // Searching - The actual search query value to submit to the API.
  const [status, setStatus] = useRecoilState(dataPointFilterStatusState);
  // const [isHeartRate, setIsHeartRate] = useRecoilState(dataPointFilterIsHeartRateState);
  // const [isStepsCounter, setIsStepsCounter] = useRecoilState(dataPointFilterIsStepsCounterState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [listRank, setListRank] = useState("");
  const [selectedFitnessPlanForDeletion, setSelectedFitnessPlanForDeletion] =
    useState("");
  const [isFetching, setFetching] = useState(false);
  const [pageSize, setPageSize] = useState(100); // Pagination
  const [previousCursors, setPreviousCursors] = useState([]); // Pagination
  const [nextCursor, setNextCursor] = useState(""); // Pagination
  const [currentCursor, setCurrentCursor] = useState(""); // Pagination
  const [isHeartRate, setIsHeartRate] = useState(true);
  const [isStepsCounter, setIsStepsCounter] = useState(false);
  const [period, setPeriod] = useState(RANK_POINT_PERIOD_DAY);
  const [calcFunction, setCalcFunction] = useState(RANK_POINT_FUNCTION_AVERAGE);

  const { id } = useParams();

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

  const fetchList = (
    user,
    cur,
    limit,
    keywords,
    stat,
    sbv,
    isHeartRate,
    isStepsCounter,
    p,
    cf
  ) => {
    setFetching(true);
    setErrors({});

    // Make the submission to the API backend.
    getFitnessChalengeLeaderboard(
      id,
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

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.
      fetchList(
        currentUser,
        currentCursor,
        pageSize,
        actualSearchText,
        status,
        sort,
        isHeartRate,
        isStepsCounter,
        period,
        calcFunction
      );
    }

    return () => {
      mounted = false;
    };
  }, [
    currentUser,
    currentCursor,
    pageSize,
    actualSearchText,
    status,
    sort,
    isHeartRate,
    isStepsCounter,
    period,
    calcFunction,
  ]);

  ////
  //// Component rendering.
  ////

  return (
    <>
      <div className="container is-fluid">
        <section className="section">
          {/* Desktop Breadcrumbs */}
          <nav class="breadcrumb is-hidden-touch" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to="/dashboard" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faGauge} />
                  &nbsp;Dashboard
                </Link>
              </li>
              <li class="">
                <Link to="/fitness-challenge" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faBolt} />
                  &nbsp;Challenges
                </Link>
              </li>
              <li class="">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faEye} />
                  &nbsp;Detail
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to="/fitness-challenge" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Challenges
                </Link>
              </li>
            </ul>
          </nav>

          {/* Page */}
          <nav className="box">
            <div className="columns">
              <div className="column">
                <h1 className="title is-4">
                  <FontAwesomeIcon className="fas" icon={faRankingStar} />
                  &nbsp;Leaderboard - Global
                </h1>
              </div>
              <div className="column has-text-right">
                <button
                  onClick={() =>
                    fetchList(
                      currentUser,
                      currentCursor,
                      pageSize,
                      actualSearchText,
                      status,
                      sort,
                      isHeartRate,
                      isStepsCounter,
                      period,
                      calcFunction
                    )
                  }
                  class="is-fullwidth-mobile button is-link is-small"
                  type="button"
                >
                  <FontAwesomeIcon className="mdi" icon={faRefresh} />
                  &nbsp;
                  <span class="is-hidden-desktop is-hidden-tablet">
                    Refresh
                  </span>
                </button>
                &nbsp;
              </div>
            </div>

            {isFetching ? (
              <PageLoadingContent displayMessage={"Please wait..."} />
            ) : (
              <>
                <FormErrorBox errors={errors} />
                {/* Tab Navigation */}
                <div class="tabs is-medium is-size-7-mobile">
                  <ul>
                    <li class="">
                      <Link to={`/fitness-challenge/${id}`}>
                        <strong>Detail</strong>
                      </Link>
                    </li>
                    <li class="is-active">
                      <Link>
                        <strong>LeaderBoard</strong>
                      </Link>
                    </li>
                  </ul>
                </div>

                {/* Section for selecting `metric type` */}
                <div class="column has-text-right">
                  {/*
                                DEVELOPERS NOTE:
                                - As we add more sensors, add your new sensors here...
                            */}
                </div>
                {/* Section for selecting `function` */}
                <div class="column has-text-right">
                  {/*
                                DEVELOPERS NOTE:
                                - Some functions are not available for some metric types because it makes sense. Why would you keep a summation of heart rate? This
                                is the reason for code restrictions below.
                                - The only functions you should use are `RANK_POINT_FUNCTION_AVERAGE` and `RANK_POINT_FUNCTION_SUM`.
                          */}
                </div>

                {listRank &&
                listRank.results &&
                (listRank.results.length > 0 || previousCursors.length > 0) ? (
                  <div>
                    {/*
                            ##################################################################
                            EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A DESKTOP SCREEN.
                            ##################################################################
                        */}
                    <div class="is-hidden-touch">
                      <MemberLeaderboardGlobalTabularListDesktop
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
                      <MemberLeaderboardGlobalTabularListMobile
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
                      <p className="subtitle">No biometrics data. </p>
                    </div>
                  </section>
                )}
              </>
            )}

            <div class="columns pt-5">
              <div class="column is-half">
                <Link
                  class="button is-fullwidth-mobile"
                  to={`/fitness-challenge`}
                >
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Challenges
                </Link>
              </div>
            </div>
          </nav>
        </section>
      </div>
    </>
  );
}

export default MemberLeaderboardGlobalTabularListForChallenge;
