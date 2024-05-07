import { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTrash,
  faClock,
  faArrowLeft,
  faGauge,
  faPencil,
  faTrophy,
  faEye,
  faIdCard,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import { useParams } from "react-router-dom";

import {
  getFitnessPlanDetailAPI,
  deleteFitnessPlanAPI,
} from "../../../API/FitnessPlan";
import FormErrorBox from "../../Reusable/FormErrorBox";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import { topAlertMessageState, topAlertStatusState } from "../../../AppState";
import DataDisplayRowText from "../../Reusable/DataDisplayRowText";
import FitnessPlanDisplay from "../../Reusable/FitnessPlanDisplay";
import ExerciseDisplay from "../../Reusable/ExerciseDisplay";
import {
  FITNESS_GOAL_STATUS_QUEUED,
  FITNESS_GOAL_STATUS_ACTIVE,
  FITNESS_GOAL_STATUS_ARCHIVED,
  FITNESS_GOAL_STATUS_ERROR,
  FITNESS_GOAL_STATUS_IN_PROGRESS,
  FITNESS_GOAL_STATUS_PENDING,
} from "../../../Constants/App";
import Layout from "../../Menu/Layout";

function MemberFitnessPlanDetail() {
  ////
  //// URL Parameters.
  ////

  const { id } = useParams();

  ////
  //// Global state.
  ////

  const [topAlertMessage, setTopAlertMessage] =
    useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] =
    useRecoilState(topAlertStatusState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [datum, setDatum] = useState({});
  const [tabIndex, setTabIndex] = useState(1);
  const [selectedFitnessPlanForDeletion, setSelectedFitnessPlanForDeletion] =
    useState(null);

  ////
  //// Event handling.
  ////

  const onDeleteConfirmButtonClick = () => {
    deleteFitnessPlanAPI(
      selectedFitnessPlanForDeletion.id,
      onFitnessPlanDeleteSuccess,
      onFitnessPlanDeleteError,
      onFitnessPlanDeleteDone
    );
    setSelectedFitnessPlanForDeletion(null);
  };

  ////
  //// API.
  ////

  // --- Detail --- //

  function onFitnessPlanDetailSuccess(response) {
    setDatum(response);
  }

  function onFitnessPlanDetailError(apiErr) {
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onFitnessPlanDetailDone() {
    setFetching(false);
  }

  // --- Delete --- //

  function onFitnessPlanDeleteSuccess(response) {
    // Update notification.
    setTopAlertStatus("success");
    setTopAlertMessage("Fitness plan deleted");
    setTimeout(() => {
      setTopAlertMessage("");
    }, 2000);

    // Redirect back to the video categories page.
    setForceURL("/fitness-plans");
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
      {
        text: "Fitness Plans",
        link: "/fitness-plans",
        icon: faTrophy,
        isActive: false,
      },
      { text: "Detail", link: "#", icon: faEye, isActive: true },
    ],
    mobileBackLinkItems: {
      link: "/fitness-plans",
      text: "Back to Fitness Plans",
      icon: faArrowLeft,
    },
  };

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.

      setFetching(true);
      getFitnessPlanDetailAPI(
        id,
        onFitnessPlanDetailSuccess,
        onFitnessPlanDetailError,
        onFitnessPlanDetailDone
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
    <Layout breadcrumbItems={breadcrumbItems}>
      {/* Modal */}
      <nav>
        {/* Delete modal */}
        <div
          class={`modal ${
            selectedFitnessPlanForDeletion !== null ? "is-active" : ""
          }`}
        >
          <div class="modal-background"></div>
          <div class="modal-card">
            <header class="modal-card-head">
              <p class="modal-card-title">Are you sure?</p>
              <button
                class="delete"
                aria-label="close"
                onClick={(e, ses) => setSelectedFitnessPlanForDeletion(null)}
              ></button>
            </header>
            <section class="modal-card-body">
              You are about to delete this fitness plan and all the data
              associated with it. This action is cannot be undone. Are you sure
              you would like to continue?
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
                onClick={(e, ses) => setSelectedFitnessPlanForDeletion(null)}
              >
                Cancel
              </button>
            </footer>
          </div>
        </div>
      </nav>

      {/* Page */}
      <div class="box">
        {datum && (
          <div class="columns">
            <div class="column">
              <p class="title is-4">
                <FontAwesomeIcon className="fas" icon={faTrophy} />
                &nbsp;Fitness Plan
              </p>
            </div>
            {datum.status === FITNESS_GOAL_STATUS_ACTIVE && (
              <div class="column has-text-right">
                <Link
                  to={`/fitness-plan/${id}/update`}
                  class="button is-warning is-small is-fullwidth-mobile"
                  type="button"
                >
                  <FontAwesomeIcon className="mdi" icon={faPencil} />
                  &nbsp;Edit & Re-request
                </Link>
                &nbsp;
                <Link
                  onClick={(e, s) => {
                    setSelectedFitnessPlanForDeletion(datum);
                  }}
                  class="button is-danger is-small is-fullwidth-mobile"
                  type="button"
                >
                  <FontAwesomeIcon className="mdi" icon={faTrash} />
                  &nbsp;Delete
                </Link>
              </div>
            )}
          </div>
        )}
        <FormErrorBox errors={errors} />

        {/* <p class="pb-4">Please fill out all the required fields before submitting this form.</p> */}

        {isFetching ? (
          <PageLoadingContent displayMessage={"Please wait..."} />
        ) : (
          <>
            {datum && (
              <div key={datum.id}>
                {/*
                  ---------------------------------------------
                  Queue Status GUI
                  ---------------------------------------------
                */}
                {(datum.status === FITNESS_GOAL_STATUS_QUEUED ||
                  datum.status === FITNESS_GOAL_STATUS_IN_PROGRESS ||
                  datum.status === FITNESS_GOAL_STATUS_PENDING) && (
                  <>
                    <section className="hero is-medium has-background-white-ter">
                      <div className="hero-body">
                        <p className="title">
                          <FontAwesomeIcon className="fas" icon={faClock} />
                          &nbsp;Fitness Plan Submitted
                        </p>
                        <p className="subtitle">
                          You have successfully submitted this fitness plan to
                          our team. The estimated time until our team completes
                          your fitness plan will take about <b>1 or 2 days</b>.
                          Please check back later.
                        </p>
                      </div>
                    </section>
                  </>
                )}

                {/*
                  ---------------------------------------------
                  Active Status GUI
                  ---------------------------------------------
                */}
                {datum.status === FITNESS_GOAL_STATUS_ACTIVE && (
                  <>
                    {/* Tab navigation */}

                    <div class="tabs is-medium is-size-7-mobile">
                      <ul>
                        <li class="is-active">
                          <Link>
                            <strong>Detail</strong>
                          </Link>
                        </li>
                        <li>
                          <Link
                            to={`/fitness-plan/${datum.id}/submission-form`}
                          >
                            Submission Form
                          </Link>
                        </li>
                      </ul>
                    </div>

                    <p class="title is-6">META</p>
                    <hr />

                    <DataDisplayRowText label="Name" value={datum.name} />

                    <p class="title is-6 pt-5">
                      <FontAwesomeIcon className="fas" icon={faIdCard} />
                      &nbsp;DETAIL
                    </p>
                    <hr />

                    <DataDisplayRowText
                      label="Instructions"
                      value={datum.instructions}
                      type="text_with_linebreaks"
                    />

                    {datum.weeklyFitnessPlans !== null &&
                      datum.weeklyFitnessPlans.length > 0 && (
                        <FitnessPlanDisplay
                          weeklyFitnessPlans={datum.weeklyFitnessPlans}
                        />
                      )}

                    {datum.exercises !== null && datum.exercises.length > 0 && (
                      <ExerciseDisplay
                        exercises={datum.exercises}
                        label="Main Exercises"
                      />
                    )}
                  </>
                )}

                <div class="columns pt-5">
                  <div class="column is-half">
                    <Link class="button is-hidden-touch" to={`/fitness-plans`}>
                      <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                      &nbsp;Back to fitness plans
                    </Link>
                    <Link
                      class="button is-fullwidth is-hidden-desktop"
                      to={`/fitness-plans`}
                    >
                      <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                      &nbsp;Back to fitness plans
                    </Link>
                  </div>
                  <div class="column is-half has-text-right">
                    <Link
                      class="button is-success is-hidden-touch"
                      to={`/fitness-plan/${id}/update`}
                    >
                      <FontAwesomeIcon className="fas" icon={faPencil} />
                      &nbsp;Edit & Re-request
                    </Link>
                    <Link
                      class="button is-success is-fullwidth is-hidden-desktop"
                      to={`/fitness-plan/${id}/update`}
                    >
                      <FontAwesomeIcon className="fas" icon={faPencil} />
                      &nbsp;Edit & Re-request
                    </Link>
                  </div>
                </div>

                {/*
                  ---------------------------------------------
                  Archived Status GUI
                  ---------------------------------------------
                */}
                {datum.status === FITNESS_GOAL_STATUS_ARCHIVED && (
                  <>
                    <section className="hero is-medium has-background-white-ter">
                      <div className="hero-body">
                        <p className="title">
                          <FontAwesomeIcon className="fas" icon={faClock} />
                          &nbsp;Fitness Plan Archived
                        </p>
                        <p className="subtitle">
                          This fitness plan has been archived.
                        </p>
                      </div>
                    </section>
                  </>
                )}

                {/*
                  ---------------------------------------------
                  Error Status GUI
                  ---------------------------------------------
                */}
                {datum.status === FITNESS_GOAL_STATUS_ERROR && (
                  <>
                    <section className="hero is-medium has-background-white-ter">
                      <div className="hero-body">
                        <p className="title">
                          <FontAwesomeIcon className="fas" icon={faClock} />
                          &nbsp;Fitness Plan Problem
                        </p>
                        <p className="subtitle">
                          There appears to be an problem with your fitness plan
                          submission. We are investigating and working through
                          the issue. Please check in again in another day.
                        </p>
                      </div>
                    </section>
                  </>
                )}
              </div>
            )}
          </>
        )}
      </div>
    </Layout>
  );
}

export default MemberFitnessPlanDetail;
