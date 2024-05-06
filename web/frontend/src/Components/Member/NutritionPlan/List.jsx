import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faFilterCircleXmark,
  faLeaf,
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

import FormErrorBox from "../../Reusable/FormErrorBox";
import { getNutritionPlanListAPI, deleteNutritionPlanAPI } from "../../../API/NutritionPlan";
import {
  topAlertMessageState,
  topAlertStatusState,
  currentUserState,
  videoCategoryFilterShowState,
  videoCategoryFilterTemporarySearchTextState,
  videoCategoryFilterActualSearchTextState,
  videoCategoryFilterSortState,
  videoCategoryFilterStatusState
} from "../../../AppState";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import FormInputFieldWithButton from "../../Reusable/FormInputFieldWithButton";
import { PAGE_SIZE_OPTIONS } from "../../../Constants/FieldOptions";
import AdminNutritionPlanListDesktop from "./ListDesktop";
import AdminNutritionPlanListMobile from "./ListMobile";
import Layout from "../../Menu/Layout";


function AdminNutritionPlanList() {
  ////
  //// Global state.
  ////

  const [topAlertMessage, setTopAlertMessage] = useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] = useRecoilState(topAlertStatusState);
  const [currentUser] = useRecoilState(currentUserState);
  const [showFilter, setShowFilter] = useRecoilState(videoCategoryFilterShowState); // Filtering + Searching
  const [sort, setSort] = useRecoilState(videoCategoryFilterSortState); // Sorting
  const [temporarySearchText, setTemporarySearchText] = useRecoilState(videoCategoryFilterTemporarySearchTextState); // Searching - The search field value as your writes their query.
  const [actualSearchText, setActualSearchText] = useRecoilState(videoCategoryFilterActualSearchTextState); // Searching - The actual search query value to submit to the API.
  const [status, setStatus] = useRecoilState(videoCategoryFilterStatusState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [listData, setListData] = useState("");
  const [selectedNutritionPlanForDeletion, setSelectedNutritionPlanForDeletion] = useState("");
  const [isFetching, setFetching] = useState(false);
  const [pageSize, setPageSize] = useState(10); // Pagination
  const [previousCursors, setPreviousCursors] = useState([]); // Pagination
  const [nextCursor, setNextCursor] = useState(""); // Pagination
  const [currentCursor, setCurrentCursor] = useState(""); // Pagination

  ////
  //// API.
  ////

  function onNutritionPlanListSuccess(response) {
    console.log("onNutritionPlanListSuccess: Starting...");
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

  function onNutritionPlanListError(apiErr) {
    console.log("onNutritionPlanListError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onNutritionPlanListDone() {
    console.log("onNutritionPlanListDone: Starting...");
    setFetching(false);
  }

  function onNutritionPlanDeleteSuccess(response) {
    console.log("onNutritionPlanDeleteSuccess: Starting..."); // For debugging purposes only.

    // Update notification.
    setTopAlertStatus("success");
    setTopAlertMessage("NutritionPlan deleted");
    setTimeout(() => {
      console.log(
        "onDeleteConfirmButtonClick: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus
      );
      setTopAlertMessage("");
    }, 2000);

    // Fetch again an updated list.
    fetchList(currentCursor, pageSize, actualSearchText, status, sort);
  }

  function onNutritionPlanDeleteError(apiErr) {
    console.log("onNutritionPlanDeleteError: Starting..."); // For debugging purposes only.
    setErrors(apiErr);

    // Update notification.
    setTopAlertStatus("danger");
    setTopAlertMessage("Failed deleting");
    setTimeout(() => {
      console.log(
        "onNutritionPlanDeleteError: topAlertMessage, topAlertStatus:",
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

  function onNutritionPlanDeleteDone() {
    console.log("onNutritionPlanDeleteDone: Starting...");
    setFetching(false);
  }

  ////
  //// BREADCRUMB
  ////
  const breadcrumbItems = {
    items: [
      { text: 'Dashboard', link: '/dashboard', isActive: false, icon: faGauge },
      { text: 'Nutrition Plans', link: '#', icon: faLeaf, isActive: true }
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

    getNutritionPlanListAPI(
      params,
      onNutritionPlanListSuccess,
      onNutritionPlanListError,
      onNutritionPlanListDone
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

  const onSelectNutritionPlanForDeletion = (e, datum) => {
    console.log("onSelectNutritionPlanForDeletion", datum);
    setSelectedNutritionPlanForDeletion(datum);
  };

  const onDeselectNutritionPlanForDeletion = (e) => {
    console.log("onDeselectNutritionPlanForDeletion");
    setSelectedNutritionPlanForDeletion("");
  };

  const onDeleteConfirmButtonClick = (e) => {
    console.log("onDeleteConfirmButtonClick"); // For debugging purposes only.

    deleteNutritionPlanAPI(
      selectedNutritionPlanForDeletion.id,
      onNutritionPlanDeleteSuccess,
      onNutritionPlanDeleteError,
      onNutritionPlanDeleteDone
    );
    setSelectedNutritionPlanForDeletion("");
  };

  // Function resets the filter state to its default state.
  const onClearFilterClick = (e) => {
    setShowFilter(false);
    setActualSearchText("");
    setTemporarySearchText("");
    setSort("no,1");
    setStatus(0);
  }

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
      <div className="box">
        <div
          className={`modal ${selectedNutritionPlanForDeletion ? "is-active" : ""
            }`}
        >
          <div className="modal-background"></div>
          <div className="modal-card">
            <header className="modal-card-head">
              <p className="modal-card-title">Are you sure?</p>
              <button
                className="delete"
                aria-label="close"
                onClick={onDeselectNutritionPlanForDeletion}
              ></button>
            </header>
            <section className="modal-card-body">
              You are about to <b>archive</b> this Nutrition Plan; it will no longer
              appear on your dashboard nor will the Nutrition Plan be able to log
              into their account. This action can be undone but you'll need
              to contact the system administrator. Are you sure you would
              like to continue?
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
                onClick={onDeselectNutritionPlanForDeletion}
              >
                Cancel
              </button>
            </footer>
          </div>
        </div>

        <div className="columns">
          <div className="column">
            <h1 className="title is-4">
              <FontAwesomeIcon className="fas" icon={faLeaf} />
              &nbsp;Nutrition Plans
            </h1>
          </div>
          <div className="column has-text-right">
            <button onClick={() => fetchList(currentCursor, pageSize, actualSearchText, status, sort)} class="is-fullwidth-mobile button is-link is-small" type="button">
              <FontAwesomeIcon className="mdi" icon={faRefresh} />&nbsp;<span class="is-hidden-desktop is-hidden-tablet">Refresh</span>
            </button>
            &nbsp;
            <button onClick={(e) => setShowFilter(!showFilter)} class="is-fullwidth-mobile button is-small is-primary" type="button">
              <FontAwesomeIcon className="mdi" icon={faFilter} />&nbsp;Filter
            </button>
            &nbsp;
            <Link to={`/nutrition-plans/add`} className="is-fullwidth-mobile button is-small is-success" type="button">
              <FontAwesomeIcon className="mdi" icon={faPlus} />&nbsp;Request Plan
            </Link>
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
              <div>

                {/*
                            ##################################################################
                            EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A DESKTOP SCREEN.
                            ##################################################################
                        */}
                <div class="is-hidden-touch" >
                  <AdminNutritionPlanListDesktop
                    listData={listData}
                    setPageSize={setPageSize}
                    pageSize={pageSize}
                    previousCursors={previousCursors}
                    onPreviousClicked={onPreviousClicked}
                    onNextClicked={onNextClicked}
                    onSelectNutritionPlanForDeletion={onSelectNutritionPlanForDeletion}
                  />
                </div>

                {/*
                            ###########################################################################
                            EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A TABLET OR MOBILE SCREEN.
                            ###########################################################################
                        */}
                <div class="is-fullwidth is-hidden-desktop">
                  <AdminNutritionPlanListMobile
                    listData={listData}
                    setPageSize={setPageSize}
                    pageSize={pageSize}
                    previousCursors={previousCursors}
                    onPreviousClicked={onPreviousClicked}
                    onNextClicked={onNextClicked}
                    onSelectNutritionPlanForDeletion={onSelectNutritionPlanForDeletion}
                  />
                </div>

              </div>
            ) : (
              <section className="hero is-medium has-background-white-ter">
                <div className="hero-body">
                  <p className="title">
                    <FontAwesomeIcon className="fas" icon={faTable} />
                    &nbsp;No Nutrition Plans
                  </p>
                  <p className="subtitle">
                    You currently have no nutrition plans.{" "}
                    <b>
                      <Link to="/nutrition-plans/add">
                        Click here&nbsp;
                        <FontAwesomeIcon
                          className="mdi"
                          icon={faArrowRight}
                        />
                      </Link>
                    </b>{" "}
                    to get started requesting your first nutrition plan from our team!
                  </p>
                </div>
              </section>
            )}
          </>
        )}

        <div class="columns pt-5">
					<div class="column is-half">
						<Link class="button is-hidden-touch" to={`/dashboard`}>
							<FontAwesomeIcon className="fas" icon={faArrowLeft} />
							&nbsp;Back to Dashboard
						</Link>
						<Link
							class="button is-fullwidth is-hidden-desktop"
							to={`/dashboard`}>
							<FontAwesomeIcon className="fas" icon={faArrowLeft} />
							&nbsp;Back to Dashboard
						</Link>
					</div>
					<div class="column is-half has-text-right">
						<Link
							class="button is-success is-hidden-touch"
							to={`/nutrition-plans/add`}
>
              <FontAwesomeIcon className="fas" icon={faPlus} />&nbsp;Request Plan
            </Link>
						<Link
							class="button is-success is-fullwidth is-hidden-desktop"
							to={`/nutrition-plans/add`}
>
              <FontAwesomeIcon className="fas" icon={faPlus} />&nbsp;Request Plan
            </Link>
					</div>
				</div>

      </div>
    </Layout>
  );
}

export default AdminNutritionPlanList;
