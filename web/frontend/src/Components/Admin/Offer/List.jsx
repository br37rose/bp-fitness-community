import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faHandHolding,
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
  faFilterCircleXmark
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import FormErrorBox from "../../Reusable/FormErrorBox";
import { getOfferListAPI, deleteOfferAPI } from "../../../API/Offer";
import {
  topAlertMessageState,
  topAlertStatusState,
  currentUserState,
  offersFilterShowState,
  offersFilterTemporarySearchTextState,
  offersFilterActualSearchTextState,
  offersFilterStatusState,
} from "../../../AppState";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import FormSelectField from "../../Reusable/FormSelectField";
import FormInputFieldWithButton from "../../Reusable/FormInputFieldWithButton";
import { PAGE_SIZE_OPTIONS, OFFER_STATUS_OPTIONS } from "../../../Constants/FieldOptions";
import AdminOfferListDesktop from "./ListDesktop";
import AdminOfferListMobile from "./ListMobile";


function AdminOfferList() {
  ////
  //// Global state.
  ////

  const [topAlertMessage, setTopAlertMessage] = useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] = useRecoilState(topAlertStatusState);
  const [currentUser] = useRecoilState(currentUserState);
  const [showFilter, setShowFilter] = useRecoilState(offersFilterShowState);  // Filtering + Searching
  const [temporarySearchText, setTemporarySearchText] = useRecoilState(offersFilterTemporarySearchTextState);  // Searching - The search field value as your writes their query.
  const [actualSearchText, setActualSearchText] = useRecoilState(offersFilterActualSearchTextState); // Searching - The actual search query value to submit to the API.
  const [status, setStatus] = useRecoilState(offersFilterStatusState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [listData, setListData] = useState("");
  const [selectedOfferForDeletion, setSelectedOfferForDeletion] = useState("");
  const [isFetching, setFetching] = useState(false);
  const [pageSize, setPageSize] = useState(10); // Pagination
  const [previousCursors, setPreviousCursors] = useState([]); // Pagination
  const [nextCursor, setNextCursor] = useState(""); // Pagination
  const [currentCursor, setCurrentCursor] = useState(""); // Pagination
  const [sortField, setSortField] = useState("created"); // Sorting

  ////
  //// API.
  ////

  function onOfferListSuccess(response) {
    console.log("onOfferListSuccess: Starting...");
    if (response.results !== null) {
      setListData(response);
      if (response.hasNextPage) {
        setNextCursor(response.nextCursor); // For pagination purposes.
      }
    } else {
      setListData([]);
      setNextCursor("");}
  }

  function onOfferListError(apiErr) {
    console.log("onOfferListError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onOfferListDone() {
    console.log("onOfferListDone: Starting...");
    setFetching(false);
  }

  function onOfferDeleteSuccess(response) {
    console.log("onOfferDeleteSuccess: Starting..."); // For debugging purposes only.

    // Update notification.
    setTopAlertStatus("success");
    setTopAlertMessage("Offer deleted");
    setTimeout(() => {
      console.log(
        "onDeleteConfirmButtonClick: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus
      );
      setTopAlertMessage("");
    }, 2000);

    // Fetch again an updated list.
    fetchList(currentCursor, pageSize, actualSearchText, status);
  }

  function onOfferDeleteError(apiErr) {
    console.log("onOfferDeleteError: Starting..."); // For debugging purposes only.
    setErrors(apiErr);

    // Update notification.
    setTopAlertStatus("danger");
    setTopAlertMessage("Failed deleting");
    setTimeout(() => {
      console.log(
        "onOfferDeleteError: topAlertMessage, topAlertStatus:",
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

  function onOfferDeleteDone() {
    console.log("onOfferDeleteDone: Starting...");
    setFetching(false);
  }

  ////
  //// Event handling.
  ////

  const fetchList = (cur, limit, keywords, status) => {
    setFetching(true);
    setErrors({});

    let params = new Map();
    params.set("page_size", limit); // Pagination
    params.set("sort_field", "created"); // Sorting
    params.set("sort_order", -1);         // Sorting - descending, meaning most recent start date to oldest start date.
    params.set("status", status);

    if (cur !== "") {
      // Pagination
      params.set("cursor", cur);
    }

    // Filtering
    if (keywords !== undefined && keywords !== null && keywords !== "") {
      // Searhcing
      params.set("search", keywords);
    }

    getOfferListAPI(
      params,
      onOfferListSuccess,
      onOfferListError,
      onOfferListDone
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

  const onSelectOfferForDeletion = (e, datum) => {
    console.log("onSelectOfferForDeletion", datum);
    setSelectedOfferForDeletion(datum);
  };

  const onDeselectOfferForDeletion = (e) => {
    console.log("onDeselectOfferForDeletion");
    setSelectedOfferForDeletion("");
  };

  const onDeleteConfirmButtonClick = (e) => {
    console.log("onDeleteConfirmButtonClick"); // For debugging purposes only.

    deleteOfferAPI(
      selectedOfferForDeletion.id,
      onOfferDeleteSuccess,
      onOfferDeleteError,
      onOfferDeleteDone
    );
    setSelectedOfferForDeletion("");
  };

  // Function resets the filter state to its default state.
  const onClearFilterClick = (e) => {
    setShowFilter(false);
    setActualSearchText("");
    setTemporarySearchText("");
    setStatus(2); // 1=Pending, 2=Active, 3=Archived
  }

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.
      fetchList(currentCursor, pageSize, actualSearchText, status);
    }

    return () => {
      mounted = false;
    };
  }, [currentCursor, pageSize, actualSearchText, status]);

  ////
  //// Component rendering.
  ////

  return (
    <>
      <div className="container">
        <section className="section">

          {/* Desktop Breadcrumbs */}
          <nav className="breadcrumb is-hidden-touch" aria-label="breadcrumbs">
            <ul>
              <li className="">
                <Link to="/admin/dashboard" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faGauge} />
                  &nbsp;Dashboard
                </Link>
              </li>
              <li className="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faHandHolding} />
                  &nbsp;Offers
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class=""><Link to="/admin/dashboard" aria-current="page"><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to Dashboard</Link></li>
            </ul>
          </nav>

          {/* Page */}
          <nav className="box">
            <div
              className={`modal ${
                selectedOfferForDeletion ? "is-active" : ""
              }`}
            >
              <div className="modal-background"></div>
              <div className="modal-card">
                <header className="modal-card-head">
                  <p className="modal-card-title">Are you sure?</p>
                  <button
                    className="delete"
                    aria-label="close"
                    onClick={onDeselectOfferForDeletion}
                  ></button>
                </header>
                <section className="modal-card-body">
                  You are about to <b>archive</b> this Video Category; it will no longer
                  appear on your dashboard nor will the Video Category be able to log
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
                    onClick={onDeselectOfferForDeletion}
                  >
                    Cancel
                  </button>
                </footer>
              </div>
            </div>

            <div className="columns">
              <div className="column">
                <h1 className="title is-4">
                  <FontAwesomeIcon className="fas" icon={faHandHolding} />
                  &nbsp;Offers
                </h1>
              </div>
              <div className="column has-text-right">
                  <button onClick={()=>fetchList(currentCursor, pageSize, actualSearchText)} class="is-fullwidth-mobile button is-link is-small" type="button">
                      <FontAwesomeIcon className="mdi" icon={faRefresh} />&nbsp;<span class="is-hidden-desktop is-hidden-tablet">Refresh</span>
                  </button>
                  &nbsp;
                  <button onClick={(e)=>setShowFilter(!showFilter)} class="is-fullwidth-mobile button is-small is-primary" type="button">
                      <FontAwesomeIcon className="mdi" icon={faFilter} />&nbsp;Filter
                  </button>
                  &nbsp;
                  <Link to={`/admin/offers/add`} className="is-fullwidth-mobile button is-small is-success" type="button">
                      <FontAwesomeIcon className="mdi" icon={faPlus} />&nbsp;New
                  </Link>
              </div>
            </div>

            {/* FILTER */}
            {showFilter && (
              <div class="has-background-white-bis" style={{ borderRadius: "15px", padding: "20px" }}>
                {/* Filter Title + Clear Button */}
                <div class="columns">
                    <div class="column is-half">
                        <strong><u><FontAwesomeIcon className="mdi" icon={faFilter} />&nbsp;Filter</u></strong>
                    </div>
                    <div class="column is-half has-text-right">
                        <Link onClick={onClearFilterClick}><FontAwesomeIcon className="mdi" icon={faFilterCircleXmark} />&nbsp;Clear Filter</Link>
                    </div>
                </div>
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
                        type="number"
                        placeholder="#"
                        selectedValue={status}
                        errorText={errors && errors.status}
                        helpText={
                          <ul class="content">
                            <li>pending - will not show up for members</li>
                            <li>active - will show up for everyone</li>
                            <li>archived - will be hidden from everyone</li>
                          </ul>
                        }
                        onChange={(e)=>setStatus(parseInt(e.target.value))}
                        isRequired={true}
                        options={OFFER_STATUS_OPTIONS}
                        maxWidth="80px"
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
                        <AdminOfferListDesktop
                            listData={listData}
                            setPageSize={setPageSize}
                            pageSize={pageSize}
                            previousCursors={previousCursors}
                            onPreviousClicked={onPreviousClicked}
                            onNextClicked={onNextClicked}
                            onSelectOfferForDeletion={onSelectOfferForDeletion}
                        />
                    </div>

                    {/*
                        ###########################################################################
                        EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A TABLET OR MOBILE SCREEN.
                        ###########################################################################
                    */}
                    <div class="is-fullwidth is-hidden-desktop">
                        <AdminOfferListMobile
                            listData={listData}
                            setPageSize={setPageSize}
                            pageSize={pageSize}
                            previousCursors={previousCursors}
                            onPreviousClicked={onPreviousClicked}
                            onNextClicked={onNextClicked}
                            onSelectOfferForDeletion={onSelectOfferForDeletion}
                        />
                    </div>

                  </div>
                ) : (
                  <section className="hero is-medium has-background-white-ter">
                    <div className="hero-body">
                      <p className="title">
                        <FontAwesomeIcon className="fas" icon={faTable} />
                        &nbsp;No Offers
                      </p>
                      <p className="subtitle">
                        No offers.{" "}
                        <b>
                          <Link to="/admin/offers/add">
                            Click here&nbsp;
                            <FontAwesomeIcon
                              className="mdi"
                              icon={faArrowRight}
                            />
                          </Link>
                        </b>{" "}
                        to get started creating your first offer.
                      </p>
                    </div>
                  </section>
                )}
              </>
            )}

            <div class="columns pt-5">
                <div class="column is-half">
                    <Link class="button is-fullwidth-mobile" to={`/admin/dashboard`}><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to Dashboard</Link>
                </div>
                <div class="column is-half has-text-right">
                    <Link to={`/admin/offers/add`} class="button is-success is-fullwidth-mobile"><FontAwesomeIcon className="fas" icon={faPlus} />&nbsp;New</Link>
                </div>
            </div>

          </nav>
        </section>
      </div>
    </>
  );
}

export default AdminOfferList;
