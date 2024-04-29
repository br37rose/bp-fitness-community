import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faFilterCircleXmark,
  faTrophy,
  faArrowLeft,
  faPlus,
  faGauge,
  faArrowRight,
  faTable,
  faRefresh,
  faFilter,
  faSearch,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import FormErrorBox from "../../Reusable/FormErrorBox";
import {
  getFitnessPlanListAPI,
  deleteFitnessPlanAPI,
} from "../../../API/FitnessPlan";
import {
  topAlertMessageState,
  topAlertStatusState,
  currentUserState,
  videoCategoryFilterShowState,
  videoCategoryFilterTemporarySearchTextState,
  videoCategoryFilterActualSearchTextState,
  videoCategoryFilterSortState,
  videoCategoryFilterStatusState,
} from "../../../AppState";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import FormInputFieldWithButton from "../../Reusable/FormInputFieldWithButton";
import AdminFitnessPlanListDesktop from "./ListDesktop";
import AdminFitnessPlanListMobile from "./ListMobile";
import Layout from "../../Menu/Layout";

function AdminFitnessPlanList() {
  ////
  //// Global state.
  ////

  const [topAlertMessage, setTopAlertMessage] =
    useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] =
    useRecoilState(topAlertStatusState);
  const [currentUser] = useRecoilState(currentUserState);
  const [showFilter, setShowFilter] = useRecoilState(
    videoCategoryFilterShowState
  ); // Filtering + Searching
  const [sort, setSort] = useRecoilState(videoCategoryFilterSortState); // Sorting
  const [temporarySearchText, setTemporarySearchText] = useRecoilState(
    videoCategoryFilterTemporarySearchTextState
  ); // Searching - The search field value as your writes their query.
  const [actualSearchText, setActualSearchText] = useRecoilState(
    videoCategoryFilterActualSearchTextState
  ); // Searching - The actual search query value to submit to the API.
  const [status, setStatus] = useRecoilState(videoCategoryFilterStatusState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [listData, setListData] = useState("");
  const [selectedFitnessPlanForDeletion, setSelectedFitnessPlanForDeletion] =
    useState("");
  const [isFetching, setFetching] = useState(false);
  const [pageSize, setPageSize] = useState(10); // Pagination
  const [previousCursors, setPreviousCursors] = useState([]); // Pagination
  const [nextCursor, setNextCursor] = useState(""); // Pagination
  const [currentCursor, setCurrentCursor] = useState(""); // Pagination

  ////
  //// API.
  ////

  function onFitnessPlanListSuccess(response) {
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
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onFitnessPlanListDone() {
    setFetching(false);
  }

  function onFitnessPlanDeleteSuccess(response) {
    // Update notification.
    setTopAlertStatus("success");
    setTopAlertMessage("FitnessPlan deleted");
    setTimeout(() => {
      setTopAlertMessage("");
    }, 2000);

    // Fetch again an updated list.
    fetchList(currentCursor, pageSize, actualSearchText, status, sort);
  }

  function onFitnessPlanDeleteError(apiErr) {
    setErrors(apiErr);

    // Update notification.
    setTopAlertStatus("danger");
    setTopAlertMessage("Failed deleting");
    setTimeout(() => {
      setTopAlertMessage("");
    }, 2000);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onFitnessPlanDeleteDone() {
    setFetching(false);
  }

  ////
  //// BREADCRUMB
  ////
  const breadcrumbItems = {
    items: [
      { text: "Dashboard", link: "/dashboard", isActive: false, icon: faGauge },
      { text: "Fitness Plans", link: "#", icon: faTrophy, isActive: true },
    ],
    mobileBackLinkItems: {
      link: "/dashboard",
      text: "Back to Dashboard",
      icon: faArrowLeft,
    },
  };

  ////
  //// Event handling.
  ////

  const fetchList = (cur, limit, keywords, stat, sbv) => {
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

    getFitnessPlanListAPI(
      params,
      onFitnessPlanListSuccess,
      onFitnessPlanListError,
      onFitnessPlanListDone
    );
  };

  const onNextClicked = (e) => {
    let arr = [...previousCursors];
    arr.push(currentCursor);
    setPreviousCursors(arr);
    setCurrentCursor(nextCursor);
  };

  const onPreviousClicked = (e) => {
    let arr = [...previousCursors];
    const previousCursor = arr.pop();
    setPreviousCursors(arr);
    setCurrentCursor(previousCursor);
  };

  const onSearchButtonClick = (e) => {
    // Searching
    setActualSearchText(temporarySearchText);
  };

  const onSelectFitnessPlanForDeletion = (e, datum) => {
    setSelectedFitnessPlanForDeletion(datum);
  };

  const onDeselectFitnessPlanForDeletion = (e) => {
    setSelectedFitnessPlanForDeletion("");
  };

  const onDeleteConfirmButtonClick = (e) => {
    deleteFitnessPlanAPI(
      selectedFitnessPlanForDeletion.id,
      onFitnessPlanDeleteSuccess,
      onFitnessPlanDeleteError,
      onFitnessPlanDeleteDone
    );
    setSelectedFitnessPlanForDeletion("");
  };

  // Function resets the filter state to its default state.
  const onClearFilterClick = (e) => {
    setShowFilter(false);
    setActualSearchText("");
    setTemporarySearchText("");
    setSort("no,1");
    setStatus(0);
  };

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.
      fetchList(currentCursor, pageSize, actualSearchText, status, sort);
    }

    return () => {
      mounted = false;
    };
  }, [currentCursor, pageSize, actualSearchText, status, sort]);

  ////
  //// Component rendering.
  ////

  return (
    <Layout breadcrumbItems={breadcrumbItems}>
      {/* Page */}
      <nav className="box">
        <div
          className={`modal ${
            selectedFitnessPlanForDeletion ? "is-active" : ""
          }`}
        >
          <div className="modal-background"></div>
          <div className="modal-card">
            <header className="modal-card-head">
              <p className="modal-card-title">Are you sure?</p>
              <button
                className="delete"
                aria-label="close"
                onClick={onDeselectFitnessPlanForDeletion}
              ></button>
            </header>
            <section className="modal-card-body">
              You are about to <b>archive</b> this Fitness Plan; it will no
              longer appear on your dashboard. This action cannot be undone. Are
              you sure you want to continue?
            </section>
            <footer className="modal-card-foot">
              <button
                className="button is-success"
                onClick={onDeleteConfirmButtonClick}
              >
                Confirm
              </button>
              <button
                className="button"
                onClick={onDeselectFitnessPlanForDeletion}
              >
                Cancel
              </button>
            </footer>
          </div>
        </div>

        <div className="columns">
          <div className="column">
            <h1 className="title is-4">
              <FontAwesomeIcon className="fas" icon={faTrophy} />
              &nbsp;Fitness Plans
            </h1>
          </div>
          <div className="column has-text-right">
            <button
              onClick={() =>
                fetchList(
                  currentCursor,
                  pageSize,
                  actualSearchText,
                  status,
                  sort
                )
              }
              class="is-fullwidth-mobile button is-link is-small"
              type="button"
            >
              <FontAwesomeIcon className="mdi" icon={faRefresh} />
              &nbsp;
              <span class="is-hidden-desktop is-hidden-tablet">Refresh</span>
            </button>
            &nbsp;
            <button
              onClick={(e) => setShowFilter(!showFilter)}
              class="is-fullwidth-mobile button is-small is-primary"
              type="button"
            >
              <FontAwesomeIcon className="mdi" icon={faFilter} />
              &nbsp;Filter
            </button>
            &nbsp;
            <Link
              to={`/fitness-plans/add`}
              className="is-fullwidth-mobile button is-small is-success"
              type="button"
            >
              <FontAwesomeIcon className="mdi" icon={faPlus} />
              &nbsp;Request Plan
            </Link>
          </div>
        </div>

        {/* FILTER */}
        {showFilter && (
          <div
            class="has-background-white-bis"
            style={{ borderRadius: "15px", padding: "20px" }}
          >
            {/* Filter Title + Clear Button */}
            <div class="columns is-mobile">
              <div class="column is-half">
                <strong>
                  <u>
                    <FontAwesomeIcon className="mdi" icon={faFilter} />
                    &nbsp;Filter
                  </u>
                </strong>
              </div>
              <div class="column is-half has-text-right">
                <Link onClick={onClearFilterClick}>
                  <FontAwesomeIcon className="mdi" icon={faFilterCircleXmark} />
                  &nbsp;Clear Filter
                </Link>
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
                <div class="is-hidden-touch">
                  <AdminFitnessPlanListDesktop
                    listData={listData}
                    setPageSize={setPageSize}
                    pageSize={pageSize}
                    previousCursors={previousCursors}
                    onPreviousClicked={onPreviousClicked}
                    onNextClicked={onNextClicked}
                    onSelectFitnessPlanForDeletion={
                      onSelectFitnessPlanForDeletion
                    }
                  />
                </div>

                {/*
                            ###########################################################################
                            EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A TABLET OR MOBILE SCREEN.
                            ###########################################################################
                        */}
                <div class="is-fullwidth is-hidden-desktop">
                  <AdminFitnessPlanListMobile
                    listData={listData}
                    setPageSize={setPageSize}
                    pageSize={pageSize}
                    previousCursors={previousCursors}
                    onPreviousClicked={onPreviousClicked}
                    onNextClicked={onNextClicked}
                    onSelectFitnessPlanForDeletion={
                      onSelectFitnessPlanForDeletion
                    }
                  />
                </div>
              </div>
            ) : (
              <section className="hero is-medium has-background-white-ter">
                <div className="hero-body">
                  <p className="title">
                    <FontAwesomeIcon className="fas" icon={faTable} />
                    &nbsp;No Fitness Plans
                  </p>
                  <p className="subtitle">
                    You currently have no fitness plans.{" "}
                    <b>
                      <Link to="/fitness-plans/add">
                        Click here&nbsp;
                        <FontAwesomeIcon className="mdi" icon={faArrowRight} />
                      </Link>
                    </b>{" "}
                    to get started requesting your first fitness plan from our
                    team!
                  </p>
                </div>
              </section>
            )}
          </>
        )}

        <div class="columns pt-5">
          <div class="column is-half">
            <Link class="button is-fullwidth-mobile" to={`/dashboard`}>
              <FontAwesomeIcon className="fas" icon={faArrowLeft} />
              &nbsp;Back to Dashboard
            </Link>
          </div>
          <div class="column is-half has-text-right">
            <Link
              to={`/fitness-plans/add`}
              class="button is-success is-fullwidth-mobile"
            >
              <FontAwesomeIcon className="fas" icon={faPlus} />
              &nbsp;Request Plan
            </Link>
          </div>
        </div>
      </nav>
    </Layout>
  );
}

export default AdminFitnessPlanList;
