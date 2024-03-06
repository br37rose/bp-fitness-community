import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import {
    faCreditCard, faImage, faFile, faTasks, faTachometer, faPlus, faTimesCircle,
    faCheckCircle, faUserCircle, faGauge, faPencil, faUsers, faIdCard,
    faAddressBook, faContactCard, faChartPie, faCogs, faEye, faArrowLeft,
    faLock, faCheck
} from '@fortawesome/free-solid-svg-icons'
import { useRecoilState } from 'recoil';

import { postAccountAvatarAPI, getAccountDetailAPI } from "../../../../../API/Account";
import { getOfferListAPI } from "../../../../../API/Offer";
import { postCreateStripeSubscriptionCheckoutSessionAPI } from "../../../../../API/PaymentProcessor";
import FormErrorBox from "../../../../Reusable/FormErrorBox";
import FormInputField from "../../../../Reusable/FormInputField";
import FormTextareaField from "../../../../Reusable/FormTextareaField";
import FormRadioField from "../../../../Reusable/FormRadioField";
import FormMultiSelectField from "../../../../Reusable/FormMultiSelectField";
import FormSelectField from "../../../../Reusable/FormSelectField";
import FormCheckboxField from "../../../../Reusable/FormCheckboxField";
import FormCountryField from "../../../../Reusable/FormCountryField";
import FormRegionField from "../../../../Reusable/FormRegionField";
import PageLoadingContent from "../../../../Reusable/PageLoadingContent";
import { topAlertMessageState, topAlertStatusState } from "../../../../../AppState";
import Layout from "../../../../Menu/Layout";


function AccountMoreOperationSubscribe() {
    ////
    //// Global state.
    ////

    const [topAlertMessage, setTopAlertMessage] = useRecoilState(topAlertMessageState);
    const [topAlertStatus, setTopAlertStatus] = useRecoilState(topAlertStatusState);

    ////
    //// Component states.
    ////

    const [errors, setErrors] = useState({});
    const [isFetching, setFetching] = useState(false);
    const [forceURL, setForceURL] = useState("");
    const [selectedFile, setSelectedFile] = useState(null);
    const [currentUser, setCurrentUser] = useState({});
    const [listData, setListData] = useState("");

    ////
    //// Event handling.
    ////

    // --- Select Offer --- //

    const onOfferClick = (priceID) => {
        // action={`${process.env.REACT_APP_API_PROTOCOL}://${process.env.REACT_APP_API_DOMAIN}/api/v1/stripe/create-subscription-checkout-session`}
        setFetching(true);
        postCreateStripeSubscriptionCheckoutSessionAPI(
            priceID,
            onCreateStripeSubscriptionCheckoutSessionSuccess,
            onCreateStripeSubscriptionCheckoutSessionError,
            onCreateStripeSubscriptionCheckoutSessionDone
        );
    }

    // --- Offer List --- //

    const fetchList = async () => {
        setFetching(true);
        setErrors({});

        getOfferListAPI(new Map(), onListSuccess, onListError, onListDone);
    };

    ////
    //// API.
    ////

    // --- Offer List --- //

    function onListSuccess(response) {
        console.log("onMemberListSuccess: Starting...");
        if (response.results !== null) {
            setListData(response);
        } else {
            setListData([]);
        }
    }

    function onListError(apiErr) {
        console.log("onMemberListError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onListDone() {
        console.log("onMemberListDone: Starting...");
        setFetching(false);
    }

    // --- Checkout Session --- //

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

    // --- Avatar Operation --- //

    function onOperationSuccess(response) {
        // For debugging purposes only.
        console.log("onOperationSuccess: Starting...");
        console.log(response);

        // Add a temporary banner message in the app and then clear itself after 2 seconds.
        setTopAlertMessage("Photo changed");
        setTopAlertStatus("success");
        setTimeout(() => {
            console.log("onOperationSuccess: Delayed for 2 seconds.");
            console.log("onOperationSuccess: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
            setTopAlertMessage("");
        }, 2000);

        // Redirect the user to the user attachments page.
        setForceURL("/account/more");
    }

    function onOperationError(apiErr) {
        console.log("onOperationError: Starting...");
        setErrors(apiErr);

        // Add a temporary banner message in the app and then clear itself after 2 seconds.
        setTopAlertMessage("Failed submitting");
        setTopAlertStatus("danger");
        setTimeout(() => {
            console.log("onOperationError: Delayed for 2 seconds.");
            console.log("onOperationError: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
            setTopAlertMessage("");
        }, 2000);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onOperationDone() {
        console.log("onOperationDone: Starting...");
        setFetching(false);
    }

    // --- Checkout Session --- //

    function onCreateStripeSubscriptionCheckoutSessionSuccess(response) {
        console.log("onCreateStripeSubscriptionCheckoutSessionSuccess: Starting...");
        console.log("onCreateStripeSubscriptionCheckoutSessionSuccess: checkoutUrl:", response.checkoutUrl);
        document.location.href = response.checkoutUrl;
    }

    function onCreateStripeSubscriptionCheckoutSessionError(apiErr) {
        console.log("onCreateStripeSubscriptionCheckoutSessionError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onCreateStripeSubscriptionCheckoutSessionDone() {
        console.log("onCreateStripeSubscriptionCheckoutSessionDone: Starting...");
        setFetching(false);
    }

    // --- Account Detail --- //

    function onAccountSuccess(response) {
        console.log("onAccountSuccess: Starting...");
        setCurrentUser(response);
    }

    function onAccountError(apiErr) {
        console.log("onAccountError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onAccountDone() {
        console.log("onAccountDone: Starting...");
        setFetching(false);
    }

    ////
    //// BREADCRUMB
    ////
    const breadcrumbItems = {
        items: [
            { text: 'Dashboard', link: '/dashboard', isActive: false, icon: faGauge },
            { text: 'Account', link: '/account', icon: faUserCircle, isActive: false },
            { text: 'Subscribe', link: '#', icon: faCreditCard, isActive: true }
        ],
        mobileBackLinkItems: {
            link: "/account",
            text: "Back to Account",
            icon: faArrowLeft
        }
    }

    // --- All --- //

    const onUnauthorized = () => {
        setForceURL("/login?unauthorized=true"); // If token expired or user is not logged in, redirect back to login.
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
            setFetching(true);
            getAccountDetailAPI(
                onAccountSuccess,
                onAccountError,
                onAccountDone,
                onUnauthorized
            );
            fetchList();
        }

        return () => { mounted = false; }
    }, [,]);
    ////
    //// Component rendering.
    ////

    if (forceURL !== "") {
        return <Navigate to={forceURL} />
    }

    return (
        <Layout breadcrumbItems={breadcrumbItems}>
            {/* Page Title */}
            <h1 className="title is-2"><FontAwesomeIcon className="fas" icon={faUserCircle} />&nbsp;Account</h1>
            <h4 className="subtitle is-4"><FontAwesomeIcon className="fas" icon={faEye} />&nbsp;Detail</h4>
            <hr />

            {/* Page */}
            <div className="box">

                {/* Title + Options */}
                <p className="title is-4"><FontAwesomeIcon className="fas" icon={faCreditCard} />&nbsp;Subscribe</p>

                <FormErrorBox errors={errors} />

                {/* <p className="pb-4 has-text-grey">Please fill out all the required fields before submitting this form.</p> */}

                {isFetching
                    ?
                    <PageLoadingContent displayMessage={"Submitting..."} />
                    :
                    <>
                        <div className="container">

                            <div className="columns is-multiline">
                                {listData.results && listData.results.length !== 0 ? (
                                    listData.results.map((content, index) => (
                                        <div className="column is-half is-flex">
                                            <PricingTable
                                                key={index}
                                                title={content.name}
                                                highlightedText={content.recommended ? "Recommended" : null}
                                                customClasses={content.recommended ? "highlighted-border" : ""}
                                                price={content.price}
                                                period={null}
                                                description={content.description}
                                                url={content.thumbnailUrl ? content.thumbnailUrl : "https://images.pexels.com/photos/841130/pexels-photo-841130.jpeg?cs=srgb&dl=pexels-victor-freitas-841130.jpg&fm=jpg"}
                                                isLocked={false}
                                                customComponent={
                                                    <CustomComponent
                                                        onClick={onOfferClick}
                                                        id={content.id}
                                                        isLocked={false}
                                                    />}
                                            />
                                        </div>
                                    ))
                                ) : (
                                    <div className="column is-size-6 subtitle has-background-light is-fullwidth has-text-centered">
                                        {`No offers found. Please check back later.`}
                                    </div>
                                )}
                            </div>


                            <div className="columns pt-5">
                                <div className="column is-half">
                                    <Link to={`/account`} className="button is-fullwidth-mobile"><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to Account</Link>
                                </div>
                                <div className="column is-half has-text-right">
                                    {/*
                                            <button className="button is-medium is-success is-fullwidth-mobile" onClick={onSubmitClick}><FontAwesomeIcon className="fas" icon={faCheckCircle} />&nbsp;Save</button>
                                            */}
                                </div>
                            </div>

                        </div>
                    </>
                }
            </div>
        </Layout>
    );
}


const PricingTableList = ({ data, title, subtitle, onClick }) => {
    return (
        <>
            <h2 className="title is-size-5 has-text-centered">{title}</h2>
            <p className="subtitle is-size-6 has-text-centered">{subtitle}</p>
            <div className="columns is-multiline">
                {data && data.length !== 0 ? (
                    data.map((content, index) => (
                        <div className="column is-one-third is-flex">
                            <PricingTable
                                key={index}
                                title={content.name}
                                highlightedText={content.recommended ? "Recommended" : null}
                                customClasses={content.recommended ? "highlighted-border" : ""}
                                price={content.price}
                                period={null}
                                description={content.description}
                                url={content.thumbnailUrl ? content.thumbnailUrl : "https://images.pexels.com/photos/841130/pexels-photo-841130.jpeg?cs=srgb&dl=pexels-victor-freitas-841130.jpg&fm=jpg"}
                                isLocked={!content.currentUserHasAccessGranted}
                                customComponent={
                                    <CustomComponent
                                        onClick={onClick}
                                        id={content.id}
                                        isLocked={!content.currentUserHasAccessGranted}
                                    />}
                            />
                        </div>
                    ))
                ) : (
                    <div className="column is-size-6 subtitle has-background-light is-fullwidth has-text-centered">
                        {`No ${title.toLowerCase()} found. Please check back later.`}
                    </div>
                )}
            </div>
        </>
    );
};

const CustomComponent = ({ onClick, id, isLocked }) => {
    return (
        <div className="buttons is-centered mt-5">
            <Link
                onClick={(e, p) => onClick(id)}
                className="button is-primary is-medium"
                type="button"
                disabled={isLocked}
            >
                <span>Buy Now</span>
            </Link>
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
        <figure class="image is-5by5">
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

export default AccountMoreOperationSubscribe;
