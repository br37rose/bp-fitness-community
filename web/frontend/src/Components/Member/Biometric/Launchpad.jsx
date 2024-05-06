import React, { useEffect } from "react";
import { Link } from "react-router-dom";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faStar, faArrowLeft, faHeartbeat, faChartLine, faRankingStar, faTrophy, faVideoCamera, faDumbbell, faTasks, faGauge, faArrowRight, faUsers, faBarcode } from '@fortawesome/free-solid-svg-icons';
import { useRecoilState } from 'recoil';

import { topAlertMessageState, topAlertStatusState } from "../../../AppState";
import Footer from "../../Menu/Footer";
import Layout from "../../Menu/Layout";


function MemberBiometricLaunchpad() {

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
      { text: 'Dashboard', link: '/dashboard', isActive: false, icon: faGauge },
      { text: 'Biometrics', link: '#', icon: faHeartbeat, isActive: true }
    ],
    mobileBackLinkItems: {
      link: "/dashboard",
      text: "Back to Dashboard",
      icon: faArrowLeft
    }
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
            <h1 class="title is-4 is-hidden-touch"><FontAwesomeIcon className="fas" icon={faHeartbeat} />
              &nbsp;Biometrics
            </h1>
            <h1 class="title is-6 is-hidden-desktop mt-2"><FontAwesomeIcon className="fas" icon={faHeartbeat} />
              &nbsp;Biometrics
            </h1>
          </div>
        </div>

        <section class="hero is-medium is-dark">
          <div class="hero-body">
            <p class="title">
              <FontAwesomeIcon className="fas" icon={faRankingStar} />&nbsp;Leaderboard
            </p>
            <p class="subtitle">
              View your ranking in the world:
              <br />
              <br />
              <Link to={"/biometrics/leaderboard/global"}>View&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} /></Link>
            </p>
          </div>
        </section>
        <section class="hero is-medium is-dark">
          <div class="hero-body">
            <p class="title">
              <FontAwesomeIcon className="fas" icon={faStar} />&nbsp;My Summary
            </p>
            <p class="subtitle">
              View the summary of your data:
              <br />
              <br />
              <Link to={"/biometrics/summary"}>View&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} /></Link>
            </p>
          </div>
        </section>
        <section class="hero is-medium is-dark">
              <div class="hero-body">
                <p class="title">
                  <FontAwesomeIcon className="fas" icon={faChartLine} />&nbsp;My History
                </p>
                <p class="subtitle">
                  View all your data:
                  <br />
                  <br />
                  <Link to={"/biometrics/history/tableview"}>View&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} /></Link>
                </p>
              </div>
            </section>

      </div>
    </Layout>

  );
}

export default MemberBiometricLaunchpad;
