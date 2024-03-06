import React, { useEffect } from "react";
import { Link } from "react-router-dom";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faDumbbell, faGauge, faArrowRight, faUsers, faChalkboardTeacher } from '@fortawesome/free-solid-svg-icons';
import { useRecoilState } from 'recoil';

import { topAlertMessageState, topAlertStatusState } from "../../AppState";
import Layout from "../Menu/Layout";

function TrainerDashboard() {

  ////
  //// Global state.
  ////

  const [topAlertMessage, setTopAlertMessage] = useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] = useRecoilState(topAlertStatusState);

  ////
  //// Component states.
  ////

  ////
  //// API.
  ////

  ////
  //// BREADCRUMB
  ////
  const breadcrumbItems = {
    items: [
      { text: 'Dashboard', link: '#', isActive: true, icon: faGauge }
    ]
  }

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
      window.scrollTo(0, 0);  // Start the page at the top of the page.
    }

    return () => { mounted = false; }
  }, []);


  ////
  //// Component rendering.
  ////

  return (
    <Layout breadcrumbItems={breadcrumbItems}>
      <div class="box">
        <div class="columns">
          <div class="column">
            <h1 class="title is-4"><FontAwesomeIcon className="fas" icon={faGauge} />&nbsp;Admin Dashboard</h1>
          </div>
        </div>

        <section class="hero is-medium is-link">
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
              <Link to={"/admin/sessions/calendar"}>View Calendar&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} /></Link>
            </p>
          </div>
        </section>

        <section class="hero is-medium is-info">
          <div class="hero-body">
            <p class="title">
              <FontAwesomeIcon className="fas" icon={faChalkboardTeacher} />&nbsp;Trainers
            </p>
            <p class="subtitle">
              Manage the trainers by clicking below:
              <br />
              <br />
              <Link to={"/admin/trainers"}>View&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} /></Link>
              <br />
              <br />
              <Link to={"/admin/trainers/add"}>Add&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} /></Link>
            </p>
          </div>
        </section>

        <section class="hero is-medium is-success">
          <div class="hero-body">
            <p class="title">
              <FontAwesomeIcon className="fas" icon={faUsers} />&nbsp;Members
            </p>
            <p class="subtitle">
              Manage the members that belong to your system.
              <br />
              <br />
              <Link to={"/admin/members"}>View&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} /></Link>
              <br />
              <br />
              <Link to={"/admin/members/add"}>Add&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} /></Link>
            </p>
          </div>
        </section>

      </div>
    </Layout>
  );
}

export default TrainerDashboard;
