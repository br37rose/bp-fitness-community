import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faArrowLeft,
  faDumbbell,
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
import { getExerciseListAPI, deleteExerciseAPI } from "../../../API/Exercise";
import {
  topAlertMessageState,
  topAlertStatusState,
  currentUserState,
  exercisesFilterShowState,
  exercisesFilterTemporarySearchTextState,
  exercisesFilterActualSearchTextState,
  exercisesFilterCategoryState,
  exercisesFilterMovementTypeState,
  exercisesFilterStatusState,
  exercisesFilterGenderState,
  exercisesFilterVideoTypeState,
  exercisesFilterSortState
} from "../../../AppState";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import FormInputFieldWithButton from "../../Reusable/FormInputFieldWithButton";
import FormSelectField from "../../Reusable/FormSelectField";
import FormMultiSelectFieldForTags from "../../Reusable/FormMultiSelectFieldForTags";
import {
    PAGE_SIZE_OPTIONS,
    EXERCISE_CATEGORY_OPTIONS_WITH_EMPTY_OPTION,
    EXERCISE_MOMENT_TYPE_OPTIONS_WITH_EMPTY_OPTION,
    EXERCISE_STATUS_OPTIONS_WITH_EMPTY_OPTION,
    EXERCISE_GENDER_OPTIONS_WITH_EMPTY_OPTION,
    EXERCISE_VIDEO_FILE_TYPE_OPTIONS_WITH_EMPTY_OPTION
} from "../../../Constants/FieldOptions";
import AdminExerciseListDesktop from "./ListDesktop";
import AdminExerciseListMobile from "./ListMobile";


function AdminExerciseList() {
    ////
    //// Global state.
    ////

    const [topAlertMessage, setTopAlertMessage] = useRecoilState(topAlertMessageState);
    const [topAlertStatus, setTopAlertStatus] = useRecoilState(topAlertStatusState);
    const [currentUser] = useRecoilState(currentUserState);
    const [showFilter, setShowFilter] = useRecoilState(exercisesFilterShowState);                                  // Filtering + Searching
    const [sort, setSort] = useRecoilState(exercisesFilterSortState);                                              // Sorting
    const [temporarySearchText, setTemporarySearchText] = useRecoilState(exercisesFilterTemporarySearchTextState); // Searching - The search field value as your writes their query.
    const [actualSearchText, setActualSearchText] = useRecoilState(exercisesFilterActualSearchTextState);          // Searching - The actual search query value to submit to the API.
    const [category, setCategory] = useRecoilState(exercisesFilterCategoryState);
    const [movementType, setMovementType] = useRecoilState(exercisesFilterMovementTypeState);
    const [status, setStatus] = useRecoilState(exercisesFilterStatusState);
    const [gender, setGender] = useRecoilState(exercisesFilterGenderState);
    const [videoType, setVideoType] = useRecoilState(exercisesFilterVideoTypeState);

    ////
    //// Component states.
    ////

    const [errors, setErrors] = useState({});
    const [listData, setListData] = useState("");
    const [selectedExerciseForDeletion, setSelectedExerciseForDeletion] = useState("");
    const [isFetching, setFetching] = useState(false);
    const [pageSize, setPageSize] = useState(10); // Pagination
    const [previousCursors, setPreviousCursors] = useState([]); // Pagination
    const [nextCursor, setNextCursor] = useState(""); // Pagination
    const [currentCursor, setCurrentCursor] = useState(""); // Pagination
    const [tags, setTags] = useState([]);

    ////
    //// API.
    ////

    function onExerciseListSuccess(response) {
        console.log("onExerciseListSuccess: Starting...");
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

    function onExerciseListError(apiErr) {
        console.log("onExerciseListError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onExerciseListDone() {
        console.log("onExerciseListDone: Starting...");
        setFetching(false);
    }

    function onExerciseDeleteSuccess(response) {
        console.log("onExerciseDeleteSuccess: Starting..."); // For debugging purposes only.

        // Update notification.
        setTopAlertStatus("success");
        setTopAlertMessage("Exercise deleted");
        setTimeout(() => {
          console.log(
            "onDeleteConfirmButtonClick: topAlertMessage, topAlertStatus:",
            topAlertMessage,
            topAlertStatus
          );
          setTopAlertMessage("");
        }, 2000);

        // Fetch again an updated list.
        fetchList(currentCursor, pageSize, actualSearchText, category, movementType, status, gender, videoType, sort, tags);
    }

    function onExerciseDeleteError(apiErr) {
        console.log("onExerciseDeleteError: Starting..."); // For debugging purposes only.
        setErrors(apiErr);

        // Update notification.
        setTopAlertStatus("danger");
        setTopAlertMessage("Failed deleting");
        setTimeout(() => {
          console.log(
            "onExerciseDeleteError: topAlertMessage, topAlertStatus:",
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

    function onExerciseDeleteDone() {
        console.log("onExerciseDeleteDone: Starting...");
        setFetching(false);
    }

    ////
    //// Event handling.
    ////

    const fetchList = (cur, limit, keywords, cat, mt, st, g, vt, sbv, t) => {
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
        if (t !== undefined && t !== null && t !== "") {
            params.set("tags", t);
        }

        getExerciseListAPI(
          params,
          onExerciseListSuccess,
          onExerciseListError,
          onExerciseListDone
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

    const onSelectExerciseForDeletion = (e, datum) => {
        console.log("onSelectExerciseForDeletion", datum);
        setSelectedExerciseForDeletion(datum);
    };

    const onDeselectExerciseForDeletion = (e) => {
        console.log("onDeselectExerciseForDeletion");
        setSelectedExerciseForDeletion("");
    };

    const onDeleteConfirmButtonClick = (e) => {
        console.log("onDeleteConfirmButtonClick"); // For debugging purposes only.

        deleteExerciseAPI(
          selectedExerciseForDeletion.id,
          onExerciseDeleteSuccess,
          onExerciseDeleteError,
          onExerciseDeleteDone
        );
        setSelectedExerciseForDeletion("");
    };

    // Function resets the filter state to its default state.
    const onClearFilterClick = (e) => {
        setShowFilter(false);
        setActualSearchText("");
        setTemporarySearchText("");
        setStatus(0);
    }

    ////
    //// Misc.
    ////

    useEffect(() => {
        let mounted = true;

        if (mounted) {
          window.scrollTo(0, 0); // Start the page at the top of the page.
          fetchList(currentCursor, pageSize, actualSearchText, category, movementType, status, gender, videoType, sort, tags);
        }

        return () => {
          mounted = false;
        };
    }, [currentCursor, pageSize, actualSearchText, category, movementType, status, gender, videoType, sort, tags]);

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
                            <FontAwesomeIcon className="fas" icon={faDumbbell} />
                            &nbsp;Exercises
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
                    selectedExerciseForDeletion ? "is-active" : ""
                  }`}
                >
                  <div className="modal-background"></div>
                  <div className="modal-card">
                    <header className="modal-card-head">
                      <p className="modal-card-title">Are you sure?</p>
                      <button
                        className="delete"
                        aria-label="close"
                        onClick={onDeselectExerciseForDeletion}
                      ></button>
                    </header>
                    <section className="modal-card-body">
                      You are about to <b>archive</b> this exercise; it will no longer
                      appear on your dashboard nor will the exercise be able to log
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
                        onClick={onDeselectExerciseForDeletion}
                      >
                        Cancel
                      </button>
                    </footer>
                  </div>
                </div>

                <div className="columns">
                  <div className="column">
                    <h1 className="title is-4">
                      <FontAwesomeIcon className="fas" icon={faDumbbell} />
                      &nbsp;Exercises
                    </h1>
                  </div>
                  <div className="column has-text-right">
                      <button onClick={()=>fetchList(currentCursor, pageSize, actualSearchText, category, movementType, status, gender, videoType, sort, tags)} class="is-fullwidth-mobile button is-link is-small" type="button">
                          <FontAwesomeIcon className="mdi" icon={faRefresh} />&nbsp;<span class="is-hidden-desktop is-hidden-tablet">Refresh</span>
                      </button>
                      &nbsp;
                      <button onClick={(e)=>setShowFilter(!showFilter)} class="is-fullwidth-mobile button is-small is-primary" type="button">
                          <FontAwesomeIcon className="mdi" icon={faFilter} />&nbsp;Filter
                      </button>
                      &nbsp;
                      <Link to={`/admin/exercises/add`} className="is-fullwidth-mobile button is-small is-success" type="button">
                          <FontAwesomeIcon className="mdi" icon={faPlus} />&nbsp;New
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
                            <AdminExerciseListDesktop
                                listData={listData}
                                setPageSize={setPageSize}
                                pageSize={pageSize}
                                previousCursors={previousCursors}
                                onPreviousClicked={onPreviousClicked}
                                onNextClicked={onNextClicked}
                                onSelectExerciseForDeletion={onSelectExerciseForDeletion}
                            />
                        </div>

                        {/*
                            ###########################################################################
                            EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A TABLET OR MOBILE SCREEN.
                            ###########################################################################
                        */}
                        <div class="is-fullwidth is-hidden-desktop">
                            <AdminExerciseListMobile
                                listData={listData}
                                setPageSize={setPageSize}
                                pageSize={pageSize}
                                previousCursors={previousCursors}
                                onPreviousClicked={onPreviousClicked}
                                onNextClicked={onNextClicked}
                                onSelectExerciseForDeletion={onSelectExerciseForDeletion}
                            />
                        </div>

                      </div>
                    ) : (
                      <section className="hero is-medium has-background-white-ter">
                        <div className="hero-body">
                          <p className="title">
                            <FontAwesomeIcon className="fas" icon={faTable} />
                            &nbsp;No Exercises
                          </p>
                          <p className="subtitle">
                            No class types.{" "}
                            <b>
                              <Link to="/admin/exercises/add">
                                Click here&nbsp;
                                <FontAwesomeIcon
                                  className="mdi"
                                  icon={faArrowRight}
                                />
                              </Link>
                            </b>{" "}
                            to get started creating your first exercise location type.
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
                        <Link to={`/admin/exercises/add`} class="button is-success is-fullwidth-mobile"><FontAwesomeIcon className="fas" icon={faPlus} />&nbsp;New Exercise</Link>
                    </div>
                </div>

                </nav>
                </section>
            </div>
        </>
    );
}

export default AdminExerciseList;
