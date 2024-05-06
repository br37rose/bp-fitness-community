import { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faArrowLeft,
  faGauge,
  faEye,
  faBolt,
  faPencilRuler,
  faUsers,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import { useParams } from "react-router-dom";

import FormErrorBox from "../../Reusable/FormErrorBox";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import { topAlertMessageState, topAlertStatusState } from "../../../AppState";
import {
  deletefitnessChallengeAPI,
  getfitnessChallengeDetailAPI,
} from "../../../API/FitnessChallenge";
import DataDisplayRowText from "../../Reusable/DataDisplayRowText";

function MemberFitnessChallengeDetail() {
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
  const [selectedChallengeForDeletion, setselectedChallengeForDeletion] =
    useState(null);

  ////
  //// Event handling.
  ////

  const onDeleteConfirmButtonClick = () => {
    deletefitnessChallengeAPI(id, ondeleteSuccess, ondeleteError, onDeleteDone);
    setselectedChallengeForDeletion(null);
  };

  ////
  //// API.
  ////

  // --- Detail --- //

  function onVideoCollectionDetailSuccess(response) {
    setDatum(response);
  }

  function onVideoCollectionDetailError(apiErr) {
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onVideoCollectionDetailDone() {
    setFetching(false);
  }

  // --- Delete --- //

  function ondeleteSuccess(response) {
    // Update notification.
    setTopAlertStatus("success");
    setTopAlertMessage("challenge deleted");
    setTimeout(() => {
      setTopAlertMessage("");
    }, 2000);

    // Redirect back to the members page.
    setForceURL("/fitness-challenge");
  }

  function ondeleteError(apiErr) {
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

  function onDeleteDone() {
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
      getfitnessChallengeDetailAPI(
        id,
        onVideoCollectionDetailSuccess,
        onVideoCollectionDetailError,
        onVideoCollectionDetailDone
      );
    }

    return () => {
      mounted = false;
    };
  }, [id]);
  ////
  //// Component rendering.
  ////

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  return (
    <>
      <div class="container is-fluid">
        <section class="section">
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
                <Link to="/fitness-challenge" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Challenges
                </Link>
              </li>
            </ul>
          </nav>

          {/* Modal */}
          <nav>
            {/* Delete modal */}
            <div
              class={`modal ${
                selectedChallengeForDeletion !== null ? "is-active" : ""
              }`}
            >
              <div class="modal-background"></div>
              <div class="modal-card">
                <header class="modal-card-head">
                  <p class="modal-card-title">Are you sure?</p>
                  <button
                    class="delete"
                    aria-label="close"
                    onClick={(e, ses) => setselectedChallengeForDeletion(null)}
                  ></button>
                </header>
                <section class="modal-card-body">
                  You are about to delete this challenge and all the data
                  associated with it. This action is cannot be undone. Are you
                  sure you would like to continue?
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
                    onClick={(e, ses) => setselectedChallengeForDeletion(null)}
                  >
                    Cancel
                  </button>
                </footer>
              </div>
            </div>
          </nav>

          {/* Page */}
          <nav class="box">
            {/* Title + Options */}
            {datum && (
              <div class="columns">
                <div class="column">
                  <p class="title is-4">
                    <FontAwesomeIcon className="fas" icon={faBolt} />
                    &nbsp;Challenges
                  </p>
                </div>
              </div>
            )}

            {/* Tab Navigation */}
            <div class="tabs is-medium is-size-7-mobile">
              <ul>
                <li class="is-active">
                  <Link>
                    <strong>Detail</strong>
                  </Link>
                </li>
                <li class="">
                  <Link to={`/fitness-challenge/${id}/leaderboard`}>
                    <strong>LeaderBoard</strong>
                  </Link>
                </li>
              </ul>
            </div>

            {/* <p class="pb-4">Please fill out all the required fields before submitting this form.</p> */}

            {isFetching ? (
              <PageLoadingContent displayMessage={"Please wait..."} />
            ) : (
              <>
                <div className="columns">
                  <div className="column ">
                    <h2 className="title is-4 mb-3">{datum?.name}</h2>
                    <DataDisplayRowText
                      label="Description"
                      value={datum.description}
                    />
                    <DataDisplayRowText
                      label="Duration in Weeks"
                      value={datum.durationInWeeks}
                    />
                    <DataDisplayRowText
                      label="Starts on"
                      value={datum.startTime}
                      type="datetime"
                    />
                    <p class="subtitle is-6 mt-3">
                      <FontAwesomeIcon className="fas" icon={faPencilRuler} />
                      &nbsp;Rules
                    </p>
                    {datum?.rules?.map((rule, index) => (
                      <div className="box" key={index}>
                        <h3 className="title is-5 mb-3">
                          Rule {index + 1}: {rule.name}
                        </h3>
                        <DataDisplayRowText
                          label="Description"
                          value={rule.description}
                        />
                      </div>
                    ))}
                    <p class="subtitle is-6 mt-3">
                      <FontAwesomeIcon className="fas" icon={faUsers} />
                      &nbsp;Joined Members
                    </p>
                    <div className="content">
                      <ul>
                        {datum?.userNames?.map(
                          (user, index) => user && <li key={index}>{user}</li>
                        )}
                      </ul>
                    </div>
                  </div>
                </div>
                <FormErrorBox errors={errors} />
              </>
            )}
          </nav>
        </section>
      </div>
    </>
  );
}

export default MemberFitnessChallengeDetail;
