import { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faArrowLeft,
  faPlus,
  faGauge,
  faArrowRight,
  faRefresh,
  faFilter,
  faFilterCircleXmark,
  faQuestionCircle,
  faList,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import FormErrorBox from "../../Reusable/FormErrorBox";
import {
  topAlertMessageState,
  topAlertStatusState,
  currentUserState,
  questionnaireFilterStatus,
  questionnaireFilterShowState,
  questionnaireFilterSortState,
} from "../../../AppState";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import FormSelectField from "../../Reusable/FormSelectField";
import { QUESTIONNAIRE_STATUS_OPTIONS } from "../../../Constants/FieldOptions";
import {
  deleteQuestionnaireAPI,
  getQuestionnaireListApi,
} from "../../../API/questionnaire";
import AdminQuestionnaireListDesktop from "./ListDesktop";
import AdminQuestionnaireListMobile from "./ListMobile";

function AdminQuestionnaireList() {
  ////
  //// Global state.
  ////

  const [topAlertMessage, setTopAlertMessage] =
    useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] =
    useRecoilState(topAlertStatusState);
  const [currentUser] = useRecoilState(currentUserState);
  const [showFilter, setShowFilter] = useRecoilState(
    questionnaireFilterShowState
  ); // Filtering + Searching
  const [status, setStatus] = useRecoilState(questionnaireFilterStatus);
  const [sort, setSort] = useRecoilState(questionnaireFilterSortState);
  const [selectedQuestionForDeletion, setSelectedQuestionForDeletion] =
    useState("");

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [listData, setListData] = useState("");

  const [isFetching, setFetching] = useState(false);
  const [pageSize, setPageSize] = useState(10); // Pagination
  const [previousCursors, setPreviousCursors] = useState([]); // Pagination
  const [nextCursor, setNextCursor] = useState(""); // Pagination
  const [currentCursor, setCurrentCursor] = useState(""); // Pagination
  const [showModal, setShowModal] = useState(false);

  ////
  //// API.
  ////

  function OnListSuccess(response) {
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

  function OnListErr(apiErr) {
    setErrors(apiErr);
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onListDone() {
    setFetching(false);
  }

  ////
  //// Event handling.
  ////

  const fetchList = (cur, limit, keywords, st, sbv) => {
    setFetching(true);
    setErrors({});

    let params = new Map();
    params.set("page_size", limit); // Pagination

    // DEVELOPERS NOTE: Our `sortByValue` is string with the sort field
    // and sort order combined with a comma seperation. Therefore we
    // need to split as follows.
    if (sbv !== undefined && sbv !== null && sbv !== "") {
      const sortArray = sbv.split(",");
      params.set("sort_field", sortArray[0]);
      params.set("sort_order", sortArray[1]);
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
    if (st !== undefined && st !== null && st !== "") {
      params.set("status", status === 1 ? true : false);
    }

    getQuestionnaireListApi(params, OnListSuccess, OnListErr, onListDone);
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

  // Function resets the filter state to its default state.
  const onClearFilterClick = (e) => {
    setShowFilter(false);
    setStatus(0);
    setSort("created,-1");
  };

  ////
  //// Misc.
  ////

  function onQuestionDeleteError(apiErr) {
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

  function onQuestionDeleteSuccess(response) {
    // Update notification.
    setTopAlertStatus("success");
    setTopAlertMessage("Question deleted");
    setTimeout(() => {
      setTopAlertMessage("");
    }, 2000);

    // Fetch again an updated list.
    fetchList(currentCursor, pageSize, status, sort);
  }
  function onQuestionDeleteDone() {
    setFetching(false);
    setShowModal(false);
  }

  const onSelectQuestionForDeletion = (e, datum) => {
    setSelectedQuestionForDeletion(datum);
    setShowModal(true);
  };
  const onDeselectQuestionForDeletion = (e) => {
    setSelectedQuestionForDeletion("");
    setShowModal(false);
  };

  const onDeleteConfirmButtonClick = (e) => {
    deleteQuestionnaireAPI(
      selectedQuestionForDeletion.id,
      onQuestionDeleteSuccess,
      onQuestionDeleteError,
      onQuestionDeleteDone
    );
    setSelectedQuestionForDeletion("");
  };

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.
      fetchList(currentCursor, pageSize, status, sort);
    }

    return () => {
      mounted = false;
    };
  }, [currentCursor, pageSize, status, sort]);

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
                  <FontAwesomeIcon className="fas" icon={faList} />
                  &nbsp;Questions
                </Link>
              </li>
            </ul>
          </nav>

          {showModal && (
            <nav>
              <div
                className="modal is-active"
                // className={`modal ${
                //   selectedQuestionForDeletion ? "is-active" : ""
                // }`}
              >
                <div className="modal-background"></div>
                <div className="modal-card">
                  <header className="modal-card-head">
                    <p className="modal-card-title">Are you sure?</p>
                    <button
                      className="delete"
                      aria-label="close"
                      onClick={onDeselectQuestionForDeletion}
                    ></button>
                  </header>
                  <section className="modal-card-body">
                    You are about to <b>Delete</b> this Question; it will no
                    longer appear on your dashboard. This action cannot be
                    undone. Are you sure you would like to continue?
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
                      onClick={onDeselectQuestionForDeletion}
                    >
                      Cancel
                    </button>
                  </footer>
                </div>
              </div>
            </nav>
          )}

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to="/admin/dashboard" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Dashboard
                </Link>
              </li>
            </ul>
          </nav>

          {/* Page */}
          <nav className="box">
            <div className="columns">
              <div className="column">
                <h1 className="title is-4">
                  <FontAwesomeIcon className="fas" icon={faList} />
                  &nbsp;Questions
                </h1>
              </div>
              <div className="column has-text-right">
                <button
                  onClick={() =>
                    fetchList(currentCursor, pageSize, status, sort)
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
                  to={`/admin/questions/add`}
                  className="is-fullwidth-mobile button is-small is-success"
                  type="button"
                >
                  <FontAwesomeIcon className="mdi" icon={faPlus} />
                  &nbsp;New
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
                      <FontAwesomeIcon
                        className="mdi"
                        icon={faFilterCircleXmark}
                      />
                      &nbsp;Clear Filter
                    </Link>
                  </div>
                </div>

                {/* Filter Options */}
                <div class="columns">
                  <div class="column">
                    <FormSelectField
                      label="Status"
                      name="status"
                      placeholder="Pick"
                      selectedValue={status}
                      errorText={errors && errors.status}
                      helpText=""
                      onChange={(e) => setStatus(parseInt(e.target.value))}
                      options={QUESTIONNAIRE_STATUS_OPTIONS}
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
                      <AdminQuestionnaireListDesktop
                        listData={listData}
                        setPageSize={setPageSize}
                        pageSize={pageSize}
                        previousCursors={previousCursors}
                        onPreviousClicked={onPreviousClicked}
                        onNextClicked={onNextClicked}
                        onSelectQuestionForDeletion={
                          onSelectQuestionForDeletion
                        }
                      />
                    </div>

                    {/*
                            ###########################################################################
                            EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A TABLET OR MOBILE SCREEN.
                            ###########################################################################
                        */}
                    <div class="is-fullwidth is-hidden-desktop">
                      <AdminQuestionnaireListMobile
                        listData={listData}
                        setPageSize={setPageSize}
                        pageSize={pageSize}
                        previousCursors={previousCursors}
                        onPreviousClicked={onPreviousClicked}
                        onNextClicked={onNextClicked}
                        onSelectQuestionForDeletion={
                          onSelectQuestionForDeletion
                        }
                      />
                    </div>
                  </div>
                ) : (
                  <section className="hero is-medium has-background-white-ter">
                    <div className="hero-body">
                      <p className="title">
                        <FontAwesomeIcon
                          className="fas"
                          icon={faQuestionCircle}
                        />
                        &nbsp;No Questions Available
                      </p>
                      <p className="subtitle">
                        There are currently no questions available.&nbsp;
                        <b>
                          <Link to="/admin/questions/add">
                            Click here&nbsp;
                            <FontAwesomeIcon
                              className="mdi"
                              icon={faArrowRight}
                            />
                          </Link>
                        </b>{" "}
                        to get started creating your questions.
                      </p>
                    </div>
                  </section>
                )}
              </>
            )}

            <div class="columns pt-5">
              <div class="column is-half">
                <Link
                  class="button is-fullwidth-mobile"
                  to={`/admin/dashboard`}
                >
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Dashboard
                </Link>
              </div>
              <div class="column is-half has-text-right">
                <Link
                  to={`/admin/questions/add`}
                  class="button is-success is-fullwidth-mobile"
                >
                  <FontAwesomeIcon className="fas" icon={faPlus} />
                  &nbsp;New
                </Link>
              </div>
            </div>
          </nav>
        </section>
      </div>
    </>
  );
}

export default AdminQuestionnaireList;
