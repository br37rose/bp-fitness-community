import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
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
  faFilterCircleXmark,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import FormErrorBox from "../../Reusable/FormErrorBox";
import { getMemberListAPI } from "../../../API/member";
import {
  topAlertMessageState,
  topAlertStatusState,
  currentUserState,
  membersFilterShowState,
  membersFilterTemporarySearchTextState,
  membersFilterActualSearchTextState,
  membersFilterOfferIDState,
  membersFilterStatusState,
  membersFilterSortState
} from "../../../AppState";
import FormSelectFieldForOffer from "../../Reusable/FormSelectFieldForOffer";
import FormRadioField from "../../Reusable/FormRadioField";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import FormInputFieldWithButton from "../../Reusable/FormInputFieldWithButton";
import { PAGE_SIZE_OPTIONS } from "../../../Constants/FieldOptions";
import AdminMemberListDesktop from "./ListDesktop";
import AdminMemberListMobile from "./ListMobile";


function AdminMemberList() {
    ////
    //// Global state.
    ////

    const [topAlertMessage, setTopAlertMessage] = useRecoilState(topAlertMessageState);
    const [topAlertStatus, setTopAlertStatus] = useRecoilState(topAlertStatusState);
    const [currentUser] = useRecoilState(currentUserState);
    const [showFilter, setShowFilter] = useRecoilState(membersFilterShowState);                                   // Filtering + Searching
    const [actualSearchText, setActualSearchText] = useRecoilState(membersFilterActualSearchTextState);           // Searching - The actual search query value to submit to the API.
    const [temporarySearchText, setTemporarySearchText] = useRecoilState(membersFilterTemporarySearchTextState);  // Searching - The search field value as your writes their query.
    const [offerID, setOfferID] = useRecoilState(membersFilterOfferIDState);                                      // Filtering
    const [status, setStatus] = useRecoilState(membersFilterStatusState);                                         // Filtering
    const [sort, setSort] = useRecoilState(membersFilterSortState);                                               // Sorting

    ////
    //// Component states.
    ////

    const [isOfferOther, setIsOfferOther] = useState("");
    const [errors, setErrors] = useState({});
    const [listData, setListData] = useState("");
    const [selectedMemberForDeletion, setSelectedMemberForDeletion] = useState("");
    const [isFetching, setFetching] = useState(false);
    const [pageSize, setPageSize] = useState(10); // Pagination
    const [previousCursors, setPreviousCursors] = useState([]); // Pagination
    const [nextCursor, setNextCursor] = useState(""); // Pagination
    const [currentCursor, setCurrentCursor] = useState(""); // Pagination

    ////
    //// API.
    ////

    function onMemberListSuccess(response) {
        console.log("onMemberListSuccess: Starting...");
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

    function onMemberListError(apiErr) {
        console.log("onMemberListError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onMemberListDone() {
        console.log("onMemberListDone: Starting...");
        setFetching(false);
    }

    function onMemberDeleteSuccess(response) {
        console.log("onMemberDeleteSuccess: Starting..."); // For debugging purposes only.

        // Update notification.
        setTopAlertStatus("success");
        setTopAlertMessage("Member deleted");
        setTimeout(() => {
          console.log(
            "onDeleteConfirmButtonClick: topAlertMessage, topAlertStatus:",
            topAlertMessage,
            topAlertStatus
          );
          setTopAlertMessage("");
        }, 2000);

        // Fetch again an updated list.
        fetchList(currentCursor, pageSize, actualSearchText, offerID, status, sort);
    }

    function onMemberDeleteError(apiErr) {
    console.log("onMemberDeleteError: Starting..."); // For debugging purposes only.
        setErrors(apiErr);

        // Update notification.
        setTopAlertStatus("danger");
        setTopAlertMessage("Failed deleting");
        setTimeout(() => {
          console.log(
            "onMemberDeleteError: topAlertMessage, topAlertStatus:",
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

    function onMemberDeleteDone() {
        console.log("onMemberDeleteDone: Starting...");
        setFetching(false);
    }

    ////
    //// Event handling.
    ////

    // Function resets the filter state to its default state.
    const onClearFilterClick = (e) => {
        setShowFilter(false);
        setActualSearchText("");
        setTemporarySearchText("");
        setOfferID(null);
        setStatus(0);
        setSort("name,-1");
    }

    const fetchList = (cur, limit, keywords, o, s, sbv) => {
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

        if (s !== undefined && s !== null && s !== "" && s !== 0) {
          params.set("status", status);
        }

        // Filtering
        if (keywords !== undefined && keywords !== null && keywords !== "") {
            // Searhcing
            params.set("search", keywords);
        }
        if (o !== undefined && o !== null && o !== "") {
            params.set("subscription_offer_id", o);
        }

        getMemberListAPI(
            params,
            onMemberListSuccess,
            onMemberListError,
            onMemberListDone
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
          fetchList(currentCursor, pageSize, actualSearchText, offerID, status, sort);
        }

        return () => {
          mounted = false;
        };
    }, [currentCursor, pageSize, actualSearchText, offerID, status, sort]);

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
                          <FontAwesomeIcon className="fas" icon={faUsers} />
                          &nbsp;Members
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
                    <div className="columns">
                      <div className="column">
                        <h1 className="title is-4">
                          <FontAwesomeIcon className="fas" icon={faUsers} />
                          &nbsp;Members
                        </h1>
                      </div>
                      <div className="column has-text-right">
                          <button onClick={()=>fetchList(currentCursor, pageSize, actualSearchText, offerID, status, sort)} class="is-fullwidth-mobile button is-link is-small" type="button">
                              <FontAwesomeIcon className="mdi" icon={faRefresh} />&nbsp;<span class="is-hidden-desktop is-hidden-tablet">Refresh</span>
                          </button>
                          &nbsp;
                          <button onClick={(e)=>setShowFilter(!showFilter)} class="is-fullwidth-mobile button is-small is-primary" type="button">
                              <FontAwesomeIcon className="mdi" icon={faFilter} />&nbsp;Filter
                          </button>
                          &nbsp;
                          <Link to={`/admin/members/add`} className="is-fullwidth-mobile button is-small is-success" type="button">
                              <FontAwesomeIcon className="mdi" icon={faPlus} />&nbsp;New Member
                          </Link>
                      </div>
                    </div>

                    {/* FILTER */}
                    {showFilter && (
                        <div class="has-background-white-bis" style={{borderRadius:"15px", padding:"20px"}}>

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
                                <div class="column">
                                    <FormSelectFieldForOffer
                                        label={`Enrollment`}
                                        isSubscription={true}
                                        offerID={offerID}
                                        setOfferID={setOfferID}
                                        isOfferOther={isOfferOther}
                                        setIsOfferOther={setIsOfferOther}
                                        errorText={errors && errors.offerId}
                                    />
                                </div>
                                <div class="column">
                                    <FormRadioField
                                        label="Status"
                                        name="status"
                                        placeholder="Pick"
                                        value={status}
                                        opt2Value={1}
                                        opt2Label="Active"
                                        opt4Value={2}
                                        opt4Label="Archived"
                                        errorText={errors && errors.status}
                                        onChange={(e) => setStatus(parseInt(e.target.value))}
                                        maxWidth="180px"
                                        disabled={false}
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
                                <AdminMemberListDesktop
                                    listData={listData}
                                    setPageSize={setPageSize}
                                    pageSize={pageSize}
                                    previousCursors={previousCursors}
                                    onPreviousClicked={onPreviousClicked}
                                    onNextClicked={onNextClicked}
                                />
                            </div>

                            {/*
                                ###########################################################################
                                EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A TABLET OR MOBILE SCREEN.
                                ###########################################################################
                            */}
                            <div class="is-fullwidth is-hidden-desktop">
                                <AdminMemberListMobile
                                    listData={listData}
                                    setPageSize={setPageSize}
                                    pageSize={pageSize}
                                    previousCursors={previousCursors}
                                    onPreviousClicked={onPreviousClicked}
                                    onNextClicked={onNextClicked}
                                />
                            </div>

                          </div>
                        ) : (
                          <section className="hero is-medium has-background-white-ter">
                            <div className="hero-body">
                              <p className="title">
                                <FontAwesomeIcon className="fas" icon={faTable} />
                                &nbsp;No Members
                              </p>
                              <p className="subtitle">
                                No class types.{" "}
                                <b>
                                  <Link to="/admin/members/add">
                                    Click here&nbsp;
                                    <FontAwesomeIcon
                                      className="mdi"
                                      icon={faArrowRight}
                                    />
                                  </Link>
                                </b>{" "}
                                to get started creating your first member location type.
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
                            <Link to={`/admin/members/add`} class="button is-success is-fullwidth-mobile"><FontAwesomeIcon className="fas" icon={faPlus} />&nbsp;New Member</Link>
                        </div>
                    </div>

                  </nav>
                </section>
            </div>
        </>
    );
}

export default AdminMemberList;
