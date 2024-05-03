import React, { useState, useEffect } from "react";
import { Link, Navigate, useNavigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faPlus,
  faGauge,
  faArrowRight,
  faTable,
  faRefresh,
  faArrowUp,
  faChild,
  faRunning,
  faBolt,
  faHeart,
  faWeight,
  faMapMarkedAlt,
  faShoePrints,
  faUser,
  faTrophy,
  faDumbbell,
  faVideoCamera,
  faLeaf,
} from "@fortawesome/free-solid-svg-icons";
import { faFreeCodeCamp } from "@fortawesome/free-brands-svg-icons";
import { useRecoilState } from "recoil";

import FormErrorBox from "../Reusable/FormErrorBox";
import { getMySummaryAPI } from "../../API/Biometric";
import {
  topAlertMessageState,
  topAlertStatusState,
  currentUserState,
} from "../../AppState";
import PageLoadingContent from "../Reusable/PageLoadingContent";
import { formatDateStringWithTimezone } from "../../Helpers/timeUtility";
import AverageAndTimeComparison from "../Reusable/AverageDateAndTimeComparison";
import Layout from "../Menu/Layout";

function MemberDashboard() {
  ////
  //// Global state.
  ////

  const navigate = useNavigate();

  const [topAlertMessage, setTopAlertMessage] =
    useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] =
    useRecoilState(topAlertStatusState);
  const [currentUser] = useRecoilState(currentUserState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [datum, setDatum] = useState("");
  const [isFetching, setFetching] = useState(false);
  // const [isComingSoon, setComingSoon] = useState(true);
  const [forceURL, setForceURL] = useState("");
  ////
  //// Event handling.
  ////

  const handleNavigateToAccount = (e) => {
    e.preventDefault();
    navigate("/account", { state: { activeTabProp: "wearableTech" } });
  };

  const getDatum = (user) => {
    if (user !== undefined && user !== null && user !== "") {
      if (user.primaryHealthTrackingDeviceType !== 0) {
        setFetching(true);
        setErrors({});

        let params = new Map();
        params.set("user_id", user.id);

        getMySummaryAPI(
          params,
          onSummarySuccess,
          onSummaryError,
          onSummaryDone
        );
      } else {
        console.log("user does not have a device, prevented pulling data.");
      }
    }
  };

  ////
  //// API.
  ////

  // --- SUMMARY --- //

  function onSummarySuccess(response) {
    setDatum(response);
  }

  function onSummaryError(apiErr) {
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onSummaryDone() {
    setFetching(false);
  }

  ////
  //// BREADCRUMB
  ////
  const breadcrumbItems = {
    items: [{ text: "Dashboard", link: "#", isActive: true, icon: faGauge }],
  };

  ////
  //// Misc.
  ////
  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.

      getDatum(currentUser);
    }

    return () => {
      mounted = false;
    };
  }, [currentUser]);

  if (forceURL !== undefined && forceURL !== null && forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  return (
    <Layout breadcrumbItems={breadcrumbItems}>
      {/* Wizard Component */}
      {/* <OnBoardingQuestionWizard questions={questionsData} /> */}

      <div className="box">
        <div className="columns">
          <div className="column">
            <h1 className="title is-4">
              <FontAwesomeIcon className="fas" icon={faGauge} />
              &nbsp;Dashboard
            </h1>
          </div>
          <div className="column has-text-right">
            <button
              onClick={() => getDatum(currentUser)}
              class="is-fullwidth-mobile button is-link is-small"
              type="button"
            >
              <FontAwesomeIcon className="mdi" icon={faRefresh} />
              &nbsp;
              <span class="is-hidden-desktop is-hidden-tablet">Refresh</span>
            </button>
            {/*
                    DEVELOPERS NOTE:
                    - If the logged in user doesn't have a device registered then
                      show the following button to encourage them to register.
                */}
            {currentUser !== undefined &&
              currentUser !== null &&
              currentUser !== "" && (
                <>
                  {currentUser.primaryHealthTrackingDeviceType === 0 && (
                    <>
                      &nbsp;
                      <Link
                        to={`/account/wearable-tech`}
                        className="is-fullwidth-mobile button is-small is-success"
                        type="button"
                      >
                        <FontAwesomeIcon className="mdi" icon={faPlus} />
                        &nbsp;Register Wearable
                      </Link>
                    </>
                  )}
                </>
              )}
          </div>
        </div>
        {isFetching ? (
          <PageLoadingContent displayMessage={"Please wait..."} />
        ) : (
          <>
            <FormErrorBox errors={errors} />
            {datum !== undefined && datum !== null && datum !== "" ? (
              <section class="main_dashboard">
                <div class="container">
                  <div class="columns">
                    <div class="column">
                      <div class="box has-background-dark has-text-white">
                        <div class="is-flex is-align-items-center ">
                          <div class="">
                            <span>
                              <FontAwesomeIcon
                                className="fas is-size-1"
                                icon={faUser}
                              />
                            </span>
                          </div>
                          <div class="ml-6">
                            <h5 class="is-size-3 has-text-primary has-text-weight-semibold is-size-5-mobile">{`Hi, ${currentUser.firstName}`}</h5>
                            <p class="is-size-5 has-text-white is-size-6-mobile">
                              Here is your state as per{" "}
                              {formatDateStringWithTimezone(
                                new Date().toISOString()
                              )}
                            </p>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div class="columns">
                    <div class="column">
                      <div class="box bg_brand">
                        <div class="is-flex is-justify-content-center is-align-items-center">
                          <div class="">
                            <span>
                              <FontAwesomeIcon
                                className="fas px-3 has-text-primary is-size-1 is-size-3-mobile"
                                icon={faShoePrints}
                              />
                            </span>
                          </div>
                          <div class="ml-6">
                            <h5 class="is-size-2 is-size-4-mobile  has-text-centered has-text-weight-semibold">
                              {datum &&
                                datum.stepCountDeltaThisDaySummary &&
                                datum.stepCountDeltaThisDaySummary.sum}
                              <span class="is-size-5 has-text-weight-semibold is-size-6-mobile">
                                &nbsp;Total Steps
                              </span>
                            </h5>
                            <AverageAndTimeComparison
                              lastDay={
                                datum && datum.stepCountDeltaLastDaySummary
                              }
                              thisDay={
                                datum && datum.stepCountDeltaThisDaySummary
                              }
                              mode={2}
                            />
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div class="columns">
                    <div class="column">
                      <div class="box bg_brand is-flex is-align-items-center is-flex-direction-column	">
                        <span>
                          <FontAwesomeIcon
                            className="fas px-3 has-text-primary is-size-1"
                            icon={faMapMarkedAlt}
                          />
                        </span>
                        <h5 class="mt-2 is-size-4 has-text-centered has has-text-weight-semibold is-size-5-mobile">
                          Coming Soon
                        </h5>
                        <p class="is-size-5 has-text-weight-semibold is-size-6-mobile">
                          Total Miles
                        </p>
                        <p class="is-size-6 is-size-6-mobile">
                          <FontAwesomeIcon
                            className="fas px-3 is-size-6-mobile"
                            icon={faArrowUp}
                          />
                          &nbsp;Coming Soon
                        </p>
                        <p></p>
                      </div>
                    </div>
                    <div class="column">
                      <div class="box is-flex is-align-items-center is-flex-direction-column	">
                        <span>
                          <FontAwesomeIcon
                            className="fas px-3 has-text-primary is-size-1"
                            icon={faFreeCodeCamp}
                          />
                        </span>
                        <h5 class="mt-2 is-size-4 has-text-centered has has-text-weight-semibold is-size-5-mobile">
                          Coming Soon
                        </h5>
                        <p class="is-size-5 is-size-6-mobile has-text-weight-semibold">
                          avg calories burn
                        </p>
                        <p class="is-size-6 is-size-6-mobile">
                          <FontAwesomeIcon
                            className="fas px-3 is-size-6-mobile"
                            icon={faArrowUp}
                          />
                          &nbsp;Coming Soon
                        </p>
                        <p></p>
                      </div>
                    </div>
                    <div class="column">
                      <div class="box is-flex is-align-items-center is-flex-direction-column	">
                        <span>
                          <FontAwesomeIcon
                            className="fas px-3 has-text-primary is-size-1"
                            icon={faBolt}
                          />
                        </span>
                        <h5 class="mt-2 is-size-4 has-text-centered has has-text-weight-semibold is-size-5-mobile">
                          Coming Soon
                        </h5>
                        <p class="is-size-5 has-text-weight-semibold is-size-6-mobile">
                          Active minutes
                        </p>
                        <p class="is-size-6 is-size-6-mobile">
                          <FontAwesomeIcon
                            className="fas px-3 is-size-6-mobile"
                            icon={faArrowUp}
                          />
                          &nbsp;Coming Soon
                        </p>
                        <p></p>
                      </div>
                    </div>
                  </div>
                  <div class="columns">
                    <div class="column">
                      <div class="box bg_brand is-flex is-align-items-center is-flex-direction-column	">
                        <span>
                          <FontAwesomeIcon
                            className="fas px-3 has-text-primary is-size-1"
                            icon={faRunning}
                          />
                        </span>
                        <h5 class="mt-2 is-size-4 has-text-centered has has-text-weight-semibold is-size-5-mobile">
                          Coming Soon
                        </h5>
                        <p class="is-size-5 has-text-weight-semibold is-size-6-mobile">
                          exercising this week
                        </p>
                        <p class="is-size-6 is-size-6-mobile">
                          <FontAwesomeIcon
                            className="fas px-3"
                            icon={faArrowUp}
                          />
                          &nbsp;Coming Soon
                        </p>
                        <p></p>
                      </div>
                    </div>
                    <div class="column">
                      <div class="box bg_brand is-flex is-align-items-center is-flex-direction-column	">
                        <span>
                          <FontAwesomeIcon
                            className="fas px-3 has-text-primary is-size-1"
                            icon={faChild}
                          />
                        </span>
                        <h5 class="mt-2 is-size-4 has-text-centered has has-text-weight-semibold is-size-5-mobile">
                          Coming Soon
                        </h5>
                        <p class="is-size-5 has-text-weight-semibold is-size-6-mobile">
                          avg hrs 250+ steps
                        </p>
                        <p class="is-size-6 is-size-6-mobile">
                          <FontAwesomeIcon
                            className="fas px-3"
                            icon={faArrowUp}
                          />
                          &nbsp;Coming Soon
                        </p>
                        <p></p>
                      </div>
                    </div>
                    <div class="column">
                      <div class="box is-flex is-align-items-center is-flex-direction-column	">
                        <span>
                          <FontAwesomeIcon
                            className="fas px-3 has-text-primary is-size-1"
                            icon={faHeart}
                          />
                        </span>
                        <h5 class="mt-2 is-size-4 has-text-centered has has-text-weight-semibold is-size-5-mobile">
                          {datum &&
                            datum.heartRateThisDaySummary &&
                            datum.heartRateThisDaySummary.average.toFixed(
                              2
                            )}{" "}
                        </h5>
                        <p class="is-size-5 has-text-weight-semibold is-size-6-mobile">
                          avg. resting heart rate
                        </p>
                        <AverageAndTimeComparison
                          lastDay={datum && datum.heartRateLastDaySummary}
                          thisDay={datum && datum.heartRateThisDaySummary}
                          iconState={true}
                          mode={1}
                        />
                      </div>
                    </div>
                  </div>
                  <div class="columns">
                    <div class="column">
                      <div class="box is-flex is-align-items-center is-flex-direction-column	">
                        <span>
                          <FontAwesomeIcon
                            className="fas px-3 has-text-primary is-size-1"
                            icon={faWeight}
                          />
                        </span>
                        <h5 class="mt-2 is-size-4 has-text-centered has has-text-weight-semibold is-size-5-mobile">
                          Coming Soon
                        </h5>
                        <p class="is-size-5 has-text-weight-semibold is-size-6-mobile">
                          Weight
                        </p>
                        <p class="is-size-6 is-size-6-mobile">
                          <FontAwesomeIcon
                            className="fas px-3 is-size-6-mobile"
                            icon={faArrowUp}
                          />
                          &nbsp;Coming Soon
                        </p>
                        <p></p>
                      </div>
                    </div>
                  </div>
                  <div class="columns">
                    <div class="column is-half">
                      <div class="box has-text-centered hero is-medium is-dark custom-hero">
                        <div class="hero-body">
                          <p class="title">
                            <FontAwesomeIcon
                              className="fas has-text-primary"
                              icon={faLeaf}
                            />
                            <br />
                            Nutrition Plans
                          </p>
                          <p class="subtitle">
                            View and generate nutrition plans that work for you:
                            <br />
                            <br />
                            <Link
                              className="has-text-white"
                              to={"/nutrition-plans"}
                            >
                              View&nbsp;
                              <FontAwesomeIcon
                                className="fas"
                                icon={faArrowRight}
                              />
                            </Link>
                          </p>
                        </div>
                      </div>
                    </div>
                    <div class="column is-half">
                      <div class="box has-text-centered hero is-medium is-dark custom-hero">
                        <div class="hero-body">
                          <p class="title">
                            <FontAwesomeIcon
                              className="fas has-text-primary"
                              icon={faTrophy}
                            />
                            <br />
                            Fitness Plans
                          </p>
                          <p class="subtitle">
                            View and generate fitness plans that work for you:
                            <br />
                            <br />
                            <Link
                              className="has-text-white"
                              to={"/fitness-plans"}
                            >
                              View&nbsp;
                              <FontAwesomeIcon
                                className="fas"
                                icon={faArrowRight}
                              />
                            </Link>
                          </p>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div class="columns">
                    <div class="column is-half">
                      <div class="box has-text-centered hero is-medium is-dark custom-hero">
                        <div class="hero-body">
                          <p class="title">
                            <FontAwesomeIcon
                              className="fas has-text-primary"
                              icon={faDumbbell}
                            />
                            <br />
                            Exercises
                          </p>
                          <p class="subtitle">
                            View all the exercises to help you at the gym:
                            <br />
                            <br />
                            <Link className="has-text-white" to={"/exercises"}>
                              View&nbsp;
                              <FontAwesomeIcon
                                className="fas"
                                icon={faArrowRight}
                              />
                            </Link>
                          </p>
                        </div>
                      </div>
                    </div>
                    <div class="column is-half">
                      <div class="box has-text-centered hero is-medium is-dark custom-hero">
                        <div class="hero-body">
                          <p class="title">
                            <FontAwesomeIcon
                              className="fas has-text-primary"
                              icon={faVideoCamera}
                            />
                            <br />
                            Videos
                          </p>
                          <p class="subtitle">
                            View the videos of BP8 Fitness:
                            <br />
                            <br />
                            <Link
                              className="has-text-white"
                              to={"/video-categories"}
                            >
                              View&nbsp;
                              <FontAwesomeIcon
                                className="fas"
                                icon={faArrowRight}
                              />
                            </Link>
                          </p>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </section>
            ) : (
              <>
                <section className="hero is-medium has-background-white-ter">
                  <div className="hero-body">
                    <p className="title">
                      <FontAwesomeIcon className="fas" icon={faTable} />
                      &nbsp;No Biometrics
                    </p>
                    <p className="subtitle">
                      You currently have no biometrics data.{" "}
                      <b>
                        <Link onClick={(e) => handleNavigateToAccount(e)}>
                          Click here&nbsp;
                          <FontAwesomeIcon
                            className="mdi"
                            icon={faArrowRight}
                          />
                        </Link>
                      </b>
                      to get started registering your wearable tech!
                    </p>
                  </div>
                </section>
                <section className="hero mt-5 is-medium has-background-white-ter">
                  <div className="is-medium has-background-white">
                    <div class="columns">
                      <div class="column is-half">
                        <div class="box has-text-centered hero is-medium is-dark custom-hero">
                          <div class="hero-body">
                            <p class="title">
                              <FontAwesomeIcon
                                className="fas has-text-primary"
                                icon={faLeaf}
                              />
                              <br />
                              Nutrition Plans
                            </p>
                            <p class="subtitle">
                              View and generate nutrition plans that work for
                              you:
                              <br />
                              <br />
                              <Link
                                className="has-text-white"
                                to={"/nutrition-plans"}
                              >
                                View&nbsp;
                                <FontAwesomeIcon
                                  className="fas"
                                  icon={faArrowRight}
                                />
                              </Link>
                            </p>
                          </div>
                        </div>
                      </div>
                      <div class="column is-half">
                        <div class="box has-text-centered hero is-medium is-dark custom-hero">
                          <div class="hero-body">
                            <p class="title">
                              <FontAwesomeIcon
                                className="fas has-text-primary"
                                icon={faTrophy}
                              />
                              <br />
                              Fitness Plans
                            </p>
                            <p class="subtitle">
                              View and generate fitness plans that work for you:
                              <br />
                              <br />
                              <Link
                                className="has-text-white"
                                to={"/fitness-plans"}
                              >
                                View&nbsp;
                                <FontAwesomeIcon
                                  className="fas"
                                  icon={faArrowRight}
                                />
                              </Link>
                            </p>
                          </div>
                        </div>
                      </div>
                    </div>
                    <div class="columns">
                      <div class="column is-half">
                        <div class="box has-text-centered hero is-medium is-dark custom-hero">
                          <div class="hero-body">
                            <p class="title">
                              <FontAwesomeIcon
                                className="fas has-text-primary"
                                icon={faDumbbell}
                              />
                              <br />
                              Exercises
                            </p>
                            <p class="subtitle">
                              View all the exercises to help you at the gym:
                              <br />
                              <br />
                              <Link
                                className="has-text-white"
                                to={"/exercises"}
                              >
                                View&nbsp;
                                <FontAwesomeIcon
                                  className="fas"
                                  icon={faArrowRight}
                                />
                              </Link>
                            </p>
                          </div>
                        </div>
                      </div>
                      <div class="column is-half">
                        <div class="box has-text-centered hero is-medium is-dark custom-hero">
                          <div class="hero-body">
                            <p class="title">
                              <FontAwesomeIcon
                                className="fas has-text-primary"
                                icon={faVideoCamera}
                              />
                              <br />
                              Videos
                            </p>
                            <p class="subtitle">
                              View the videos of BP8 Fitness:
                              <br />
                              <br />
                              <Link
                                className="has-text-white"
                                to={"/video-categories"}
                              >
                                View&nbsp;
                                <FontAwesomeIcon
                                  className="fas"
                                  icon={faArrowRight}
                                />
                              </Link>
                            </p>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </section>
              </>
            )}
          </>
        )}
      </div>
    </Layout>
  );
}

export default MemberDashboard;
