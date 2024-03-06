import React, { useEffect, useState } from "react";
import { postLogoutAPI } from "../../API/gateway";
import Scroll from 'react-scroll';
// import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
// import { faSpinner } from '@fortawesome/free-solid-svg-icons';


function LogoutRedirector() {

    ////
    //// Component states.
    ////

    const [errors, setErrors] = useState({});

    ////
    //// API.
    ////

    function onLogoutnSuccess(response){
        console.log("onLogoutnSuccess: Starting...");
    }

    function onLogoutnError(apiErr) {
        console.log("onLogoutnError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onLogoutnDone() {
        console.log("onLogoutnDone: Starting...");


        setTimeout(onRedirect, 250);
    }

    ////
    //// Event handling.
    ////

    const onRedirect = (e) => {
        // Clear the entire local storage.
        localStorage.clear();

        // Do not use `Link` but instead use the `window.location` change
        // to fix the issue with the `TopNavigation` component to restart.
        // If you use use `Link` then when you redirect to the navigation then
        // the menu will not update.
        window.location.href = "/login";
    }

    ////
    //// Misc.
    ////

    useEffect(() => {
        let mounted = true;

        if (mounted) {
            postLogoutAPI(
                onLogoutnSuccess,
                onLogoutnError,
                onLogoutnDone
            );

            // DEFENSIVE CODE: In case we stay too long on this page,
            // automatically redirect the user regardless of the results
            // from the API.
            setTimeout(onRedirect, 750);
        }

        return () => mounted = false;
    }, []);

    return (
        <div class="container column is-12">
            <div class="section">
                <section class="hero is-fullheight">
                    <div class="hero-body">

                        <div class="container">
                            <div class="columns is-centered">
                                <div class="column is-one-third-tablet">
                                    <h1 className="is-size-1">LOGGING OFF...</h1>
                                </div>
                            </div>
                        </div>

                    </div>
                </section>
            </div>
        </div>
    );
}

export default LogoutRedirector;
