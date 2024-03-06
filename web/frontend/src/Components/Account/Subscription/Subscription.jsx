import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faCheck, faLock, faEllipsis, faXmark, faCartShopping, faArrowRight, faFileInvoiceDollar, faTasks, faTachometer, faKey, faArrowLeft, faCheckCircle, faUserCircle, faGauge, faPencil, faIdCard, faAddressBook, faContactCard, faChartPie } from '@fortawesome/free-solid-svg-icons'
import { useRecoilState } from 'recoil';

import { getAccountDetailAPI } from "../../../API/Account";
import { postSubscriptionCancelAPI } from "../../../API/PaymentProcessor";
import { getMemberDetailAPI } from "../../../API/member";
import FormErrorBox from "../../Reusable/FormErrorBox";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import { topAlertMessageState, topAlertStatusState, currentUserState } from "../../../AppState";
// import Footer from "../Menu/Footer";
// import PricingTable from "../Reusable/Pricing";
import { PAY_FREQUENCY } from "../../../Constants/FieldOptions";
import Layout from "../../Menu/Layout";


function AccountSubscriptionDetailAndCancel() {
    ////
    //// Global state.
    ////

    const [topAlertMessage, setTopAlertMessage] = useRecoilState(topAlertMessageState);
    const [topAlertStatus, setTopAlertStatus] = useRecoilState(topAlertStatusState);
    const [currentUser, setCurrentUser] = useRecoilState(currentUserState);

    ////
    //// Component states.
    ////

    const [errors, setErrors] = useState({});
    const [isFetching, setFetching] = useState(false);
    const [forceURL, setForceURL] = useState("");
    const [showModal, setShowModal] = useState(false);

    ////
    //// Event handling.
    ////

    const getMemberDetailUpdated = () => {
        getMemberDetailAPI(
            currentUser.id,
            getMemberDetailUpdatedSuccess,
            getMemberDetailUpdatedError,
            getMemberDetailUpdatedDone
        );
    };

    const onConfirmMySubscriptionCancelClick = () => {
        setFetching(true);
        setShowModal(false);
        setErrors({});

        postSubscriptionCancelAPI(
            "",
            onMySubscriptionCancelSuccess,
            onMySubscriptionCancelError,
            onMySubscriptionCancelDone
        );
    };

    ////
    //// API.
    ////

    const onMySubscriptionCancelSuccess = (response) => {
        console.log("onMySubscriptionCancelSuccess: Starting...");
        setTopAlertStatus("success");
        setTopAlertMessage("MySubscription Cancelled");

        // Invokes member detail API to update the recoil state
        getMemberDetailUpdated();

        setTimeout(() => {
            setTopAlertMessage("");
        }, 2000);
    };

    const onMySubscriptionCancelError = (apiErr) => {
        console.log("onMySubscriptionCancelError: Starting...");
        setTopAlertMessage("MySubscription Cancellation is unsuccessful");
        setTopAlertStatus("danger");
        setTimeout(() => {
            setTopAlertMessage("");
        }, 2000);
        setErrors(apiErr);
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    };

    const onMySubscriptionCancelDone = () => {
        console.log("onMySubscriptionCancelDone: Starting...");
        setFetching(false);
    };

    const getMemberDetailUpdatedSuccess = (response) => {
        console.log("getMemberDetailUpdatedSuccess: Starting...");
        setCurrentUser(response);
    };

    const getMemberDetailUpdatedError = (apiErr) => {
        console.log("getMemberDetailUpdatedError: Starting...");
        setErrors(apiErr);

        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    };

    const getMemberDetailUpdatedDone = () => {
        console.log("onMySubscriptionCancelDone: Starting...");
        setFetching(false);
    };

    function onAccountDetailSuccess(response) {
        console.log("onAccountDetailSuccess: Starting...");
        setCurrentUser(response);
    }

    function onAccountDetailError(apiErr) {
        console.log("onAccountDetailError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onAccountDetailDone() {
        console.log("onAccountDetailDone: Starting...");
        setFetching(false);
    }

    ////
    //// Misc.
    ////

    useEffect(() => {
        let mounted = true;

        if (mounted) {
            window.scrollTo(0, 0);  // Start the page at the top of the page.
            setFetching(true);
            setErrors({});
            getAccountDetailAPI(
                onAccountDetailSuccess,
                onAccountDetailError,
                onAccountDetailDone
            );
        }

        return () => { mounted = false; }
    }, []);
    ////
    //// Component rendering.
    ////

    if (forceURL !== "") {
        return <Navigate to={forceURL} />
    }

    // For convienence variables...
    const { isSubscriber, stripeSubscription } = currentUser;

    return (
        <div>
            {/* Modal(s) */}
            {showModal && (
                <div className={`modal ${showModal ? "is-active" : ""}`}>
                    <div className="modal-background"></div>
                    <div className="modal-card">
                        <header className="modal-card-head">
                            <p className="is-size-4 modal-card-title">Are you sure?</p>
                            <button
                                className="delete"
                                aria-label="cancel"
                                onClick={(i) => setShowModal(!showModal)}
                            ></button>
                        </header>
                        <p className="is-size-6 modal-card-body">
                            You are about to <b>cancel</b> your{" "}
                            {stripeSubscription !== undefined &&
                                stripeSubscription !== null &&
                                stripeSubscription !== "" && (
                                    <span>{`${stripeSubscription.interval}ly`}</span>
                                )}{" "}
                            subscription; This action cannot be reversed. Are you sure you
                            would like to continue?
                        </p>
                        <footer className="modal-card-foot">
                            <button
                                className="button is-danger"
                                onClick={onConfirmMySubscriptionCancelClick}
                            >
                                Continue
                            </button>
                            <button
                                className="button"
                                onClick={(i) => setShowModal(!showModal)}
                            >
                                Cancel
                            </button>
                        </footer>
                    </div>
                </div>
            )}
            {isSubscriber !== undefined && isSubscriber !== null && isSubscriber !== "" && isSubscriber === true
                ? <>
                    {stripeSubscription !== undefined && stripeSubscription !== null && stripeSubscription !== "" &&
                        <>
                            {stripeSubscription.offerPurchase !== undefined && stripeSubscription.offerPurchase !== null && stripeSubscription.offerPurchase !== ""
                                ?
                                <>
                                    {/*
                                                                DEVELOPERS NOTE:
                                                                THE FOLLOWING CODE IS TO RENDER MORE SPECIFIC SUBSCRIPTION DATA.
                                                            */}

                                    <div className="columns">

                                        <div className="column"></div>
                                        <div className="column">
                                            <PricingTable
                                                title={stripeSubscription.offerPurchase.offerName}
                                                price={stripeSubscription.offerPurchase.offerPrice}
                                                period={`${stripeSubscription.offerPurchase.offerPayFrequency === PAY_FREQUENCY[1].value ? 'month' : 'year'}`}
                                                description={stripeSubscription.offerPurchase.offerDescription}
                                                url={"https://images.pexels.com/photos/841130/pexels-photo-841130.jpeg?cs=srgb&dl=pexels-victor-freitas-841130.jpg&fm=jpg"}
                                                highlightedText="Current Plan"
                                                customComponent={
                                                    <CustomComponent
                                                        showModal={showModal}
                                                        setShowModal={setShowModal}
                                                    />
                                                }
                                            />
                                        </div>
                                        <div className="column"></div>
                                    </div>

                                </>
                                :
                                <>
                                    {/*
                                                                DEVELOPERS NOTE:
                                                                THE FOLLOWING CODE IS TO RENDER GENERIC SUBSCRIPTION DATA. RUN THIS CODE IF ERROR.
                                                            */}
                                    <p className="is-size-6 has-text-secondary pb-1">
                                        You are subscribed to our{" "}
                                        {`${stripeSubscription.interval === 'year' ? "annual" : "monthly"}`} premium plan.
                                    </p>
                                    <div className="columns is-centered mt-4">
                                        <div className="column is-narrow">
                                            <Link to="/account/subscription/invoices" className="button is-dark">
                                                <FontAwesomeIcon className="fas" icon={faCartShopping} />
                                                &nbsp;View Past Invoices
                                            </Link>
                                        </div>
                                        <div className="column is-narrow">
                                            <button className="button is-danger" onClick={(i) => setShowModal(!showModal)}><FontAwesomeIcon className="fas" icon={faXmark} />&nbsp;Cancel My Subscription</button>
                                        </div>
                                    </div>
                                </>
                            }
                        </>
                    }
                </>
                : <section class="hero has-background-white-ter">
                    <div class="hero-body">
                        <p class="title is-hidden-mobile">
                            <FontAwesomeIcon className="fas" icon={faFileInvoiceDollar} />&nbsp;No Subscription Plan
                        </p>
                        <p class="title is-6 is-hidden-tablet">
                            <FontAwesomeIcon className="fas" icon={faFileInvoiceDollar} />&nbsp;No Subscription Plan
                        </p>
                        <p class="subtitle">
                            No subscription enrollment. <b><Link to="/account/more/subscribe">Click here&nbsp;<FontAwesomeIcon className="mdi" icon={faArrowRight} /></Link></b>{" "}to get started enrolling in a new subscription.
                        </p>
                    </div>
                </section>
            }

            <div class="columns pt-5">
                <div class="column is-half">
                    <Link class="button is-hidden-touch" to={"/dashboard"}><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to Dashboard</Link>
                    <Link class="button is-fullwidth is-hidden-desktop" to={"/dashboard"}><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to Dashboard</Link>
                </div>
                <div class="column is-half has-text-right">
                    {/*
                                            <button class="button is-primary is-hidden-touch" onClick={onSubmitClick}><FontAwesomeIcon className="fas" icon={faCheckCircle} />&nbsp;Save</button>
                                            <button class="button is-primary is-fullwidth is-hidden-desktop" onClick={onSubmitClick}><FontAwesomeIcon className="fas" icon={faCheckCircle} />&nbsp;Save</button>
                                            */}
                </div>
            </div>

        </div>
    );
}

export default AccountSubscriptionDetailAndCancel;

const CustomComponent = ({ setShowModal, showModal }) => {
    return (
        <div className="columns is-centered mt-4">
            <div className="column is-narrow">
                <Link to="/account/subscription/invoices" className="button is-dark">
                    <FontAwesomeIcon className="fas" icon={faCartShopping} />
                    &nbsp;View Past Invoices
                </Link>
            </div>
            <div className="column is-narrow">
                <button className="button is-danger" onClick={(i) => setShowModal(!showModal)}><FontAwesomeIcon className="fas" icon={faXmark} />&nbsp;Cancel My Subscription</button>
            </div>
        </div>
    );
}



const PricingTable = props => (
    <PricingOption {...props} />
);


// Main PricingOption Component
const PricingOption = ({
    title,
    url,
    subtitle,
    description,
    price,
    period,
    features,
    highlightedText,
    isLocked,
    customComponent,
    customClasses = "",
    customStyles = {}
}) => (
    <div class={`card is-flex is-flex-direction-column ${customClasses} ${isLocked ? 'is-locked' : ''}`} style={{ ...customStyles, minWidth: "100%", position: "relative" }}>
        {isLocked && (
            <div class="is-size-4 has-text-centered lock-icon">
                <FontAwesomeIcon icon={faLock} /><br />Purchased
            </div>
        )}
        {highlightedText && (
            <div style={{ position: "absolute", top: "0", left: "0", right: "0", margin: "auto", zIndex: "1" }} class="highlightedText-badge is-uppercase has-text-centered has-background-primary has-text-white has-text-weight-semibold py-2">{highlightedText}</div>
        )}
        {url && <CardImage url={url} title={title} />}
        <div class="card-content is-flex is-flex-direction-column is-flex-grow-1">
            {title && <h3 class="title is-5 has-text-centered is-primary mb-4">{title}</h3>}
            {subtitle && <p class="subtitle has-text-centered is-6 mb-5">{subtitle}</p>}
            <Pricing price={price} period={period} />
            <div class="is-flex-grow-1">
                {description && <p class="subtitle has-text-centered is-6 my-5">{description}</p>}
                {features && <FeaturesList features={features} />}
            </div>
            <div>
                {customComponent}
            </div>
        </div>
    </div>
);

// CardImage Component
const CardImage = ({ url, title }) => (
    <div class="card-image">
        {/* `is-4by5` can be found via https://bulma.io/documentation/elements/image/#responsive-images-with-ratios */}
        <figure class="image is-4by5">
            <img src={url} alt={title} />
        </figure>
    </div>
);

// Pricing Component
const Pricing = ({ price, period }) => (
    <div class="pricing has-text-centered my-4">
        <span class="price is-size-2 has-text-weight-bold">${price}</span>
        {period && <span class='text-highlight'>{` / ${period}`}</span>}
    </div>
);

// FeaturesList Component
const FeaturesList = ({ features }) => (
    <ul class="features-list subtitle has-text-centered">
        {features.map((feature, index) => (
            <li key={index} class="has-text-dark">
                <span class="icon is-large has-text-primary mr-2">
                    <FontAwesomeIcon icon={faCheck} />
                </span>
                <span>{feature}</span>
            </li>
        ))}
    </ul>
);
