import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faArrowLeft,
  faVideoCamera,
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
import { getVideoCollectionListAPI, deleteVideoCollectionAPI } from "../../../API/VideoCollection";
import {
  topAlertMessageState,
  topAlertStatusState,
  currentUserState,
} from "../../../AppState";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import FormInputFieldWithButton from "../../Reusable/FormInputFieldWithButton";
import FormSelectField from "../../Reusable/FormSelectField";
import {
    PAGE_SIZE_OPTIONS,
    EXERCISE_CATEGORY_OPTIONS_WITH_EMPTY_OPTION,
    EXERCISE_MOMENT_TYPE_OPTIONS_WITH_EMPTY_OPTION,
    EXERCISE_STATUS_OPTIONS_WITH_EMPTY_OPTION,
    EXERCISE_GENDER_OPTIONS_WITH_EMPTY_OPTION,
    EXERCISE_VIDEO_FILE_TYPE_OPTIONS_WITH_EMPTY_OPTION
} from "../../../Constants/FieldOptions";
import AdminVideoCollectionListDesktop from "./ListDesktop";
import AdminVideoCollectionListMobile from "./ListMobile";


function AdminVideoCollectionList() {
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
  const [selectedVideoCollectionForDeletion, setSelectedVideoCollectionForDeletion] = useState("");
  const [isFetching, setFetching] = useState(false);
  const [pageSize, setPageSize] = useState(10); // Pagination
  const [previousCursors, setPreviousCursors] = useState([]); // Pagination
  const [nextCursor, setNextCursor] = useState(""); // Pagination
  const [currentCursor, setCurrentCursor] = useState(""); // Pagination
  const [showFilter, setShowFilter] = useState(false); // Filtering + Searching
  const [sortField, setSortField] = useState("created"); // Sorting
  const [temporarySearchText, setTemporarySearchText] = useState(""); // Searching - The search field value as your writes their query.
  const [actualSearchText, setActualSearchText] = useState(""); // Searching - The actual search query value to submit to the API.
  const [category, setCategory] = useState("");
  const [movementType, setMovementType] = useState("");
  const [status, setStatus] = useState("");
  const [gender, setGender] = useState("");
  const [videoType, setVideoType] = useState("");

  ////
  //// API.
  ////

  function onVideoCollectionListSuccess(response) {
    console.log("onVideoCollectionListSuccess: Starting...");
    if (response.results !== null) {
      setListData(response);
      if (response.hasNextPage) {
        setNextCursor(response.nextCursor); // For pagination purposes.
      }
    } else {
      setListData([]);
      setNextCursor("");}
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

  function onVideoCollectionDeleteSuccess(response) {
    console.log("onVideoCollectionDeleteSuccess: Starting..."); // For debugging purposes only.

    // Update notification.
    setTopAlertStatus("success");
    setTopAlertMessage("VideoCollection deleted");
    setTimeout(() => {
      console.log(
        "onDeleteConfirmButtonClick: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus
      );
      setTopAlertMessage("");
    }, 2000);

    // Fetch again an updated list.
    fetchList(currentCursor, pageSize, actualSearchText, category, movementType, status, gender, videoType);
  }

  function onVideoCollectionDeleteError(apiErr) {
    console.log("onVideoCollectionDeleteError: Starting..."); // For debugging purposes only.
    setErrors(apiErr);

    // Update notification.
    setTopAlertStatus("danger");
    setTopAlertMessage("Failed deleting");
    setTimeout(() => {
      console.log(
        "onVideoCollectionDeleteError: topAlertMessage, topAlertStatus:",
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

  function onVideoCollectionDeleteDone() {
    console.log("onVideoCollectionDeleteDone: Starting...");
    setFetching(false);
  }

  ////
  //// Event handling.
  ////

  const fetchList = (cur, limit, keywords, cat, mt, st, g, vt) => {
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
    if (cat !== undefined && cat !== null && cat !== "") {
        params.set("category", cat);
    }
    if (mt !== undefined && mt !== null && mt !== "") {
        params.set("movement_type", mt);
    }
    if (st !== undefined && st !== null && st !== "") {
        params.set("status", st);
    }
    if (g !== undefined && g !== null && g !== "") {
        params.set("gender", g);
    }
    if (vt !== undefined && vt !== null && vt !== "") {
        params.set("video_type", vt);
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

  const onSelectVideoCollectionForDeletion = (e, datum) => {
    console.log("onSelectVideoCollectionForDeletion", datum);
    setSelectedVideoCollectionForDeletion(datum);
  };

  const onDeselectVideoCollectionForDeletion = (e) => {
    console.log("onDeselectVideoCollectionForDeletion");
    setSelectedVideoCollectionForDeletion("");
  };

  const onDeleteConfirmButtonClick = (e) => {
    console.log("onDeleteConfirmButtonClick"); // For debugging purposes only.

    deleteVideoCollectionAPI(
      selectedVideoCollectionForDeletion.id,
      onVideoCollectionDeleteSuccess,
      onVideoCollectionDeleteError,
      onVideoCollectionDeleteDone
    );
    setSelectedVideoCollectionForDeletion("");
  };

    ////
    //// Misc.
    ////

    useEffect(() => {
        let mounted = true;

        if (mounted) {
          window.scrollTo(0, 0); // Start the page at the top of the page.
          fetchList(currentCursor, pageSize, actualSearchText, category, movementType, status, gender, videoType);
        }

        return () => {
          mounted = false;
        };
    }, [currentCursor, pageSize, actualSearchText, category, movementType, status, gender, videoType]);

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
                            <FontAwesomeIcon className="fas" icon={faVideoCamera} />
                            &nbsp;Video Collections
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
                    selectedVideoCollectionForDeletion ? "is-active" : ""
                  }`}
                >
                  <div className="modal-background"></div>
                  <div className="modal-card">
                    <header className="modal-card-head">
                      <p className="modal-card-title">Are you sure?</p>
                      <button
                        className="delete"
                        aria-label="close"
                        onClick={onDeselectVideoCollectionForDeletion}
                      ></button>
                    </header>
                    <section className="modal-card-body">
                      You are about to <b>archive</b> this video collection; it will no longer
                      appear on your dashboard nor will the video collection be able to log
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
                        onClick={onDeselectVideoCollectionForDeletion}
                      >
                        Cancel
                      </button>
                    </footer>
                  </div>
                </div>

                <div className="columns">
                  <div className="column">
                    <h1 className="title is-4">
                      <FontAwesomeIcon className="fas" icon={faVideoCamera} />
                      &nbsp;Video Collections
                    </h1>
                  </div>
                  <div className="column has-text-right">
                      <button onClick={()=>fetchList(currentCursor, pageSize, actualSearchText, category, movementType, status, gender, videoType)} class="is-fullwidth-mobile button is-link is-small" type="button">
                          <FontAwesomeIcon className="mdi" icon={faRefresh} />&nbsp;<span class="is-hidden-desktop is-hidden-tablet">Refresh</span>
                      </button>
                      &nbsp;
                      <button onClick={(e)=>setShowFilter(!showFilter)} class="is-fullwidth-mobile button is-small is-primary" type="button">
                          <FontAwesomeIcon className="mdi" icon={faFilter} />&nbsp;Filter
                      </button>
                      &nbsp;
                      <Link to={`/admin/video-collections/add`} className="is-fullwidth-mobile button is-small is-success" type="button">
                          <FontAwesomeIcon className="mdi" icon={faPlus} />&nbsp;New
                      </Link>
                  </div>
                </div>

                {showFilter && (
                  <div class="columns has-background-white-bis" style={{ borderRadius: "15px", padding: "20px" }}>
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
                            label="Category"
                            name="category"
                            placeholder="Pick"
                            selectedValue={category}
                            errorText={errors && errors.category}
                            helpText=""
                            onChange={(e)=>setCategory(parseInt(e.target.value))}
                            options={EXERCISE_CATEGORY_OPTIONS_WITH_EMPTY_OPTION}
                        />
                    </div>
                    <div class="column">
                        <FormSelectField
                            label="Movement Type"
                            name="movementType"
                            placeholder="Pick"
                            selectedValue={movementType}
                            errorText={errors && errors.movementType}
                            helpText=""
                            onChange={(e)=>setMovementType(parseInt(e.target.value))}
                            options={EXERCISE_MOMENT_TYPE_OPTIONS_WITH_EMPTY_OPTION}
                        />
                    </div>
                    <div class="column">
                        <FormSelectField
                            label="Status"
                            name="status"
                            placeholder="Pick"
                            selectedValue={status}
                            errorText={errors && errors.status}
                            helpText=""
                            onChange={(e)=>setStatus(parseInt(e.target.value))}
                            options={EXERCISE_STATUS_OPTIONS_WITH_EMPTY_OPTION}
                        />
                    </div>
                    <div class="column">
                        <FormSelectField
                            label="Gender"
                            name="gender"
                            placeholder="Pick"
                            selectedValue={gender}
                            errorText={errors && errors.gender}
                            helpText=""
                            onChange={(e)=>setGender(e.target.value)}
                            options={EXERCISE_GENDER_OPTIONS_WITH_EMPTY_OPTION}
                        />
                    </div>
                    <div class="column">
                        <FormSelectField
                            label="Video Type"
                            name="videoType"
                            placeholder="Pick"
                            selectedValue={videoType}
                            errorText={errors && errors.videoType}
                            helpText=""
                            onChange={(e)=>setVideoType(e.target.value)}
                            options={EXERCISE_VIDEO_FILE_TYPE_OPTIONS_WITH_EMPTY_OPTION}
                        />
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
                            <AdminVideoCollectionListDesktop
                                listData={listData}
                                setPageSize={setPageSize}
                                pageSize={pageSize}
                                previousCursors={previousCursors}
                                onPreviousClicked={onPreviousClicked}
                                onNextClicked={onNextClicked}
                                onSelectVideoCollectionForDeletion={onSelectVideoCollectionForDeletion}
                            />
                        </div>

                        {/*
                            ###########################################################################
                            EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A TABLET OR MOBILE SCREEN.
                            ###########################################################################
                        */}
                        <div class="is-fullwidth is-hidden-desktop">
                            <AdminVideoCollectionListMobile
                                listData={listData}
                                setPageSize={setPageSize}
                                pageSize={pageSize}
                                previousCursors={previousCursors}
                                onPreviousClicked={onPreviousClicked}
                                onNextClicked={onNextClicked}
                                onSelectVideoCollectionForDeletion={onSelectVideoCollectionForDeletion}
                            />
                        </div>

                      </div>
                    ) : (
                      <section className="hero is-medium has-background-white-ter">
                        <div className="hero-body">
                          <p className="title">
                            <FontAwesomeIcon className="fas" icon={faTable} />
                            &nbsp;No Video Collections
                          </p>
                          <p className="subtitle">
                            No class types.{" "}
                            <b>
                              <Link to="/admin/video-collections/add">
                                Click here&nbsp;
                                <FontAwesomeIcon
                                  className="mdi"
                                  icon={faArrowRight}
                                />
                              </Link>
                            </b>{" "}
                            to get started creating your first video collections.
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
                        <Link to={`/admin/video-collections/add`} class="button is-success is-fullwidth-mobile"><FontAwesomeIcon className="fas" icon={faPlus} />&nbsp;New</Link>
                    </div>
                </div>

                </nav>
                </section>
            </div>
        </>
    );
}

export default AdminVideoCollectionList;
