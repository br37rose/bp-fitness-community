import React, { useEffect } from "react";
import { Link } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faVideo,
  faDumbbell,
  faGauge,
  faArrowRight,
  faUsers,
  faChalkboardTeacher,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import { topAlertMessageState, topAlertStatusState } from "../../AppState";

function AdminDashboard() {
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

  ////
  //// API.
  ////

  ////
  //// Event handling.
  ////

  ////
  //// Misc.
  ////

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.
    }

    return () => {
      mounted = false;
    };
  }, []);

  ////
  //// Component rendering.
  ////

  return (
    <>
      <div class="container is-fluid">
        <section class="section">
          <nav class="breadcrumb" aria-label="breadcrumbs">
            <ul>
              <li class="is-active">
                <Link to="/dashboard" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faGauge} />
                  &nbsp;Dashboard
                </Link>
              </li>
            </ul>
          </nav>
          <nav class="box">
            <div class="columns">
              <div class="column">
                <h1 class="title is-4">
                  <FontAwesomeIcon className="fas" icon={faGauge} />
                  &nbsp;Admin Dashboard
                </h1>
              </div>
            </div>

            {/*

            <section class="hero is-medium is-primary">
              <div class="hero-body">
                <p class="title">
                  <FontAwesomeIcon className="fas" icon={faDumbbell} />&nbsp;Workouts
                </p>
                <p class="subtitle">
                  Manage the classes and schedule by clicking below:
                  <br />
                  <br />
                  <Link to={"/admin/classes"}>View Classes&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} /></Link>
                  <br />
                  <br />
                  <Link to={"/admin/sessions/calendar"}>View Sessions Calendar&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} /></Link>
                </p>
              </div>
            </section>



            */}

            <section class="hero is-medium is-dark">
              <div class="hero-body">
                <p class="title">
                  <FontAwesomeIcon className="fas" icon={faUsers} />
                  &nbsp;Members
                </p>
                <p class="subtitle">
                  Manage the members that belong to your system.
                  <br />
                  <br />
                  <Link to={"/admin/members"}>
                    View&nbsp;
                    <FontAwesomeIcon className="fas" icon={faArrowRight} />
                  </Link>
                  <br />
                  <br />
                  <Link to={"/admin/members/add"}>
                    Add&nbsp;
                    <FontAwesomeIcon className="fas" icon={faArrowRight} />
                  </Link>
                </p>
              </div>
            </section>

            <section class="hero is-medium is-dark">
              <div class="hero-body">
                <p class="title">
                  <FontAwesomeIcon className="fas" icon={faChalkboardTeacher} />
                  &nbsp;Exercises
                </p>
                <p class="subtitle">
                  Manage the exercises by clicking below:
                  <br />
                  <br />
                  <Link to={"/admin/exercises"}>
                    View&nbsp;
                    <FontAwesomeIcon className="fas" icon={faArrowRight} />
                  </Link>
                  <br />
                  <br />
                  <Link to={"/admin/exercises/add"}>
                    Add&nbsp;
                    <FontAwesomeIcon className="fas" icon={faArrowRight} />
                  </Link>
                </p>
              </div>
            </section>

            <section class="hero is-medium is-dark">
              <div class="hero-body">
                <p class="title">
                  <FontAwesomeIcon className="fas" icon={faVideo} />
                  &nbsp;Video Collections
                </p>
                <p class="subtitle">
                  Manage the videos that belong to your system.
                  <br />
                  <br />
                  <Link to={"/admin/video-collections"}>
                    View&nbsp;
                    <FontAwesomeIcon className="fas" icon={faArrowRight} />
                  </Link>
                  <br />
                  <br />
                  <Link to={"/admin/video-collections/add"}>
                    Add&nbsp;
                    <FontAwesomeIcon className="fas" icon={faArrowRight} />
                  </Link>
                </p>
              </div>
            </section>
          </nav>
        </section>
      </div>
    </>
  );
}

export default AdminDashboard;
