import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faArrowLeft,
  faGauge,
  faPencil,
  faDumbbell,
  faEye,
  faTrash,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import { useParams } from "react-router-dom";

import FormErrorBox from "../../Reusable/FormErrorBox";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import { topAlertMessageState, topAlertStatusState } from "../../../AppState";
import DataDisplayRowText from "../../Reusable/DataDisplayRowText";
import {
  deleteQuestionnaireAPI,
  getQuestionnaireDetailAPI,
} from "../../../API/questionnaire";

function AdminQuestionnaireDetail() {
  ////
  //// URL Parameters.
  ////

  const { id } = useParams();

  const [topAlertMessage, setTopAlertMessage] =
    useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] =
    useRecoilState(topAlertStatusState);

  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [datum, setDatum] = useState({});
  const [selectedQuestionForDeletion, setSelectedQuestionForDeletion] =
    useState(null);

  const onDeleteConfirmButtonClick = () => {
    deleteQuestionnaireAPI(
      selectedQuestionForDeletion.id,
      OnQuestionDeleteSuccess,
      OnQuestionDeleteError,
      OnQuestionDeleteDone
    );
    setSelectedQuestionForDeletion(null);
  };

  // --- Detail --- //

  function OnQuestionSuccess(response) {
    setDatum(response);
  }

  function OnQuestionError(apiErr) {
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function OnQuestionDone() {
    setFetching(false);
  }

  // --- Delete --- //

  function OnQuestionDeleteSuccess(response) {
    // Update notification.
    setTopAlertStatus("success");
    setTopAlertMessage("Question deleted successfully");
    setTimeout(() => {
      setTopAlertMessage("");
    }, 2000);

    // Redirect back to the video categories page.
    setForceURL("/admin/questions");
  }

  function OnQuestionDeleteError(apiErr) {
    setErrors(apiErr);

    // Update notification.
    setTopAlertStatus("danger");
    setTopAlertMessage("Failed deleting");
    setTimeout(() => {
      console.log(topAlertMessage, topAlertStatus);
      setTopAlertMessage("");
    }, 2000);

    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function OnQuestionDeleteDone() {
    setFetching(false);
  }

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.

      setFetching(true);
      getQuestionnaireDetailAPI(
        id,
        OnQuestionSuccess,
        OnQuestionError,
        OnQuestionDone
      );
    }

    return () => {
      mounted = false;
    };
  }, []);
  ////
  //// Component rendering.
  ////

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  return (
    <>
      <div class="container">
        <section class="section">
          {/* Desktop Breadcrumbs */}
          <nav class="breadcrumb is-hidden-touch" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to="/admin/dashboard" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faGauge} />
                  &nbsp;Dashboard
                </Link>
              </li>
              <li class="">
                <Link to="/admin/questions" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faDumbbell} />
                  &nbsp;Questions
                </Link>
              </li>
              <li class="is-active">
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
                <Link to="/admin/video-categories" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Questions
                </Link>
              </li>
            </ul>
          </nav>

          {/* Modal */}
          <nav>
            {/* Delete modal */}
            <div
              class={`modal ${
                selectedQuestionForDeletion !== null ? "is-active" : ""
              }`}
            >
              <div class="modal-background"></div>
              <div class="modal-card">
                <header class="modal-card-head">
                  <p class="modal-card-title">Are you sure?</p>
                  <button
                    class="delete"
                    aria-label="close"
                    onClick={(e, ses) => setSelectedQuestionForDeletion(null)}
                  ></button>
                </header>
                <section class="modal-card-body">
                  You are about to delete this questionnaire and all associated
                  data. This action cannot be undone. Are you sure you want to
                  continue?
                </section>
                <footer class="modal-card-foot">
                  <button
                    class="button is-success"
                    onClick={onDeleteConfirmButtonClick}
                  >
                    Confirm
                  </button>
                  <button
                    class="button"
                    onClick={(e, ses) => setSelectedQuestionForDeletion(null)}
                  >
                    Cancel
                  </button>
                </footer>
              </div>
            </div>
          </nav>

          {/* Page */}
          <nav class="box">
            {datum && (
              <div class="columns">
                <div class="column">
                  <p class="title is-4">
                    <FontAwesomeIcon className="fas" icon={faDumbbell} />
                    &nbsp;Question
                  </p>
                </div>
                <div class="column has-text-right">
                  <Link
                    to={`/admin/questions/${id}/update`}
                    class="button is-warning is-small is-fullwidth-mobile"
                    type="button"
                  >
                    <FontAwesomeIcon className="mdi" icon={faPencil} />
                    &nbsp;Edit
                  </Link>
                  &nbsp;
                  <Link
                    onClick={(e, s) => {
                      setSelectedQuestionForDeletion(datum);
                    }}
                    class="button is-danger is-small is-fullwidth-mobile"
                    type="button"
                  >
                    <FontAwesomeIcon className="mdi" icon={faTrash} />
                    &nbsp;Delete
                  </Link>
                </div>
              </div>
            )}
            <FormErrorBox errors={errors} />

            {/* <p class="pb-4">Please fill out all the required fields before submitting this form.</p> */}

            {isFetching ? (
              <PageLoadingContent displayMessage={"Please wait..."} />
            ) : (
              <>
                {datum && (
                  <div class="container" key={datum.id}>
                    <p class="subtitle is-6">
                      <FontAwesomeIcon className="fas" icon={faEye} />
                      &nbsp;Detail
                    </p>
                    <hr />
                    <DataDisplayRowText label="Title" value={datum.title} />
                    {datum.subtitle && (
                      <DataDisplayRowText
                        label="Subtitle"
                        value={datum.subtitle}
                      />
                    )}
                    <DataDisplayRowText
                      label="Status"
                      value={datum.Status ? "Active" : "Archived"}
                    />

                    <DataDisplayRowText
                      label="IsMultiSelect?"
                      value={datum.isMultiselect ? "Yes" : "No"}
                    />

                    {datum && datum.options && (
                      <div className="content">
                        <label className="label">Options</label>
                        <ul>
                          {datum.options.map((op, i) => (
                            <li key={i}>{op}</li>
                          ))}
                        </ul>
                      </div>
                    )}

                    <div class="columns pt-5">
                      <div class="column is-half">
                        <Link
                          class="button is-fullwidth-mobile"
                          to={`/admin/video-categories`}
                        >
                          <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                          &nbsp;Back to Questions
                        </Link>
                      </div>
                      <div class="column is-half has-text-right">
                        <Link
                          to={`/admin/video-category/${id}/update`}
                          class="button is-warning is-fullwidth-mobile"
                        >
                          <FontAwesomeIcon className="fas" icon={faPencil} />
                          &nbsp;Edit
                        </Link>
                      </div>
                    </div>
                  </div>
                )}
              </>
            )}
          </nav>
        </section>
      </div>
    </>
  );
}

export default AdminQuestionnaireDetail;
